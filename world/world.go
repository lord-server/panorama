package world

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Backend interface {
	GetBlockData(x, y, z int) ([]byte, error)
	Close()
}

type PostgresBackend struct {
	conn *pgxpool.Pool
}

func NewPostgresBackend(dsn string) (*PostgresBackend, error) {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresBackend{
		conn: conn,
	}, nil
}

func (p *PostgresBackend) Close() {
	p.conn.Close()
}

func (p *PostgresBackend) GetBlockData(x, y, z int) ([]byte, error) {
	var data []byte
	err := p.conn.QueryRow(context.Background(), "SELECT data FROM blocks WHERE posx=$1 and posy=$2 and posz=$3", x, y, z).Scan(&data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

type World struct {
	backend Backend
}

func NewWorldWithBackend(backend Backend) World {
	return World{
		backend: backend,
	}
}

func (w *World) GetBlock(x, y, z int) (*MapBlock, error) {
	data, err := w.backend.GetBlockData(x, y, z)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	block, err := DecodeMapBlock(data)
	if err != nil {
		return nil, err
	}

	return block, nil
}
