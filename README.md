# 🏠 localias

![Latest Version](https://badgers.space/badge/latest%20version/v2.0.0/blueviolet?corner_radius=m)
![Golang](https://badgers.space/badge/golang/1.18+/blue?corner_radius=m)

Localias is a tool for developers to securely manage local aliases for development servers.

Use Localias to redirect `https://server.test` &rarr; `http://localhost:3000` in your browser and on your command line. 

<img width="464" alt="iTerm showing the most basic usage of Localias" src="https://github.com/peterldowns/localias/assets/824173/5b0121df-237e-47e7-92b8-d09017fcf95f.png">

### Major Features
- Use convenient names, without ports, in your URLs
- Serve your development website behind TLS, minimizing differences between development and production.
  - No more CORS problems!
  - Set secure cookies!
- Works on MacOS, Linux, and even WSL2 (!)
- Automatically provisions and installs TLS certificates for all of your aliases
  by default.
- Automatically updates `/etc/hosts` as you add and remove aliases, so that they
  work with all of your tools.
- Runs in the foreground or as a background daemon process, your choice.
- Supports shared configuration files so your whole team can use the same
  aliases for your development services.
- Proxies requests and generates TLS certs with
  [`caddy`](https://caddyserver.com/) so it's fast and secure by default.
- Serves `.local` domains over mDNS, so you can visit your development
  servers from your phone or any other device connected to the same network.


# Install

#### Homebrew:
```bash
# install it
brew install peterldowns/tap/localias
```

#### Golang:
```bash
# run it
go run github.com/peterldowns/localias/cmd/localias@latest --help
# install it
go install github.com/peterldowns/localias/cmd/localias@latest
```

#### Nix (flakes):
```bash
# run it
nix run github:peterldowns/localias -- --help
# install it
nix profile install --refresh github:peterldowns/localias
```

#### Manually download binaries
Visit [the latest Github release](https://github.com/peterldowns/localias/releases/latest) and pick the appropriate binary. Or, click one of the shortcuts here:
- [darwin-amd64](https://github.com/peterldowns/localias/releases/latest/download/localias-darwin-amd64)
- [darwin-arm64](https://github.com/peterldowns/localias/releases/latest/download/localias-darwin-arm64)
- [linux-amd64](https://github.com/peterldowns/localias/releases/latest/download/localias-linux-amd64)
- [linux-arm64](https://github.com/peterldowns/localias/releases/latest/download/localias-linux-arm64)

# How does it work?

Localias has two parts:
- the configuration file
- the proxy server

The configuration file is where Localias keeps track of your aliases, and to which local ports they should be pointing. The proxy server then runs and actually proxies requests based on the configuration.

### configuration file
Every time you run `localias`, it looks for a config file in the following places, using the first one that it finds:

- If you pass an explicit `--configfile <path>`, it will attempt to use `<path>`
- If you set an environment variable `LOCALIAS_CONFIGFILE=<path>`, it will attempt to use `<path>`
- If your current directory has `.localias.yaml`, it will use `$pwd/.localias.yaml`
- If you are in a git repository and there is a `.localias.yaml` at the root of the repository, use `$repo_root/.localias.yaml`
- Otherwise, use `$XDG_CONFIG_HOME/localias.yaml`, creating it if necessary.
  - On MacOS, this defaults to `~/Library/Application\ Support/localias.yaml`
  - On Linux or on WSL, this defaults to `~/.config/localias.yaml`

This means that your whole dev team can share the same aliases by adding `.localias.yaml` to the root of your git repository.

To show the configuration file currently in use, you can run

```bash
# Print the path to the current configuration file
localias debug config
# Print the contents of the current configuration file
localias debug config --print
```

The following commands all interact directly with the configuration file:

```shell
# add or edit an alias
localias set <alias> <port>
# clear all aliases
localias clear
# list all aliases
localias list
# remove an alias
localias remove <alias>
```

The configuration file is just a YAML map of `<alias>: <port>`! For example, this is a valid configuration file:

```yaml
bareTLD: 9003 # serves over https and http
implicitly_secure.test: 9002 # serves over https and http
https://explicit_secure.test: 9000 # serves over https and http
http://explicit_insecure.test: 9001 # serves over http only
```

### proxy server

When you execute `localias run` or `localias start` to run the proxy server, Localias performs the following operations:

- Reads the current Localias configuration file to find all the current aliases and the ports to which they're pointing.
- Checks the `/etc/hosts` file to make sure that every alias is present
  - Adds any new aliases that aren't already present
  - Removes any old aliases that are no longer in the Localias config
  - Only updates the file if any changes were made, since this requires `sudo` privileges.
- Runs the Caddy proxy server
  - If Caddy has not already generated a local root certificate:
    - Generate a local root certificate to sign TLS certificates
    - Install the local root certificate to the system's trust stores, and the Firefox certificate store if it exists and an be accessed.
  - Generate a Caddy configuration telling it how to redirect each alias to the correct local port.
  - Generate and sign TLS certificates for each of the aliases currently in use
  - Bind to ports 80/443 in order to proxy requests

Localias requires elevated privileges to perform these actions as part of running the proxy server:
- Edit `/etc/hosts`
- Install the locally generated root certificate to your system store
- Bind to ports 80/443 in order to run the proxy server

When you run Localias, each time it needs to do these things, it will open a subshell using `sudo` to perform these actions, and this will prompt you for your password. Localias *does not read or interact with your password*.

Localias is entirely local and performs no telemetry.

# Quickstart

## Running the server for the first time

After installing `localias`, you will need to configure some aliases. For this quickstart example, we'll assume that you're running a local http frontend devserver on `http://localhost:3000`, and that you'd like to be able to access it at `https://frontend.test` in your browser and via tools like `curl`.

First, create the alias:

```console
$ localias set frontend.test 3000
[added] frontend.test -> 3000
```

You can check to see that it was added correctly:

```console
$ localias list
frontend.test -> 3000
```

That's it in terms of configuration!

Now, start the proxy server. You can do this in the foreground with `localias run` (and stop it with `ctrl-c`) or you can start the server in the background with `localias start`. For the purposes of this quickstart, we'll do it in the foreground.

```console
$ localias run
# some prompts to authenticate as root
# ... lots of server logs like this:
2023/05/02 23:12:58.218 INFO    tls.obtain      acquiring lock  {"identifier": "frontend.test"}
2023/05/02 23:12:58.229 INFO    tls.obtain      lock acquired   {"identifier": "frontend.test"}
2023/05/02 23:12:58.230 INFO    tls.obtain      obtaining certificate   {"identifier": "frontend.test"}
2023/05/02 23:12:58.230 INFO    tls.obtain      certificate obtained successfully       {"identifier": "frontend.test"}
2023/05/02 23:12:58.230 INFO    tls.obtain      releasing lock  {"identifier": "frontend.test"}
# process is now waiting for requests
```

This will prompt you to authenticate at least once. Each time Localias runs, it will

- Automatically edit your `/etc/hosts` file and add entries for each of your aliases.
- Sign TLS certificates for your aliases, and generate+install a custom root certificate to your system if it hasn't done so already. 

Each of these steps requires sudo access. But starting/stopping Localias will only prompt for sudo when it needs to, so if you hit `control-C` and restart the process you won't get prompted again:

```console
^C
$ localias run
# ... lots of server logs
# ... but no sudo prompts!
```

Congratulations, you're done! Start your development servers (or just one of them) in another console. You should be able to visit [`https://frontend.test`](https://frontend.test) in your browser, or make a request with `curl`, and see everything work perfectly\*.

\* *are you using Firefox, or are you on WSL? See the notes below for how to do the one-time install of the localias root certificate*

## Running as a daemon

Instead of explicitly running the proxy server as a foreground process with `localias run`, you can also run Localias in the background with `localias start`. You can interact with this daemon with the following commands:

```shell
# Start the proxy server as a daemon process
localias start
# Show the status of the daemon process
localias status
# Apply the latest configuration to the proxy server in the daemon process
localias reload
# Stop the daemon process
localias stop
```

When running as a daemon process, if you make any changes to your configuration you
will need to explicitly reload the daemon:

```shell
# Start with frontend.test -> 3000
localias set frontend.test 3000
localias start
# Update frontend.test -> 4004. 
localias set frontend.test 4004
# The daemon will still be running with frontend.test -> 3000, so
# to apply the new changes you'll need to reload it
localias reload
```

# Using the CLI 

`localias` has many different subcommands, each of which is documented
(including usage examples). To see the available subcommands, run `localias`. To
see help on any command, you can run `localias help $command` or
`localias $command --help`. 

```console
$ localias
securely manage local aliases for development servers

Usage:
  localias [flags]
  localias [command]

Examples:
  # Add an alias forwarding https://secure.test to http://127.0.0.1:9000
  localias set secure.test 9000
  # Update an existing alias to forward to a different port
  localias set secure.test 9001
  # Remove an alias
  localias rm secure.test
  # List all aliases
  localias list
  # Clear all aliases
  localias clear
  
  # Start the proxy server as a daemon process
  localias start
  # Show the status of the daemon process
  localias status
  # Apply the latest configuration to the proxy server in the daemon process
  localias reload
  # Stop the daemon process
  localias stop
  # Run the proxy server in the foreground
  localias run

Available Commands:
  clear       clear all aliases
  help        Help about any command
  list        list all aliases
  reload      apply the latest configuration to the proxy server in the daemon process
  rm          remove an alias
  run         run the proxy server in the foreground
  set         add or edit an alias
  start       start the proxy server as a daemon process
  status      show the status of the daemon process
  stop        stop the daemon process
  version     show the version of this binary

Flags:
  -c, --configfile string   path to the configuration file to edit
  -h, --help                help for localias
  -v, --version             version for localias

Use "localias [command] --help" for more information about a command.
```

# Errata

## Why build this?

Localias is the tool I've always wanted to use for local web development. After years of just visiting `localhost:8080`, I finally got around to looking for a solution, and came across [hotel](https://github.com/typicode/hotel) (unmaintained) and its fork [chalet](https://github.com/jeansaad/chalet) (maintained). These are wonderful projects that served as inspiration for Localias, but I think Localias is implemented in a better and more useful way.

Finally, [my friend Justin wanted this to exist, too](https://twitter.com/jmduke/status/1628034461605539840?s=20):

> I swear there's a tool that lets me do:
> 
> localhost:8000 → application.local  
> localhost:3000 → marketing.local  
> localhost:3002 → docs.local  
> 
> But I can't for the life of me remember the name of it. Does anyone know what I'm talking about?

## Why not hotel/chalet?
Localias is designed to replace alternative tools like [hotel](https://github.com/typicode/hotel)/[chalet](https://github.com/jeansaad/chalet). Hotel is no longer maintained, and Chalet is a fork of Hotel with basically the same features. I think Localias compares favorably:

  - Localias is a single binary. Hotel requires a working NodeJS runtime.
  - Localias works by modifying `/etc/hosts` (and the windows equivalent), which makes it easy to observe and debug. Hotel requires you to configure itself as a proxy in your browser or in your operating system.
    - Aliases configured with Localias will also work in command-line scipts or requests sent by programs like `curl`. Hotel aliases only work in your browser.
  - Localias allows you to create any number of aliases on different TLDs at the same time. Hotel only allows you to use one TLD.
  - Localias will generate a root certificate and any necessary certificates for each alias, and install the root certificate in your system store so you do not see any warnings about invalid self-signed certificates. Hotel does not do any TLS signing.
  - Localias will automatically discover configuration files committed to your git repository, which makes it easy to share a configuration with you development team. Hotel does not allow for shared configuration files.
  - Localias does not attempt to do any kind of process management or launching, leaving that entirely up to you. Hotel attemps to run and manage processes for you.


## Domain conflicts and HSTS

When using Localias, you **should not** create aliases with the same name as existing websites. For instance, if you're working on a website hosted in production at `https://example.com`, you really do not want to create a local alias for `example.com` to point to your development server. If you do,
your browser may do things you don't expect:

- Your development cookies will be included in requests to production, and vice-versa. If you are turning localias off/on and switching between
  development and production, these cookies will conflict with each other and generally make you and your website extremely confused.
- If your production website uses [HSTS / certificate pinning](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security), you will see
  very scary errors when trying to use it as a local alias for a development server. This is because localias will be serving content with a different
  private key, but HSTS explicitly tells your browser to disallow this.

In general, it's best to avoid this problem entirely and use aliases that end in [`.test`](https://en.wikipedia.org/wiki/.test), [`.example`](https://en.wikipedia.org/wiki/.example), [`.localhost`](https://en.wikipedia.org/wiki/.localhost), or some other TLD that is not in use.

## `.local` domains
Thanks to ["mDNS", or "multicast
dns"](https://en.wikipedia.org/wiki/Multicast_DNS), any aliases that you create
that end in `.local` will be broadcast to your entire network. This makes it
easy to visit a development server from any other device, including your phone,
which makes testing responsive websites really easy. All you need to do is create
an alias ending in `.local`:

```console
$ localias add frontend.local 8080
[added] frontend.local -> 8080
$ localias add http://insecure.local 8080
[added] http://insecure.local 8080
```

When you visit a secure local alias from another device, you may be prompted with
a certificate warning. You should feel free to accept the warning and continue
on to the site, which is just your local development site.

## The Localias Root Certificate and System Trust Stores
Localias's proxy server, Caddy, automatically generates certificates for any
secure aliases you'd like to make. When Localias runs it will make sure that
its root signing certificate is installed in the system store on Mac and Linux.
If your browser reads from the system store to determine which certificate
authorities to trust, this means that everything will work nicely for you out
of the box.

This means that if you're using Safari/Edge/Chrome on MacOS/Linux, you're good
to go, and you will see a nice "verified" or "secure" status when you visit one
of your secure aliases in your browser.

### WSL
When you run Localias inside of WSL, so basically inside of a Linux virtual
machine with a Windows host, Caddy will generate certificates and install them
to the Linux VM's trust store, but not to the parent Windows host. This means
that if you're using a browser running in Windows, you will see a certificate
warning if you visit a secure alias.

You can fix this by explicitly installing the Localias root certificate to your
Windows machine's certificate store. You can do this with the following
command, which will prompt you to authorize it as an administrator:

```bash
localias debug cert --install
```

### Firefox
Firefox [does not trust the system certificate store by
default](https://blog.mozilla.org/security/2019/02/14/why-does-mozilla-maintain-our-own-root-certificate-store/).
This means that unfortunately, if you visit you secure alias, you will see a
warning that the certificate is invalid.

On MacOS/Linux, Firefox can be configured to trust the system store by changing
a configuration setting.

1. Open Firefox
1. Visit `about:config`
1. Set
   ```
   security.enterprise_roots.enabled = true
   ```
1. Quit and re-open Firefox

Altenately, or if you're using Firefox on Windows to try to browse to a server
running in WSL, you can manually add the Localias root certificate to Firefox.
You will need to do this if you're using WSL, since Firefox on Windows does not
read from the system trust store.

1. Find the path to the root certificate being used by Localias. If you're on MacOS or Linux, run:

   ```console
   $ localias debug cert
   /Users/pd/Library/Application Support/localias/caddy/pki/authorities/local/root.crt
   ```
   to print the path to the certificate.

   In WSL, you'll need to convert this to a Windows file path using the `wslpath` tool:

   ```console
   $ wslpath -w $(localias debug cert)
   \\wsl$\Ubuntu-20.04\home\pd\.local\state\localias\caddy\pki\authorities\local\root.crt
   ```
   Copy this path to the clipboard.
1. In Firefox, visit *Settings > Privacy & Security > Security > Certificates*,
   or visit *Settings* and search for "certificates".
1. Click *View Certificates*
1. Under the *Authorities* tab, click *Import...*. This will open a filepicker dialog.

   - On MacOS: hit "Cmd+Shift+G" to open a filepath dialog. Paste the path you copied earlier to select the `root.crt`.
   - On Windows: in the "Name" field, paste the path to the root certificate that you copied earlier.

   Click *Open*.
1. Check the box next to *Trust this CA to identify websites.* then click *OK*.

You should now see "localias" listed as a certificate authority. If you visit a
secure alias, you should see that the certificate is trusted and no errors or
warnings are displayed.


## Allow Localias to bind to ports 443/80 on Linux
Localias works by proxying requests from ports 80 and 443 to your development
servers. When you run Localias, it therefore will attempt to listen on ports 80
and 443. On Linux you may not be allowed to do this by default -- you may see an
error like:

```console
$ localias run
# ... some informational output
error: loading new config: http app module: start: listening on :443: listen tcp :443: bind: permission denied
```

or you may notice that starting the daemon does not result in a running daemon
```console
$ localias start
$ localias status
daemon is not running
```

To fix this, after installing or upgrading Localias, you can use capabilities
to grant the `localias` binary permission to bind on these privileged ports:

```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip $(which localias)
```

For more information, view the [arch man pages for `capabilities`](https://man.archlinux.org/man/capabilities.7#CAP_NET_BIND_SERVICE) and [this Stackoverflow answer](https://stackoverflow.com/a/414258).

## error: localias could not start successfully

If you've tried running `localias run` and see this error:

```console
$ localias run
error: localias could not start successfully. Most likely there is another instance of
localias or some other kind of proxy or server listening to ports 443/80, which
is preventing another instance from starting. Common causes:

- You have another instance of localias running in a different terminal
- You have a proxy server like Caddy, Nginx, or Apache running
- There is a bug in localias

Please see the https://github.com/peterldowns/localias README for some
diagnostics and ideas for how to debug this.
```

Or you've tried to start the daemon `localias start` but no daemon gets started:

```console
$ localias start
$ localias status
daemon is not running
```

Then most likely some other process is bound to ports 443/80, preventing
localias from starting up correctly. The only way localias will start is if it
is able to bind to these ports, which it needs to do to act as a proxy.

To find out if there are any other instances of localias running, use `ps`. In
this example, the first result is an instance of localias, and the second result
is the `grep` process itself.

```console
$ ps aux | grep -i localias
pd               39020   0.0  0.1 409289408  38736 s003  S+    1:42PM   0:00.09 localias run
pd               39198   0.0  0.0 407965536    624 s005  R+    1:47PM   0:00.00 grep -i localias
```

You can find out what services are listening on your ports by using `lsof`. In this example,
there the results show that there is an instance of localias bound to both port 80 and port 443:

```console
$  lsof -Pn | grep -E '\*:443|\*:80'
localias  39020   pd    9u     IPv6 0xb3abbd50442d943f       0t0                 TCP *:443 (LISTEN)
localias  39020   pd   11u     IPv6 0xb3abbd4b78f6ba3f       0t0                 UDP *:443
localias  39020   pd   12u     IPv6 0xb3abbd50442da23f       0t0                 TCP *:80 (LISTEN)
```

In order for localias to start, you'll have to kill the process that is
interfering and binding to these ports.

# General reading / links / sources

- https://blog.mozilla.org/security/2019/02/14/why-does-mozilla-maintain-our-own-root-certificate-store/
- https://support.mozilla.org/en-US/kb/setting-certificate-authorities-firefox
- https://wiki.mozilla.org/CA/AddRootToFirefox#Windows_Enterprise_Support
- https://adamtheautomator.com/windows-certificate-manager/
- https://stackoverflow.com/a/49553299
- https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/certutil
- https://github.com/christian-korneck/firefox_add-certs
- https://rud.is/b/2021/04/24/making-macos-universal-apps-with-universal-golang-static-libraries/
- https://caddyserver.com/docs/automatic-https#overview


## Future Work

- [ ] Daemon config command for dumping running config
- [ ] `--json` formatting for command line controller + caddy logs as well
- [ ] Helper for doing explicit certificate installation
  - [ ] Handle firefox if `certutil` is available?
  - [ ] automatically install localias root certs using powershell script when
        running in wsl2
- [ ] Daemonized server errors are reported if it fails to start
- [ ] Better helpers for getting access to logs
- [ ] General code cleanup and tests
