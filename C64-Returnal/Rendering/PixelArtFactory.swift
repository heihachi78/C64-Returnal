//
//  PixelArtFactory.swift
//  C64-Returnal
//

import SpriteKit

enum PixelArtFactory {
    static func makeGrassTextures(tileSize: CGFloat) -> [SKTexture] {
        (0..<8).map { variant in
            nearestTexture(makeGrassTexture(size: Int(tileSize), variant: variant))
        }
    }

    static func makeMageTextures() -> [SKTexture] {
        (0..<2).map { variant in
            nearestTexture(makeMageImageTexture(variant: variant))
        }
    }

    static func makeMageTexture() -> SKTexture {
        makeMageTextures()[0]
    }

    static func makeSkeletonTextures() -> [SKTexture] {
        (0..<2).map { variant in
            nearestTexture(makeSkeletonImageTexture(variant: variant))
        }
    }

    static func makeSkeletonTexture() -> SKTexture {
        makeSkeletonTextures()[0]
    }

    static func makeFireballTextures() -> [SKTexture] {
        (0..<2).map { variant in
            nearestTexture(makeFireballTexture(variant: variant))
        }
    }

    static func makeLightningTexture() -> SKTexture {
        nearestTexture(makeLightningImageTexture())
    }

    static func makeOrbitalOrbTextures() -> [SKTexture] {
        (0..<2).map { variant in
            nearestTexture(makeOrbitalOrbTexture(variant: variant))
        }
    }

    static func makeBeamTexture() -> SKTexture {
        nearestTexture(makeBeamImageTexture())
    }

    static func makeMeteorTextures() -> [SKTexture] {
        (0..<2).map { variant in
            nearestTexture(makeMeteorTexture(variant: variant))
        }
    }

    static func makeLifeTexture() -> SKTexture {
        nearestTexture(makeLifeImageTexture())
    }

    static func makeChestTexture(tier: ChestTier) -> SKTexture {
        nearestTexture(makeChestImageTexture(tier: tier))
    }

    static func makeCoinTextures() -> [SKTexture] {
        (0..<4).map { variant in
            nearestTexture(makeCoinImageTexture(variant: variant))
        }
    }

    private static func nearestTexture(_ texture: SKTexture) -> SKTexture {
        texture.filteringMode = .nearest
        return texture
    }

