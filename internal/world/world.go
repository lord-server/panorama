package world

import (
	"context"
	"errors"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lord-server/panorama/internal/spatial"
)

type Backend interface {
	GetBlockData(pos spatial.BlockPosition) ([]byte, error)
	Close()
}

type PostgresBackend struct {
	conn *pgxpool.Pool
}

func NewPostgresBackend(dsn string) (*PostgresBackend, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
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

func (p *PostgresBackend) GetBlockData(pos spatial.BlockPosition) ([]byte, error) {
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
	backend           Backend
	decodedBlockCache *lru.Cache[spatial.BlockPosition, *MapBlock]
}

func NewWorldWithBackend(backend Backend) World {
	decodedBlockCache, err := lru.New[spatial.BlockPosition, *MapBlock](1024 * 16)
	if err != nil {
		panic(err)
	}
	return World{
		backend:           backend,
		decodedBlockCache: decodedBlockCache,
	}
}

func (w *World) GetBlock(pos spatial.BlockPosition) (*MapBlock, error) {
	cachedBlock, ok := w.decodedBlockCache.Get(pos)

	if ok {
		if cachedBlock == nil {
			return nil, nil
		}
		return cachedBlock, nil
	}

	data, err := w.backend.GetBlockData(pos)
	if err != nil {
		return nil, err
	}

	if data == nil {
		w.decodedBlockCache.Add(pos, nil)
		return nil, nil
	}

	block, err := DecodeMapBlock(data)
	if err != nil {
		return nil, err
	}

	w.decodedBlockCache.Add(pos, block)

	return block, nil
}
