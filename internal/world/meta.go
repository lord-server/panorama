package world

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Meta map[string]string

func ParseMeta(path string) (Meta, error) {
	meta := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't read world metadata: %w", err)
	}

	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		meta[key] = value
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return meta, nil
}
