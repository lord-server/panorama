package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Vertex struct {
	position Vector3
	texcoord Vector2
	normal   Vector3
}

type Mesh struct {
	vertices []Vertex
}

func NewMesh() Mesh {
	return Mesh{
		vertices: []Vertex{},
	}
}

func parseVector3(fields []string) (Vector3, error) {
	if len(fields) < 3 {
		return Vector3{}, fmt.Errorf("expected at least 3 vector elements, found %d", len(fields))
	}

	x, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return Vector3{}, err
	}
	y, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		return Vector3{}, err
	}
	z, err := strconv.ParseFloat(fields[2], 32)
	if err != nil {
		return Vector3{}, err
	}

	return NewVector3(float32(x), float32(y), float32(z)), nil
}

func parseVector2(fields []string) (Vector2, error) {
	if len(fields) < 2 {
		return Vector2{}, fmt.Errorf("expected at least 3 vector elements, found %d", len(fields))
	}

	x, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return Vector2{}, err
	}
	y, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		return Vector2{}, err
	}

	return NewVector2(float32(x), float32(y)), nil
}

type Triplet struct {
	positionIndex int
	texcoordIndex *int
	normalIndex   *int
}

func parseFace(fields []string) ([]Triplet, error) {
	if len(fields) < 3 {
		return []Triplet{}, fmt.Errorf("face needs at least 3 fields, got %v", len(fields))
	}

	triplets := []Triplet{}
	for _, field := range fields {
		parts := strings.SplitN(field, "/", 3)

		var err error
		triplet := Triplet{}

		triplet.positionIndex, err = strconv.Atoi(parts[0])
		if err != nil {
			return []Triplet{}, err
		}

		if len(parts) > 1 && len(parts[1]) != 0 {
			texcoordIndex, err := strconv.Atoi(parts[1])
			if err != nil {
				return []Triplet{}, err
			}
			triplet.texcoordIndex = &texcoordIndex
		}

		if len(parts) > 2 && len(parts[2]) != 0 {
			normalIndex, err := strconv.Atoi(parts[2])
			if err != nil {
				return []Triplet{}, err
			}
			triplet.normalIndex = &normalIndex
		}
		triplets = append(triplets, triplet)
	}

	return triplets, nil
}

type objParser struct {
	positions []Vector3
	texcoords []Vector2
	normals   []Vector3

	mesh Mesh
}

func (o *objParser) vertexAt(triplet Triplet) Vertex {
	texcoord := Vector2{}
	normal := Vector3{}

	if triplet.texcoordIndex != nil {
		texcoord = o.texcoords[*triplet.texcoordIndex-1]
	}

	if triplet.normalIndex != nil {
		normal = o.normals[*triplet.normalIndex-1]
	}

	return Vertex{
		position: o.positions[triplet.positionIndex-1],
		texcoord: texcoord,
		normal:   normal,
	}
}

func (o *objParser) triangulatePolygon(triplets []Triplet) []Vertex {
	vertices := []Vertex{}

	origin := o.vertexAt(triplets[0])

	for i := 2; i < len(triplets); i++ {
		vertices = append(vertices, origin, o.vertexAt(triplets[i-1]), o.vertexAt(triplets[i]))
	}

	return vertices
}

func (o *objParser) processLine(line string) error {
	// Skip comments
	if strings.HasPrefix(line, "#") {
		return nil
	}

	fields := strings.Fields(line)

	// Skip empty lines
	if len(fields) == 0 {
		return nil
	}

	switch fields[0] {
	case "v":
		position, err := parseVector3(fields[1:])
		if err != nil {
			return err
		}

		o.positions = append(o.positions, position)
	case "vt":
		texcoord, err := parseVector2(fields[1:])
		if err != nil {
			return err
		}

		o.texcoords = append(o.texcoords, texcoord)
	case "vn":
		normal, err := parseVector3(fields[1:])
		if err != nil {
			return err
		}

		o.normals = append(o.normals, normal)
	case "f":
		triplets, err := parseFace(fields[1:])
		if err != nil {
			return err
		}

		vertices := o.triangulatePolygon(triplets)

		o.mesh.vertices = append(o.mesh.vertices, vertices...)

	default:
		// log.Printf("unknown attribute %v; ignoring\n", fields[0])
	}

	return nil
}

func loadOBJ(path string) (Mesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return Mesh{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	parser := objParser{
		positions: []Vector3{},
		texcoords: []Vector2{},
		normals:   []Vector3{},
		mesh:      NewMesh(),
	}

	lineNumber := 1
	for scanner.Scan() {
		lineNumber += 1
		err := parser.processLine(scanner.Text())
		if err != nil {
			return Mesh{}, err
		}
	}

	if err := scanner.Err(); err != nil {
		return Mesh{}, err
	}

	return parser.mesh, nil
}

func Cube() *Mesh {
	mesh, err := loadOBJ("untitled.obj")
	if err != nil {
		log.Panic(err)
	}

	return &mesh
}
