package config

type Config struct {
	Server struct {
		Env       string `mapstructure:"env" yaml:"env"`
		Port      string `mapstructure:"port" yaml:"port"`
		ApiPrefix string `mapstructure:"apiPrefix" yaml:"apiPrefix"`
		Storage   string `mapstructure:"storage" yaml:"storage"`
	} `mapstructure:"server" yaml:"server"`
	Log struct {
		Level      string `mapstructure:"level" yaml:"level"`
		RootDir    string `mapstructure:"rootDir" yaml:"rootDir"`
		Filename   string `mapstructure:"filename" yaml:"filename"`
		MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
		MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize"` // MB
		MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge"`   // day
		Compress   bool   `mapstructure:"compress" yaml:"compress"`
	} `mapstructure:"log" yaml:"log"`
	Web struct {
		StaticDir   string `mapstructure:"staticDir" yaml:"staticDir"`
		DefaultHtml string `mapstructure:"defaultHtml" yaml:"defaultHtml"`
	} `mapstructure:"web" yaml:"web"`
	Security struct {
		SecretKey string `mapstructure:"secretKey" yaml:"secretKey"`
	} `mapstructure:"security" yaml:"security"`
	DB DBConfig `mapstructure:"db" yaml:"db"`
}

type DBConfig struct {
	Driver   string         `mapstructure:"driver" yaml:"driver"` // sqlite | mysql | postgres
	SQLite   SQLiteConfig   `mapstructure:"sqlite" yaml:"sqlite"`
	MySQL    MySQLConfig    `mapstructure:"mysql" yaml:"mysql"`
	Postgres PostgresConfig `mapstructure:"postgres" yaml:"postgres"`
	Migrate  struct {
		Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
		Dir     string `mapstructure:"dir" yaml:"dir"`
	} `mapstructure:"migrate" yaml:"migrate"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path" yaml:"path"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	Params   string `mapstructure:"params" yaml:"params"` // charset=utf8mb4&parseTime=true
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}
