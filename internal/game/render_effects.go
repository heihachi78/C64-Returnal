package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

func (g *Game) drawEffect(screen *ebiten.Image, effect Effect) {
	alpha := effectFadeAlpha(effect.TTL, effect.MaxTTL)
	switch effect.Kind {
	case EffectLightning:
		innerPoints := effect.InnerPoints
		if len(innerPoints) == 0 {
			innerPoints = effect.Points
		}
		g.drawBolt(screen, effect.Points, g.tuning.LightningBranchWidth+5, color.RGBA{33, 163, 255, alpha / 2})
		g.drawBolt(screen, effect.Points, g.tuning.LightningBranchWidth, color.RGBA{33, 163, 255, alpha})
		g.drawBolt(screen, innerPoints, 2, color.RGBA{255, 255, 255, alpha / 2})
		g.drawBolt(screen, innerPoints, 1, color.RGBA{255, 255, 255, alpha})
		endX, endY := g.worldToScreen(effect.End)
		g.drawSpriteScreen(screen, g.assets.Lightning, endX, endY, 24, 24, false, color.RGBA{255, 255, 255, alpha})
	case EffectLightningHit:
		x, y := g.worldToScreen(effect.Pos)
		hitAlpha := lightningHitEffectAlpha(effect.TTL, effect.MaxTTL)
		g.drawSpriteRotatedBlend(
			screen,
			g.assets.Skeleton[effect.Frame%len(g.assets.Skeleton)],
			x,
			y,
			30,
			42,
			0,
			effect.Facing < 0,
			color.RGBA{89, 219, 255, hitAlpha},
			0.8,
		)
	case EffectBeam:
		startX, startY := g.worldToScreen(effect.Start)
		endX, endY := g.worldToScreen(effect.End)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 16, color.RGBA{255, 184, 20, alpha / 4}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 9, color.RGBA{255, 184, 20, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 7, color.RGBA{255, 240, 56, alpha / 3}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 4, color.RGBA{255, 240, 56, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 1, color.RGBA{255, 255, 255, alpha}, false)
	case EffectMeteorImpact:
		x, y := g.worldToScreen(effect.Pos)
		style := meteorImpactRenderStyle(effect)
		vector.DrawFilledCircle(screen, float32(x), float32(y), float32(style.Radius), color.RGBA{74, 43, 20, scaleAlpha(115, style.Alpha)}, false)
		vector.DrawFilledCircle(screen, float32(x), float32(y), float32(style.CoreRadius), color.RGBA{156, 97, 43, scaleAlpha(166, style.Alpha)}, false)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(style.Radius), float32(style.GlowWidth), color.RGBA{194, 128, 61, scaleAlpha(64, style.Alpha)}, false)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(style.Radius), float32(style.StrokeWidth), color.RGBA{194, 128, 61, scaleAlpha(217, style.Alpha)}, false)
	}
}

type meteorImpactStyle struct {
	Scale float64
	Alpha float64
}
type meteorImpactRenderMetrics struct {
	Radius      float64
	CoreRadius  float64
	GlowWidth   float64
	StrokeWidth float64
	Alpha       float64
}

func meteorImpactPresentation(effect Effect) meteorImpactStyle {
	if effect.MaxTTL <= 0 {
		return meteorImpactStyle{Scale: 1.25, Alpha: 0}
	}
	age := Clamp(effect.MaxTTL-effect.TTL, 0, effect.MaxTTL)
	switch {
	case age < 0.08:
		return meteorImpactStyle{Scale: 0.25 + 0.75*(age/0.08), Alpha: 1}
	case age < 0.16:
		return meteorImpactStyle{Scale: 1, Alpha: 1}
	default:
		fade := Clamp((age-0.16)/0.16, 0, 1)
		return meteorImpactStyle{Scale: 1 + 0.25*fade, Alpha: 1 - fade}
	}
}
func meteorImpactRenderStyle(effect Effect) meteorImpactRenderMetrics {
	presentation := meteorImpactPresentation(effect)
	radius := effect.Radius * presentation.Scale
	return meteorImpactRenderMetrics{
		Radius:      radius,
		CoreRadius:  radius * 0.35,
		GlowWidth:   3 * presentation.Scale,
		StrokeWidth: 2 * presentation.Scale,
		Alpha:       presentation.Alpha,
	}
}
func scaleAlpha(base uint8, alpha float64) uint8 {
	return uint8(math.Round(float64(base) * Clamp(alpha, 0, 1)))
}
func flashActionAlpha(elapsed, total, fadeDown, fadeUp float64) uint8 {
	if elapsed < 0 || elapsed >= total || fadeDown <= 0 || fadeUp <= 0 {
		return 255
	}
	cycle := fadeDown + fadeUp
	phase := math.Mod(elapsed, cycle)
	const minAlpha = 0.35
	if phase < fadeDown {
		progress := phase / fadeDown
		return uint8(math.Round(255 * (1 - (1-minAlpha)*progress)))
	}
	progress := (phase - fadeDown) / fadeUp
	return uint8(math.Round(255 * (minAlpha + (1-minAlpha)*progress)))
}
func coinFloatOffset(phase float64) float64 {
	return 5 * linearPingPong(phase, 0.42)
}
func coinShimmerAlpha(phase float64) uint8 {
	alpha := 1 - 0.28*linearPingPong(phase, 0.28)
	return uint8(math.Round(255 * alpha))
}
func linearPingPong(phase, halfPeriod float64) float64 {
	if halfPeriod <= 0 {
		return 0
	}
	t := math.Mod(phase, halfPeriod*2)
	if t < 0 {
		t += halfPeriod * 2
	}
	if t <= halfPeriod {
		return t / halfPeriod
	}
	return 1 - (t-halfPeriod)/halfPeriod
}
func (g *Game) drawBolt(screen *ebiten.Image, points []Vec2, width float32, clr color.RGBA) {
	if len(points) < 2 {
		return
	}
	for i := 1; i < len(points); i++ {
		startX, startY := g.worldToScreen(points[i-1])
		endX, endY := g.worldToScreen(points[i])
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), width, clr, false)
	}
}
