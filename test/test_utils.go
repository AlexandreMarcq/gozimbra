package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func NewFakeServer(t testing.TB, behavior http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(behavior)
}

func NewXMLServer(t testing.TB, wantedRequest, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Fatalf("Excepted 'Content-Type: accplication/xml' header, got: %s", r.Header.Get("Content-Type"))
		}

		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error while reading request body: %s", err)
		}

		diff := cmp.Diff(wantedRequest, string(got))
		if diff != "" {
			t.Fatalf("wrong request body: %v", diff)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

func AssertXML(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong XML, wanted %s got %s", want, got)
	}
}

func AssertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("an error should have been raised")
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("an error have been raised: %v", err)
	}
}

func ExecuteCommandC(root *cobra.Command, args ...string) (*cobra.Command, string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err := root.ExecuteC()

	return c, buf.String(), err
}
