package types

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type DurationString struct {
	time.Duration
}

func NewDurationString(td time.Duration) DurationString {
	return DurationString{td}
}

func (d *DurationString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	d.Duration = dur
	return nil
}

func (d DurationString) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d DurationString) Schema(r huma.Registry) *huma.Schema {
	t := reflect.TypeFor[DurationString]()
	r.RegisterTypeAlias(t, reflect.TypeFor[string]())
	return r.Schema(t, true, "")
}
