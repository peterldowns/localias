# Localias

| :warning: Work In Progress |
|----------------------------|

Localias is a tool for developers to securely manage local aliases for development servers.

You can use Localias to make it possible to visit `https://server.test` in your browser, and have that request served by a local devserver running at `http://localhost:3000`.

Major features:
- Works perfectly on MacOS, Linux, and even WSL2 (!)
- Automatically provisions and installs TLS certificates for all of your aliases by default.
- Automatically updates `/etc/hosts` as you add and remove aliases.
- Runs in the foreground or as a background daemon process.
- Uses a shared configuration file if your team puts one in your git repository.
- Built with [`caddy`](https://caddyserver.com/) so it's fast and secure by default.

# Quickstart

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

# Install

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

# Configuration
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

## Syntax

The configuration file is just a YAML map of `<alias>: <port>`! For example, this is a valid configuration file:

```yaml
https://secure.test: 9000
http://insecure.test: 9001
insecure2.test: 9002
bareTLD: 9003
```

# Errata

## Why build this?

Localias is the tool I've always wanted to use for local web development. After years of just visiting `localhost:8080`, I finally got around to looking for a solution, and came across [hotel](https://github.com/typicode/hotel) (unmaintained) and its fork [chalet](https://jeansaad/chalet) (maintained). These are wonderful projects that served as inspiration for Localias, but I think that Localias compares favorably:

- Localias is a single binary, whereas Hotel require a working NodeJS runtime
- Localias works by modifying /etc/hosts (and the windows equivalent), which makes it easy to observe and debug. Hotel requires you to configure itself as a proxy in your browser or in your operating system.
  - As a consequence, aliases configured with Localias will also work in command-line scipts or requests sent by progarms like `curl`, whereas aliases managed by Hotel will not work.
- Localias allows you to create any number of aliases on different TLDs at the same time, but Hotel only allows you to use one TLD.
- Localias will install its root certificate to your system store so that you do not see any warnings about invalid self-signed certificates.
- Localias will automatically discover configuration files committed to your git repository, which makes it easy to share a configuration with you development team. 
- Localias does not attempt to do any kind of process management or launching, leaving that entirely up to you.

I also wanted an excuse to play around with building a MacOS app, and this seemed like a small and well-defined problem that would be amenable to learning Swift.

Finally, [my friend Justin wanted this to exist, too](https://twitter.com/jmduke/status/1628034461605539840?s=20):

> I swear there's a tool that lets me do:
> 
> localhost:8000 → application.local  
> localhost:3000 → marketing.local  
> localhost:3002 → docs.local  
> 
> But I can't for the life of me remember the name of it. Does anyone know what I'm talking about?


## Domain conflicts and HSTS

When using Localias, you **should not** create aliases with the same name as existing websites. For instance, if you're working on a website hosted in production at `https://example.com`, you really do not want to create a local alias for `example.com` to point to your development server. If you do,
your browser may do things you don't expect:

- Your development cookies will be included in requests to production, and vice-versa. If you are turning localias off/on and switching between
  development and production, these cookies will conflict with each other and generally make you and your website extremely confused.
- If your production website uses [HSTS / certificate pinning](https://en.wikipedia.org/wiki/HTTP_Strict_Transport_Security), you will see
  very scary errors when trying to use it as a local alias for a development server. This is because localias will be serving content with a different
  private key, but HSTS explicitly tells your browser to disallow this.

In general, it's best to avoid this problem entirely and use aliases that end in [`.test`](https://en.wikipedia.org/wiki/.test), [`.example`](https://en.wikipedia.org/wiki/.example), [`.localhost`](https://en.wikipedia.org/wiki/.localhost), or some other TLD that is not in use.

## `.local` domains on MacOS
If you add an alias to a `.local` domain on a Mac, resolving the domain for the first time [will take add ~5-10s to every
request thanks to Bonjour](https://superuser.com/questions/1596225/dns-resolution-delay-for-entries-in-etc-hosts). The workaround would be to set `127.0.0.1 domain.local` as well as `::1 domain.local` but that's tricky with the way that the `hostctl` package is currently implemented. 

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
warning that the certificate is invalid:

(TODO: image)

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


## Allow Caddy to bind to ports 443/80 on Linux
Localias works by proxying requests from ports 80 and 443 to your development servers. When you run Localias, it therefore will attempt to listen on ports 80 and 443. On Linux you may not be allowed to do this by default -- you may see an error like:

```console
$ localias run
# ... some informational output
error: loading new config: http app module: start: listening on :443: listen tcp :443: bind: permission denied
```

or you may notice that starting the daemon does not result in a running daemon
```console
$ localias daemon start
$ localias daemon status
daemon is not running
```

To fix this, after installing or upgrading Localias, you can use capabilities
to grant the `localias` binary permission to bind on these privileged ports:

```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip $(which localias)
```

For more information, view the [arch man pages for `capabilities`](https://man.archlinux.org/man/capabilities.7#CAP_NET_BIND_SERVICE) and [this Stackoverflow answer](https://stackoverflow.com/a/414258).


## General reading / links / sources

- https://blog.mozilla.org/security/2019/02/14/why-does-mozilla-maintain-our-own-root-certificate-store/
- https://support.mozilla.org/en-US/kb/setting-certificate-authorities-firefox
- https://wiki.mozilla.org/CA/AddRootToFirefox#Windows_Enterprise_Support
- https://adamtheautomator.com/windows-certificate-manager/
- https://stackoverflow.com/a/49553299
- https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/certutil
- https://github.com/christian-korneck/firefox_add-certs

## TODOS

- [ ] Docs
  - [ ] How it works section
  - [ ] WSL2 details
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
