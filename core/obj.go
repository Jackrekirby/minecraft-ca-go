package core

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// Vertex represents a 3D vertex.
type Vertex3 Point3D

// TextureCoord represents a texture coordinate.
type TextureCoord struct {
	U, V float64
}

// Normal represents a normal vector.
type Normal Point3D

// Face represents a single face in the OBJ file.
type Face struct {
	VertexIndices  []int
	TextureIndices []int
	NormalIndices  []int
}

// OBJData stores parsed data from an OBJ file.
type OBJData struct {
	Vertices      []Vertex3
	TextureCoords []TextureCoord
	Normals       []Normal
	Faces         []Face
}

// ParseOBJFile reads an OBJ file and returns parsed data.
func ParseOBJFile(filePath string) (*OBJData, error) {
	file, err := OpenAsset(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	objData := &OBJData{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// Skip empty lines and comments
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "v": // Vertex
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			objData.Vertices = append(objData.Vertices, Vertex3{x, y, z})

		case "vt": // Texture coordinate
			u, _ := strconv.ParseFloat(fields[1], 64)
			v, _ := strconv.ParseFloat(fields[2], 64)
			objData.TextureCoords = append(objData.TextureCoords, TextureCoord{u, v})

		case "vn": // Normal
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			objData.Normals = append(objData.Normals, Normal{x, y, z})

		case "f": // Face
			face := Face{}
			for _, part := range fields[1:] {
				indices := strings.Split(part, "/")

				// Parse vertex index
				if len(indices) > 0 && indices[0] != "" {
					vIdx, _ := strconv.Atoi(indices[0])
					if vIdx < 0 {
						vIdx = len(objData.Vertices) + vIdx + 1
					}
					face.VertexIndices = append(face.VertexIndices, vIdx-1)
				}

				// Parse texture index
				if len(indices) > 1 && indices[1] != "" {
					vtIdx, _ := strconv.Atoi(indices[1])
					if vtIdx < 0 {
						vtIdx = len(objData.TextureCoords) + vtIdx + 1
					}
					face.TextureIndices = append(face.TextureIndices, vtIdx-1)
				}

				// Parse normal index
				if len(indices) > 2 && indices[2] != "" {
					vnIdx, _ := strconv.Atoi(indices[2])
					if vnIdx < 0 {
						vnIdx = len(objData.Normals) + vnIdx + 1
					}
					face.NormalIndices = append(face.NormalIndices, vnIdx-1)
				}
			}
			objData.Faces = append(objData.Faces, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return objData, nil
}

// PrintOBJData prints parsed OBJ data for debugging.
func PrintOBJData(objData *OBJData) {
	fmt.Println("Vertices:")
	for _, v := range objData.Vertices {
		fmt.Printf("  %f %f %f\n", v.X, v.Y, v.Z)
	}

	fmt.Println("\nTexture Coordinates:")
	for _, vt := range objData.TextureCoords {
		fmt.Printf("  %f %f\n", vt.U, vt.V)
	}

	fmt.Println("\nNormals:")
	for _, vn := range objData.Normals {
		fmt.Printf("  %f %f %f\n", vn.X, vn.Y, vn.Z)
	}

	fmt.Println("\nFaces:")
	for _, f := range objData.Faces {
		fmt.Printf("  Vertices: %v, Textures: %v, Normals: %v\n", f.VertexIndices, f.TextureIndices, f.NormalIndices)
	}
}

func RunOBJTest() {
	// Replace with your OBJ file path
	filePath := "model.obj"

	objData, err := ParseOBJFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	PrintOBJData(objData)
}
