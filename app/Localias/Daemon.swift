import SwiftUI

class Daemon: ObservableObject {
    @Published var status: String = "stopped"
    @Published var error: String = ""

    func Start() -> Bool {
        if let raw = daemon_start() {
            self.status = "stopped"
            self.error = String(cString: raw)
            print("failed to start:", self.error)
            return false
        } else {
            self.status = "running"
            self.error = ""
            return true
        }
    }

    func Stop() -> Bool {
        self.status = "stopped" // TODO: enum
        daemon_stop() // TODO: handle errors?
        return false
    }

    func IsOn() -> Bool {
        return self.status == "running"
    }
}
