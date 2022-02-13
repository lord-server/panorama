package world

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Backend interface {
	GetBlockData(x, y, z int) ([]byte, error)
	Close()
}

type PgBackend struct {
	conn *pgxpool.Pool
}

func NewPgBackend(dsn string) (*PgBackend, error) {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PgBackend{
		conn: conn,
	}, nil
}

func (p *PgBackend) Close() {
	p.conn.Close()
}

func (p *PgBackend) GetBlockData(x, y, z int) ([]byte, error) {
	var data []byte
	err := p.conn.QueryRow(context.Background(), "SELECT data FROM blocks WHERE posx=$1 and posy=$2 and posz=$3", x, y, z).Scan(&data)
	if err != nil {
		return []byte{}, err
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

	block, err := DecodeMapBlock(data)
	if err != nil {
		return nil, err
	}

	return block, nil
}
