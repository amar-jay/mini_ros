package msgs

// Vector3 represents a 3D vector.
type Vector3 struct {
	ROS_MSG
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Quaternion struct {
	ROS_MSG
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

// Twist represents the velocity of a robot in free space broken into its linear and angular parts.
type Twist struct {
	ROS_MSG
	Linear  Vector3 `json:"linear"`
	Angular Vector3 `json:"angular"`
}

type Pose struct {
	ROS_MSG
	Position    Vector3    `json:"position"`
	Orientation Quaternion `json:"orientation"`
}

type Transform struct {
	ROS_MSG
	Translation Vector3    `json:"translation"`
	Rotation    Quaternion `json:"rotation"`
}
