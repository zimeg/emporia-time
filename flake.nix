{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    zimeg.url = "github:zimeg/nur-packages";
  };
  outputs =
    { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gon = if pkgs.stdenv.isDarwin then inputs.zimeg.packages.${pkgs.system}.gon else null;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            gnumake
            go
            go-junit-report
            gocyclo
            gofumpt
            golangci-lint
            gopls
            goreleaser
          ];
          shellHook = ''
            go mod tidy
          '';
        };
        devShells.gon =
          if pkgs.stdenv.isDarwin then
            pkgs.mkShell {
              packages = with pkgs; [
                go
                gon
                goreleaser
              ];
              shellHook = ''
                export PATH=/usr/bin:$PATH # https://github.com/zimeg/nur-packages/issues/4
              '';
            }
          else
            null;
        devShells.tom = pkgs.mkShell {
          packages = with pkgs; [
            time
          ];
        };
        packages.default = pkgs.buildGoModule rec {
          pname = "emporia-time";
          version = "unversioned";
          src = ./.;
          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];
          doCheck = true;
          vendorHash = "sha256-lQIFcGH87XIv4DAI7yzA6gV8AEJ2ENLtoGyyGW+c9Ig=";
          meta = {
            description = "an energy aware time command";
            homepage = "https://github.com/zimeg/emporia-time";
            license = pkgs.lib.licenses.mit;
          };
        };
      }
    );
}
