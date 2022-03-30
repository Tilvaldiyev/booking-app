package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
 } {
	{"home", "/", http.MethodGet, []postData{}, http.StatusOK},
	{"about", "/about", http.MethodGet, []postData{}, http.StatusOK},
	{"gq", "/rooms/general-quarters", http.MethodGet, []postData{}, http.StatusOK},
	{"ms", "/rooms/major-suite", http.MethodGet, []postData{}, http.StatusOK},
	{"sa", "/search-availability", http.MethodGet, []postData{}, http.StatusOK},
	{"mr", "/make-reservation", http.MethodGet, []postData{}, http.StatusOK},
	{"contact", "/contact", http.MethodGet, []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", http.MethodPost, []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", http.MethodPost, []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/search-availability-json", http.MethodPost, []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "Smith@here.com"},
		{key: "phone", value: "+77778885522"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == http.MethodGet {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL + e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}

}
