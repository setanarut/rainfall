package main

import "github.com/setanarut/rainfall"

func main() {
	dem := rainfall.ImageToSlice(rainfall.LoadImage("noise.png"))
	rf := rainfall.NewRainfall(dem, int64(666555))

	for i := 0; i < 200; i++ {
		rf.Simulate(1000)

	}
	dem = rf.GetDem()
	rainfall.SaveImage("noise_out.png", rainfall.SliceToImage(dem))
}
