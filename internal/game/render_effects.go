package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const lightningBoltStrokeScale = 0.25

func (g *Game) drawEffect(screen *ebiten.Image, effect Effect) {
	alpha := effectFadeAlpha(effect.TTL, effect.MaxTTL)
	switch effect.Kind {
	case EffectLightning:
		innerPoints := effect.InnerPoints
		if len(innerPoints) == 0 {
			innerPoints = effect.Points
		}
		g.drawBolt(screen, effect.Points, (g.tuning.LightningBranchWidth+5)*lightningBoltStrokeScale, color.RGBA{33, 163, 255, alpha / 2})
		g.drawBolt(screen, effect.Points, g.tuning.LightningBranchWidth*lightningBoltStrokeScale, color.RGBA{33, 163, 255, alpha})
		g.drawBolt(screen, innerPoints, 2*lightningBoltStrokeScale, color.RGBA{255, 255, 255, alpha / 2})
		g.drawBolt(screen, innerPoints, 1*lightningBoltStrokeScale, color.RGBA{255, 255, 255, alpha})
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
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 4, color.RGBA{255, 184, 20, alpha / 4}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 2.25, color.RGBA{255, 184, 20, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 1.75, color.RGBA{255, 240, 56, alpha / 3}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 1, color.RGBA{255, 240, 56, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 0.25, color.RGBA{255, 255, 255, alpha}, false)
	case EffectFireballImpact:
		g.drawFireballImpact(screen, effect)
	}
}

func (g *Game) drawFireballImpact(screen *ebiten.Image, effect Effect) {
	if effect.MaxTTL <= 0 {
		return
	}
	progress := Clamp((effect.MaxTTL-effect.TTL)/effect.MaxTTL, 0, 1)
	alpha := effectFadeAlpha(effect.TTL, effect.MaxTTL)
	const particles = 7
	for i := 0; i < particles; i++ {
		angle := float64(i)*math.Pi*2/particles + 0.28
		distance := 4 + 18*progress
		pos := Vec2{
			X: effect.Pos.X + math.Cos(angle)*distance,
			Y: effect.Pos.Y + math.Sin(angle)*distance,
		}
		x, y := g.worldToScreen(pos)
		radius := float32(math.Max(1, 3.2*(1-progress)))
		clr := color.RGBA{255, 145, 25, scaleAlpha(alpha, 0.85)}
		if i%3 == 0 {
			clr = color.RGBA{255, 238, 74, alpha}
		}
		vector.DrawFilledCircle(screen, float32(x), float32(y), radius, clr, false)
	}
}

func (g *Game) drawDeathWave(screen *ebiten.Image, wave DeathWave) {
	if wave.MaxRadius <= 0 || wave.Radius <= 0 {
		return
	}
	progress := Clamp(wave.Radius/wave.MaxRadius, 0, 1)
	alpha := uint8(math.Round(255 * (1 - progress) * 0.86))
	if alpha == 0 {
		return
	}
	x, y := g.worldToScreen(wave.Origin)
	width := float32(max(1, g.tuning.DeathWaveWidth*0.18))
	pulse := 0.5 + 0.5*math.Sin(wave.Radius*0.08+g.totalTime*12)
	outerAlpha := scaleAlpha(alpha, 0.28+0.18*pulse)
	echoAlpha := scaleAlpha(alpha, 0.32*(1-progress))
	rimAlpha := scaleAlpha(alpha, 0.96)

	vector.StrokeCircle(screen, float32(x), float32(y), float32(wave.Radius), width+9, color.RGBA{42, 5, 58, outerAlpha}, false)
	vector.StrokeCircle(screen, float32(x), float32(y), float32(wave.Radius), width+5, color.RGBA{116, 18, 156, scaleAlpha(alpha, 0.72)}, false)
	vector.StrokeCircle(screen, float32(x), float32(y), float32(wave.Radius), width+2, color.RGBA{226, 43, 217, rimAlpha}, false)
	vector.StrokeCircle(screen, float32(x), float32(y), float32(wave.Radius), max(1, width-1), color.RGBA{246, 243, 231, rimAlpha}, false)
	if wave.Radius > g.tuning.DeathWaveWidth {
		vector.StrokeCircle(screen, float32(x), float32(y), float32(wave.Radius-g.tuning.DeathWaveWidth*0.65), max(1, width-2), color.RGBA{71, 187, 255, echoAlpha}, false)
	}
	g.drawDeathWaveSparks(screen, wave, alpha)
}

func (g *Game) drawDeathWaveSparks(screen *ebiten.Image, wave DeathWave, alpha uint8) {
	if wave.Radius < 18 {
		return
	}
	const sparks = 18
	for i := 0; i < sparks; i++ {
		angle := float64(i)*math.Pi*2/sparks + math.Sin(wave.Radius*0.045+float64(i))*0.16
		length := 9 + 7*linearPingPong(wave.Radius*0.018+float64(i)*0.19, 0.5)
		radius := wave.Radius + math.Sin(wave.Radius*0.06+float64(i)*1.7)*g.tuning.DeathWaveWidth*0.28
		start := wave.Origin.Add(Vec2{X: math.Cos(angle) * (radius - length*0.35), Y: math.Sin(angle) * (radius - length*0.35)})
		end := wave.Origin.Add(Vec2{X: math.Cos(angle) * (radius + length), Y: math.Sin(angle) * (radius + length)})
		startX, startY := g.worldToScreen(start)
		endX, endY := g.worldToScreen(end)
		clr := color.RGBA{255, 72, 214, scaleAlpha(alpha, 0.62)}
		if i%3 == 0 {
			clr = color.RGBA{246, 243, 231, scaleAlpha(alpha, 0.78)}
		} else if i%3 == 1 {
			clr = color.RGBA{71, 187, 255, scaleAlpha(alpha, 0.48)}
		}
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 1.4, clr, false)
	}
}

