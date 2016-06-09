// Package urlgen provides generator, which can generate a URL from a list of
// named routes with parmeters.
package urlgen

import (
    "fmt"
    "net/url"
    "regexp"
    "strings"
)

const (
    paramRegex  = `/:([a-zA-Z\d\-_]+)`
    paramPrefix = ":"
)

// Routes is a map of named routes. The key is the name of the route and the value
// is a URL, which can have named parameters in the path. The parameters must
// have the format :<parameter name>, i.e. /:foo/bar, where foo is the name of
// the parameter.
type Routes map[string]*url.URL

// Params is a map of named route parameters and their corresponding values.
// The values must not be encoded.
type Params map[string]string

// URLGen is a URL generator interface.
type URLGen interface {
    URL(name string, params Params) (*url.URL, error)
}

// URLGen is the URL generator.
type gen struct {
    routes Routes
}

// New returns a new URL generator for the routes r.
func New(r Routes) *gen {
    return &gen{routes: r}
}

// URL returns the URL and nil error, if it successfully finds a route with
// name and replaces the route params. It returns an error on failure.
func (g *gen) URL(name string, params Params) (*url.URL, error) {
    u, ok := g.routes[name]
    if !ok {
        var names []string
        for n := range g.routes {
            names = append(names, n)
        }
        return new(url.URL), fmt.Errorf("got route %q, want \"%s\"", name, strings.Join(names, `", "`))
    }

    path := u.EscapedPath()
    uparams := paramNames(path)

    if want, got := len(uparams), len(params); want == 0 && got == 0 {
        return u, nil
    } else if got != want {
        return new(url.URL), fmt.Errorf("route %q wants %d parameter(s), got %d", name, want, got)
    }
    for _, p := range uparams {
        if _, ok := params[p]; !ok {
            return new(url.URL), fmt.Errorf("route %q wants parameter %q", name, p)
        }
    }

    for _, name := range uparams {
        path = strings.Replace(path, paramPrefix+name, url.QueryEscape(params[name]), -1)
    }
    r := strings.Replace(u.String(), u.EscapedPath(), path, 1)
    ret, err := url.Parse(r)
    if err != nil {
        return new(url.URL), fmt.Errorf("can not parse route %q: %s", r, err)
    }

    return ret, nil
}

func paramNames(path string) []string {
    var names []string
    for _, m := range regexp.MustCompile(paramRegex).FindAllStringSubmatch(path, -1) {
        if len(m) != 2 {
            panic("paramNames(): param regex must have a single match")
        }
        names = append(names, m[1])
    }
    return names
}
