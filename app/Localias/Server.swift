import SwiftUI

class Server: ObservableObject {
    @Published var status: String = "stopped"
    @Published var error: String = ""

    func Start() -> Bool {
        if let raw = server_start() {
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
        server_stop() // TODO: handle errors?
        return false
    }

    func IsOn() -> Bool {
        return self.status == "running"
    }
}
