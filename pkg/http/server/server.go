package server

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/papidb/aqua/pkg/config"
)

type Server struct {
	App *config.App
}
