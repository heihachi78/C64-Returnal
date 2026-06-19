package game

import (
	"golang.org/x/image/font"
	"image/color"
	"math"
	"os"
	"slices"
	"testing"
)

func TestInitialWindowSizeMatchesOriginalApp(t *testing.T) {
	if ScreenWidth != 800 || ScreenHeight != 600 {
		t.Fatalf("initial screen size = %dx%d, want 800x600", ScreenWidth, ScreenHeight)
	}
}

func TestHUDFontSizesMatchOriginalLabels(t *testing.T) {
	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{name: "status", got: statusFontSize, want: 16},
		{name: "combat", got: combatFontSize, want: 14},
		{name: "level title", got: levelUpTitleFontSize, want: 40},
		{name: "level option", got: levelUpOptionFontSize, want: 22},
		{name: "level key", got: levelUpKeyFontSize, want: 18},
		{name: "chest title", got: chestTitleFontSize, want: 34},
		{name: "chest item", got: chestItemFontSize, want: 20},
		{name: "chest continue", got: chestContinueFontSize, want: 18},
		{name: "game over title", got: gameOverTitleFontSize, want: 40},
		{name: "game over option", got: gameOverOptionFontSize, want: 22},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("%s font size = %v, want %v", tt.name, tt.got, tt.want)
		}
	}

	if got := fontFaceForSize(levelUpTitleFontSize).Metrics().Height.Ceil(); got < 36 {
		t.Fatalf("level title font metric height = %d, want at least 36", got)
	}
}

func TestHUDFontPrefersOriginalMenloBoldWhenAvailable(t *testing.T) {
	if !fontNameMatches("Menlo Bold", []string{"menlo", "bold"}) {
		t.Fatal("fontNameMatches rejected Menlo Bold")
	}
	if fontNameMatches("Menlo Regular", []string{"menlo", "bold"}) {
		t.Fatal("fontNameMatches accepted Menlo Regular")
	}

	const menloPath = "/System/Library/Fonts/Menlo.ttc"
	if _, err := os.Stat(menloPath); err != nil {
		t.Skip("Menlo.ttc is not available on this platform")
	}

	font, name := loadSystemFontByFullName([]string{menloPath}, []string{"menlo", "bold"})
	if font == nil {
		t.Fatalf("could not load Menlo Bold from %s", menloPath)
	}
	if !fontNameMatches(name, []string{"menlo", "bold"}) {
		t.Fatalf("loaded HUD font name = %q, want Menlo Bold", name)
	}
}

func TestHUDIntervalFormattingMatchesOriginalStatusPanel(t *testing.T) {
	if got := formattedSeconds(3); got != "3.0" {
		t.Fatalf("formatted seconds for whole interval = %q, want 3.0", got)
	}
	if got := formattedSeconds(0.91); got != "0.91" {
		t.Fatalf("formatted seconds below one = %q, want 0.91", got)
	}
	if got := formattedSeconds(0.955); got != "0.95" {
		t.Fatalf("formatted seconds rounding below one = %q, want original 0.95", got)
	}
}

func TestGrassTintBlendFactorMatchesOriginalField(t *testing.T) {
	if math.Abs(grassTintBlendFactor-0.22) > 0.0001 {
		t.Fatalf("grassTintBlendFactor = %v, want 0.22", grassTintBlendFactor)
	}
}

func TestGrassGridMatchesOriginalInfiniteFieldLayout(t *testing.T) {
	startColumn, startRow, columns, rows := grassGrid(ScreenWidth, ScreenHeight, 64, Vec2{})
	if columns != 17 || rows != 14 {
		t.Fatalf("grass grid size = %dx%d, want 17x14", columns, rows)
	}
	if startColumn != -8 || startRow != -7 {
		t.Fatalf("grass grid start = (%d,%d), want (-8,-7)", startColumn, startRow)
	}

	startColumn, startRow, columns, rows = grassGrid(ScreenWidth, ScreenHeight, 64, Vec2{X: 130, Y: -130})
	if columns != 17 || rows != 14 {
		t.Fatalf("shifted grass grid size = %dx%d, want 17x14", columns, rows)
	}
	if startColumn != -6 || startRow != -10 {
		t.Fatalf("shifted grass grid start = (%d,%d), want (-6,-10)", startColumn, startRow)
	}
}

func TestGrassHashMatchesOriginalIntMinGuard(t *testing.T) {
	g := New()
	minInt := -int(^uint(0)>>1) - 1
	if got := g.grassHash(0, 0, minInt); got != 0 {
		t.Fatalf("grassHash with Int.min-equivalent salt = %d, want 0", got)
	}
}

