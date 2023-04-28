import SwiftUI

struct Wrapper : View {
    @State var isOn: Bool
    @State var path: String
    @State var entries: [Entry]
    @State var daemonProcess: DispatchWorkItem?
    
    init() {
        let cfg = loadConfig()!
        _daemonProcess = .init(initialValue: nil)
        _isOn = .init(initialValue: false)
        _path = .init(initialValue: cfg.Path)
        _entries = .init(initialValue:
                                (cfg.Entries ?? []).sorted())
    }
    
    var body: some View {
        ContentView(
                isOn: $isOn,
                path: $path,
                entries: $entries,
                daemonProcess: $daemonProcess
            )
    }
}

struct Wrapper_Previews: PreviewProvider {
    static var previews: some View {
        Wrapper()
    }
}

struct ContentView: View {
    @StateObject var daemon = Server()
    @Binding var isOn: Bool;
    @Binding var path: String
    @Binding var entries: [Entry]
    @Binding var daemonProcess: DispatchWorkItem?
    
    @State var newEntry: Entry =  Entry(
        Alias:"",
        Port: 0
    )
    
    func add() {
        // Validate the new rule
        var entry = newEntry
        entry.Alias = entry.Alias.trim()
        if entry.Alias == "" {
            return
        }
        if entry.Port == 0 {
            return
        }
        // If there is an existing rule with the same alias,
        // replace it with the new one. Otherwise, add to the list.
        if var existing = entries.first(where: {d in
            d.Alias == entry.Alias
        }) {
            existing.Port = entry.Port
        } else {
            entries.append(entry)
        }
        // Sort them and clear the input field.
        newEntry = Entry(
            Alias: "",
            Port: 0
        )
        entries = entries.sorted()
        focusedField = .port
    }
    
    func remove(_ entry: Entry) {
        entries.removeAll(where: {d in
            d == entry
        })
        entries = entries.sorted()
    }
    
    func minHeight() -> CGFloat {
        var min = CGFloat(140 + 36 * entries.count)
        if min > 800 {
            min = 800
        }
        if min < 240 {
            min = 240
        }
        return min
    }
    func save() {
        let encoder = JSONEncoder()
        let config = Config(Path: path, Entries: entries)
        let data = try! encoder.encode(config)
        let string = strdup(String(data: data, encoding: .utf8))
        if let raw = config_save(string) {
            let error = String(cString: raw)
            print("failed to save config:", error)
        }
        if daemon.IsOn() {
            _ = daemon.Start()
        }
        entries = entries.sorted()
    }
    
    func reload() {
        let cfg = loadConfig()!
        path = cfg.Path
        entries = (cfg.Entries ?? []).sorted()
    }
    
    
    @Environment(\.colorScheme) var colorScheme
    @State private var serverRunning: Bool = false
    
    
    var dynamicItems: some View {
        VStack(
            alignment:.leading,
            spacing: 16
        ) {
            ForEach($entries) { $entry in
                HStack(alignment: .top){
                    TextField(
                        "port",
                        value: $entry.Port,
                        formatter: NumberFormatter(),
                        prompt: Text("port:")
                    )
                    .textFieldStyle(.plain)
                    .multilineTextAlignment(.leading)
                    .frame(
                        width: 60,
                        alignment: .leading
                    )
                    TextField(
                        "alias",
                        text: $entry.Alias,
                        prompt: Text("alias")
                    )
                    .textFieldStyle(.plain)
                    .multilineTextAlignment(.leading)
                    Spacer()
                    Button(action: {
                        remove(entry)
                    }) {
                        Image(systemName: "trash")
                    }
                    .buttonStyle(.plain)
                }.contextMenu {
                    Button(action: {
                        remove(entry)
                    }){
                        Text("Delete")
                    }
                }.padding(.trailing, 20)
            }
        }
    }
    
    private enum Field: Int, Hashable {
        case port, alias
    }
    @FocusState private var focusedField: Field?
    
    var body: some View {
        let serverRunningToggle = Binding {
            serverRunning
        } set: { x in
            if x {
                serverRunning = daemon.Start()
            } else {
                self.save()
                serverRunning = daemon.Stop()
            }
        }
        
        VStack(alignment:.leading){
            HStack(alignment:.top){
                Toggle("Server", isOn: serverRunningToggle)
                    .toggleStyle(SwitchToggleStyle(tint:.accentColor))
                    .labelsHidden()
                    .help("Turn the proxy server on or off")
                Spacer()
                Button(action: save) {
                    // system(size: 20, weight: .light)
                    Label("Save", systemImage: "square.and.arrow.down")
                }.help("Save the current config")
                Button(action: reload) {
                    Label("Undo", systemImage: "arrow.counterclockwise.circle")
                }.help("Reset to the last saved config")
                Spacer()
                Button(action: {
                    exit(0)
                }) {
                    Label("Quit", systemImage: "multiply")
                }
                .help("Quit Localias")
            }.padding(.bottom, 10).padding(.top, 10).padding(.trailing, 20)
            Divider().padding([.bottom,.trailing], 10)
            Form {
                HStack(alignment: .top){
                    TextField(
                        "port",
                        value: $newEntry.Port,
                        formatter: HiddenZeroFormatter,
                        prompt: Text("port")
                    )
                        .textFieldStyle(.plain)
                        .font(.body)
                        .multilineTextAlignment(.leading)
                        .frame(width: 60, alignment: .leading)
                        .focused($focusedField, equals:.port)
                    TextField(
                        "alias",
                        text: $newEntry.Alias,
                        prompt: Text("alias")
                    )
                        .textFieldStyle(.plain)
                        .font(.body)
                        .multilineTextAlignment(.leading)
                        .focused($focusedField, equals:.alias)
                    Spacer()
                    Button(action: add) {
                        Image(systemName: "plus")
                    }.buttonStyle(.plain)
                    
                }.padding(.trailing, 20)
            }.onSubmit(add)
                .labelsHidden()
                .padding(.top, 2)
                .padding(.bottom, 10)
            Divider().padding([.bottom,.trailing], 10)
            
            
            ScrollView {
                dynamicItems
            }
            
            
        }.padding()
            .padding(.leading, 10)
            .padding(.bottom, 10)
            .frame(
                minWidth: 360,
                minHeight: self.minHeight(),
                maxHeight: 800,
                alignment:.top
            )
            .background(
                colorScheme == .light ?
                    .white.opacity(0.5) :
                        .black.opacity(0.5)
            )
    }
}
