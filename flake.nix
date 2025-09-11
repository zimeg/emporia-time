{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    zimeg = {
      url = "github:zimeg/nur-packages";
      inputs.nixpkgs.follows = "nixpkgs";
    };
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
        tom = pkgs.mkShell {
          packages = with pkgs; [
            go # https://github.com/golang/go
            goreleaser # https://github.com/goreleaser/goreleaser
            inputs.zimeg.packages.${pkgs.system}.quill # https://github.com/anchore/quill
            sops # https://github.com/getsops/sops
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
          vendorHash = "sha256-co+4DhzaNyswF1MgMQNlHDXkYVPgsINhL77BufiG7Xo=";
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