func TestSpriteCullingMatchesOriginalNonVisibleNodeIntent(t *testing.T) {
	if spriteBoundsVisible(800, 600, -20, 300, 30, 42, 0) {
		t.Fatal("fully offscreen sprite was visible, want culled")
	}
	if !spriteBoundsVisible(800, 600, -14, 300, 30, 42, 0) {
		t.Fatal("edge-overlapping sprite was culled, want visible")
	}
	if !spriteBoundsVisible(800, 600, 808, 300, 24, 24, math.Pi/4) {
		t.Fatal("rotated edge-overlapping sprite was culled, want conservative visible")
	}
	if spriteBoundsVisible(800, 600, 900, 300, 24, 24, math.Pi/4) {
		t.Fatal("far offscreen rotated sprite was visible, want culled")
	}
}

func TestWorldRenderLayerOrderMatchesOriginalZPositions(t *testing.T) {
	want := []float64{-20, 8.5, 8.75, 9, 9.5, 10, 11, 12, 13, 14}
	got := worldRenderLayerOrder()
	if !slices.Equal(got, want) {
		t.Fatalf("world render layer order = %v, want original z order %v", got, want)
	}
	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Fatalf("world render layer order is not strictly ascending at %d: %v", i, got)
		}
	}
}

func TestWorldRotationConvertsToEbitenScreenSpaceLikeOriginalVisuals(t *testing.T) {
	clockwiseDeathRotation := -math.Pi / 2
	if got, want := worldRotationToScreen(clockwiseDeathRotation), math.Pi/2; math.Abs(got-want) > 0.0001 {
		t.Fatalf("screen death rotation = %v, want %v", got, want)
	}

	upwardProjectileRotation := math.Pi / 2
	if got, want := worldRotationToScreen(upwardProjectileRotation), -math.Pi/2; math.Abs(got-want) > 0.0001 {
		t.Fatalf("screen projectile rotation = %v, want %v", got, want)
	}
}

func TestPanelCornerRadiusMatchesOriginalHUD(t *testing.T) {
	if math.Abs(panelCornerRadius-6) > 0.0001 {
		t.Fatalf("panelCornerRadius = %v, want 6", panelCornerRadius)
	}
}

func TestScaledTextUsesOriginalLabelScaleInsteadOfScaledFontSize(t *testing.T) {
	baseFace := fontFaceForSize(levelUpTitleFontSize)
	scaledFace := fontFaceForSize(levelUpTitleFontSize * modalTitleScale(0))
	text := "LEVEL 3"
	baseWidth := font.MeasureString(baseFace, text).Ceil()
	scaledWidth := font.MeasureString(scaledFace, text).Ceil()

	layout := baseScaledTextLayout(baseFace, text, true)

	if layout.Width != baseWidth+8 {
		t.Fatalf("scaled text backing width = %d, want base label width %d plus padding", layout.Width, baseWidth)
	}
	if layout.Width == scaledWidth+8 {
		t.Fatalf("scaled text backing width matched scaled font width %d; want scaled base label image", scaledWidth)
	}
	if math.Abs(layout.AnchorX-(4+float64(baseWidth)/2)) > 0.0001 {
		t.Fatalf("centered scaled text anchor X = %v, want label center", layout.AnchorX)
	}
	leftLayout := baseScaledTextLayout(baseFace, text, false)
	if leftLayout.AnchorX != 4 {
		t.Fatalf("left scaled text anchor X = %v, want left label origin plus padding", leftLayout.AnchorX)
	}

	g := New()
	first := g.scaledTextImage(text, levelUpTitleFontSize, true)
	second := g.scaledTextImage(text, levelUpTitleFontSize, true)
	left := g.scaledTextImage(text, levelUpTitleFontSize, false)
	if first.Image == nil || second.Image == nil {
		t.Fatal("scaled text cache returned nil image")
	}
	if first.Image != second.Image {
		t.Fatal("scaled text cache did not reuse the base label image")
	}
	if first.Image == left.Image {
		t.Fatal("scaled text cache reused centered label image for left-aligned label")
	}
}

func TestGameOverLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := gameOverPanelRect(800, 600)
	if x != 90 || y != 210 || w != 620 || h != 190 {
		t.Fatalf("game over panel rect = (%v,%v,%v,%v), want (90,210,620,190)", x, y, w, h)
	}
	if y+gameOverTitleOffsetY != 258 {
		t.Fatalf("game over title y = %v, want 258", y+gameOverTitleOffsetY)
	}
	if y+gameOverRestartOffsetY != 322 || y+gameOverExitOffsetY != 370 {
		t.Fatalf("game over option y values = %v,%v; want 322,370", y+gameOverRestartOffsetY, y+gameOverExitOffsetY)
	}

	g := New()
	if got := g.gameOverOptionAt(400, 322); got != "restart" {
		t.Fatalf("game over option at restart center = %q, want restart", got)
	}
	if got := g.gameOverOptionAt(400, 370); got != "exit" {
		t.Fatalf("game over option at exit center = %q, want exit", got)
	}

	face := fontFaceForSize(gameOverOptionFontSize)
	restartOutsideX := 400 + float64(font.MeasureString(face, "RESTART").Ceil())/2 + 27
	if got := g.gameOverOptionAt(restartOutsideX, 322); got != "" {
		t.Fatalf("game over option outside restart label hit area = %q, want none like original behavior", got)
	}
}

func TestHUDStatusPanelLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := topStatusPanelRect(3)
	if x != 8 || y != 9 || w != 210 || h != 104 {
		t.Fatalf("three-life top panel rect = (%v,%v,%v,%v), want (8,9,210,104)", x, y, w, h)
	}
	x, y, w, h = topStatusPanelRect(13)
	if x != 8 || y != 9 || w != 210 || h != 120 {
		t.Fatalf("thirteen-life top panel rect = (%v,%v,%v,%v), want (8,9,210,120)", x, y, w, h)
	}

	x, y = lifeIconScreenPosition(0)
	if x != 25 || y != 98 {
		t.Fatalf("first life icon position = (%v,%v), want (25,98)", x, y)
	}
	x, y = lifeIconScreenPosition(11)
	if x != 201 || y != 98 {
		t.Fatalf("twelfth life icon position = (%v,%v), want (201,98)", x, y)
	}
	x, y = lifeIconScreenPosition(12)
	if x != 25 || y != 114 {
		t.Fatalf("thirteenth life icon position = (%v,%v), want (25,114)", x, y)
	}

	x, y, w, h = combatStatusPanelRect(600)
	if x != 8 || y != 267 || w != 176 || h != 330 {
		t.Fatalf("combat panel rect = (%v,%v,%v,%v), want (8,267,176,330)", x, y, w, h)
	}
}

func TestChestRewardLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := chestOverlayPanelRect(800, 600, 1)
	if x != 90 || y != 196 || w != 620 || h != 208 {
		t.Fatalf("one-item chest panel rect = (%v,%v,%v,%v), want (90,196,620,208)", x, y, w, h)
	}
	if chestOverlayTitleY(y) != 250 {
		t.Fatalf("one-item chest title y = %v, want 250", chestOverlayTitleY(y))
	}
	if chestOverlayContinueY(y, h) != 368 {
		t.Fatalf("one-item chest continue y = %v, want 368", chestOverlayContinueY(y, h))
	}
	if chestRewardItemY(y, h, 1, 0) != 312 {
		t.Fatalf("one-item chest item y = %v, want 312", chestRewardItemY(y, h, 1, 0))
	}

	x, y, w, h = chestOverlayPanelRect(800, 600, 2)
	if x != 90 || y != 179 || w != 620 || h != 242 {
		t.Fatalf("two-item chest panel rect = (%v,%v,%v,%v), want (90,179,620,242)", x, y, w, h)
	}
	if chestOverlayTitleY(y) != 233 {
		t.Fatalf("two-item chest title y = %v, want 233", chestOverlayTitleY(y))
	}
	if chestOverlayContinueY(y, h) != 385 {
		t.Fatalf("two-item chest continue y = %v, want 385", chestOverlayContinueY(y, h))
	}
	if chestRewardItemY(y, h, 2, 0) != 297 || chestRewardItemY(y, h, 2, 1) != 327 {
		t.Fatalf("two-item chest item y values = %v,%v; want 297,327", chestRewardItemY(y, h, 2, 0), chestRewardItemY(y, h, 2, 1))
	}
}

