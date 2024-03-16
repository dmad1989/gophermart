package config

import (
	"flag"
	"os"
)

const (
	defHost = "localhost:8080"
)

type Config struct {
	ApiURL     string
	AccrualURL string
	DbConnName string
}

var LoggerCtxKey = &ContextKey{"logger"}

type ContextKey struct {
	name string
}

var conf = Config{}

func init() {
	flag.StringVar(&conf.ApiURL, "a", defHost, "server URL format host:port, :port")
	flag.StringVar(&conf.AccrualURL, "r", defHost, "URL for accrual system format host:port, :port")
	flag.StringVar(&conf.DbConnName, "d", "", "database connection addres, format host=? port=? user=? password=? dbname=? sslmode=?")
}

func ParseConfig() Config {
	flag.Parse()
	if os.Getenv("RUN_ADDRESS") != "" {
		conf.ApiURL = os.Getenv("RUN_ADDRESS")
	}

	if os.Getenv("DATABASE_URI") != "" {
		conf.DbConnName = os.Getenv("DATABASE_URI")
	}

	if os.Getenv("ACCRUAL_SYSTEM_ADDRESS") != "" {
		conf.AccrualURL = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	}

	return conf
}
