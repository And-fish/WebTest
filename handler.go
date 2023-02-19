package main

import "net/http"

// 作为多个接口都实现的方法，同样可以将其抽象出来作为接口来表示，Routable接口只需要专注于实现路由和注册
type Routable interface {
	Route(method string, pattern string, handlerFunc MyhandlerFunc)
}

// 我们自己实现的接口应该要实现http.Handler和路由功能
type Handler interface {
	ServeHTTP(*Context)
	Routable
}

// 这个Handler应该要能够实现http.Handler接口，也就是必须要实现ServeHTTP(ResponseWriter, *Request)方法
type HandlerBaseOnMap struct {
	// v1版本直接使用map作为路由(我就问你能不能找到8)
	// key使用method + url就能实现将不同method相同url路由到不同的方法上
	handlers map[string]MyhandlerFunc
}

// 需要实现的接口函数
func (h *HandlerBaseOnMap) ServeHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	if handle, ok := h.handlers[key]; ok {
		// 如果找到了就直接调用
		handle(ctx)
	} else {
		// 如果没有找到路由返回404
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("NOT FOUND"))
	}
}

// 获取key
func (h *HandlerBaseOnMap) key(method, path string) string {
	return method + "#" + path
}

// 注册路由
func (h *HandlerBaseOnMap) Route(method string, pattern string, handlerFunc MyhandlerFunc) {
	// 注册现在是作为Handler自身的功能存在了
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

// 创建Handler
func NewHandlerBaseOnMap() Handler {
	return &HandlerBaseOnMap{
		handlers: make(map[string]MyhandlerFunc), // 像一个成熟的框架可以在这里添加option来设置默认路由容量
	}
}
