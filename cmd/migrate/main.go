package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"task-api/pkg/config"
)

func main() {
	var cfg config.AppConfig
	if err := cfg.ReadEnvConfig(); err != nil {
		log.Println("read env config failed", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Println("validate config failed", err)
	}

	connect := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s",
		cfg.MainStorage.Postgres.User,
		cfg.MainStorage.Postgres.Password,
		cfg.MainStorage.Postgres.Host,
		cfg.MainStorage.Postgres.Port,
		cfg.MainStorage.Postgres.DBName,
		cfg.MainStorage.Postgres.SSLMode,
	)

	m, err := migrate.New("file://migrations", connect)
	if err != nil {
		log.Println(err)
	}

	if err := Up(m); err != nil {
		log.Println(err)
	}

	//if err := Down(m); err != nil {
	//	log.Println("Ошибка при откате миграций", err)
	//}

	errSource, errDB := m.Close()
	if errSource != nil {
		log.Println("Ошибка при закритии источника миграций:", errSource)
	}
	if errDB != nil {
		log.Println("Ошибка при закрытии подключения к БД:", errDB)
	}
}

func Up(m *migrate.Migrate) error {
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Миграций для применения не обнаружено")
			return err
		}
		log.Printf("Ошибка при применении миграций: %v", err)
		return err
	}
	log.Println("Миграции успешно применены")
	return nil
}

func Down(m *migrate.Migrate) error {
	if err := m.Down(); err != nil {
		if err == migrate.ErrNoChange {
			return err
		}
		return err
	}
	log.Println("Миграции успешно отменены")
	return nil
}
