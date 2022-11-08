package storage

type Config struct {
	Host     string
	Port     string
	Possword string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection()
