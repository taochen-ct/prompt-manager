package config

import "time"

type Config struct {
	Server struct {
		Env       string `mapstructure:"env" yaml:"env"`
		Port      string `mapstructure:"port" yaml:"port"`
		ApiPrefix string `mapstructure:"apiPrefix" yaml:"apiPrefix"`
		Storage   string `mapstructure:"storage" yaml:"storage"`
		Buffer    int    `mapstructure:"buffer" yaml:"buffer"`
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
		SecretKey       string `mapstructure:"secretKey" yaml:"secretKey"`
		TokenExpireHour int    `mapstructure:"tokenExpireHour" yaml:"tokenExpireHour"`
	} `mapstructure:"security" yaml:"security"`
	DB    DBConfig `mapstructure:"db" yaml:"db"`
	Proxy Proxy    `mapstructure:"proxy" yaml:"proxy"`
	Else  Else     `mapstructure:"else" yaml:"else"`
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

type Proxy struct {
	Server struct {
		Host         string        `mapstructure:"host" yaml:"host"`
		Port         int           `mapstructure:"port" yaml:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
	} `mapstructure:"server" yaml:"server"`

	Limits struct {
		NonStreamConcurrency int `mapstructure:"non_stream_concurrency" yaml:"non_stream_concurrency"`
		StreamConcurrency    int `mapstructure:"stream_concurrency" yaml:"stream_concurrency"`
	} `mapstructure:"limits" yaml:"limits"`

	HttpClient struct {
		MaxIdleConns        int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
		MaxIdleConnsPerHost int           `mapstructure:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
		MaxConnsPerHost     int           `mapstructure:"max_conns_per_host" yaml:"max_conns_per_host"`
		IdleConnTimeout     time.Duration `mapstructure:"idle_conn_timeout" yaml:"idle_conn_timeout"`
	} `mapstructure:"http_client" yaml:"http_client"`

	Models map[string]struct {
		Endpoints []struct {
			Name    string `mapstructure:"name" yaml:"name"`
			ApiBase string `mapstructure:"api_base" yaml:"api_base"`
			ApiKey  string `mapstructure:"api_key" yaml:"api_key"`
		} `mapstructure:"endpoints" yaml:"endpoints"`
	} `mapstructure:"models" yaml:"models"`
}

type Else struct {
	ScSend struct {
		Enable bool   `mapstructure:"enable" yaml:"enable"`
		Key    string `mapstructure:"key" yaml:"key"`
	} `mapstructure:"scSend" yaml:"scSend"`
}
