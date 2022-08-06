package world

import (
	"context"
	"errors"

	lru "github.com/hashicorp/golang-lru"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/weqqr/panorama/pkg/spatial"
)

type Backend interface {
	GetBlockData(pos spatial.BlockPos) ([]byte, error)
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

func (p *PostgresBackend) GetBlockData(pos spatial.BlockPos) ([]byte, error) {
	var data []byte
	err := p.conn.QueryRow(context.Background(), "SELECT data FROM blocks WHERE posx=$1 and posy=$2 and posz=$3", pos.X, pos.Y, pos.Z).Scan(&data)
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

func (w *World) GetBlock(pos spatial.BlockPos) (*MapBlock, error) {
	cachedBlock, ok := w.blockCache.Get(pos)

	if ok {
		if cachedBlock == nil {
			return nil, nil
		}
		return cachedBlock.(*MapBlock), nil
	}

	data, err := w.backend.GetBlockData(pos)
	if err != nil {
		return nil, err
	}

	if data == nil {
		w.blockCache.Add(pos, nil)
		return nil, nil
	}

	block, err := DecodeMapBlock(data)
	if err != nil {
		return nil, err
	}

	w.blockCache.Add(pos, block)

	return block, nil
}
