{
  description = "lm-suggester - LLM提案をreviewdog JSON形式へ変換するGoライブラリ";

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

            # 基本的なシェルツール
            bashInteractive
            coreutils
          ];

          # 環境変数を明示的に設定し、外部ツールの干渉を防ぐ
          GOROOT = "${pkgs.go}/share/go";

          shellHook = ''
            # Nixが管理するパッケージのみを使用
            export PATH="${pkgs.lib.makeBinPath [
              pkgs.go
              pkgs.gotools
              pkgs.golangci-lint
              pkgs.gopls
              pkgs.git
              pkgs.bashInteractive
              pkgs.coreutils
            ]}"

            # GOPATHをプロジェクトローカルに設定
            export GOPATH="$PWD/.go"
          '';
        };

        # Note: This project uses a multi-module setup with replace directives
        # which is complex to build with Nix. Use `nix develop` for development,
        # or build with `go build ./cmd/lm-suggester` inside the dev shell.
      }
    );
}
