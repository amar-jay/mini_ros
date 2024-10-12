package msgs

// Vector3 represents a 3D vector.
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Quaternion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

// Twist represents the velocity of a robot in free space broken into its linear and angular parts.
type Twist struct {
	Linear  Vector3 `json:"linear"`
	Angular Vector3 `json:"angular"`
}

type Pose struct {
	Position    Vector3    `json:"position"`
	Orientation Quaternion `json:"orientation"`
}

type Transform struct {
	Translation Vector3    `json:"translation"`
	Rotation    Quaternion `json:"rotation"`
}
