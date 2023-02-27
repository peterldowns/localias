| :warning: Work In Progress          |
|---------------------------|


# todo - rename me

`pfpro` is a CLI utility for developers to control local test domains. You can use it to map arbitrary domains to local processes and ports. Built on `caddy`, you get automatic TLS configuration and good performance.

A simple example would be to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.

```shell
$ ./bin/pfpro
securely proxy domains to local development servers

Usage:
  pfpro [command]

Available Commands:
  add         add an alias
  help        Help about any command
  hostctl     modify an /etc/hosts-type file
  list        list all aliases
  remove      remove aliases
  run         run the caddy server

Flags:
  -h, --help     help for pfpro
  -t, --toggle   Help message for toggle
```

## TODO
- clear out old managed etc/hosts entries, only current set is ever active
- cli daemon
  - actually daemonize
  - allow installing the daemon with plist? status commands, etc
- tui / gui / admin controls of some sort
  - set it up on pfpro.local?

```
# to make firefox use the default trust stores that caddy edits:
# open firefox about:config on macos and set
security.enterprise_roots.enabled = true
```
