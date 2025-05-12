package connectors

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"log"
	"task-api/pkg/config"
	"time"
)

type PostgresConnect struct {
	Pool *pgxpool.Pool
}

func NewPostgresConnect(cfg *config.PostgresConfig) (*PostgresConnect, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=10 pool_max_conn_lifetime=1h30m",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
	parseConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Println("Ошибка при парсенге конфига:", err)
		return nil, err
	}

	parseConfig.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithTrimSQLInSpanName(),
	)

	pool, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		log.Println("Ошибка при создание пула соединений:", err)
		return nil, err
	}

	driver, err := postgres.WithInstance(stdlib.OpenDB(*parseConfig.ConnConfig), &postgres.Config{
		StatementTimeout: time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("%s:%w", "postgres driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", cfg.DBName, driver)
	if err != nil {
		log.Println(err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Миграций для применения не обнаружено")
		}
		log.Printf("Ошибка при применении миграций: %v", err)
	}
	log.Println("Миграции успешно применены")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		log.Println("Ошибка при тестовом соединение с БД")
		return nil, err
	}
	log.Println("Подключение к БД через pgxpool установлено успешно")
	return &PostgresConnect{Pool: pool}, nil

}

func (pool *PostgresConnect) GetConnect() *pgxpool.Pool {
	return pool.Pool
}

func (pool *PostgresConnect) Close() error {
	pool.Pool.Close()
	return nil
}
