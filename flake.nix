{
  description = "lm-suggester";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            go
            gotools
            golangci-lint
            gopls
            git
            syft

            # base utilities
            bashInteractive
            coreutils
          ];

          GOROOT = "${pkgs.go}/share/go";

          shellHook = ''
            export PATH="${pkgs.lib.makeBinPath [
              pkgs.go
              pkgs.gotools
              pkgs.golangci-lint
              pkgs.gopls
              pkgs.git
              pkgs.syft
              pkgs.bashInteractive
              pkgs.coreutils
            ]}"

            export GOPATH="$PWD/.go"
          '';
        };
      }
    );
}
