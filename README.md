| :warning: Work In Progress |
|----------------------------|
# Localias

Localias is a tool for developers to securely manage local aliases for development servers.

You can use Localias to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.

Major features:
- Works perfectly on MacOS, Linux, and even WSL2 (!)
- Automatically provisions and installs TLS certificates for all of your aliases by default.
- Automatically updates `/etc/hosts` as you add and remove aliases.
- Runs in the foreground or as a background daemon process.
- Uses a shared configuration file if your team puts one in your git repository.
- Built with [`caddy`](https://caddyserver.com/) so it's fast and secure by default.

## Quickstart

After installing `localias`, getting started is easy. First, add some aliases. We'll assume you're running a standard http devserver on `http://localhost:3000`:

```console
$ localias set frontend.test 3000
[added] frontend.test -> 3000
```

You can verify that the rules were added correctly:

```console
$ localias list
frontend.test -> 3000
```

Now, run the proxy server. You can do this in the foreground with `localias run` or in the background with `localias daemon start`. For the purposes of this example, we'll do it in the foreground:

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
- Sign TLS certificates for your aliases, and install its root certificate to your system if it hasn't done so already. 

Each of these steps requires sudo access. But starting/stopping Localias will only prompt for sudo when it needs to, so if you hit `control-C` and restart the process you won't get prompted again:

```console
^C
$ localias run
# ... lots of server logs
# ... but no sudo prompts!
```

Congratulations, you're done!  Start your development servers (or just one of them) in another console. You should be able to visit [`https://frontend.test`](https://frontend.test) in your browser, or make a request with `curl`, and see everything work perfectly\*.

\* *are you using Firefox, or are you on WSL? See the notes below for how to do the one-time install of the localias root certificate*

## Install

Golang:
```bash
# run it
go run github.com/peterldowns/localias/cmd/localias@latest --help
# install it
go install github.com/peterldowns/localias/cmd/localias@latest
```

Homebrew:
```bash
# install it
brew tap peterldowns/tap
brew install localias
```

Nix (flakes):
```bash
# run it
nix run github:peterldowns/localias --help
# install it
nix profile install github:peterldowns/localias --refresh
```

Manual:
- Visit [the latest Github release](https://github.com/peterldowns/localias/releases/latest)
- Download the appropriate binary: `localias-$os-$arch`

## Configuration
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

### Syntax

The configuration file is just a YAML map of `<alias>: <port>`! For example, this is a valid configuration file:

```yaml
https://secure.test: 9000
http://insecure.test: 9001
insecure2.test: 9002
bareTLD: 9003
```

## TODOS

- [ ] Docs
  - [ ] How it works section
  - [ ] WSL2 details
  - [ ] MacOS app details
- [ ] MacOS App
  - [ ] Allow picking the config file
  - [ ] Instructions for installing / using the macos app
- [ ] Installation
  - [ ] homebrew bottles working correctly
- [ ] Improvements
  - [ ] Daemon config command for dumping running config
  - [ ] `--json` formatting for command line controller + caddy logs as well
  - [ ] Helper for doing explicit certificate installation
    - [ ] Handle firefox as well
    - [ ] automatically install localias root certs using powershell script when
          running in wsl2 
- [ ] Code review + cleanup
  - [ ] golang
  - [ ] swift
  - [ ] infra/scripts

## Errata

### Why build this?

Localias is the tool I've always wanted to use for local web development. After years of just visiting `localhost:8080`, I finally got around to looking for a solution, and came across [hotel](https://github.com/typicode/hotel) (unmaintained) and its fork [chalet](https://jeansaad/chalet) (maintained). These projects (basically the same) have the following drawbacks compared to Localias:

- They require NodeJS
- They require configuring a proxy auto-config file
- They do things I don't want (proxying requests to external domains, process management)
- They don't install their self-signed SSL certificates
- They don't allow teams to easily share configuration files

I also wanted an excuse to play around with building a MacOS app, and this seemed like a small and well-defined problem that would be amenable to learning Swift.

Finally, [my friend Justin wanted this to exist, too](https://twitter.com/jmduke/status/1628034461605539840?s=20):

> I swear there's a tool that lets me do:
> 
> localhost:8000 → application.local
> localhost:3000 → marketing.local
> localhost:3002 → docs.local
> 
> But I can't for the life of me remember the name of it. Does anyone know what I'm talking about?


### Domain conflicts and HSTS

When using Localias, you **should not** create aliases with the same name as existing websites. For instance, if you're working on a website hosted in production at `https://example.com`, you really do not want to create a local alias for `example.com` to point to your development server. If you do,
your browser may do things you don't expect:

- Your development cookies will be included in requests to production, and vice-versa. If you are turning localias off/on and switching between
  development and production, these cookies will conflict with each other and generally make you and your website extremely confused.
- If your production website uses [HSTS / certificate pinning](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security), you will see
  very scary errors when trying to use it as a local alias for a development server. This is because localias will be serving content with a different
  private key, but HSTS explicitly tells your browser to disallow this.

In general, it's best to avoid this problem entirely and use aliases that end in [`.test`](https://en.wikipedia.org/wiki/.test), [`.example`](https://en.wikipedia.org/wiki/.example), [`.localhost`](https://en.wikipedia.org/wiki/.localhost), or some other TLD that is not in use.

### `.local` domains on MacOS
If you add an alias to a `.local` domain on a Mac, resolving the domain for the first time [will take add ~5-10s to every
request thanks to Bonjour](https://superuser.com/questions/1596225/dns-resolution-delay-for-entries-in-etc-hosts). The workaround would be to set `127.0.0.1 domain.local` as well as `::1 domain.local` but that's tricky with the way that the `hostctl` package is currently implemented. 

### Firefox and the System Trust Store
Localias's proxy server, Caddy, automatically generates certificates for any
secure aliases you'd like to make. When Localias runs, or the daemon starts, it
will install its root signing certificate into your system store on Mac and
Windows. This means that when you browse to one of your aliases using one of the
following setups, everything works great and your browser will not show any
errors.

- MacOS + Safari
- MacOS + Edge
- MacOS + Chrome
- WSL2 + Chrome
- WSL2 + Edge

Firefox, though, [does not trust the system certificate store by default](https://wiki.mozilla.org/CA/AddRootToFirefox).
This means that when you start Localias and try to browse to an alias for the first time, you will see a warning about an untrusted certificate. You can safely proceed, but it's annoying and looks bad.

But! You can fix this problem, and here's how:

#### MacOS

To make Firefox use the default trust stores that caddy edits: open Firefox,
visit `about:config`, and set

```
security.enterprise_roots.enabled = true
```

This will tell Firefox to trust the system certificate store, and should immediately fix the problem.

#### Windows

Setting `security.enterprise_roots.enabled = true` like on MacOS unfortunately does not work on Windows. The best way is to open Firefox's security settings and manually add the root certificate as a trusted authority. You can do this with the following steps:

1. Find the path to the root certificate being used by Localias. Inside your WSL terminal, run:
   ```console
   $ wslpath -w $(localias debug cert)
   \\wsl$\Ubuntu-20.04\home\pd\.local\state\localias\caddy\pki\authorities\local\root.crt
   ```
   Copy this path to the clipboard.
1. In Firefox, visit *Settings > Privacy & Security > Security > Certificates*,
   or visit *Settings* and search for "certificates".
1. Click *View Certificates*
1. Under the *Authorities* tab, click *Import...*. This will open a filepicker dialog. In the "Name" field, paste the path to the root certificate that you copied earlier. Click *Open*.
1. Check the box next to *Trust this CA to identify websites.* then click *OK*.
1. You should now see "localias" listed as a certificate authority.

Once localias has been added as a certificate authority, you should be able to visit your aliases in Firefox without any security warnings.


### Allow Caddy to bind to ports 443/80 on Linux
Localias is built by wrapping the Caddy webserver, and when you `localias run` or `localias daemon start`, that webserver will attempt to listen on ports 443/80. On Linux you may not be allowed to do this by default. You will see an error like:

```shell
$ localias run
# ... some informational output
error: loading new config: http app module: start: listening on :443: listen tcp :443: bind: permission denied
```

or you may notice that starting the daemon does not result in a running daemon
```shell
$ localias daemon start
$ localias daemon status
daemon is not running
```

To fix this, after installing or upgrading localias, you can use capabilities to
grant `localias` permission to bind on these privileged ports:

```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip $(which localias)
```

For more information, view the [arch man pages for `capabilities`](https://man.archlinux.org/man/capabilities.7#CAP_NET_BIND_SERVICE) and [this Stackoverflow answer](https://stackoverflow.com/a/414258).


### General reading / links / sources

- https://blog.mozilla.org/security/2019/02/14/why-does-mozilla-maintain-our-own-root-certificate-store/
- https://support.mozilla.org/en-US/kb/setting-certificate-authorities-firefox
- https://wiki.mozilla.org/CA/AddRootToFirefox#Windows_Enterprise_Support
- https://adamtheautomator.com/windows-certificate-manager/
- https://stackoverflow.com/a/49553299
- https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/certutil
- https://github.com/christian-korneck/firefox_add-certs
