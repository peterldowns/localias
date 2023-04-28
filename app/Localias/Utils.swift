import SwiftUI

extension String {
    func trim() -> String {
        return self.trimmingCharacters(in: NSCharacterSet.whitespaces)
    }
}

// When the number being formatted is 0, show nothing (empty string) instead of
// 0.
let HiddenZeroFormatter = {
    var fmt = NumberFormatter()
    fmt.zeroSymbol = ""
    return fmt
}()
