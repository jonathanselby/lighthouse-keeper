package game

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveRoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "save.json")
	s := NewSave()
	s.Book[Whale].Count = 12
	s.Log = append(s.Log, NightRecord{Date: "2026-11-11", Prediction: 3, Result: 3, GoodGuess: true})
	if err := WriteSave(p, &s); err != nil {
		t.Fatal(err)
	}
	got, err := LoadSave(p)
	if err != nil {
		t.Fatal(err)
	}
	if got.Book[Whale].Count != 12 || len(got.Log) != 1 {
		t.Fatal(got)
	}
	info, _ := os.Stat(p)
	if info.Mode().Perm() != 0600 {
		t.Fatal(info.Mode())
	}
}
func TestMissingSaveStartsEmpty(t *testing.T) {
	s, err := LoadSave(filepath.Join(t.TempDir(), "nothing.json"))
	if err != nil || len(s.Book) != 6 {
		t.Fatal(s, err)
	}
}
func TestMalformedSaveReportedAndKept(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.json")
	os.WriteFile(path, []byte("not json"), 0600)
	s, err := LoadSave(path)
	if err == nil || len(s.Book) != 6 {
		t.Fatal(err)
	}
	b, e := os.ReadFile(path + ".bad")
	if e != nil || string(b) != "not json" {
		t.Fatal(e)
	}
}
func TestUnknownFieldsRejected(t *testing.T) {
	s := NewSave()
	data, _ := json.Marshal(s)
	data = bytes.Replace(data, []byte(`"version":2`), []byte(`"version":2,"unexpected":true`), 1)
	if _, err := DecodeSave(data); err == nil {
		t.Fatal("unknown accepted")
	}
}
func TestOversizedSaveRejectedBeforeDecode(t *testing.T) {
	if _, err := DecodeSave([]byte(strings.Repeat("x", MaxSaveRead+1))); err == nil {
		t.Fatal("oversized accepted")
	}
}
func TestInvalidSaveFieldsRejected(t *testing.T) {
	s := NewSave()
	s.NightsToday = -2
	b, _ := json.Marshal(s)
	if _, err := DecodeSave(b); err == nil {
		t.Fatal("invalid fields")
	}
	s = NewSave()
	s.Book = append(s.Book, BookEntry{Kind: Whale})
	b, _ = json.Marshal(s)
	if _, err := DecodeSave(b); err == nil {
		t.Fatal("extra book kind")
	}
	s = NewSave()
	s.LastDate = "2026-02-31"
	b, _ = json.Marshal(s)
	if _, err := DecodeSave(b); err == nil {
		t.Fatal("invalid date")
	}
}
func TestAtomicSaveInterruptionKeepsOldFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "save.json")
	s := NewSave()
	s.Book[Whale].Count = 5
	if err := WriteSave(p, &s); err != nil {
		t.Fatal(err)
	}
	// An invalid replacement must be rejected before the original file is touched.
	s.Book[Whale].Count = -10
	if err := WriteSave(p, &s); err == nil {
		t.Fatal("invalid save written")
	}
	restored, err := LoadSave(p)
	if err != nil || restored.Book[Whale].Count != 5 {
		t.Fatal(restored, err)
	}
}
func TestLogCaps(t *testing.T) {
	s := NewSave()
	for i := 0; i < 500; i++ {
		s.Log = append(s.Log, NightRecord{Prediction: 3, Result: 3, GoodGuess: true})
	}
	p := filepath.Join(t.TempDir(), "save.json")
	if err := WriteSave(p, &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Log) != 400 {
		t.Fatal(len(s.Log))
	}
}
func TestTrailingJSONRejected(t *testing.T) {
	b, _ := json.Marshal(NewSave())
	if _, e := DecodeSave(append(b, []byte("{}")...)); e == nil {
		t.Fatal("trailing JSON")
	}
}
