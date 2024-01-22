package dbase

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type TLS struct {
	ClientCert string
	ClientKey  string
	ServerCA   string
	ServerName string
}

type Config struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	TLS      TLS
}

func (cfg Config) Open() (*sqlx.DB, error) {
	connectString := ConnectString(cfg)
	//for logging database interactions globally
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	fmt.Println(newLogger)

	db, err := sqlx.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}

	// Set Max Open Connections and Max Idle Connections
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Ping the database to check the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectString(cfg Config) string {
	return fmt.Sprintf(`host=%v port=%v user=%v dbname=%v password=%v sslmode=disable`,
		cfg.Hostname,
		cfg.Port,
		cfg.Username,
		cfg.Database,
		cfg.Password,
	)
}

func RunMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Println("intilize error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"postgres", driver)
	if err != nil {
		log.Println("error while opening the file")
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("err91", err)
		return err
	}

	fmt.Println("*************Migrations ran successfully*************")
	return nil
}
