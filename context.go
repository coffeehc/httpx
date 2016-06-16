package web

type RequestContext map[string]interface{}

func (this RequestContext) Get(key string) interface{} {
	return this[key]
}

func (this RequestContext) Set(key string, value interface{}) {
	this[key] = value
}
