package gzc

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func ok(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

func equals(t *testing.T, expected, received interface{}) {
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("Expected '%+v' was not same as received '%+v'", expected, received)
	}
}

func TestGettingEpic(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/p1/repositories/1/epics/1")
		equals(t, req.Header.Get("X-Authentication-Token"), "token")
		rw.Write([]byte(`{}`))
	}))
	defer server.Close()

	api := CreateAPI(server.Client(), server.URL)
	client := CreateClient(api, "token")
	e, err := client.RequestEpic(1, 1)

	ok(t, err)
	equals(t, &Epic{}, e)
}

func TestGettingDependencies(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		equals(t, req.URL.String(), "/p1/repositories/1/dependencies")
		equals(t, req.Header.Get("X-Authentication-Token"), "token")
		rw.Write([]byte(`{}`))
	}))
	defer server.Close()

	api := CreateAPI(server.Client(), server.URL)
	client := CreateClient(api, "token")
	d, err := client.RequestDependencies(1)

	ok(t, err)
	equals(t, &Dependencies{}, d)
}
