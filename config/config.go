package config

import (
	"fmt"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Listen      string `yaml:"listen" mapstructure:"listen"`
	Target      string `yaml:"target" mapstructure:"target"`
	Log         string `yaml:"log" mapstructure:"log"`
	CertDir     string `yaml:"cert-dir" mapstructure:"cert-dir"`
	SslStriping bool   `yaml:"ssl-striping" mapstructure:"ssl-striping"`
}

func LoadConfig(arguments []string) (config Config, err error) {
	flags := flag.NewFlagSet("http-proxy-config", pflag.ExitOnError)

	configPath := flags.StringP("config", "c", "", "configuration file path")

	/* if config file not specified get config from cmd line arguments */
	cfg := Config{}
	flags.StringVarP(&cfg.Listen, "listen", "l", "0.0.0.0:80", "listen bind address")
	flags.StringVarP(&cfg.Target, "target", "t", "", "proxy target address")
	flags.StringVarP(&cfg.Log, "log", "o", "", "log file path")
	flags.StringVarP(&cfg.CertDir, "cert-dir", "d", "./", "CA certificate and private key directory")
	flags.BoolVarP(&cfg.SslStriping, "ssl-striping", "s", false, "HTTPS connection between client and proxy server")

	flags.Parse(arguments)

	if configPath != nil && *configPath != "" {
		/* load config from config file */
		viper.SetConfigFile(*configPath)
		if err := viper.ReadInConfig(); err != nil {
			return cfg, err
		}

		err = viper.Unmarshal(&cfg)
		if err != nil {
			return cfg, err
		}
	}

	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if c.Target == "" {
		return fmt.Errorf("target address must be specified")
	}

	return nil
}
