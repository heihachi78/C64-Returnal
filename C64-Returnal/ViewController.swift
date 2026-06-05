//
//  ViewController.swift
//  C64-Returnal
//
//  Created by Tóth István on 2026. 06. 02..
//

import Cocoa
import SpriteKit

class ViewController: NSViewController {

    @IBOutlet var skView: SKView!

    override func viewDidLoad() {
        super.viewDidLoad()

        let scene = GameScene(size: skView.bounds.size)
        scene.scaleMode = .resizeFill

        skView.presentScene(scene)
        skView.ignoresSiblingOrder = true

        #if DEBUG
        skView.showsFPS = true
        skView.showsNodeCount = true
        #endif
    }
}
