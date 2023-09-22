{
  description = "Localias is a tool for developers to securely manage local aliases for development servers.";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";

    flake-utils.url = "github:numtide/flake-utils";

    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;

    nix-filter.url = "github:numtide/nix-filter";
  };

  outputs = { self, ... }@inputs:
    inputs.flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ ];
        pkgs = import inputs.nixpkgs {
          inherit system overlays;
        };
        lib = pkgs.lib;
        version = (builtins.readFile ./VERSION);
        commit = if (builtins.hasAttr "rev" self) then (builtins.substring 0 7 self.rev) else "unknown";
      in
      rec {
        packages = rec {
          localias = pkgs.buildGo120Module {
            pname = "localias";
            version = version;
            vendorHash = "sha256-L81PJ1MpXFfcZ/BPYaYlr2rS549i6Lle9l9IRIhh2iE=";
            src =
              let
                # Set this to `true` in order to show all of the source files
                # that will be included in the module build.
                debug-tracing = false;
                source-files = inputs.nix-filter.lib.filter {
                  root = ./.;
                };
              in
              (
                if (debug-tracing) then
                  pkgs.lib.sources.trace source-files
                else
                  source-files
              );
            # Add any extra packages required to build the binaries should go here.
            buildInputs = [ ];
            ldflags = [
              "-X github.com/peterldowns/localias/cmd/localias/shared.Version=${version}"
              "-X github.com/peterldowns/localias/cmd/localias/shared.Commit=${commit}"
            ];
            modRoot = "./.";
            subPackages = ["cmd/localias"];
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
            packages = with pkgs;
              [
                # golang
                delve
                go-outline
                go_1_20
                golangci-lint
                gopkgs
                gopls
                gotools
                # nix
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
      });
}
