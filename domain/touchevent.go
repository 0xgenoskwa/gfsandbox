package domain

type Touch struct {
	Identifier    string      `json:"identifier"`
	ScreenX       float64     `json:"screenX"`
	ScreenY       float64     `json:"screenY"`
	ClientX       float64     `json:"clientX"`
	ClientY       float64     `json:"clientY"`
	PageX         float64     `json:"pageY"`
	PageY         float64     `json:"pageX"`
	Target        interface{} `json:"target"`
	RadiusX       float64     `json:"radiusY"`
	RadiusY       float64     `json:"radiusX"`
	RorationAngle float64     `json:"rorationAgle"`
	Force         float64     `json:"force"`
}
