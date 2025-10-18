{
  description = "lm-suggester - LLM suggestion to reviewdog JSON converter";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        lm-suggester-lib = pkgs.buildGoModule {
          pname = "lm-suggester-lib";
          version = "0.0.0-dev";

          src = ./.;

          vendorHash = "sha256-Z+o55vhuhHe2kJbPjCKVmTY8gzXR6TjbaIUdaRaFEho=";

          subPackages = [ ];

          doInstallCheck = false;

          meta = with pkgs.lib; {
            description = "Go library for converting LLM suggestions to reviewdog JSON format";
            homepage = "https://github.com/HikaruEgashira/lm-suggester";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        lm-suggester-cli = pkgs.buildGoModule {
          pname = "lm-suggester";
          version = "0.0.0-dev";

          src = ./.;

          modRoot = "./cmd/lm-suggester";

          vendorHash = "sha256-IwBucCfotipyZKw0uil9gChEkG0KeZNRUr2ggm08wX4=";

          subPackages = [ "." ];

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${self.rev or "dev"}"
          ];

          meta = with pkgs.lib; {
            description = "CLI tool for converting LLM suggestions to reviewdog JSON format";
            homepage = "https://github.com/HikaruEgashira/lm-suggester";
            license = licenses.mit;
            maintainers = [ ];
            mainProgram = "lm-suggester";
          };
        };
      in
      {
        packages = {
          lib = lm-suggester-lib;
          cli = lm-suggester-cli;
          default = lm-suggester-cli;
        };

        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            go_1_25
            gotools
            golangci-lint
            gopls
            git
            syft

            bashInteractive
            coreutils
          ];
        };
      }
    );
}
