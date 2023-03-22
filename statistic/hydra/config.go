package hydra

import "github.com/p4gefau1t/trojan-go/config"

type HydraConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	BaseUrl   string `json:"base_url" yaml:"base-url"`
	CheckRate int    `json:"check_rate" yaml:"check-rate"`
}

type Config struct {
	Hydra HydraConfig `json:"hydra" yaml:"hydra"`
}

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return &Config{
			Hydra: HydraConfig{
				CheckRate: 30,
			},
		}
	})
}
