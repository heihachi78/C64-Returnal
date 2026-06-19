package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

func (g *Game) drawSprite(screen, img *ebiten.Image, x, y, w, h float64, flipX bool, tint color.RGBA) {
	g.drawSpriteScreen(screen, img, x, y, w, h, flipX, tint)
}
func (g *Game) drawSpriteScreen(screen, img *ebiten.Image, x, y, w, h float64, flipX bool, tint color.RGBA) {
	g.drawSpriteRotated(screen, img, x, y, w, h, 0, flipX, tint)
}
func (g *Game) drawSpriteRotated(screen, img *ebiten.Image, x, y, w, h, rotation float64, flipX bool, tint color.RGBA) {
	g.drawSpriteRotatedBlend(screen, img, x, y, w, h, rotation, flipX, tint, 0)
}
func (g *Game) drawSpriteRotatedBlend(screen, img *ebiten.Image, x, y, w, h, rotation float64, flipX bool, tint color.RGBA, blendFactor float64) {
	if !spriteBoundsVisible(g.screenW, g.screenH, x, y, w, h, rotation) {
		return
	}
	bounds := img.Bounds()
	scaleX := w / float64(bounds.Dx())
	scaleY := h / float64(bounds.Dy())
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
	if flipX {
		op.GeoM.Scale(-scaleX, scaleY)
	} else {
		op.GeoM.Scale(scaleX, scaleY)
	}
	if rotation != 0 {
		op.GeoM.Rotate(rotation)
	}
	if blendFactor > 0 {
		op.ColorM.Scale(1-blendFactor, 1-blendFactor, 1-blendFactor, 1)
		op.ColorM.Translate(float64(tint.R)/255*blendFactor, float64(tint.G)/255*blendFactor, float64(tint.B)/255*blendFactor, 0)
		if tint.A != 255 {
			op.ColorScale.ScaleAlpha(float32(tint.A) / 255)
		}
	} else if tint != (color.RGBA{255, 255, 255, 255}) {
		op.ColorScale.ScaleWithColor(tint)
	}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
}
func spriteBoundsVisible(screenW, screenH int, x, y, w, h, rotation float64) bool {
	if screenW <= 0 || screenH <= 0 || w <= 0 || h <= 0 {
		return false
	}
	halfW := w / 2
	halfH := h / 2
	if rotation != 0 {
		radius := math.Hypot(w, h) / 2
		halfW = radius
		halfH = radius
	}
	return x+halfW >= 0 && x-halfW <= float64(screenW) &&
		y+halfH >= 0 && y-halfH <= float64(screenH)
}
func skeletonTintBlend(kind SkeletonKind) (color.RGBA, float64) {
	switch kind {
	case SkeletonRed:
		return color.RGBA{242, 13, 10, 255}, 0.68
	case SkeletonPurple:
		return color.RGBA{148, 31, 242, 255}, 0.72
	case SkeletonBlack:
		return color.RGBA{5, 5, 5, 255}, 0.86
	default:
		return color.RGBA{255, 255, 255, 255}, 0
	}
}
func (g *Game) panel(screen *ebiten.Image, x, y, w, h float64) {
	g.panelWithAlpha(screen, x, y, w, h, c64Panel.A)
}
func (g *Game) panelWithAlpha(screen *ebiten.Image, x, y, w, h float64, alpha uint8) {
	panelColor := c64Panel
	panelColor.A = alpha
	drawFilledRoundedRect(screen, x, y, w, h, panelCornerRadius, panelColor)
	if c64PanelEdge.A > 0 {
		vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, c64PanelEdge, false)
	}
}
func drawFilledRoundedRect(screen *ebiten.Image, x, y, w, h, radius float64, clr color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	r := math.Min(radius, math.Min(w/2, h/2))
	if r <= 0 {
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), clr, false)
		return
	}
	vector.DrawFilledRect(screen, float32(x+r), float32(y), float32(w-2*r), float32(h), clr, false)
	vector.DrawFilledRect(screen, float32(x), float32(y+r), float32(r), float32(h-2*r), clr, false)
	vector.DrawFilledRect(screen, float32(x+w-r), float32(y+r), float32(r), float32(h-2*r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+r), float32(y+r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+w-r), float32(y+r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+r), float32(y+h-r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+w-r), float32(y+h-r), float32(r), clr, false)
}
