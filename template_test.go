package urlgen

import (
	"html/template"
	"net/url"
	"os"
	"testing"
)

func ExampleTemplateFunc() {
	r, _ := url.Parse("/hello/:name")
	f := TemplateFunc(New(Routes{"hello": r}))
	c := `<a href="{{url "hello" "name" "Jane & Jon"}}">Hello</a>`
	t := template.Must(template.New("foo").Funcs(template.FuncMap{"url": f}).Parse(c))
	t.Execute(os.Stdout, nil)
	// Output:
	// <a href="/hello/Jane&#43;%26&#43;Jon">Hello</a>
}

func TestTemplateFunc(t *testing.T) {
	var tests = [...]struct {
		path   string
		params []string
		url    string
		err    bool
	}{
		{
			"",
			[]string{},
			"",
			false,
		},
		{
			"/",
			[]string{},
			"/",
			false,
		},
		{
			"/a/:foo",
			[]string{"foo", "FOO"},
			"/a/FOO",
			false,
		},
		{
			"/a/:foo/b/:bar",
			[]string{"bar", "BAR", "foo", "FOO"},
			"/a/FOO/b/BAR",
			false,
		},
		{
			"/a/:foo/b/:bar/:baz",
			[]string{"bar", "BAR", "baz", "BAZ", "foo", "FOO"},
			"/a/FOO/b/BAR/BAZ",
			false,
		},
	}
	for i, test := range tests {
		r, _ := url.Parse(test.path)
		f := TemplateFunc(New(Routes{"foo": r}))
		u, err := f("foo", test.params...)
		if u != test.url {
			t.Errorf("#%d TemplateFunc(): got URL %q, want %q", i, u, test.url)
		}
		if (err == nil) && test.err {
			t.Errorf("#%d TemplateFunc(): got error %v, want %v", i, u, test.url)
		}
	}
}
