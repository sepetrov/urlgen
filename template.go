package urlgen

import (
	"fmt"
)

type templateHelperFunc func(string, ...string) (string, error)

// TemplateFunc is an adapter, which wraps a URL generator and returns a
// template helper function.
func TemplateFunc(g *URLGen) templateHelperFunc {
	return func(name string, params ...string) (string, error) {
		l := len(params)
		if l%2 != 0 {
			return "", fmt.Errorf("got %d parameters, want even number", l)
		}
		var p = make(Params)
		for i := 0; i < l; i += 2 {
			p[params[i]] = params[i+1]
		}
		u, err := g.URL(name, p)
		if err != nil {
			return "", err
		}
		return u.String(), nil
	}
}