    private static func makeGrassTexture(size: Int, variant: Int) -> SKTexture {
        let image = NSImage(size: CGSize(width: size, height: size))

        image.lockFocus()
        grassColor(red: 0.22, green: 0.47, blue: 0.17, variant: variant, shift: 0).setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: size, height: size)).fill()

        drawGrassGroundFlecks(size: size, variant: variant)
        drawGrassDetails(size: size, variant: variant)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func drawGrassGroundFlecks(size: Int, variant: Int) {
        let flecks = [
            (4, 7, 3, 1, 0), (13, 23, 2, 1, 1), (27, 9, 4, 1, 2), (39, 18, 2, 1, 3),
            (52, 5, 3, 1, 4), (7, 39, 2, 1, 5), (20, 53, 3, 1, 6), (34, 41, 4, 1, 7),
            (49, 55, 2, 1, 8), (58, 29, 3, 1, 9), (2, 58, 2, 1, 10), (45, 34, 3, 1, 11)
        ]

        for (index, fleck) in flecks.enumerated() {
            let color = fleck.4.isMultiple(of: 2)
                ? grassColor(red: 0.30, green: 0.56, blue: 0.20, variant: variant, shift: index)
                : grassColor(red: 0.16, green: 0.36, blue: 0.14, variant: variant, shift: index)
            drawPixelRect(
                x: wrappedPixel(fleck.0 + grassOffset(variant: variant, index: index, axis: 0), size: size),
                y: wrappedPixel(fleck.1 + grassOffset(variant: variant, index: index, axis: 1), size: size),
                width: fleck.2,
                height: fleck.3,
                color: color
            )
        }
    }

    private static func drawGrassDetails(size: Int, variant: Int) {
        let darkBlade = grassColor(red: 0.10, green: 0.29, blue: 0.11, variant: variant, shift: 17)
        let midBlade = grassColor(red: 0.28, green: 0.61, blue: 0.20, variant: variant, shift: 19)
        let lightBlade = grassColor(red: 0.50, green: 0.82, blue: 0.31, variant: variant, shift: 23)
        let rootColor = grassColor(red: 0.08, green: 0.24, blue: 0.10, variant: variant, shift: 29)

        let tufts = [
            (7, 9, 8), (22, 5, 11), (43, 8, 9), (56, 17, 10),
            (13, 29, 12), (33, 25, 8), (49, 36, 13), (5, 49, 9),
            (25, 51, 10), (39, 55, 7), (59, 53, 11)
        ]

        for (index, tuft) in tufts.enumerated() {
            let baseX = wrappedPixel(tuft.0 + grassOffset(variant: variant, index: index, axis: 4), size: size)
            let baseY = wrappedPixel(tuft.1 + grassOffset(variant: variant, index: index, axis: 5), size: size)
            let height = max(5, tuft.2 + grassOffset(variant: variant, index: index, axis: 6) / 2)

            drawGrassTuft(
                baseX: baseX,
                baseY: baseY,
                height: height,
                darkColor: darkBlade,
                midColor: midBlade,
                lightColor: lightBlade,
                rootColor: rootColor
            )
        }

        let singleBlades = [
            (4, 21, 6, 1), (17, 42, 5, -1), (30, 15, 7, 0), (41, 44, 6, 1),
            (54, 28, 5, -1), (61, 6, 6, 0), (9, 60, 4, 1), (31, 36, 5, -1)
        ]

        for (index, blade) in singleBlades.enumerated() {
            let x = wrappedPixel(blade.0 + grassOffset(variant: variant, index: index, axis: 7), size: size)
            let y = wrappedPixel(blade.1 + grassOffset(variant: variant, index: index, axis: 8), size: size)
            let height = max(3, blade.2 + grassOffset(variant: variant, index: index, axis: 9) / 2)
            drawGrassBlade(baseX: x, baseY: y, height: height, lean: blade.3, color: lightBlade, baseWidth: 1)
        }
    }

    private static func drawGrassTuft(
        baseX: Int,
        baseY: Int,
        height: Int,
        darkColor: NSColor,
        midColor: NSColor,
        lightColor: NSColor,
        rootColor: NSColor
    ) {
        drawPixelRect(x: baseX - 2, y: baseY, width: 6, height: 2, color: rootColor)
        drawPixelRect(x: baseX - 1, y: baseY + 1, width: 4, height: 1, color: darkColor)

        drawGrassBlade(baseX: baseX - 2, baseY: baseY + 1, height: height - 2, lean: -1, color: darkColor, baseWidth: 1)
        drawGrassBlade(baseX: baseX, baseY: baseY + 1, height: height, lean: 0, color: midColor, baseWidth: 2)
        drawGrassBlade(baseX: baseX + 2, baseY: baseY + 1, height: height - 1, lean: 1, color: midColor, baseWidth: 1)
        drawGrassBlade(baseX: baseX + 1, baseY: baseY + 2, height: max(3, height - 4), lean: -1, color: lightColor, baseWidth: 1)

        drawPixelRect(x: baseX, y: baseY + height + 1, width: 1, height: 1, color: lightColor)
    }

    private static func drawGrassBlade(baseX: Int, baseY: Int, height: Int, lean: Int, color: NSColor, baseWidth: Int) {
        guard height > 0 else {
            return
        }

        for step in 0..<height {
            let leanOffset = lean == 0 ? 0 : (step * lean) / 3
            let width = step < 2 ? baseWidth : 1
            drawPixelRect(x: baseX + leanOffset, y: baseY + step, width: width, height: 1, color: color)
        }
    }

    private static func makeMageImageTexture(variant: Int) -> SKTexture {
        let pixelSize = 16
        let image = NSImage(size: CGSize(width: pixelSize, height: 22))
        let leftFootX = variant == 0 ? 4 : 3
        let rightFootX = variant == 0 ? 9 : 10
        let staffX = variant == 0 ? 13 : 14
        let crystalX = variant == 0 ? 12 : 13
        let crystalY = variant == 0 ? 17 : 16

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: 22)).fill()

        drawPixelRect(x: 5, y: 20, width: 6, height: 1, color: NSColor(calibratedRed: 0.13, green: 0.07, blue: 0.24, alpha: 1))
        drawPixelRect(x: 6, y: 19, width: 4, height: 1, color: NSColor(calibratedRed: 0.38, green: 0.18, blue: 0.65, alpha: 1))
        drawPixelRect(x: 4, y: 17, width: 8, height: 2, color: NSColor(calibratedRed: 0.27, green: 0.12, blue: 0.49, alpha: 1))
        drawPixelRect(x: 3, y: 16, width: 10, height: 1, color: NSColor(calibratedRed: 0.13, green: 0.07, blue: 0.24, alpha: 1))
        drawPixelRect(x: 4, y: 11, width: 8, height: 6, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
        drawPixelRect(x: 6, y: 12, width: 4, height: 4, color: NSColor(calibratedRed: 0.88, green: 0.65, blue: 0.47, alpha: 1))
        drawPixelRect(x: 7, y: 14, width: 1, height: 1, color: NSColor(calibratedRed: 0.12, green: 0.08, blue: 0.08, alpha: 1))
        drawPixelRect(x: 10, y: 14, width: 1, height: 1, color: NSColor(calibratedRed: 0.12, green: 0.08, blue: 0.08, alpha: 1))

        drawPixelRect(x: 4, y: 4, width: 8, height: 8, color: NSColor(calibratedRed: 0.14, green: 0.24, blue: 0.66, alpha: 1))
        drawPixelRect(x: 3, y: 3, width: 10, height: 2, color: NSColor(calibratedRed: 0.09, green: 0.15, blue: 0.39, alpha: 1))
        drawPixelRect(x: 5, y: 5, width: 2, height: 6, color: NSColor(calibratedRed: 0.22, green: 0.35, blue: 0.83, alpha: 1))
        drawPixelRect(x: 9, y: 5, width: 2, height: 6, color: NSColor(calibratedRed: 0.09, green: 0.15, blue: 0.39, alpha: 1))
        if variant == 0 {
            drawPixelRect(x: 2, y: 6, width: 3, height: 3, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
            drawPixelRect(x: 11, y: 7, width: 2, height: 2, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
        } else {
            drawPixelRect(x: 3, y: 7, width: 2, height: 2, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
            drawPixelRect(x: 11, y: 5, width: 3, height: 3, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
        }
        drawPixelRect(x: leftFootX, y: 1, width: 3, height: 2, color: NSColor(calibratedRed: 0.08, green: 0.07, blue: 0.16, alpha: 1))
        drawPixelRect(x: rightFootX, y: 1, width: 3, height: 2, color: NSColor(calibratedRed: 0.08, green: 0.07, blue: 0.16, alpha: 1))

        drawPixelRect(x: staffX, y: 4, width: 1, height: 13, color: NSColor(calibratedRed: 0.37, green: 0.20, blue: 0.09, alpha: 1))
        drawPixelRect(x: crystalX, y: crystalY, width: 3, height: 3, color: NSColor(calibratedRed: 0.32, green: 0.86, blue: 0.95, alpha: 1))
        drawPixelRect(x: crystalX + 1, y: crystalY + 1, width: 1, height: 1, color: NSColor.white)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeSkeletonImageTexture(variant: Int) -> SKTexture {
        let pixelSize = 16
        let image = NSImage(size: CGSize(width: pixelSize, height: 22))
        let bone = NSColor(calibratedRed: 0.82, green: 0.84, blue: 0.76, alpha: 1)
        let brightBone = NSColor(calibratedRed: 0.94, green: 0.96, blue: 0.86, alpha: 1)
        let shadowBone = NSColor(calibratedRed: 0.48, green: 0.51, blue: 0.46, alpha: 1)
        let dark = NSColor(calibratedRed: 0.06, green: 0.07, blue: 0.08, alpha: 1)
        let rust = NSColor(calibratedRed: 0.42, green: 0.28, blue: 0.13, alpha: 1)

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: 22)).fill()

        drawPixelRect(x: 5, y: 14, width: 7, height: 6, color: bone)
        drawPixelRect(x: 6, y: 18, width: 5, height: 2, color: brightBone)
        drawPixelRect(x: 5, y: 14, width: 1, height: 4, color: shadowBone)
        drawPixelRect(x: 7, y: 17, width: 1, height: 1, color: dark)
        drawPixelRect(x: 10, y: 17, width: 1, height: 1, color: dark)
        drawPixelRect(x: 9, y: 15, width: 1, height: 1, color: dark)

        drawPixelRect(x: 7, y: 11, width: 3, height: 3, color: bone)
        drawPixelRect(x: 5, y: 8, width: 7, height: 4, color: bone)
        drawPixelRect(x: 6, y: 9, width: 5, height: 1, color: dark)
        drawPixelRect(x: 5, y: 7, width: 1, height: 4, color: shadowBone)
        drawPixelRect(x: 11, y: 7, width: 1, height: 4, color: shadowBone)

        if variant == 0 {
            drawPixelRect(x: 3, y: 8, width: 2, height: 1, color: bone)
            drawPixelRect(x: 2, y: 5, width: 1, height: 4, color: bone)
            drawPixelRect(x: 12, y: 8, width: 2, height: 1, color: bone)
            drawPixelRect(x: 13, y: 5, width: 1, height: 4, color: bone)

            drawPixelRect(x: 6, y: 4, width: 2, height: 4, color: bone)
            drawPixelRect(x: 10, y: 4, width: 2, height: 4, color: bone)
            drawPixelRect(x: 5, y: 2, width: 3, height: 2, color: bone)
            drawPixelRect(x: 10, y: 2, width: 3, height: 2, color: bone)
            drawPixelRect(x: 5, y: 1, width: 4, height: 1, color: shadowBone)
            drawPixelRect(x: 10, y: 1, width: 4, height: 1, color: shadowBone)
        } else {
            drawPixelRect(x: 3, y: 7, width: 2, height: 1, color: bone)
            drawPixelRect(x: 2, y: 4, width: 1, height: 4, color: bone)
            drawPixelRect(x: 12, y: 9, width: 2, height: 1, color: bone)
            drawPixelRect(x: 14, y: 6, width: 1, height: 4, color: bone)

            drawPixelRect(x: 5, y: 4, width: 2, height: 4, color: bone)
            drawPixelRect(x: 11, y: 4, width: 2, height: 4, color: bone)
            drawPixelRect(x: 4, y: 2, width: 3, height: 2, color: bone)
            drawPixelRect(x: 11, y: 2, width: 3, height: 2, color: bone)
            drawPixelRect(x: 4, y: 1, width: 4, height: 1, color: shadowBone)
            drawPixelRect(x: 11, y: 1, width: 4, height: 1, color: shadowBone)
        }

        drawPixelRect(x: 13, y: 7, width: 1, height: 8, color: rust)
        drawPixelRect(x: 12, y: 14, width: 3, height: 1, color: rust)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeFireballTexture(variant: Int) -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let red = NSColor(calibratedRed: 0.92, green: 0.08, blue: 0.03, alpha: 1)
        let darkRed = NSColor(calibratedRed: 0.52, green: 0.03, blue: 0.02, alpha: 1)
        let orange = NSColor(calibratedRed: 1.00, green: 0.40, blue: 0.03, alpha: 1)
        let yellow = NSColor(calibratedRed: 1.00, green: 0.88, blue: 0.18, alpha: 1)
        let core = variant == 0 ? yellow : orange
        let flame = variant == 0 ? orange : yellow

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 1, y: 4, width: 2, height: 4, color: darkRed)
        drawPixelRect(x: 2, y: 2, width: 7, height: 8, color: red)
        drawPixelRect(x: 4, y: 3, width: 5, height: 6, color: flame)
        drawPixelRect(x: 5, y: 4, width: 3, height: 4, color: core)
        drawPixelRect(x: 9, y: 5, width: 2, height: 3, color: flame)

        if variant == 0 {
            drawPixelRect(x: 3, y: 9, width: 3, height: 1, color: yellow)
            drawPixelRect(x: 2, y: 1, width: 2, height: 1, color: orange)
        } else {
            drawPixelRect(x: 7, y: 9, width: 2, height: 1, color: yellow)
            drawPixelRect(x: 1, y: 2, width: 2, height: 1, color: orange)
        }

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeLightningImageTexture() -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let white = NSColor.white
        let blue = NSColor(calibratedRed: 0.18, green: 0.82, blue: 1.0, alpha: 1)
        let darkBlue = NSColor(calibratedRed: 0.04, green: 0.25, blue: 0.68, alpha: 1)

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 6, y: 10, width: 3, height: 2, color: white)
        drawPixelRect(x: 5, y: 8, width: 4, height: 2, color: blue)
        drawPixelRect(x: 4, y: 6, width: 3, height: 2, color: white)
        drawPixelRect(x: 5, y: 4, width: 3, height: 2, color: blue)
        drawPixelRect(x: 3, y: 2, width: 3, height: 2, color: white)
        drawPixelRect(x: 2, y: 0, width: 2, height: 2, color: blue)
        drawPixelRect(x: 8, y: 8, width: 2, height: 2, color: darkBlue)
        drawPixelRect(x: 7, y: 3, width: 2, height: 2, color: darkBlue)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeOrbitalOrbTexture(variant: Int) -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let deepPurple = NSColor(calibratedRed: 0.20, green: 0.04, blue: 0.42, alpha: 1)
        let purple = NSColor(calibratedRed: 0.48, green: 0.11, blue: 0.78, alpha: 1)
        let brightPurple = NSColor(calibratedRed: 0.78, green: 0.30, blue: 1.0, alpha: 1)
        let core = NSColor(calibratedRed: 0.96, green: 0.78, blue: 1.0, alpha: 1)
        let halo = variant == 0 ? brightPurple : purple
        let inner = variant == 0 ? purple : brightPurple

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 4, y: 1, width: 4, height: 1, color: halo)
        drawPixelRect(x: 2, y: 3, width: 8, height: 6, color: deepPurple)
        drawPixelRect(x: 3, y: 2, width: 6, height: 8, color: purple)
        drawPixelRect(x: 4, y: 3, width: 5, height: 6, color: inner)
        drawPixelRect(x: 5, y: 4, width: 3, height: 4, color: core)

        if variant == 0 {
            drawPixelRect(x: 2, y: 8, width: 2, height: 2, color: brightPurple)
            drawPixelRect(x: 8, y: 2, width: 2, height: 2, color: brightPurple)
        } else {
            drawPixelRect(x: 1, y: 5, width: 2, height: 2, color: brightPurple)
            drawPixelRect(x: 7, y: 9, width: 2, height: 2, color: brightPurple)
        }

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeBeamImageTexture() -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let gold = NSColor(calibratedRed: 1.0, green: 0.72, blue: 0.08, alpha: 1)
        let yellow = NSColor(calibratedRed: 1.0, green: 0.94, blue: 0.22, alpha: 1)
        let white = NSColor.white

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 1, y: 4, width: 10, height: 4, color: gold)
        drawPixelRect(x: 0, y: 5, width: 12, height: 2, color: yellow)
        drawPixelRect(x: 3, y: 5, width: 6, height: 2, color: white)
        drawPixelRect(x: 2, y: 8, width: 2, height: 1, color: yellow)
        drawPixelRect(x: 8, y: 3, width: 2, height: 1, color: yellow)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeMeteorTexture(variant: Int) -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let darkBrown = NSColor(calibratedRed: 0.22, green: 0.12, blue: 0.06, alpha: 1)
        let brown = NSColor(calibratedRed: 0.42, green: 0.24, blue: 0.11, alpha: 1)
        let warmBrown = NSColor(calibratedRed: 0.58, green: 0.36, blue: 0.17, alpha: 1)
        let highlight = NSColor(calibratedRed: 0.76, green: 0.53, blue: 0.28, alpha: 1)
        let shadow = variant == 0 ? darkBrown : brown
        let mid = variant == 0 ? brown : warmBrown

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 4, y: 10, width: 4, height: 1, color: shadow)
        drawPixelRect(x: 2, y: 8, width: 8, height: 2, color: darkBrown)
        drawPixelRect(x: 1, y: 4, width: 10, height: 4, color: brown)
        drawPixelRect(x: 3, y: 2, width: 7, height: 2, color: shadow)
        drawPixelRect(x: 4, y: 5, width: 5, height: 3, color: mid)
        drawPixelRect(x: 5, y: 7, width: 3, height: 2, color: highlight)

        if variant == 0 {
            drawPixelRect(x: 2, y: 6, width: 2, height: 2, color: warmBrown)
            drawPixelRect(x: 8, y: 3, width: 2, height: 1, color: darkBrown)
        } else {
            drawPixelRect(x: 8, y: 6, width: 2, height: 2, color: warmBrown)
            drawPixelRect(x: 2, y: 3, width: 2, height: 1, color: darkBrown)
        }

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeLifeImageTexture() -> SKTexture {
        let pixelSize = 12
        let image = NSImage(size: CGSize(width: pixelSize, height: pixelSize))
        let red = NSColor(calibratedRed: 0.86, green: 0.05, blue: 0.10, alpha: 1)
        let brightRed = NSColor(calibratedRed: 1.0, green: 0.18, blue: 0.22, alpha: 1)
        let highlight = NSColor(calibratedRed: 1.0, green: 0.62, blue: 0.62, alpha: 1)
        let shadow = NSColor(calibratedRed: 0.45, green: 0.02, blue: 0.06, alpha: 1)

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: pixelSize, height: pixelSize)).fill()

        drawPixelRect(x: 2, y: 7, width: 3, height: 3, color: red)
        drawPixelRect(x: 7, y: 7, width: 3, height: 3, color: red)
        drawPixelRect(x: 1, y: 5, width: 10, height: 3, color: brightRed)
        drawPixelRect(x: 2, y: 3, width: 8, height: 2, color: red)
        drawPixelRect(x: 4, y: 1, width: 4, height: 2, color: shadow)
        drawPixelRect(x: 3, y: 8, width: 2, height: 1, color: highlight)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeChestImageTexture(tier: ChestTier) -> SKTexture {
        let image = NSImage(size: CGSize(width: 16, height: 14))
        let colors = chestColors(for: tier)
        let outline = NSColor(calibratedRed: 0.09, green: 0.05, blue: 0.03, alpha: 1)
        let darkStrap = NSColor(calibratedRed: 0.14, green: 0.08, blue: 0.04, alpha: 1)
        let lock = NSColor(calibratedRed: 1.0, green: 0.86, blue: 0.28, alpha: 1)

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: 16, height: 14)).fill()

        drawPixelRect(x: 3, y: 11, width: 10, height: 1, color: outline)
        drawPixelRect(x: 2, y: 9, width: 12, height: 2, color: colors.dark)
        drawPixelRect(x: 3, y: 10, width: 10, height: 1, color: colors.light)
        drawPixelRect(x: 1, y: 3, width: 14, height: 6, color: outline)
        drawPixelRect(x: 2, y: 4, width: 12, height: 5, color: colors.base)
        drawPixelRect(x: 2, y: 7, width: 12, height: 1, color: colors.light)
        drawPixelRect(x: 2, y: 4, width: 12, height: 1, color: colors.dark)
        drawPixelRect(x: 1, y: 2, width: 14, height: 1, color: outline)
        drawPixelRect(x: 7, y: 3, width: 2, height: 7, color: darkStrap)
        drawPixelRect(x: 6, y: 5, width: 4, height: 3, color: outline)
        drawPixelRect(x: 7, y: 5, width: 2, height: 2, color: lock)
        drawPixelRect(x: 3, y: 8, width: 3, height: 1, color: colors.light)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeCoinImageTexture(variant: Int) -> SKTexture {
        let image = NSImage(size: CGSize(width: 12, height: 12))
        let outline = NSColor(calibratedRed: 0.42, green: 0.24, blue: 0.03, alpha: 1)
        let dark = NSColor(calibratedRed: 0.72, green: 0.43, blue: 0.04, alpha: 1)
        let gold = NSColor(calibratedRed: 1.0, green: 0.73, blue: 0.08, alpha: 1)
        let yellow = NSColor(calibratedRed: 1.0, green: 0.92, blue: 0.24, alpha: 1)
        let white = NSColor(calibratedRed: 1.0, green: 0.98, blue: 0.72, alpha: 1)

        image.lockFocus()
        NSColor.clear.setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: 12, height: 12)).fill()

        switch variant {
        case 0:
            drawPixelRect(x: 3, y: 1, width: 6, height: 1, color: outline)
            drawPixelRect(x: 2, y: 2, width: 8, height: 1, color: outline)
            drawPixelRect(x: 1, y: 3, width: 10, height: 6, color: outline)
            drawPixelRect(x: 2, y: 9, width: 8, height: 1, color: outline)
            drawPixelRect(x: 3, y: 10, width: 6, height: 1, color: outline)
            drawPixelRect(x: 2, y: 3, width: 8, height: 6, color: gold)
            drawPixelRect(x: 3, y: 2, width: 6, height: 8, color: gold)
            drawPixelRect(x: 3, y: 7, width: 5, height: 2, color: yellow)
            drawPixelRect(x: 4, y: 4, width: 3, height: 2, color: dark)
            drawPixelRect(x: 4, y: 8, width: 2, height: 1, color: white)
        case 1, 3:
            drawPixelRect(x: 4, y: 1, width: 4, height: 1, color: outline)
            drawPixelRect(x: 3, y: 2, width: 6, height: 1, color: outline)
            drawPixelRect(x: 3, y: 3, width: 6, height: 6, color: outline)
            drawPixelRect(x: 3, y: 9, width: 6, height: 1, color: outline)
            drawPixelRect(x: 4, y: 10, width: 4, height: 1, color: outline)
            drawPixelRect(x: 4, y: 2, width: 4, height: 8, color: gold)
            drawPixelRect(x: 5, y: 3, width: 2, height: 6, color: yellow)
            drawPixelRect(x: 4, y: 4, width: 1, height: 4, color: dark)
            drawPixelRect(x: 5, y: 8, width: 1, height: 1, color: white)
        default:
            drawPixelRect(x: 5, y: 1, width: 2, height: 1, color: outline)
            drawPixelRect(x: 4, y: 2, width: 4, height: 1, color: outline)
            drawPixelRect(x: 4, y: 3, width: 4, height: 6, color: outline)
            drawPixelRect(x: 4, y: 9, width: 4, height: 1, color: outline)
            drawPixelRect(x: 5, y: 10, width: 2, height: 1, color: outline)
            drawPixelRect(x: 5, y: 2, width: 2, height: 8, color: gold)
            drawPixelRect(x: 6, y: 3, width: 1, height: 6, color: yellow)
            drawPixelRect(x: 5, y: 4, width: 1, height: 4, color: dark)
            drawPixelRect(x: 6, y: 8, width: 1, height: 1, color: white)
        }

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func chestColors(for tier: ChestTier) -> (base: NSColor, light: NSColor, dark: NSColor) {
        switch tier {
        case .bronze:
            return (
                NSColor(calibratedRed: 0.55, green: 0.30, blue: 0.13, alpha: 1),
                NSColor(calibratedRed: 0.86, green: 0.52, blue: 0.23, alpha: 1),
                NSColor(calibratedRed: 0.32, green: 0.16, blue: 0.07, alpha: 1)
            )
        case .silver:
            return (
                NSColor(calibratedRed: 0.58, green: 0.63, blue: 0.68, alpha: 1),
                NSColor(calibratedRed: 0.88, green: 0.93, blue: 0.96, alpha: 1),
                NSColor(calibratedRed: 0.34, green: 0.38, blue: 0.43, alpha: 1)
            )
        case .gold:
            return (
                NSColor(calibratedRed: 0.86, green: 0.58, blue: 0.08, alpha: 1),
                NSColor(calibratedRed: 1.0, green: 0.86, blue: 0.25, alpha: 1),
                NSColor(calibratedRed: 0.48, green: 0.30, blue: 0.03, alpha: 1)
            )
        }
    }

    private static func grassColor(red: CGFloat, green: CGFloat, blue: CGFloat, variant: Int, shift: Int) -> NSColor {
        let amount = CGFloat(((variant * 17 + shift * 11) % 9) - 4) * 0.012

        return NSColor(
            calibratedRed: clampedColor(red + amount * 0.7),
            green: clampedColor(green + amount),
            blue: clampedColor(blue + amount * 0.45),
            alpha: 1
        )
    }

    private static func clampedColor(_ value: CGFloat) -> CGFloat {
        min(1, max(0, value))
    }

    private static func grassOffset(variant: Int, index: Int, axis: Int) -> Int {
        let value = (variant * 37 + index * 19 + axis * 13) % 11
        return value - 5
    }

    private static func wrappedPixel(_ value: Int, size: Int) -> Int {
        ((value % size) + size) % size
    }

    private static func drawPixelRect(x: Int, y: Int, width: Int, height: Int, color: NSColor) {
        color.setFill()
        NSBezierPath(rect: NSRect(x: x, y: y, width: width, height: height)).fill()
    }
}
