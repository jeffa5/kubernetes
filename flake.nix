{
  description = "Kubernetes";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { system = system; };
        in
        rec
        {
          devShell = pkgs.mkShell {
            buildInputs = with pkgs;[
              protobuf
              kubectl
              k9s

              go
              etcd_3_4
              kind
              openssl

              rnix-lsp
              nixpkgs-fmt
            ];
          };
        }
      );
}
