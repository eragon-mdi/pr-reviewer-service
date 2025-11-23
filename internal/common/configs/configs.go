package configs

import (
	"log"
	"time"
)

type Config struct {
	Storages      Storages      `envconfig:"STORAGES" required:"true"`
	Servers       Servers       `envconfig:"SERVERS" required:"true"`
	Logger        Logger        `envconfig:"LOGGER" required:"true"`
	BussinesLogic BussinesLogic `envconfig:"BUSSINES_LOGIC" required:"true"`
}

func MustLoad() *Config {
	if err := loadCfg(singleConfig); err != nil {
		log.Fatalf("failed load cfgs on init stage: %v", err)
	}
	return singleConfig
}

var singleConfig = &Config{}

func Get() *Config {
	return singleConfig
}

type Storages struct {
	Postgres PsqlStore `envconfig:"POSTGRES" required:"true"`
}

type PsqlStore struct {
	HostF     string `envconfig:"HOST" required:"true"`
	PortF     string `envconfig:"PORT" default:"5432"`
	UserF     string `envconfig:"USER" required:"true"`
	PasswordF string `envconfig:"PASS" required:"true"`
	NameF     string `envconfig:"NAME" required:"true"`
	SSLmodeF  string `envconfig:"SSLM" default:"disable"`
}

type Servers struct {
	REST Server `envconfig:"REST"`
}

type Server struct {
	AddressF           string        `envconfig:"ADDR" required:"true"`
	PortF              string        `envconfig:"PORT" required:"true"`
	ReadTimeoutF       time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeoutF      time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	ReadHeaderTimeoutF time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeoutF       time.Duration `envconfig:"IDLE_TIMEOUT" default:"5s"`

	HealthCheckRoute string `envconfig:"HEALTH_CHECK_ROUTE" default:"health"`
}

type Logger struct {
	Level      string `envconfig:"LEVEL" default:"error"`
	Encoding   string `envconfig:"ENCODING" default:"json"`
	Output     string `envconfig:"OUTPUT" default:"stdout"`
	MessageKey string `envconfig:"MESSAGE_KEY" default:"message"`
}

type BussinesLogic struct {
	AllowedReuseToReasign   bool     `envconfig:"ALLOWED_REUSE_TO_REASIGN" default:"false"`
	AlloweStatusesToReasign []string `envconfig:"ALLOWE_STATUSES_TO_REASIGN"`
	AllowedRolesToReasign   []string `envconfig:"ALLOWED_ROLES_TO_REASIGN"`
}
