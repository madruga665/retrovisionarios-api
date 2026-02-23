package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DbPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar string de conexão: %w", err)
	}

	// Configurações recomendadas para o pool
	config.MaxConns = 10
	config.MinConns = 2

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar pool de conexões: %w", err)
	}

	// Verifica se a conexão está ativa
	if err := dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("erro ao pingar o banco de dados: %w", err)
	}

	return dbPool, nil
}
