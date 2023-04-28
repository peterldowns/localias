{
  description = "securely proxy domains to local development servers";
  inputs = {
    nixpkgs.url = github:nixos/nixpkgs;

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
          # standard nix definitions
          overlays = [
            inputs.gomod2nix.overlays.default
          ];
          pkgs = import inputs.nixpkgs {
            inherit system overlays;
          };
          lib = pkgs.lib;
          # localias specific
          localiasVersion = (builtins.readFile ./VERSION);
          xcodewrapper = (pkgs.callPackage ./xcodewrapper.nix {});
        in
        rec {
          packages = rec {
            # TODO: somehow pass ldflags here?
            localias = pkgs.buildGoApplication {
              ldflags = [ "-X main.Version=${localiasVersion}" ];
              pname = "localias";
              version = localiasVersion;
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
                # xcode: this wrapper symlinks to the host system.
                # other tools
                just
                cobra-cli
              ] ++ lib.optional stdenv.isDarwin [
                # When in a MacOS environment, must use this wrapper in order
                # for xcodebuild / clang / etc to all work correctly.
                # When not on MacOS, no problem, just won't be able to build
                # the app.
                #
                # TODO: the xcodeenv in nixpkgs is tuned for iOS applications,
                # not for MacOS. May as well modify it and keep something local.
                # Or create a separate repo for it as a flake + correctly
                # set LD=clang.
                (xcodewrapper { allowHigher = true; })
                # Makes the xcodebuild invocation prettier
                xcpretty
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

                # Make it easy to test while developing; add the golang and nix
                # build outputs to the path.
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
