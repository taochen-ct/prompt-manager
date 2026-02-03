package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	if strings.TrimSpace(cfg.Server.Storage) == "" {
		cfg.Server.Storage = "storage"
		cfg.Log.RootDir = fmt.Sprintf("./%s/%s", cfg.Server.Storage, cfg.Log.RootDir)
		if cfg.DB.Driver == "sqlite" {
			cfg.DB.SQLite.Path = fmt.Sprintf("./%s/%s", cfg.Server.Storage, cfg.DB.SQLite.Path)
		}
		log.Printf("default storage dir: %+v", cfg.Server.Storage)
	}

	if strings.TrimSpace(cfg.Web.StaticDir) == "" {
		cfg.Web.StaticDir = "web"
		log.Printf("default static dir: %+v", cfg.Web.StaticDir)
	}

	if strings.TrimSpace(cfg.Web.DefaultHtml) == "" {
		cfg.Web.DefaultHtml = "index.html"
		log.Printf("default html: %+v", cfg.Web.DefaultHtml)
	}

	return &cfg
}