func TestSkeletonTintBlendFactorsMatchOriginalValues(t *testing.T) {
	tests := []struct {
		kind       SkeletonKind
		wantColor  [3]uint8
		wantFactor float64
	}{
		{kind: SkeletonRegular, wantColor: [3]uint8{255, 255, 255}, wantFactor: 0},
		{kind: SkeletonRed, wantColor: [3]uint8{242, 13, 10}, wantFactor: 0.68},
		{kind: SkeletonPurple, wantColor: [3]uint8{148, 31, 242}, wantFactor: 0.72},
		{kind: SkeletonBlack, wantColor: [3]uint8{5, 5, 5}, wantFactor: 0.86},
	}

	for _, tt := range tests {
		color, factor := skeletonTintBlend(tt.kind)
		if color.R != tt.wantColor[0] || color.G != tt.wantColor[1] || color.B != tt.wantColor[2] {
			t.Fatalf("kind %v color = (%d,%d,%d), want %v", tt.kind, color.R, color.G, color.B, tt.wantColor)
		}
		if math.Abs(factor-tt.wantFactor) > 0.0001 {
			t.Fatalf("kind %v factor = %v, want %v", tt.kind, factor, tt.wantFactor)
		}
	}
}

func TestPlayerSpritePresentationMatchesOriginalDeathTint(t *testing.T) {
	presentation := playerSpritePresentation(Player{HitFlash: playerHitFlashDuration / 2}, false)
	if presentation.Tint.R != 255 || presentation.Tint.G != 255 || presentation.Tint.B != 255 {
		t.Fatalf("active hit-flash tint rgb = %+v, want white with alpha-only flash", presentation.Tint)
	}
	if presentation.BlendFactor != 0 || presentation.Rotation != 0 {
		t.Fatalf("active hit-flash presentation = %+v, want no color blend or rotation", presentation)
	}

	presentation = playerSpritePresentation(Player{HitFlash: playerHitFlashDuration, DeathRotation: -math.Pi / 2}, true)
	if presentation.Tint != (color.RGBA{217, 13, 20, 115}) {
		t.Fatalf("death tint = %+v, want original red color with 0.45 alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.65) > 0.0001 {
		t.Fatalf("death blend factor = %v, want 0.65", presentation.BlendFactor)
	}
	if math.Abs(presentation.Rotation-math.Pi/2) > 0.0001 {
		t.Fatalf("death screen rotation = %v, want pi/2", presentation.Rotation)
	}
}

func TestSkeletonSpritePresentationPreservesTintBlendDuringHitFlash(t *testing.T) {
	presentation := skeletonSpritePresentation(Skeleton{Kind: SkeletonPurple, HitFlash: skeletonDamageFlashDuration})
	if presentation.Tint.R != 148 || presentation.Tint.G != 31 || presentation.Tint.B != 242 || presentation.Tint.A != 255 {
		t.Fatalf("fresh purple hit-flash tint = %+v, want purple tint with full alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.72) > 0.0001 {
		t.Fatalf("purple blend factor = %v, want 0.72", presentation.BlendFactor)
	}

	presentation = skeletonSpritePresentation(Skeleton{Kind: SkeletonPurple, HitFlash: skeletonDamageFlashDuration - 0.06})
	if presentation.Tint.R != 148 || presentation.Tint.G != 31 || presentation.Tint.B != 242 || presentation.Tint.A != 89 {
		t.Fatalf("dim purple hit-flash tint = %+v, want purple tint with 0.35 alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.72) > 0.0001 {
		t.Fatalf("dim purple blend factor = %v, want 0.72", presentation.BlendFactor)
	}
}

func TestLockedCombatRowsKeepKillCountsLikeOriginalHUD(t *testing.T) {
	first, second, kills := combatRowLabels("x0", "3.0s", "KILLS 0", false)
	if first != "LOCKED" || second != "--" || kills != "KILLS 0" {
		t.Fatalf("locked row labels = %q, %q, %q; want LOCKED, --, KILLS 0", first, second, kills)
	}
	row := combatRowPresentation("x0", "3.0s", "KILLS 0", false)
	if row.Tint != (color.RGBA{255, 255, 255, 255}) {
		t.Fatalf("locked row icon tint = %+v, want unchanged white", row.Tint)
	}
	if row.TextColor != c64Text {
		t.Fatalf("locked row text color = %+v, want primary text color %+v", row.TextColor, c64Text)
	}
}

func TestCombatRowTextOffsetsMatchOriginalHUDLayout(t *testing.T) {
	if combatRowFirstOffsetY != -10 {
		t.Fatalf("first row offset = %v, want -10", combatRowFirstOffsetY)
	}
	if combatRowSecondOffsetY != 6 {
		t.Fatalf("second row offset = %v, want 6", combatRowSecondOffsetY)
	}
	if combatRowKillsOffsetY != 20 {
		t.Fatalf("kills row offset = %v, want 20", combatRowKillsOffsetY)
	}
}
