package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

// Migrations contém todas as migrations descobertas nos arquivos .sql
// deste diretório, seguindo o padrão <timestamp>_<nome>.up.sql / .down.sql.
var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
