package drouter

import (
	"reflect"
	"testing"
)

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func() {
		routed = true
	}
	wantParams := Params{Param{"name", "gopher"}}

	router := NewRouter()

	// try empty router first
	handle, _, tsr := router.Lookup("/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.AddRoute("/user/:name", wantHandle)
	handle, params, _ := router.Lookup("/user/gopher")
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle.(func())()
		if !routed {
			t.Fatal("Routing failed!")
		}
	}
	if !reflect.DeepEqual(params, wantParams) {
		t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, params)
	}
	routed = false

	// route without param
	router.AddRoute("/user", wantHandle)
	handle, params, _ = router.Lookup("/user")
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle.(func())()
		if !routed {
			t.Fatal("Routing failed!")
		}
	}
	if params != nil {
		t.Fatalf("Wrong parameter values: want %v, got %v", nil, params)
	}

	handle, _, tsr = router.Lookup("/user/gopher/")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, _, tsr = router.Lookup("/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}
