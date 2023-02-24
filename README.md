# Port Forward Pro

`pfpro` is a CLI utility for developers to control local test domains. You can use it to map arbitrary domains to local processes and ports. Built on `caddy`, you get automatic TLS configuration and good performance.

A simple example would be to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.


## Docs

```bash
# configuration
pfpro list
prpro add <alias> <target>
pfpro remove <alias>
pfpro enable <alias>
pfpro disable <alias>

# turning it on and off
pfpro start # start the daemon in the background
pfpro stop # stop the daemon in the background
pfpro status # show the status of the daemon

# explicitly running the daemon
pfpro daemon # run the daemon
```

## TODO
- clear out old managed etc/hosts entries, only current set is ever active
- store config in a consistent location $XDG_HOME/pfpro/config.yaml
  - lockdown to only support 'https? domain -> 127.0.0.1:port'
  - daemon reads config from the file
  - cli edits the config
- cli daemon
  - actually daemonize
  - allow installing the daemon with plist? status commands, etc
- tui / gui / admin controls of some sort
  - set it up on pfpro.local?

```
# firefox about:config on macos
security.enterprise_roots.enabled = true
```
