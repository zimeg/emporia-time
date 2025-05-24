{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    zimeg.url = "github:zimeg/nur-packages";
  };
  outputs =
    { nixpkgs, ... }@inputs:
    let
      each =
        function:
        nixpkgs.lib.genAttrs [
          "x86_64-darwin"
          "x86_64-linux"
          "aarch64-darwin"
          "aarch64-linux"
        ] (system: function nixpkgs.legacyPackages.${system});
    in
    {
      devShells = each (pkgs: {
        default = pkgs.mkShell {
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
        gon =
          if pkgs.stdenv.isDarwin then
            pkgs.mkShell {
              packages = with pkgs; [
                go
                inputs.zimeg.packages.${pkgs.system}.gon
                goreleaser
              ];
              shellHook = ''
                export PATH=/usr/bin:$PATH # https://github.com/zimeg/nur-packages/issues/4
              '';
            }
          else
            null;
        tom = pkgs.mkShell {
          packages = with pkgs; [
            time
          ];
        };
      });
      packages = each (pkgs: {
        default = pkgs.buildGoModule rec {
          pname = "etime";
          version = "unversioned";
          src = ./.;
          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];
          doCheck = true;
          vendorHash = "sha256-TIVUxxozRGrpkVRFHvSJCCVlErMHj7ThBtCSPOmYWXU=";
          nativeBuildInputs = [
            pkgs.installShellFiles
          ];
          installPhase = ''
            mkdir -p $out/bin
            mv $GOPATH/bin/emporia-time $out/bin/etime
            installManPage etime.1
          '';
          meta = {
            description = "an energy aware time command";
            homepage = "https://github.com/zimeg/emporia-time";
            license = pkgs.lib.licenses.mit;
          };
        };
      });
    };
}
