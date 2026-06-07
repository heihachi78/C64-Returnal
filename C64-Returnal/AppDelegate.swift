import Cocoa
import SpriteKit

@main
final class AppDelegate: NSObject, NSApplicationDelegate {
    private static var sharedDelegate: AppDelegate?

    private var window: NSWindow?

    static func main() {
        let app = NSApplication.shared
        let delegate = AppDelegate()

        sharedDelegate = delegate
        app.delegate = delegate
        app.setActivationPolicy(.regular)
        app.mainMenu = makeMainMenu()
        app.finishLaunching()
        delegate.showGameWindow()
        app.run()
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        showGameWindow()
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        true
    }

    private func showGameWindow() {
        if let window {
            bringGameWindowForward(window)
            return
        }

        let frame = NSRect(x: 0, y: 0, width: 800, height: 600)
        let skView = SKView(frame: frame)
        let scene = GameScene(size: frame.size)

        scene.scaleMode = .resizeFill
        skView.presentScene(scene)
        skView.ignoresSiblingOrder = true

        let window = NSWindow(
            contentRect: frame,
            styleMask: [.titled, .closable, .miniaturizable, .resizable],
            backing: .buffered,
            defer: false
        )
        window.title = "C64-Returnal"
        window.center()
        window.contentView = skView
        window.isRestorable = false
        window.isReleasedWhenClosed = false
        window.makeKeyAndOrderFront(nil)
        window.makeFirstResponder(skView)

        self.window = window
        bringGameWindowForward(window)
    }

    private func bringGameWindowForward(_ window: NSWindow) {
        NSApp.setActivationPolicy(.regular)
        NSApp.activate(ignoringOtherApps: true)
        NSRunningApplication.current.activate(options: [.activateAllWindows])
        window.makeKeyAndOrderFront(nil)
        window.orderFrontRegardless()

        DispatchQueue.main.async {
            NSApp.activate(ignoringOtherApps: true)
            NSRunningApplication.current.activate(options: [.activateAllWindows])
            window.makeKeyAndOrderFront(nil)
            window.orderFrontRegardless()
            _ = window.makeFirstResponder(window.contentView)
        }
    }

    private static func makeMainMenu() -> NSMenu {
        let mainMenu = NSMenu()
        let appMenuItem = NSMenuItem()
        let appMenu = NSMenu(title: "C64-Returnal")

        appMenu.addItem(
            withTitle: "Quit C64-Returnal",
            action: #selector(NSApplication.terminate(_:)),
            keyEquivalent: "q"
        )
        appMenuItem.submenu = appMenu
        mainMenu.addItem(appMenuItem)

        return mainMenu
    }
}
