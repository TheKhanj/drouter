package drouter

import (
	"context"
	"sync"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

// ParamsFromContext pulls the URL parameters from a request context,
// or returns nil if none are present.
func ParamsFromContext(ctx context.Context) Params {
	p, _ := ctx.Value(ParamsKey).(Params)
	return p
}

// MatchedRoutePathParam is the Param name under which the path of the matched
// route is stored, if Router.SaveMatchedRoutePath is set.
var MatchedRoutePathParam = "$matchedRoutePath"

// MatchedRoutePath retrieves the path of the matched route.
// Router.SaveMatchedRoutePath must have been enabled when the respective
// handle was added, otherwise this function always returns an empty string.
func (ps Params) MatchedRoutePath() string {
	return ps.ByName(MatchedRoutePathParam)
}

type Handle interface{}

type Router struct {
	root *node

	paramsPool sync.Pool
	maxParams  uint16
}

func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0]
	return ps
}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

func (r *Router) Lookup(path string) (Handle, Params, bool) {
	root := r.root

	if root == nil {
		return nil, nil, false
	}

	handle, ps, tsr := root.getValue(path, r.getParams)

	if handle == nil {
		r.putParams(ps)
		return nil, nil, tsr
	}

	if ps == nil {
		return handle, nil, tsr
	}

	return handle, *ps, tsr
}

func (r *Router) AddRoute(path string, handle Handle) {
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if handle == nil {
		panic("handle must not be nil")
	}

	root := r.root

	if root == nil {
		root = new(node)
		r.root = root
	}

	root.addRoute(path, handle)

	r.updateMaxParams(path)
	r.lazyInitParamsPool()
}

func (r *Router) lazyInitParamsPool() {
	if !(r.paramsPool.New == nil && r.maxParams > 0) {
		return
	}

	r.paramsPool.New = func() interface{} {
		ps := make(Params, 0, r.maxParams)
		return &ps
	}
}

func (r *Router) updateMaxParams(path string) {
	if paramsCount := countParams(path); paramsCount > r.maxParams {
		r.maxParams = paramsCount
	}
}
