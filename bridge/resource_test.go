package bridge

import (
	"reflect"
	"testing"
)

const (
	fakeID   = "id"
	fakeName = "name"
)

func TestCopyFrom(t *testing.T) {
	var (
		dest   = &Resource{}
		src    = &Resource{ID: fakeID, Metadata: &Metadata{Name: fakeName}}
		output = &Resource{ID: fakeID, Metadata: &Metadata{Name: fakeName}}
	)
	dest.CopyFrom(src)
	if !reflect.DeepEqual(dest, output) {
		t.Fatalf("%+v != %+v", dest, output)
	}
}
