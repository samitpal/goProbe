package misc

import (
	"html/template"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestFormattedTime(t *testing.T) {
	ps := NewProbesStatus([]string{"probe1"})
	unixEpoch := int64(1360287003083988472)
	rt := ps.FormattedTime(&unixEpoch)
	et := "Fri Feb  8 01:30:03 UTC 2013"
	if !reflect.DeepEqual(rt.UTC().Format(time.UnixDate), et) {
		t.Errorf("Got: %v\n Want: %v", rt.UTC().Format(time.UnixDate), et)
	}
}

func TestFormattedHttpHeaders(t *testing.T) {
	ps := NewProbesStatus([]string{"probe1"})
	h := http.Header{"My-Header": []string{"one", "two"}}
	rh := ps.FormattedHttpHeaders(h)
	eh := template.HTML("My-Header: one two<br />")
	if !reflect.DeepEqual(rh, eh) {
		t.Errorf("Got: %v\n Want: %v", rh, eh)
	}
}
