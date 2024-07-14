package config

type (
	Config struct {
		App
		HTTP
		Log
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
)

func New() *Config {
	return &Config{
		App: App{
			Name:    "realworld-fiber-sqlc",
			Version: "v1.0.0",
		},
		HTTP: HTTP{
			Port: ":3000",
		},
		Log: Log{
			Level: "info",
		},
	}
}
