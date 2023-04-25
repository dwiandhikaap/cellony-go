package util

import "github.com/jaypipes/ghw"

func Clamp(a float64, min float64, max float64) float64 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

func GpuInfo() (gpu string) {
	gpu = "Unknown GPU"
	info, err := ghw.GPU()
	if err != nil {
		return
	}

	for _, gc := range info.GraphicsCards {
		if gc.DeviceInfo != nil {
			return gc.DeviceInfo.Product.Name
		}
	}
	return
}

func RangeInterpolate(a float64, aMin float64, aMax float64, bMin float64, bMax float64) float64 {
	return bMin + (a-aMin)*(bMax-bMin)/(aMax-aMin)
}
