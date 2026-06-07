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
        if !isRunningUnitTests {
            app.setActivationPolicy(.regular)
            app.mainMenu = makeMainMenu()
        }
        app.finishLaunching()
        if !isRunningUnitTests {
            delegate.showGameWindow()
        }
        app.run()
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        guard !Self.isRunningUnitTests else {
            return
        }

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

    private static var isRunningUnitTests: Bool {
        let environment = ProcessInfo.processInfo.environment
        let arguments = ProcessInfo.processInfo.arguments.joined(separator: " ")

        return environment["C64_RETURNAL_HEADLESS"] == "1"
            || environment["XCTestConfigurationFilePath"] != nil
            || environment["XCTestBundlePath"] != nil
            || environment["XCTestSessionIdentifier"] != nil
            || environment["XCInjectBundleInto"] != nil
            || arguments.localizedCaseInsensitiveContains("xctest")
            || Bundle.allBundles.contains { bundle in
                bundle.bundlePath.localizedCaseInsensitiveContains(".xctest")
            }
            || NSClassFromString("XCTest.XCTestCase") != nil
            || NSClassFromString("XCTestCase") != nil
    }
}
