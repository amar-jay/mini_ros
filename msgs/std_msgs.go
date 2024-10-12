package msgs

import "time"

type String string
type Int32 int32
type Float64 float64
type Bool bool
type ColorRGBA struct {
	R float32
	G float32
	B float32
	A float32
}

type ColorRGB struct {
	R float32
	G float32
	B float32
}

type Header struct {
	Seq     uint32
	Stamp   time.Time
	FrameId string
}
