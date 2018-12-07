package unbox

type Pack struct {
	Frames    map[string]interface{} `json:"frames"`
	Meta      Meta                   `json:"meta"`
	FrameList []*FrameData
}

type Meta struct {
	scale string `json:"scale"`
	Image string `json:"image"`
	Size  Size   `json:"size"`
}

type FrameData struct {
	Name             string
	Frame            Rect
	Rotated          bool
	SpriteSourceSize Rect
	SourceSize       Size
	Pivot            Point
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Size struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
