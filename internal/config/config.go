package config

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	Dsn      string `mapstructure:"dsn"`
}

type JwtConfig struct {
	AccessTokenTTL   int    `mapstructure:"access_token_ttl_sec"`
	RefreshTokenTTL  int    `mapstructure:"refresh_token_ttl_sec"`
	AccessSecretKey  string `mapstructure:"access_secret_key"`
	RefreshSecretKey string `mapstructure:"refresh_secret_key"`
}

type Config struct {
	App AppConfig `mapstructure:"app"`
	Log LogConfig `mapstructure:"log"`
	DB  DBConfig  `mapstructure:"db"`
	Jwt JwtConfig `mapstructure:"jwt"`
}
