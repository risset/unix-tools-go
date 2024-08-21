{
  description = "Collection of classic Unix shell utilities implemented in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
        nativeBuildInputs = with pkgs; [
          go
          gopls
          delve
          gnumake
        ];
        buildInputs = with pkgs; [];
        lib = pkgs.lib;
      in {
        devShells.default = pkgs.mkShell {inherit nativeBuildInputs buildInputs; };
        packages.default = pkgs.buildGoModule rec {
          name = "unix-utils";
          src = ./.;
          inherit buildInputs;
          vendorHash = lib.fakeHash;
        };
      }
    );
}
