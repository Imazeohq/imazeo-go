package imazeo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	imazeo "github.com/Imazeohq/imazeo-go"
)

func TestRedirectDoesNotForwardProjectKey(t *testing.T) {
	var destinationRequests atomic.Int32
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		destinationRequests.Add(1)
		jsonOK(t, w, 200, map[string]string{"projectId": "p", "keyId": "k"})
	}))
	defer destination.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "secret-key" {
			t.Error("missing project key on original request")
		}
		http.Redirect(w, r, destination.URL, http.StatusTemporaryRedirect)
	}))
	defer origin.Close()
	client, err := imazeo.New(imazeo.WithAPIKey("secret-key"), imazeo.WithBaseURL(origin.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Context.Get(context.Background()); err == nil {
		t.Fatal("expected redirect to be rejected")
	}
	if destinationRequests.Load() != 0 {
		t.Fatal("project credentials were forwarded")
	}
}
