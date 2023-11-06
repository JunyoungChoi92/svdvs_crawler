package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

type Database struct {
	conn *sqlx.DB
}

func (database *Database) Connect(path string) error {
	config, err := loadConfig(path)
	if err != nil {
		return err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode)

	database.conn, err = sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return err
	}

	log.Println("Successfully connected!")
	return nil
}

func (database *Database) Disconnect() error {
	err := database.conn.Close()
	if err != nil {
		return err
	}

	log.Print("Successfully disconnected!")
	return nil
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
