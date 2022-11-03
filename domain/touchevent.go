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

type TouchEvent struct {
	AltKey         bool    `json:"altKey"`
	ChangedTouches []Touch `json:"changedTouches"`
	CtrlKey        bool    `json:"ctrlKey"`
	MetaKey        bool    `json:"metaKey"`
	ShiftKey       bool    `json:"shiftKey"`
	TargetTouches  []Touch `json:"targetTouches"`
	Touches        []Touch `json:"touches"`
}
