| :warning: Work In Progress          |
|---------------------------|


# todo - rename me

`localias` is a CLI utility for developers to control local test domains. You can use it to map arbitrary domains to local processes and ports. Built on `caddy`, you get automatic TLS configuration and good performance.

A simple example would be to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.

```shell
$ localias
securely proxy domains to local development servers

Usage:
  localias [command]

Examples:
  # Add an alias forwarding https://secure.local to http://127.0.0.1:9000
  localias add --alias secure.local -p 9000
  # Remove an alias
  localias remove secure.local
  # Show aliases
  localias list
  # Clear all aliases
  localias clear
  # Run the server, automatically applying all necessary rules to
  # /etc/hosts and creating any necessary TLS certificates
  localias run

Available Commands:
  add         add an alias
  clear       clear all aliases
  help        Help about any command
  list        list all aliases
  remove      remove an alias
  run         run the caddy server

Flags:
  -h, --help     help for localias
  -t, --toggle   Help message for toggle

Use "localias [command] --help" for more information about a command.
TODO
```

## TODO
- clear out old managed etc/hosts entries, only current set is ever active
- cli daemon
  - actually daemonize
  - allow installing the daemon with plist? status commands, etc
- tui / gui / admin controls of some sort
  - set it up on localias.local?

```
# to make firefox use the default trust stores that caddy edits:
# open firefox about:config on macos and set
security.enterprise_roots.enabled = true
```
