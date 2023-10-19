package rainfall

import "math"

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Add(other Vec2) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2) Mul(other Vec2) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *Vec2) Div(other Vec2) {
	v.X /= other.X
	v.Y /= other.Y
}

func (v *Vec2) AddS(scalar float64) {
	v.X += scalar
	v.Y += scalar
}
func (v *Vec2) MulS(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2) DivS(scalar float64) {
	v.X /= scalar
	v.Y /= scalar
}

func (v *Vec2) Length() float64 {
	return float64(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// VECTOR 3

type Vec3 struct {
	X, Y, Z float64
}

func (v *Vec3) Add(other Vec3) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

// Multiple other Vec3
func (v *Vec3) Mul(other Vec3) {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
}

func (v *Vec3) Div(other Vec3) {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
}

// Multiple scalar
func (v *Vec3) MulS(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
}

// MulR returns new multiple Vec3
func (v Vec3) MulR(other Vec3) Vec3 {
	return Vec3{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

// AddR returns new added Vec3
func (v Vec3) AddR(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vec3) Normalize() {
	length := v.Length()
	v.X /= length
	v.Y /= length
	v.Z /= length
}
