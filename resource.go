package hue

const (
	TypeBridgeHome   = "bridge_home"
	TypeGroupedLight = "grouped_light"
	TypeLight        = "light"
	TypeZone         = "zone"
)

type Owner struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

type Metadata struct {
	Name string `json:"name"`
}

type On struct {
	On bool `json:"on"`
}

type Dimming struct {
	Brightness float64 `json:"brightness"`
}

type Dynamics struct {
	Duration int64 `json:"duration"`
}

type ColorXY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Color struct {
	XY *ColorXY `json:"xy,omitempty"`
}

type Resource struct {
	ID       string    `json:"id,omitempty"`
	Owner    *Owner    `json:"owner,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	On       *On       `json:"on,omitempty"`
	Dimming  *Dimming  `json:"dimming,omitempty"`
	Color    *Color    `json:"color,omitempty"`
	Dynamics *Dynamics `json:"dynamics,omitempty"`
}
