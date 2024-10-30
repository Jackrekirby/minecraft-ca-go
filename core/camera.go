package core

import "math"

// Camera represents the camera in 3D space, with position, rotation, and projection parameters
type Camera struct {
	Position    Point3D
	Rotation    Point3D
	FOV         float64
	AspectRatio float64
	Near        float64
	Far         float64
}

// Convert3DTo2D projects a 3D point onto a 2D plane using perspective projection, considering the camera's position and orientation
func (c *Camera) Convert3DTo2D(point Point3D) Point2D {
	// Translate the point by the camera position
	p := point.
		Subtract(c.Position).
		RotateX(c.Rotation.X).
		RotateY(c.Rotation.Y).
		RotateZ(c.Rotation.Z)

	// Check if the point is behind the camera
	if p.Z <= 0 {
		return Point2D{X: 0, Y: 0}
	}

	// Perspective projection
	fovRad := 1.0 / math.Tan(c.FOV*0.5*math.Pi/180)
	q := c.Far / (c.Far - c.Near)

	ndcX := p.X * fovRad * c.AspectRatio
	ndcY := p.Y * fovRad
	ndcZ := p.Z * q

	// Project to screen coordinates
	screenX := ndcX / ndcZ
	screenY := ndcY / ndcZ

	return Point2D{X: screenX, Y: screenY}
}

// Convert3DTo2D projects a 3D point onto a 2D plane using perspective projection.
func Convert3DTo2D(point Point3D, fov, aspectRatio, near, far float64) Point2D {
	fovRad := 1.0 / math.Tan(fov*0.5*math.Pi/180)
	q := far / (far - near)

	ndcX := point.X * fovRad * aspectRatio
	ndcY := point.Y * fovRad
	ndcZ := point.Z * q

	if point.Z <= 0 {
		return Point2D{X: 0, Y: 0}
	}

	screenX := ndcX / ndcZ
	screenY := ndcY / ndcZ

	return Point2D{X: screenX, Y: screenY}
}
