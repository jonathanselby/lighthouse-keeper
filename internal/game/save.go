package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const SaveVersion = 2
const MaxSaveRead = 256 * 1024
const MaxSaveOutput = 200 * 1024

type BookEntry struct {
	Kind  Kind `json:"kind"`
	Count int  `json:"count"`
}
type NightRecord struct {
	Date       string `json:"date"`
	Prediction int    `json:"prediction"`
	Result     int    `json:"result"`
	GoodGuess  bool   `json:"good_guess"`
	CameBacks  int    `json:"came_backs"`
	Noticed    int    `json:"noticed"`
	CheckIns   []bool `json:"check_ins"`
}
type DeskRecord struct {
	Date            string `json:"date"`
	LengthMinutes   int    `json:"length_minutes"`
	Answers         []bool `json:"answers"`
	UnansweredAtEnd bool   `json:"unanswered_at_end"`
	EndedEarly      bool   `json:"ended_early"`
}
type Save struct {
	Version         int           `json:"version"`
	Book            []BookEntry   `json:"book"`
	NightsCompleted int           `json:"nights_completed"`
	LastDate        string        `json:"last_date"`
	NightsToday     int           `json:"nights_today"`
	ExtraNightUsed  bool          `json:"extra_night_used"`
	Log             []NightRecord `json:"log"`
	Desk            []DeskRecord  `json:"desk"`
}

func NewSave() Save {
	s := Save{Version: SaveVersion, Book: make([]BookEntry, int(KindCount))}
	for i := range s.Book {
		s.Book[i] = BookEntry{Kind: Kind(i)}
	}
	return s
}
func isDate(s string) bool {
	if s == "" {
		return true
	}
	if len(s) != 10 || s[4] != '-' || s[7] != '-' {
		return false
	}
	t, e := time.Parse("2006-01-02", s)
	return e == nil && t.Format("2006-01-02") == s
}
func (s Save) Validate() error {
	if s.Version != SaveVersion {
		return fmt.Errorf("save version %d", s.Version)
	}
	if len(s.Book) != int(KindCount) {
		return errors.New("book must have six entries")
	}
	for i, b := range s.Book {
		if b.Kind != Kind(i) || b.Count < 0 || b.Count > 100000000 {
			return fmt.Errorf("invalid book entry %d", i)
		}
	}
	if s.NightsCompleted < 0 || s.NightsCompleted > 100000000 || s.NightsToday < 0 || s.NightsToday > 3 || !isDate(s.LastDate) {
		return errors.New("invalid progression or date")
	}
	if s.ExtraNightUsed && s.LastDate == "" {
		return errors.New("override without date")
	}
	if len(s.Log) > 400 || len(s.Desk) > 400 {
		return errors.New("history beyond limit")
	}
	for _, r := range s.Log {
		if !isDate(r.Date) || r.Prediction < 0 || r.Prediction > 5 || r.Result < 0 || r.Result > 5 || r.GoodGuess != (absInt(r.Prediction-r.Result) <= 1) || r.CameBacks < 0 || r.CameBacks > 1000000 || r.Noticed < 0 || r.Noticed > 1000000 || len(r.CheckIns) > 3 {
			return errors.New("invalid night record")
		}
	}
	for _, r := range s.Desk {
		if !isDate(r.Date) || (r.LengthMinutes != 15 && r.LengthMinutes != 25) || len(r.Answers) > 9 {
			return errors.New("invalid desk record")
		}
	}
	return nil
}
func DecodeSave(data []byte) (Save, error) {
	if len(data) > MaxSaveRead {
		return Save{}, errors.New("save exceeds 256 KiB")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var s Save
	if err := dec.Decode(&s); err != nil {
		return Save{}, err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return Save{}, errors.New("trailing JSON")
		}
		return Save{}, err
	}
	if err := s.Validate(); err != nil {
		return Save{}, err
	}
	return s, nil
}
func SavePath() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", errors.New("user configuration directory unavailable")
	}
	return filepath.Join(home, "lighthouse-keeper", "save.json"), nil
}

// LoadSave preserves malformed saves in the first available .bad location.
// It deliberately returns a new playable save plus an error for the adapter to report.
func LoadSave(path string) (Save, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewSave(), nil
	}
	if err != nil {
		return NewSave(), err
	}
	data, readErr := io.ReadAll(io.LimitReader(f, MaxSaveRead+1))
	closeErr := f.Close()
	if readErr == nil {
		readErr = closeErr
	}
	var s Save
	if readErr == nil {
		s, readErr = DecodeSave(data)
	}
	if readErr == nil {
		return s, nil
	}
	// Never overwrite an earlier recovery copy.
	for i := 0; i < 1000; i++ {
		dst := path + ".bad"
		if i > 0 {
			dst = fmt.Sprintf("%s.bad.%d", path, i)
		}
		_, statErr := os.Lstat(dst)
		if errors.Is(statErr, os.ErrNotExist) {
			if renameErr := os.Rename(path, dst); renameErr != nil {
				return NewSave(), fmt.Errorf("%w (recovery copy unavailable: %v)", readErr, renameErr)
			}
			return NewSave(), fmt.Errorf("old journal preserved at %s: %w", dst, readErr)
		}
		if statErr != nil {
			return NewSave(), fmt.Errorf("journal recovery check failed: %w", statErr)
		}
	}
	return NewSave(), fmt.Errorf("%w: too many recovery copies", readErr)
}
func (s *Save) capHistory() {
	if len(s.Log) > 400 {
		s.Log = s.Log[len(s.Log)-400:]
	}
	if len(s.Desk) > 400 {
		s.Desk = s.Desk[len(s.Desk)-400:]
	}
}
func (s *Save) bytesForWrite() ([]byte, error) {
	s.capHistory()
	for {
		if err := s.Validate(); err != nil {
			return nil, err
		}
		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return nil, err
		}
		if len(data) <= MaxSaveOutput {
			return append(data, '\n'), nil
		}
		// Keep the latest record from each log before considering further trimming.
		if len(s.Log) <= 1 && len(s.Desk) <= 1 {
			return nil, errors.New("save exceeds bound despite history trimming")
		}
		if len(s.Log) > len(s.Desk) && len(s.Log) > 1 {
			s.Log = s.Log[1:]
		} else if len(s.Desk) > 1 {
			s.Desk = s.Desk[1:]
		} else if len(s.Log) > 1 {
			s.Log = s.Log[1:]
		}
	}
}
func WriteSave(path string, s *Save) error {
	data, err := s.bytesForWrite()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".save-*.tmp")
	if err != nil {
		return err
	}
	temporary := f.Name()
	defer os.Remove(temporary)
	if err = f.Chmod(0600); err != nil {
		f.Close()
		return err
	}
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = os.Rename(temporary, path); err != nil {
		return err
	}
	if d, e := os.Open(dir); e == nil {
		_ = d.Sync()
		_ = d.Close()
	}
	return nil
}
