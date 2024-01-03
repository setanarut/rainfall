package rainfall

import (
	"image"
	"math"
	"math/rand"
)

// constants
var root2 = math.Sqrt(2)
var v015 = vec3{0.15, 0.15, 0.15}
var v01 = vec3{0.1, 0.1, 0.1}

type Options struct {
	// Scale is used as a multiplier for determining the surface normal at any given point. Increasing this value would cause a raindrop to move a farther distance, whereas decreasing it would cause it to move a shorter distance. This is because the scale directly affects the speed of the raindrop. Acceleration = force / mass, and this variable determines the force in that equation.
	Scale float64
	// Density is used to calculate the acceleration of the raindrop as it flows across the heightmap. Acceleration = force / mass, and mass is calculated by volume * density. Therefore, increasing the density causes the raindrop to accelerate more slowly.
	Density float64
	// Friction is a multiplier which effects the speed of the raindrop after it moves across a surface. Higher friction means lower speeds.
	Friction float64
	// Deposition rate is a multiplier which controls how much sediment is deposited to the terrain.
	DepositionRate float64
	// Evaporation rate is how many times the raindrop can move to a new position before it is terminated. By default, this number is set to the 1 / the length of the X axis.
	EvaporationRate float64
	// The seed of randomness in each raindrop.
	RaindropRandomSeed int64
}

// DefaultOptions returns default simulation options
func DefaultOptions() *Options {
	return &Options{
		Scale:          100.0,
		Density:        1.0,
		Friction:       0.1,
		DepositionRate: 0.3,
		// 1.0 / height-map width
		EvaporationRate:    1.0 / 512,
		RaindropRandomSeed: 1923,
	}
}

type Rainfall struct {
	// Terrain is a 2D array height map in range -1.0~1.0
	Terrain [][]float64
	Opt     *Options

	// size of terrain (width and height)
	terrainWidth, terrainHeight int
	rand                        *rand.Rand
}

// New returns new Rainfall from 2D slice in range [-1 ~ 1]
func New(terrain [][]float64, opt *Options) *Rainfall {
	return &Rainfall{
		Terrain:       terrain,
		Opt:           opt,
		terrainWidth:  len(terrain),
		terrainHeight: len(terrain[0]),
		rand:          rand.New(rand.NewSource(opt.RaindropRandomSeed)),
	}
}

// NewFromImageFile returns new Rainfall from image file
func NewFromImageFile(filePath string, opt *Options) *Rainfall {
	return NewFromImage(openImage(filePath), opt)

}

// NewFromImage returns new Rainfall from image.Image
func NewFromImage(img image.Image, opt *Options) *Rainfall {
	return &Rainfall{
		Terrain:       imageToSlice(img),
		Opt:           opt,
		terrainWidth:  img.Bounds().Size().X,
		terrainHeight: img.Bounds().Size().Y,
		rand:          rand.New(rand.NewSource(opt.RaindropRandomSeed)),
	}
}

func (rf *Rainfall) getSurfaceNormal(x, y int) vec3 {
	// create the vector and add the 4 points directly adjacent to it
	surfNorm := v015 // Surface normal
	surfNorm.Mul(vec3{
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x+1][y]),
		1,
		0})
	surfNorm.Add(v015.MulR(vec3{
		rf.Opt.Scale * (rf.Terrain[x-1][y] - rf.Terrain[x][y]),
		1,
		0}))

	surfNorm.Add(v015.MulR(vec3{
		0,
		1,
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x][y+1])}))

	surfNorm.Add(v015.MulR(vec3{
		0,
		1,
		rf.Opt.Scale * (rf.Terrain[x][y-1] - rf.Terrain[x][y])}))

	// and the 4 diagonal adjacent
	surfNorm.Add(v01.MulR(vec3{
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x+1][y+1]) / root2,
		root2,
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x+1][y+1]) / root2}))

	surfNorm.Add(v01.MulR(vec3{
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x+1][y-1]) / root2,
		root2,
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x+1][y-1]) / root2}))

	surfNorm.Add(v01.MulR(vec3{
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x-1][y+1]) / root2,
		root2,
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x-1][y+1]) / root2}))

	surfNorm.Add(v01.MulR(vec3{
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x-1][y-1]) / root2,
		root2,
		rf.Opt.Scale * (rf.Terrain[x][y] - rf.Terrain[x-1][y-1]) / root2}))

	// normalize
	surfNorm.Normalize()

	return surfNorm
}

func (rf *Rainfall) randRangeInt(min, max int) int {
	return rf.rand.Intn(max-min+1) + min

}

// Raindrop drops a random single raindrop
func (rf *Rainfall) Raindrop() {
	// initialize the random raindrop location
	loc := vec2{
		X: float64(rf.randRangeInt(1, rf.terrainWidth-2)),
		Y: float64(rf.randRangeInt(1, rf.terrainHeight-2))}

	speed := vec2{X: 0, Y: 0}
	volume := 1.0
	percentSediment := 0.0

	// loop while the raindrop still exists
	for volume > 0 {
		initPos := vec2{X: loc.X, Y: loc.Y}
		positionNormal := rf.getSurfaceNormal(int(initPos.X), int(initPos.Y))
		// accelerate the raindrop using acceleration = force / mass
		acc := vec2{X: positionNormal.X, Y: positionNormal.Z}
		acc.DivS(volume * rf.Opt.Density)
		speed.Add(acc)
		// update the position based on the new speed
		loc.Add(speed)
		// reduce the speed due to friction after the movement
		speed.MulS(1.0 - rf.Opt.Friction)

		// check to see if the raindrop went out of bounds
		if loc.X >= float64(rf.terrainWidth-1) || loc.X < 1 || loc.Y >= float64(rf.terrainHeight-1) || loc.Y < 1 {
			break
		}
		// compute the value of the maximum sediment and the difference between it and the percent sediment in the raindrop
		// positive numbers will cause erosion, negative numbers will cause deposition
		maxSediment := volume * speed.Length() * (rf.Terrain[int(initPos.X)][int(initPos.Y)] - rf.Terrain[int(loc.X)][int(loc.Y)])
		if maxSediment < 0.0 {
			maxSediment = 0.0
		}
		sedimentDifference := maxSediment - percentSediment
		// erode or deposit to the dem
		percentSediment += rf.Opt.DepositionRate * sedimentDifference
		rf.Terrain[int(initPos.X)][int(initPos.Y)] -= volume * rf.Opt.DepositionRate * sedimentDifference
		// evaporate the raindrop
		volume -= rf.Opt.EvaporationRate
	}
}

// Raindrops drops the given amount of random raindrops.
func (rf *Rainfall) Raindrops(amount int) {
	for i := 0; i < amount; i++ {
		rf.Raindrop()
	}
}

// WriteToImageFile writes terrain to image file
func (rf *Rainfall) WriteToImageFile(filePath string) {
	saveImage(filePath, sliceToImage(rf.Terrain))
}

// GetImage returns terrain as image
func (rf *Rainfall) GetImage() *image.Gray {
	return sliceToImage(rf.Terrain)
}
