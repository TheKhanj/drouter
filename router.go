package drouter

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

type Router struct{}

type Tsr bool

func (r *Router) Lookup(path string) (Handler, Tsr) {
}

func (r *Router) AddRoute(path string, handler Handler) {
}
