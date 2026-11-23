{
  description = "rssagg";
  inputs.nixpkgs.url = "https://channels.nixos.org/nixos-26.05/nixexprs.tar.zst";

  outputs = {
    self,
    nixpkgs,
  }: let
    lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
    version = builtins.substring 0 8 lastModifiedDate;
    supportedSystems = [
      "x86_64-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "aarch64-darwin"
    ];
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});
  in {
    packages = forAllSystems (
      system: let
        pkgs = nixpkgsFor.${system};
      in {
        rssagg = pkgs.buildGoModule {
          pname = "rssagg";
          inherit version;
        };
      }
    );

    devShells = forAllSystems (
      system: let
        pkgs = nixpkgsFor.${system};
      in {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            goose
            sqlc
            postgresql
          ];
          shellHook = ''
            export PGDATA="$PWD/.direnv/postgres"
            export PGSOCK="$PWD/.direnv/postgres_socket"
            export LOG_FILE="$PWD/.direnv/postgres.log"

            mkdir -p "$PGDATA" "$PGSOCK"

            if [ ! -d "$PGDATA/base" ]; then
              initdb --no-locale -D "$PGDATA"
            fi

            echo "unix_socket_directories = '$PGSOCK'" >> "$PGDATA/postgresql.conf"

            pg_ctl -o "-p 5432 -k $PGSOCK" -l "$LOG_FILE" -w start

            trap "pg_ctl stop -m fast" EXIT

            echo "Para conectar usa: psql -h \$PGSOCK"
          '';
        };
      }
    );
  };
}
