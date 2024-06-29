{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    zimeg.url = "github:zimeg/nur-packages";
  };
  outputs = { nixpkgs, flake-utils, zimeg, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gon =
          if system == "x86_64-darwin" || system == "aarch64-darwin" then
            zimeg.packages.${pkgs.system}.gon
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
            gon
            gopls
            goreleaser
          ];
          shellHook = "go mod tidy";
        };
      });
}
