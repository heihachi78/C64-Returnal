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

    static func makeMageTexture() -> SKTexture {
        nearestTexture(makeMageImageTexture())
    }

    static func makeSkeletonTexture() -> SKTexture {
        nearestTexture(makeSkeletonImageTexture())
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

    private static func nearestTexture(_ texture: SKTexture) -> SKTexture {
        texture.filteringMode = .nearest
        return texture
    }

    private static func makeGrassTexture(size: Int, variant: Int) -> SKTexture {
        let image = NSImage(size: CGSize(width: size, height: size))

        image.lockFocus()
        grassColor(red: 0.24, green: 0.50, blue: 0.18, variant: variant, shift: 0).setFill()
        NSBezierPath(rect: NSRect(x: 0, y: 0, width: size, height: size)).fill()

        let patchSeeds: [(Int, Int, Int, Int, CGFloat, CGFloat, CGFloat)] = [
            (3, 5, 19, 7, 0.34, 0.61, 0.22),
            (28, 2, 14, 13, 0.18, 0.43, 0.16),
            (48, 8, 10, 19, 0.39, 0.68, 0.25),
            (9, 28, 24, 10, 0.15, 0.38, 0.15),
            (36, 36, 21, 11, 0.31, 0.57, 0.19),
            (1, 50, 15, 9, 0.42, 0.70, 0.27),
            (22, 53, 26, 6, 0.20, 0.45, 0.17),
            (53, 49, 7, 10, 0.29, 0.54, 0.20)
        ]

        for (index, patch) in patchSeeds.enumerated() {
            grassColor(red: patch.4, green: patch.5, blue: patch.6, variant: variant, shift: index).setFill()
            NSBezierPath(
                rect: NSRect(
                    x: wrappedPixel(patch.0 + grassOffset(variant: variant, index: index, axis: 0), size: size),
                    y: wrappedPixel(patch.1 + grassOffset(variant: variant, index: index, axis: 1), size: size),
                    width: max(6, patch.2 + grassOffset(variant: variant, index: index, axis: 2) / 2),
                    height: max(4, patch.3 + grassOffset(variant: variant, index: index, axis: 3) / 2)
                )
            ).fill()
        }

        drawGrassDetails(size: size, variant: variant)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func drawGrassDetails(size: Int, variant: Int) {
        let bladeColor = grassColor(red: 0.49, green: 0.78, blue: 0.30, variant: variant, shift: 19)
        let shadowColor = grassColor(red: 0.11, green: 0.31, blue: 0.13, variant: variant, shift: 23)
        let blades = [
            (5, 18, 1, 5), (17, 43, 1, 4), (24, 18, 1, 6), (41, 24, 1, 5),
            (57, 33, 1, 4), (12, 60, 1, 3), (31, 47, 1, 4), (50, 4, 1, 5)
        ]
        let shadows = [
            (8, 12, 6, 2), (32, 21, 8, 2), (46, 58, 7, 2), (19, 7, 6, 2)
        ]

        bladeColor.setFill()
        for (index, blade) in blades.enumerated() {
            NSBezierPath(
                rect: NSRect(
                    x: wrappedPixel(blade.0 + grassOffset(variant: variant, index: index, axis: 4), size: size),
                    y: wrappedPixel(blade.1 + grassOffset(variant: variant, index: index, axis: 5), size: size),
                    width: blade.2,
                    height: max(2, blade.3 + grassOffset(variant: variant, index: index, axis: 6) / 2)
                )
            ).fill()
        }

        shadowColor.setFill()
        for (index, shadow) in shadows.enumerated() {
            NSBezierPath(
                rect: NSRect(
                    x: wrappedPixel(shadow.0 + grassOffset(variant: variant, index: index, axis: 7), size: size),
                    y: wrappedPixel(shadow.1 + grassOffset(variant: variant, index: index, axis: 8), size: size),
                    width: max(3, shadow.2 + grassOffset(variant: variant, index: index, axis: 9) / 2),
                    height: shadow.3
                )
            ).fill()
        }
    }

    private static func makeMageImageTexture() -> SKTexture {
        let pixelSize = 16
        let image = NSImage(size: CGSize(width: pixelSize, height: 22))

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
        drawPixelRect(x: 2, y: 6, width: 3, height: 3, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
        drawPixelRect(x: 11, y: 7, width: 2, height: 2, color: NSColor(calibratedRed: 0.32, green: 0.13, blue: 0.56, alpha: 1))
        drawPixelRect(x: 4, y: 1, width: 3, height: 2, color: NSColor(calibratedRed: 0.08, green: 0.07, blue: 0.16, alpha: 1))
        drawPixelRect(x: 9, y: 1, width: 3, height: 2, color: NSColor(calibratedRed: 0.08, green: 0.07, blue: 0.16, alpha: 1))

        drawPixelRect(x: 13, y: 4, width: 1, height: 13, color: NSColor(calibratedRed: 0.37, green: 0.20, blue: 0.09, alpha: 1))
        drawPixelRect(x: 12, y: 17, width: 3, height: 3, color: NSColor(calibratedRed: 0.32, green: 0.86, blue: 0.95, alpha: 1))
        drawPixelRect(x: 13, y: 18, width: 1, height: 1, color: NSColor.white)

        image.unlockFocus()
        return SKTexture(image: image)
    }

    private static func makeSkeletonImageTexture() -> SKTexture {
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
