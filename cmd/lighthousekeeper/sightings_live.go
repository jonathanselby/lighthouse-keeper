package main

import (
 "image"
 "math"
 "github.com/hajimehoshi/ebiten/v2"
 "github.com/jonathanselby/lighthouse-keeper/internal/game"
)

// Presentation-only poses: the deterministic notice/chase rules stay in game.
func emergenceProgress(s game.Sighting) float64 {
 switch s.State {
 case game.Appearing:
  p:=math.Min(1,float64(s.Age)/48)
  return 1-math.Pow(1-p,3)
 case game.Leaving,game.Chased:
  p:=math.Min(1,float64(s.Age)/30)
  return 1-p*p*(3-2*p)
 default: return 1
 }
}
func (a *app) drawLiveSighting(dst *ebiten.Image,s game.Sighting,index int,sec float64) {
 im:=a.liveSightings[s.Kind]
 if im==nil {return}
 waterY:=sceneY(s.Pos.Y)
 progress:=emergenceProgress(s)
 wobble:=math.Sin(sec*1.8+float64(index)*2.2)*1.6
 alpha:=float32(math.Min(1,.20+progress*.80))
 if s.State==game.Noticed {a.imageAt(dst,a.images["glow.png"],s.Pos.X,waterY-18,.50)}
 if s.Kind==game.Whale||s.Kind==game.Seal||s.Kind==game.Jellyfish {
  // Source cropping is an occlusion plane: the creature rises from below
  // the waterline, and its sprite contains no baked spray or water.
  const scale=.76
  const waterlineInSprite=204.0
  ytop:=waterY*renderScale-waterlineInSprite*scale+(1-progress)*110+wobble*renderScale
  cut:=int(math.Floor((waterY*renderScale-ytop)/scale))
  if cut>240 {cut=240}
  if cut>0 {
   src:=im.SubImage(image.Rect(0,0,240,cut)).(*ebiten.Image)
   op:=&ebiten.DrawImageOptions{Filter:ebiten.FilterLinear}
   op.GeoM.Scale(scale,scale)
   op.GeoM.Translate(s.Pos.X*renderScale-120*scale,ytop)
   op.ColorScale.ScaleAlpha(alpha)
   dst.DrawImage(src,op)
  }
  if progress>.14 {
   a.drawOneImpact(dst,impact{x:s.Pos.X,y:waterY+1,w:65,h:34,pose:int(s.Kind)%2,phase:.27+float64(index)*.13},sec)
  }
  return
 }
 op:=&ebiten.DrawImageOptions{Filter:ebiten.FilterLinear}
 op.GeoM.Translate(-120,-120)
 switch s.Kind {
 case game.ShootingStar:
  op.GeoM.Scale(.72,.72)
  op.GeoM.Rotate(-.04)
  op.GeoM.Translate((s.Pos.X+18*progress)*renderScale,(waterY-35)*renderScale)
 case game.Owl:
  op.GeoM.Scale(.68,.68)
  op.GeoM.Translate(s.Pos.X*renderScale,(waterY-28+wobble*.45)*renderScale)
 case game.Ferry:
  op.GeoM.Scale(.72,.72)
  op.GeoM.Translate((s.Pos.X-14+14*progress)*renderScale,(waterY-35+wobble)*renderScale)
  if progress>.1 {a.drawOneImpact(dst,impact{x:s.Pos.X-28,y:waterY+1,w:75,h:26,pose:0,phase:.53},sec)}
 }
 op.ColorScale.ScaleAlpha(alpha)
 dst.DrawImage(im,op)
}
