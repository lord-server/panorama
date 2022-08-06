package world

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/weqqr/panorama/pkg/spatial"
)

const MapBlockSize = 16
const MapBlockVolume = MapBlockSize * MapBlockSize * MapBlockSize
const NodeSizeInBytes = 4

type Node struct {
	ID     uint16
	Param1 uint8
	Param2 uint8
}

func readU8(r io.Reader) (uint8, error) {
	var value uint8
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func readU16(r io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(r, binary.BigEndian, &value)
	return value, err
}

func readString(r io.Reader) (string, error) {
	length, err := readU16(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

type MapBlock struct {
	mappings map[uint16]string
	nodeData []byte
}

type ReaderCounter struct {
	inner *bytes.Reader
	count int64
}

func NewReaderCounter(r *bytes.Reader) *ReaderCounter {
	return &ReaderCounter{
		inner: r,
		count: 0,
	}
}

func (r *ReaderCounter) Read(p []byte) (n int, err error) {
	n, err = r.inner.Read(p)
	r.count += int64(n)
	return
}

func (r *ReaderCounter) ReadByte() (byte, error) {
	b, err := r.inner.ReadByte()
	r.count += 1
	return b, err
}

func inflate(reader *bytes.Reader) ([]byte, error) {
	position, _ := reader.Seek(0, io.SeekCurrent)

	counter := NewReaderCounter(reader)
	z, err := zlib.NewReader(counter)
	if err != nil {
		panic(err)
	}
	defer z.Close()

	data, err := io.ReadAll(z)
	if err != nil {
		panic(err)
	}

	_, err = reader.Seek(position+counter.count, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return data, err
}

func readMappings(reader *bytes.Reader) (map[uint16]string, error) {
	mappingCount, err := readU16(reader)
	if err != nil {
		return nil, err
	}

	mappings := make(map[uint16]string)
	for i := 0; i < int(mappingCount); i++ {
		id, err := readU16(reader)
		if err != nil {
			return nil, err
		}
		name, err := readString(reader)
		if err != nil {
			return nil, err
		}

		mappings[id] = name
	}

	return mappings, nil
}

func decodeLegacyBlock(reader *bytes.Reader, version uint8) (*MapBlock, error) {
	if version >= 27 {
		// - uint8 flags
		// - uint16 lighting_complete
		// - uint8 content_width
		// - uint8 params_width
		_, err := reader.Seek(1+2+1+1, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	} else {
		// - uint8 flags
		// - uint8 content_width
		// - uint8 params_width
		_, err := reader.Seek(1+1+1, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}

	nodeData, err := inflate(reader)
	if err != nil {
		panic(err)
	}

	_, err = inflate(reader)
	if err != nil {
		panic(err)
	}

	// - uint8 staticObjectVersion
	_, err = reader.Seek(1, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	staticObjectCount, err := readU16(reader)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(staticObjectCount); i++ {
		// - uint8 type
		// - int32 x, y, z
		_, err = reader.Seek(1+4+4+4, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		dataSize, err := readU16(reader)
		if err != nil {
			panic(err)
		}
		_, err = reader.Seek(int64(dataSize), io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}

	// - uint32 timestamp
	// - uint8 mappingVersion
	_, err = reader.Seek(4+1, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	mappings, err := readMappings(reader)
	if err != nil {
		return nil, err
	}

	return &MapBlock{
		mappings: mappings,
		nodeData: nodeData,
	}, nil
}

func decodeBlock(reader *bytes.Reader) (*MapBlock, error) {
	z, err := zstd.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	data, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}

	reader = bytes.NewReader(data)

	// Skip:
	// - uint8 flags
	// - uint16 lighting_complete
	// - uint32 timestamp
	// - uint8 mapping version
	_, err = reader.Seek(1+2+4+1, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	mappings, err := readMappings(reader)
	if err != nil {
		return nil, err
	}

	// Skip uint8 contentWidth, uint8 paramsWidth
	_, err = reader.Seek(1+1, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	nodeData := make([]byte, MapBlockVolume*NodeSizeInBytes)
	_, err = io.ReadFull(reader, nodeData)
	if err != nil {
		return nil, err
	}

	return &MapBlock{
		mappings: mappings,
		nodeData: nodeData,
	}, nil
}

func DecodeMapBlock(data []byte) (*MapBlock, error) {
	reader := bytes.NewReader(data)

	version, err := readU8(reader)
	if err != nil {
		return nil, err
	}

	if version < 29 {
		mapblock, err := decodeLegacyBlock(reader, version)
		if err != nil {
			panic(err)
		}
		return mapblock, nil
	}

	return decodeBlock(reader)
}

func (b *MapBlock) ResolveName(id uint16) string {
	return b.mappings[id]
}

func (b *MapBlock) GetNode(pos spatial.NodePos) Node {
	index := pos.Z*MapBlockSize*MapBlockSize + pos.Y*MapBlockSize + pos.X
	idHi := uint16(b.nodeData[2*index])
	idLo := uint16(b.nodeData[2*index+1])
	param1 := b.nodeData[2*MapBlockVolume+index]
	param2 := b.nodeData[3*MapBlockVolume+index]
	return Node{
		ID:     (idHi << 8) | idLo,
		Param1: param1,
		Param2: param2,
	}
}
