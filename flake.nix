{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        gon =
          if system == "x86_64-darwin" || system == "aarch64-darwin" then
            let
              gonZip = pkgs.fetchurl {
                url = "https://github.com/Bearer/gon/releases/download/v0.0.36/gon_macos.zip";
                sha256 = "1firj23pgdfx9hybjjr91chn1jzf7lzjrbx3nm1s9h3xbpx945x2";
              };
            in
            pkgs.runCommand "gon" { nativeBuildInputs = [ pkgs.unzip ]; } ''
              mkdir -p $out/bin
              unzip ${gonZip} -d $out/bin
            ''
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
          shellHook = "go mod tidy";
        };
        devShells.gh = pkgs.mkShell {
          packages = with pkgs; [
            github-runner
          ];
        };
        devShells.gon =
          if system == "x86_64-darwin" || system == "aarch64-darwin" then
            pkgs.mkShell
              {
                buildInputs = with pkgs; [
                  go
                  gon
                  goreleaser
                ];
                shellHook = ''
                  export PATH=/usr/bin:$PATH:${gon}/bin
                '';
              }
          else
            null;
      });
}
