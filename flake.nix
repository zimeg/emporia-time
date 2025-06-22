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
            gnumake # https://github.com/mirror/make
            go # https://github.com/golang/go
            go-junit-report # https://github.com/jstemmer/go-junit-report
            gocyclo # https://github.com/fzipp/gocyclo
            gofumpt # https://github.com/mvdan/gofumpt
            golangci-lint # https://github.com/golangci/golangci-lint
            gopls # https://github.com/golang/tools
            goreleaser # https://github.com/goreleaser/goreleaser
          ];
          shellHook = ''
            go mod tidy
          '';
        };
        gon =
          if pkgs.stdenv.isDarwin then
            pkgs.mkShell {
              packages = with pkgs; [
                go # https://github.com/golang/go
                inputs.zimeg.packages.${pkgs.system}.gon # https://github.com/Bearer/gon
                goreleaser # https://github.com/goreleaser/goreleaser
              ];
              shellHook = ''
                export PATH=/usr/bin:$PATH # https://github.com/zimeg/nur-packages/issues/4
              '';
            }
          else
            null;
        tom = pkgs.mkShell {
          packages = with pkgs; [
            time # https://git.savannah.gnu.org/cgit/time.git
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
          vendorHash = "sha256-a835B4+BraqOy1GH3XBYnBhiBjKEss71DoEroV9/0gA=";
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
