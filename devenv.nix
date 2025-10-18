{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # Disable cachix to avoid warnings
  cachix.enable = false;

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    syft
  ];

  # https://devenv.sh/languages/
  languages.go.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.test.exec = ''
    go test ./...
  '';

  scripts.test-race.exec = ''
    go test -race ./...
  '';

  scripts.lint.exec = ''
    go vet ./...
  '';

  scripts.coverage.exec = ''
    go test -cover ./...
  '';

  scripts.bench.exec = ''
    go test -run none -bench . ./suggester
  '';

  scripts.example.exec = ''
    cat examples/testdata/simple_replacement.json | go run examples/simple/main.go
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    go test ./...
  '';

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
