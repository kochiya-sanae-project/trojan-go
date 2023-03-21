package postgres

import "github.com/p4gefau1t/trojan-go/config"

type PostgresConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Url     string `json:"url" yaml:"url"`
	//ServerHost string `json:"server_addr" yaml:"server-addr"`
	//ServerPort int    `json:"server_port" yaml:"server-port"`
	//Database   string `json:"database" yaml:"database"`
	//Username   string `json:"username" yaml:"username"`
	//Password   string `json:"password" yaml:"password"`
	CheckRate int `json:"check_rate" yaml:"check-rate"`
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
