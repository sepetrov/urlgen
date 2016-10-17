# urlgen

[![Build Status](https://travis-ci.org/sepetrov/urlgen.svg?branch=master)](https://travis-ci.org/sepetrov/urlgen)

[![codecov](https://codecov.io/gh/sepetrov/urlgen/branch/master/graph/badge.svg)](https://codecov.io/gh/sepetrov/urlgen)


Go URL generator.


The generator uses a predefined list of named routes to generate a URL. Each of
the routes can have a named parameters in the format `:<parameter name>`.

## Usage

```go
package main

import (
	"fmt"
	"github.com/sepetrov/urlgen"
	"net/url"
)

func main() {
	route, _ := url.Parse("/hello/:name")
	generator := urlgen.New(urlgen.Routes{"hello": route})
	helloUrl, _ := generator.URL("hello", urlgen.Params{"name": "Jane & Jon"})
	fmt.Print(helloUrl.String())
	// Output:
	// /hello/Jane+%26+Jon
}
```

Using the template helper function to generate URL.
```go
package main

import (
	"github.com/sepetrov/urlgen"
	"html/template"
	"net/url"
	"os"
)

func main() {
	route, _ := url.Parse("/hello/:name")
	generator := urlgen.New(urlgen.Routes{"hello": route})
	helper := urlgen.TemplateFunc(generator)
	markup := `<a href="{{url "hello" "name" "Jane & Jon"}}">Hello</a>`
	tmpl := template.Must(template.New("foo").Funcs(template.FuncMap{"url": helper}).Parse(markup))
	tmpl.Execute(os.Stdout, nil)
	// Output:
	// <a href="/hello/Jane&#43;%26&#43;Jon">Hello</a>
}
```

## License

See [LICENSE](LICENSE).
