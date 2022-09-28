package app

import (
	"github.com/tzvetkoff-go/pasteur/pkg/db"
	"github.com/tzvetkoff-go/pasteur/pkg/hasher"
	"github.com/tzvetkoff-go/pasteur/pkg/logger"
	"github.com/tzvetkoff-go/pasteur/pkg/webserver"
)

// DefaultConfigPath ...
const DefaultConfigPath = `./config.yml`

// Config ...
type Config struct {
	Logger    logger.Config    `yaml:"logger"`
	DB        db.Config        `yaml:"db"`
	Hasher    hasher.Config    `yaml:"hasher"`
	WebServer webserver.Config `yaml:"webserver"`
}
