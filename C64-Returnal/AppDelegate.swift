//
//  AppDelegate.swift
//  C64-Returnal
//
//  Created by Tóth István on 2026. 06. 02..
//


import Cocoa

@main
class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ aNotification: Notification) {
        NSApp.windows.forEach(disableWindowRestoration)
    }

    func applicationWillTerminate(_ aNotification: Notification) {
    }

    func applicationSupportsSecureRestorableState(_ app: NSApplication) -> Bool {
        true
    }

    func applicationShouldSaveApplicationState(_ sender: NSApplication) -> Bool {
        false
    }

    func applicationShouldRestoreApplicationState(_ sender: NSApplication) -> Bool {
        false
    }

    private func disableWindowRestoration(_ window: NSWindow) {
        window.isRestorable = false
        window.restorationClass = nil
        window.identifier = nil
    }
}
