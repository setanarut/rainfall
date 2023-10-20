package rainfall

import "math"

type vec2 struct {
	X, Y float64
}

func (v *vec2) Add(other vec2) {
	v.X += other.X
	v.Y += other.Y
}

func (v *vec2) Mul(other vec2) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *vec2) Div(other vec2) {
	v.X /= other.X
	v.Y /= other.Y
}

func (v *vec2) AddS(scalar float64) {
	v.X += scalar
	v.Y += scalar
}
func (v *vec2) MulS(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *vec2) DivS(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
}

func (v *vec2) Length() float64 {
	return float64(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// VECTOR 3

type vec3 struct {
	X, Y, Z float64
}

func (v *vec3) Add(other vec3) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

// Multiple other Vec3
func (v *vec3) Mul(other vec3) {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
}

func (v *vec3) Div(other vec3) {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
}

// Multiple scalar
func (v *vec3) MulS(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
}

// MulR multiplies by other vector and returns a new one
func (v vec3) MulR(other vec3) vec3 {
	return vec3{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

// AddR adds with other and returns the new one
func (v vec3) AddR(other vec3) vec3 {
	return vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v *vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *vec3) Normalize() {
	length := v.Length()
	v.X /= length
	v.Y /= length
	v.Z /= length
}
