# This is a shim that allows non-flake Nix users to build this project, using
# the standard compatibility tools.
# 
# From https://nixos.wiki/wiki/Flakes#Using_flakes_project_from_a_legacy_Nix
(import
  (
    let
      lock = builtins.fromJSON (builtins.readFile ./flake.lock);
    in
    fetchTarball {
      url = "https://github.com/edolstra/flake-compat/archive/${lock.nodes.flake-compat.locked.rev}.tar.gz";
      sha256 = lock.nodes.flake-compat.locked.narHash;
    }
  )
  {
    src = ./.;
  }).defaultNix