type meteorImpactShake struct {
	Radius  float64
	OffsetX float64
	OffsetY float64
}

func (g *Game) drawGroundEffects(screen *ebiten.Image) {
	for _, effect := range g.effects {
		if effect.Kind == EffectMeteorImpact {
			g.drawMeteorImpactGroundShake(screen, effect)
		}
	}
}

func (g *Game) drawMeteorImpactGroundShake(screen *ebiten.Image, effect Effect) {
	shake := meteorImpactShakePresentation(effect)
	if shake.Radius <= 0 || (shake.OffsetX == 0 && shake.OffsetY == 0) {
		return
	}

	tile := g.tuning.TileSize
	startColumn := int(math.Floor((effect.Pos.X - shake.Radius) / tile))
	endColumn := int(math.Floor((effect.Pos.X + shake.Radius) / tile))
	startRow := int(math.Floor((effect.Pos.Y - shake.Radius) / tile))
	endRow := int(math.Floor((effect.Pos.Y + shake.Radius) / tile))

	for row := startRow; row <= endRow; row++ {
		for column := startColumn; column <= endColumn; column++ {
			tileCenter := Vec2{
				X: float64(column)*tile + tile/2,
				Y: float64(row)*tile + tile/2,
			}
			distance := tileCenter.Sub(effect.Pos).Len()
			falloff := Clamp(1-distance/(shake.Radius+tile/2), 0, 1)
			if falloff <= 0 {
				continue
			}
			g.drawGrassTile(screen, column, row, math.Round(shake.OffsetX*falloff), math.Round(shake.OffsetY*falloff))
		}
	}
}

func meteorImpactShakePresentation(effect Effect) meteorImpactShake {
	if effect.MaxTTL <= 0 || effect.Radius <= 0 {
		return meteorImpactShake{}
	}
	age := Clamp(effect.MaxTTL-effect.TTL, 0, effect.MaxTTL)
	progress := Clamp(age/effect.MaxTTL, 0, 1)
	amplitude := 4 * (1 - progress) * (1 - progress)
	return meteorImpactShake{
		Radius:  effect.Radius,
		OffsetX: math.Round(math.Sin(age*92) * amplitude),
		OffsetY: math.Round(math.Cos(age*73) * amplitude * 0.65),
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
