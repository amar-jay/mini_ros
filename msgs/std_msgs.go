package msgs

import "time"

type ROS_MSG interface{}
type String interface {
	ROS_MSG
	string
}

type Int32 interface {
	ROS_MSG
	int32
}

type Float64 interface {
	ROS_MSG
	float64
}

type Bool interface {
	ROS_MSG
	bool
}

type ColorRGBA struct {
	ROS_MSG
	R float32
	G float32
	B float32
	A float32
}

type ColorRGB struct {
	ROS_MSG
	R float32
	G float32
	B float32
}

type Header struct {
	ROS_MSG
	Seq     uint32
	Stamp   time.Time
	FrameId string
}
