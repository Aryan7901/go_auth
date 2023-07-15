package constants

import (
	"embed"

	"github.com/jmoiron/sqlx"
)

type AppConfig struct {
	Db        *sqlx.DB
	Templates embed.FS
}
