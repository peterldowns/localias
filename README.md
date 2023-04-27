| :warning: Work In Progress          |
|---------------------------|

# localias

`localias` is a CLI utility for developers to control local test domains. You can use it to alias arbitrary domains to local dev servers. Built on [`caddy`](https://caddyserver.com/), you get automatic TLS configuration and good performance out of the box.

A simple example would be to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.

```shell
$ localias
securely proxy domains to local development servers

Usage:
  localias [command]

Examples:
  # Add an alias forwarding https://secure.local to http://127.0.0.1:9000
  localias set --alias secure.local -p 9000
  # Update an existing alias to forward to a different port
  localias set --alias secure.local -p 9001
  # Remove an alias
  localias remove secure.local
  # Show aliases
  localias list
  # Clear all aliases
  localias clear
  # Run the server, automatically applying all necessary rules to
  # /etc/hosts and creating any necessary TLS certificates
  localias run
  # Run the server as a daemon
  localias daemon start
  # Check whether or not the daemon is running
  localias daemon status
  # Reload the config that the daemon is using
  localias daemon reload
  # Stop the daemon if it is running
  localias daemon stop

Available Commands:
  clear       clear all aliases
  config      show the configuration file path
  daemon      interact with the daemon process
  help        Help about any command
  list        list all aliases
  remove      remove an alias
  run         run the caddy server
  set         add or edit an alias
  version     show the version of this binary

Flags:
  -c, --configfile string   path to the configuration file to edit
  -h, --help                help for localias
  -v, --version             version for localias

Use "localias [command] --help" for more information about a command.
```

## Install

TODO

## Configuration

TODO

## How it works

TODO


## Errata

#### `.local` domains
If you add an alias to a `.local` domain on a Mac, this will add ~5s to every
request thanks to DNS resolution and Bonjour. See [this Stackoverflow
question](https://superuser.com/questions/1596225/dns-resolution-delay-for-entries-in-etc-hosts)
for more information.  I might be able to work around this in the future by
also adding a `::1` entry for each alias, in addition to `localhost`.

#### Using the system trust store with firefox
To make firefox use the default trust stores that caddy edits: open firefox,
visit `about:config`, and set

```
security.enterprise_roots.enabled = true
```

If you do this, you won't have to see a warning about the certificates being self-signed.
