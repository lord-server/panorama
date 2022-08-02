package world

import (
	"context"
	"errors"

	lru "github.com/hashicorp/golang-lru"
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
	backend    Backend
	blockCache *lru.Cache
}

func NewWorldWithBackend(backend Backend) World {
	blockCache, err := lru.New(1024 * 16)
	if err != nil {
		panic(err)
	}
	return World{
		backend:    backend,
		blockCache: blockCache,
	}
}

func (w *World) GetBlock(x, y, z int) (*MapBlock, error) {
	type blockPos struct {
		x, y, z int
	}

	cachedBlock, ok := w.blockCache.Get(blockPos{x, y, z})

	if ok {
		if cachedBlock == nil {
			return nil, nil
		}
		return cachedBlock.(*MapBlock), nil
	}

	data, err := w.backend.GetBlockData(x, y, z)
	if err != nil {
		return nil, err
	}

	if data == nil {
		w.blockCache.Add(blockPos{x, y, z}, nil)
		return nil, nil
	}

	block, err := DecodeMapBlock(data)
	if err != nil {
		return nil, err
	}

	w.blockCache.Add(blockPos{x, y, z}, block)

	return block, nil
}
