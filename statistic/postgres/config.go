package postgres

import "github.com/p4gefau1t/trojan-go/config"

type PostgresConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Url       string `json:"url" yaml:"url"`
	CheckRate int    `json:"check_rate" yaml:"check-rate"`
}

type Config struct {
	Postgres PostgresConfig `json:"postgres" yaml:"postgres"`
}

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return &Config{
			Postgres: PostgresConfig{
				CheckRate: 30,
			},
		}
	})
}
