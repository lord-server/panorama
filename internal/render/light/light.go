package light

const (
	ZeroIntensity    = 0
	MapEdgeIntensity = 11
	FullIntensity    = 15
)

var lut = [16]float64{
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

func Decode(param1 uint8) float64 {
	return lut[param1&0xF]
}
