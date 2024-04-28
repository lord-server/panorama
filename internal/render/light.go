package render

const (
	ZeroIntensity    = 0
	MapEdgeIntensity = 11
	FullIntensity    = 15
)

func DecodeLight(param1 uint8) float64 {
	var LUT = [16]float64{
		0.000,
		0.024,
		0.059,
		0.118,
		0.196,
		0.286,
		0.384,
		0.471,
		0.545,
		0.608,
		0.659,
		0.710,
		0.769,
		0.835,
		0.918,
		1.000,
	}

	return LUT[param1&0xF]
}
