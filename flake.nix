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

        # Library package (suggester/)
        lm-suggester-lib = pkgs.buildGoModule {
          pname = "lm-suggester-lib";
          version = "0.0.0-dev";

          src = ./.;

          # Use lib.fakeHash initially, then replace with actual hash after first build
          vendorHash = "sha256-Z+o55vhuhHe2kJbPjCKVmTY8gzXR6TjbaIUdaRaFEho=";

          # Don't install any binaries for the library package
          subPackages = [ ];

          # Skip installation phase since this is just a library
          doInstallCheck = false;

          meta = with pkgs.lib; {
            description = "Go library for converting LLM suggestions to reviewdog JSON format";
            homepage = "https://github.com/HikaruEgashira/lm-suggester";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        # CLI package (cmd/lm-suggester/)
        lm-suggester-cli = pkgs.buildGoModule {
          pname = "lm-suggester";
          version = "0.0.0-dev";

          src = ./.;

          # CLI has its own go.mod in cmd/lm-suggester/
          modRoot = "./cmd/lm-suggester";

          vendorHash = "sha256-IwBucCfotipyZKw0uil9gChEkG0KeZNRUr2ggm08wX4=";

          # Build only the CLI binary
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
        # Packages output
        packages = {
          lib = lm-suggester-lib;
          cli = lm-suggester-cli;
          default = lm-suggester-cli;
        };

        # Development shell
        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            go_1_25
            gotools
            golangci-lint
            gopls
            git
            syft

            # base utilities
            bashInteractive
            coreutils
          ];
        };
      }
    );
}
