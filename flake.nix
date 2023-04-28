{
  description = "securely proxy domains to local development servers";
  inputs = {
    nixpkgs.url = github:nixos/nixpkgs/nixos-unstable;

    flake-utils.url = github:numtide/flake-utils;

    flake-compat.url = github:edolstra/flake-compat;
    flake-compat.flake = false;

    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, ... }@inputs:
    inputs.flake-utils.lib.eachDefaultSystem
      (system:
        let
          overlays = [
            inputs.gomod2nix.overlays.default
          ];
          pkgs = import inputs.nixpkgs {
            inherit system overlays;
          };
          version = (builtins.readFile ./VERSION);
        in
        rec {
          packages = rec {
            # TODO: somehow pass ldflags here?
            localias = pkgs.buildGoApplication {
              ldflags = [ "-X main.Version=${version}" ];
              pname = "localias";
              version = version;
              src = ./.;
              modules = ./gomod2nix.toml;
              subPackages = [
                "cmd/localias"
              ];
              doCheck = false;
            };
            default = localias;
          };

          apps = rec {
            localias = {
              type = "app";
              program = "${packages.localias}/bin/localias";
            };
            default = localias;
          };

          devShells = rec {
            default = pkgs.mkShell {
              packages = with pkgs; [
                # golang
                delve
                go-outline
                go
                golangci-lint
                gopkgs
                gopls
                gotools
                # nix
                gomod2nix # have to use pkgs. prefix or it breaks lorri
                rnix-lsp
                nixpkgs-fmt
                # other tools
                just
                cobra-cli
              ];

              shellHook = ''
                # The path to this repository
                shell_nix="''${IN_LORRI_SHELL:-$(pwd)/shell.nix}"
                workspace_root=$(dirname "$shell_nix")
                export WORKSPACE_ROOT="$workspace_root"

                # We put the $GOPATH/$GOCACHE/$GOENV in $TOOLCHAIN_ROOT,
                # and ensure that the GOPATH's bin dir is on our PATH so tools
                # can be installed with `go install`.
                #
                # Any tools installed explicitly with `go install` will take precedence
                # over versions installed by Nix due to the ordering here.
                export TOOLCHAIN_ROOT="$workspace_root/.toolchain"
                export GOROOT=
                export GOCACHE="$TOOLCHAIN_ROOT/go/cache"
                export GOENV="$TOOLCHAIN_ROOT/go/env"
                export GOPATH="$TOOLCHAIN_ROOT/go/path"
                export GOMODCACHE="$GOPATH/pkg/mod"
                export PATH=$(go env GOPATH)/bin:$PATH
                export CGO_ENABLED=1
                export PATH="$workspace_root/bin:$workspace_root/result/bin:$PATH"
              '';

              # Need to disable fortify hardening because GCC is not built with -oO,
              # which means that if CGO_ENABLED=1 (which it is by default) then the golang
              # debugger fails.
              # see https://github.com/NixOS/nixpkgs/pull/12895/files
              hardeningDisable = [ "fortify" ];
            };
          };
        }
      );
}
