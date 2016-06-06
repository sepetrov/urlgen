package urlgen

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func ExampleURL() {
	r, _ := url.Parse("/hello/:name")
	g := New(Routes{"hello": r})
	u, _ := g.URL("hello", Params{"name": "Jane & Jon"})
	fmt.Print(u.String())
	// Output:
	// /hello/Jane+%26+Jon
}

func TestUrl(t *testing.T) {
	type tRoute struct {
		name string
		url  string
	}
	type tReq struct {
		name   string
		params Params
	}
	type tRes struct {
		url string
		err bool
	}
	var tests = [...]struct {
		route tRoute
		req   tReq
		res   tRes
	}{
		{
			tRoute{"route#0", "/"},
			tReq{"invalid", nil},
			tRes{"", true},
		},
		{
			tRoute{"route#1", "/:foo"},
			tReq{"route#1", Params{"foo": "one:two"}},
			tRes{"/one%3Atwo", false},
		},
		{
			tRoute{"route#2", "/:foo"},
			tReq{"route#2", Params{"bar": "baz"}},
			tRes{"", true},
		},
		{
			tRoute{"route#3", ""},
			tReq{"route#3", nil},
			tRes{"", false},
		},
		{
			tRoute{"route#4", "http://www.google.com"},
			tReq{"route#4", nil},
			tRes{"http://www.google.com", false},
		},
		{
			tRoute{"route#5", "http://www.google.com/"},
			tReq{"route#5", nil},
			tRes{"http://www.google.com/", false},
		},
		{
			tRoute{"route#6", "http://www.google.com/:Foo-1/file%3Aone%26two/:ba_r?q=go%3Alanguage#:xyz"},
			tReq{"route#6", Params{"Foo-1": "FOO", "ba_r": "BAR"}},
			tRes{"http://www.google.com/FOO/file%3Aone%26two/BAR?q=go%3Alanguage#:xyz", false},
		},
		{
			tRoute{"route#7", "/:Foo-1/file%3Aone%26two/:ba_r?q=go%3Alanguage#:xyz"},
			tReq{"route#7", Params{"Foo-1": "FOO", "ba_r": "BAR"}},
			tRes{"/FOO/file%3Aone%26two/BAR?q=go%3Alanguage#:xyz", false},
		},
	}
	var routes = make(Routes)
	for i, test := range tests {
		u, err := url.Parse(test.route.url)
		if err != nil {
			t.Fatalf("#%d TestUrl(): can not parse test route URL %q", i, test.route.url)
		}
		routes[test.route.name] = u
	}
	var gen = New(routes)
	for i, test := range tests {
		url, err := gen.URL(test.req.name, test.req.params)
		if url.String() != test.res.url {
			t.Errorf("#%d URL(%q): got URL %q, want %q", i, test.req.name, url.String(), test.res.url)
		}
		if (err != nil) != test.res.err {
			t.Errorf("#%d URL(%q): got unexpected error value %v", i, test.req.name, err)
		}
	}
}

func TestParamNames(t *testing.T) {
	var tests = [...]struct {
		path   string
		params []string
	}{
		{
			"/",
			[]string{},
		},
		{
			"/foo",
			[]string{},
		},
		{
			"/:foo",
			[]string{"foo"},
		},
		{
			"/:foo/bar",
			[]string{"foo"},
		},
		{
			"/:foo/:bar",
			[]string{"foo", "bar"},
		},
		{
			"/:foo/file%20one%26two/:bar",
			[]string{"foo", "bar"},
		},
		{
			"/:Foo-123/file%20one%26two/:6b_aR",
			[]string{"Foo-123", "6b_aR"},
		},
	}
	for i, test := range tests {
		params := paramNames(test.path)
		if got, want := len(params), len(test.params); got != want {
			t.Errorf("#%d paramNames(%q): got %d params %v, want %d params %v", i, test.path, got, params, want, test.params)
			continue
		}
	NEXT:
		for _, want := range test.params {
			for _, got := range params {
				if got == want {
					continue NEXT
				}
			}
			t.Errorf("#%d paramNames(%q): got params \"%s\", want param %q", i, test.path, strings.Join(params, `", "`), want)
		}
	}
}
