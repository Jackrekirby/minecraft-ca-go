//go:build !js && !wasm
// +build !js,!wasm

package core

import (
	"fmt"
	"image"
	"math"
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	var totalDuration time.Duration = 0
	for i := 0; i < b.N; i++ {
		startTime := time.Now()
		elapsed := time.Since(startTime)
		totalDuration += elapsed
	}
}

func BenchmarkNanotime(b *testing.B) {
	var totalDuration int64 = 0
	for i := 0; i < b.N; i++ {
		startTime := nanotime()
		elapsed := nanotime() - startTime
		totalDuration += elapsed
	}
}

func BenchmarkDrawTriangle3D(b *testing.B) {
	width, height := 512, 512
	sceneImage := image.NewRGBA(image.Rect(0, 0, width, height))
	depthBuffer := make(DepthBuffer, width*height)
	clearDepthBuffer(&depthBuffer)
	texImage, err := loadImage("crafting_table_front.png")
	if err != nil {
		panic(fmt.Sprintf("failed to load texture: %v", err))
	}
	texImageRGBA := ImageToRGBA(texImage)
	camera := Camera{
		Position:    Point3D{0.5, 0.5, 1},
		Rotation:    Point3D{0, DegToRad(180), 0},
		FOV:         90.0,
		AspectRatio: 1.0,
		Near:        0.1,
		Far:         100.0,
	}
	clr := Red.ToRGBA()
	for i := 0; i < b.N; i++ {
		DrawTriangle3D(
			Vertex{Point3D{0, 0, 0}, 16, 16},
			Vertex{Point3D{1, 0, 0}, 0, 16},
			Vertex{Point3D{1, 1, 0}, 0, 0},
			camera, sceneImage, clr, &depthBuffer, texImageRGBA,
		)
	}
	SaveImage(sceneImage, CreateProjectRelativePath("output/benchmark.png"))
}

func benchmarkSinFunc(b *testing.B, sinFunc func(float64) float64, name string) {
	b.Run(name, func(b *testing.B) {
		var total float64 = 0
		maxAngle := DegToRad(3600.0)
		deltaAngle := DegToRad(0.1)
		for i := 0; i < b.N; i++ {
			for j := 0.0; j < maxAngle; j += deltaAngle {
				total += sinFunc(j)
			}
		}
		// Prevent compiler optimizations
		_ = total
	})
}

func BenchmarkSins(b *testing.B) {
	InitSinTable()
	benchmarkSinFunc(b, math.Sin, "math.Sin")
	benchmarkSinFunc(b, FastSin2, "FastSin2")
	benchmarkSinFunc(b, FastSin, "FastSin")

}
