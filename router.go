package drouter

type Router struct{}

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
