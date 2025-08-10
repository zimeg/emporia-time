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
        quill = if pkgs.stdenv.isDarwin then inputs.zimeg.packages.${pkgs.system}.quill else null;
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
        devShells.quill =
          if pkgs.stdenv.isDarwin then
            pkgs.mkShell {
              packages = with pkgs; [
                go
                quill
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
          vendorHash = "sha256-9uN8zdc0BEqrZmpxyvKRSKGc2SBOIn4Tp/Lx3cffyjQ=";
          meta = {
            description = "an energy aware time command";
            homepage = "https://github.com/zimeg/emporia-time";
            license = pkgs.lib.licenses.mit;
          };
        };
      }
    );
}
