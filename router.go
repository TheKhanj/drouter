package drouter

import "sync"

type Tsr bool

type Param struct {
	Key   string
	Value string
}

type Params []Param

type paramsKey struct{}

var ParamsKey = paramsKey{}

type Handler interface {
	Handle(params Params)
}

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

func (r *Router) Lookup(path string) (Handler, Params, Tsr) {
	root := r.root

	if root == nil {
		return nil, nil, false
	}

	handler, ps, tsr := root.getValue(path, r.getParams)

	if handler == nil {
		r.putParams(ps)
		return nil, nil, tsr
	}

	if ps == nil {
		return handler, nil, tsr
	}

	return handler, *ps, tsr
}

func (r *Router) AddRoute(path string, handler Handler) {
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if handler == nil {
		panic("handler must not be nil")
	}

	root := r.root

	if root == nil {
		root = new(node)
		r.root = root
	}

	root.addRoute(path, handler)

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
