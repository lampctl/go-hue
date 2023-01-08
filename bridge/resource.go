package bridge

import (
	"encoding/json"
	"reflect"
)

const (
	TypeBridgeHome   = "bridge_home"
	TypeGroupedLight = "grouped_light"
	TypeLight        = "light"
	TypeZone         = "zone"
)

type Error struct {
	Description string `json:"description"`
}

type Response struct {
	Errors []*Error        `json:"errors"`
	Data   json.RawMessage `json:"data"`
}

type RegistrationRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

type RegistrationResponseError struct {
	Description string `json:"description"`
}

type RegistrationResponseSuccess struct {
	Username string `json:"username"`
}

type RegistrationResponse struct {
	Error   *RegistrationResponseError   `json:"error"`
	Success *RegistrationResponseSuccess `json:"success"`
}

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
	Type     string    `json:"type,omitempty"`
}

func deepCopy(dest, src any) {
	var (
		tSrc  = reflect.TypeOf(src).Elem()
		vDest = reflect.ValueOf(dest).Elem()
		vSrc  = reflect.ValueOf(src).Elem()
	)
	for i := 0; i < tSrc.NumField(); i++ {
		var (
			tSrcField  = tSrc.Field(i).Type
			vDestField = vDest.Field(i)
			vSrcField  = vSrc.Field(i)
		)
		switch tSrcField.Kind() {
		case reflect.Pointer:
			if !vSrcField.IsNil() {
				if vDestField.IsNil() {
					vDestField.Set(reflect.New(tSrcField.Elem()))
				}
				deepCopy(vDestField.Interface(), vSrcField.Interface())
			}
		default:
			if !vSrcField.IsZero() {
				vDestField.Set(vSrcField)
			}
		}
	}
}

// CopyFrom populates the resource with only the values set in the provided
// parameter.
func (r *Resource) CopyFrom(v *Resource) {
	deepCopy(r, v)
}
