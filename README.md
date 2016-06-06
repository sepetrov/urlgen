# urlgen

Go URL generator.


The generator uses a predefined list of named routes to generate a URL. Each of
the routes can have a named parameters in the format `:<parameter name>`.

## Usage

```go
route, _ := url.Parse("/hello/:name")
generator := urlgen.New(Routes{"hello": route})
helloUrl, _ := generator.URL("hello", Params{"name": "Jane & Jon"}) // returns *url.URL, error
fmt.Print(helloUrl.String())
// Output:
// /hello/Jane+%26+Jon
```

The generator can also be used as a template helper.
```go
route, _ := url.Parse("/hello/:name")
generator := urlgen.New(Routes{"hello": r})
urlHelper := func(name string, params map[string]string) (string, error) {
	u, err := g.URL(name, params) // returns *url.URL, error
	return u.String(), err
}
tmpl := template.Must(template.New("hello").Funcs(template.FuncMap{"url": urlHelper}).Parse(`<a href="{{url "hello" .}}">Hello</a>`))
tmpl.Execute(os.Stdout, map[string]string{"name": "Jane & Jon"})
// Output:
// <a href="/hello/Jane&#43;%26&#43;Jon">Hello</a>
```

## License

See [LICENSE](LICENSE).
