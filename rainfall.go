package rainfall

import (
	"math"
	"math/rand"
)

var rootTwo = math.Sqrt(2)
var v015 = Vec3{0.15, 0.15, 0.15}
var v01 = Vec3{0.1, 0.1, 0.1}

type Rainfall struct {
	dem             [][]float64
	sizeX, sizeY    int
	scale, density  float64
	friction        float64
	evaporationRate float64
	depositionRate  float64
	rand            *rand.Rand
}

func NewRainfall(dem [][]float64, randomSeed int64) *Rainfall {
	sizeX := len(dem[0])
	sizeY := len(dem)
	evaporRate := 1.0 / float64(sizeX)

	return &Rainfall{
		dem:             dem,
		sizeX:           sizeX,
		sizeY:           sizeY,
		scale:           100.0,
		density:         1.0,
		friction:        0.1,
		evaporationRate: evaporRate,
		depositionRate:  0.3,
		rand:            rand.New(rand.NewSource(randomSeed)),
	}
}

func (r *Rainfall) GetDem() [][]float64 {
	return r.dem
}

func (r *Rainfall) getSurfaceNormal(x, y int) Vec3 {

	// create the vector and add the 4 points directly adjacent to it
	surfaceNormal := v015
	surfaceNormal.Mul(Vec3{r.scale * (r.dem[y][x] - r.dem[y][x+1]), 1, 0})

	// 1 x-1
	surfaceNormal.Add(v015.MulR(Vec3{r.scale * (r.dem[y][x-1] - r.dem[y][x]), 1, 0}))
	// 2 y+1
	surfaceNormal.Add(v015.MulR(Vec3{0, 1, r.scale * (r.dem[y][x] - r.dem[y+1][x])}))
	// 3 y-1
	surfaceNormal.Add(v015.MulR(Vec3{0, 1, r.scale * (r.dem[y-1][x] - r.dem[y][x])}))

	// and the 4 diagonal adjacent
	surfaceNormal.Add(v01.MulR(Vec3{r.scale * (r.dem[y][x] - r.dem[y+1][x+1]) / rootTwo, rootTwo, r.scale * (r.dem[y][x] - r.dem[y+1][x+1]) / rootTwo}))
	surfaceNormal.Add(v01.MulR(Vec3{r.scale * (r.dem[y][x] - r.dem[y-1][x+1]) / rootTwo, rootTwo, r.scale * (r.dem[y][x] - r.dem[y-1][x+1]) / rootTwo}))
	surfaceNormal.Add(v01.MulR(Vec3{r.scale * (r.dem[y][x] - r.dem[y+1][x-1]) / rootTwo, rootTwo, r.scale * (r.dem[y][x] - r.dem[y+1][x-1]) / rootTwo}))
	surfaceNormal.Add(v01.MulR(Vec3{r.scale * (r.dem[y][x] - r.dem[y-1][x-1]) / rootTwo, rootTwo, r.scale * (r.dem[y][x] - r.dem[y-1][x-1]) / rootTwo}))

	surfaceNormal.Normalize()

	return surfaceNormal
}

func (r *Rainfall) randRangeInt(min, max int) int {
	return r.rand.Intn(max-min+1) + min

}

// func (r *Rainfall) randRangeFloat(min, max float64) float64 {
// 	return min + r.rand.Float64()*(max-min)
// }

func (r *Rainfall) raindrop() {
	// initialize the raindrop
	loc := Vec2{
		X: float64(r.randRangeInt(1, r.sizeX-2)),
		Y: float64(r.randRangeInt(1, r.sizeY-2))}

	speed := Vec2{0, 0}
	volume := 1.0
	percentSediment := 0.0

	// loop while the raindrop still exists
	for volume > 0 {
		initPos := Vec2{X: loc.X, Y: loc.Y}
		positionNormal := r.getSurfaceNormal(int(initPos.X), int(initPos.Y))
		// accelerate the raindrop using acceleration = force / mass
		acc := Vec2{X: positionNormal.X, Y: positionNormal.Z}
		acc.DivS(volume * r.density)
		speed.Add(acc)
		// update the position based on the new speed
		loc.Add(speed)
		// reduce the speed due to friction after the movement
		speed.MulS(1.0 - r.friction)

		// check to see if the raindrop went out of bounds
		if loc.X >= float64(r.sizeX-1) || loc.X < 1 || loc.Y >= float64(r.sizeY-1) || loc.Y < 1 {
			break
		}
		// compute the value of the maximum sediment and the difference between it and the percent sediment in the raindrop
		// positive numbers will cause erosion, negative numbers will cause deposition
		maxSediment := volume * speed.Length() * (r.dem[int(initPos.Y)][int(initPos.X)] - r.dem[int(loc.Y)][int(loc.X)])
		if maxSediment < 0.0 {
			maxSediment = 0.0
		}
		sedimentDifference := maxSediment - percentSediment
		// erode or deposit to the dem
		percentSediment += r.depositionRate * sedimentDifference
		r.dem[int(initPos.Y)][int(initPos.X)] -= volume * r.depositionRate * sedimentDifference
		// evaporate the raindrop
		volume -= r.evaporationRate
	}
}

func (r *Rainfall) Simulate(iterations int) {
	for i := 0; i < iterations; i++ {
		r.raindrop()
	}
}
