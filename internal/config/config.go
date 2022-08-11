package config

import "fmt"

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	LogLevel       string
	PercentExclude int
	Environment    string
	DB             DB
	Server         Server
	Rmq            Rmq
}

type DB struct {
	Type string // "mem", "sql"
	SQL  SQLDatabase
}

type SQLDatabase struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type Server struct {
	Port string
	Host string
}

type Rmq struct {
	Host string
	Port string
	User string
	Pswd string
}

func GetDsn(s SQLDatabase) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", s.User, s.Password, s.Host, s.Port, s.Name)
}

func NewConfig() *Config {
	return &Config{}
}
