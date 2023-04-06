package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
)

const advertPetTable = "advert_pet"

type Config struct {
	DebugMode bool
	DB        *DB
	Listen    *Listen  `yaml:"listen"`
	Gateway   *Gateway `yaml:"gateway"`
	VersionDB int      `yaml:"db_version"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user_name"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name_db"`
}

type Listen struct {
	Port string `yaml:"port"`
	IP   string `yaml:"ip"`
}

type Gateway struct {
	Port  string `yaml:"port"`
	IP    string `yaml:"ip"`
	Label string `yaml:"label"`
}

func (db *DB) GetConnectionString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		db.Username, db.Password, db.Host, db.Port, db.DBName)
}

func NewPostgresDB(dbc DB) (*sqlx.DB, error) {
	// parse connection string
	dbConf, err := pgx.ParseConfig(dbc.GetConnectionString())
	if err != nil {
		return nil, errors.New("failed to parse config")
	}

	dbConf.Host = dbc.Host

	//register pgx conn
	dsn := stdlib.RegisterConnConfig(dbConf)

	sql.Register("wrapper", stdlib.GetDefaultDriver())
	wdb, err := sql.Open("wrapper", dsn)
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}

	const (
		maxOpenConns    = 50
		maxIdleConns    = 50
		connMaxLifetime = 3
		connMaxIdleTime = 5
	)
	db := sqlx.NewDb(wdb, "pgx")
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Minute)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetConfig() *Config {

	logger := log.Default()
	logger.Print("Read application configuration")
	instance := &Config{DB: &DB{}}

	instance.DB = &DB{
		Host:     getEnv("POSTGRES_HOST", ""),
		Port:     getEnv("POSTGRES_PORT", ""),
		DBName:   getEnv("POSTGRES_DB_NAME", ""),
		Username: getEnv("POSTGRES_USER", ""),
		Password: getEnv("POSTGRES_PASSWORD", ""),
	}

	return instance
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
