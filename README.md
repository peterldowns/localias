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
