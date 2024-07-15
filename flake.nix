{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    zimeg.url = "github:zimeg/nur-packages";
  };
  outputs = { nixpkgs, flake-utils, ... } @ inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gon =
          if system == "x86_64-darwin" || system == "aarch64-darwin" then
            inputs.zimeg.packages.${pkgs.system}.gon
          else
            null;
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            gnumake
            go
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
        devShells.gh = pkgs.mkShell {
          packages = with pkgs; [
            github-runner
          ];
        };
        devShells.gon = pkgs.mkShell {
          packages = with pkgs; [
            go
            gon
            goreleaser
          ];
          shellHook = ''
            export PATH=/usr/bin:$PATH
          '';
        };
      });
}
