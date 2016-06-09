package urlgen

import (
    "bytes"
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

func TestTemplateFuncHTML(t *testing.T) {
    route0, _ := url.Parse("/")
    route1, _ := url.Parse("/:foo")
    route2, _ := url.Parse("/:foo/:bar")
    route3, _ := url.Parse("/:foo/:bar/:baz")
    generator := New(Routes{
        "zero":  route0,
        "one":   route1,
        "two":   route2,
        "three": route3,
    })
    helper := TemplateFunc(generator)
    markup := `{{url "zero"}}
{{url "one" "foo" "FOO"}}
{{url "two" "foo" "FOO" "bar" "BAR"}}
{{url "three" "foo" "FOO" "bar" "BAR" "baz" "BAZ"}}`
    tmpl := template.Must(template.New("foo").Funcs(template.FuncMap{"url": helper}).Parse(markup))
    buf := new(bytes.Buffer)
    tmpl.Execute(buf, nil)
    want := `/
/FOO
/FOO/BAR
/FOO/BAR/BAZ`
    if got := buf.String(); got != want {
        t.Errorf("TemplateFunc(): got template output\n\n%s\n\nwant\n\n%s\n", got, want)
    }
}
