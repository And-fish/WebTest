package main

import "net/http"

// 通过V1的解析很容易就将两个方法抽象到一个接口
type Server interface {
	// 因为需要隐式的创建Context所以需要将传入handlerFunc修改成更"傻瓜"的版本
	// Route(pattern string, handlerFunc http.HandlerFunc)

	// 在当前版本下已经可以支持最基本的抽象，但是作为当前时代的web框架支持RESTFul API也是值得优化的点
	// Route(pattern string, handlerFunc MyhandlerFunc)

	// 此时的Server功能变得更加清晰，即路由能力和启动能力
	Routable
	Start(address string) error
}

type MyhandlerFunc func(*Context)

// 作为一个Server来作为实现的实体(基于http模块的Server结构体做上层)
type myHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

// 注册路由v1
func (s *myHttpServer) Route1(pattern string, handlerFunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlerFunc)
}

// 注册路由v2
func (s *myHttpServer) Route2(pattern string, handlerFunc MyhandlerFunc) {
	// 从v1变化到v2可能有点难理解，这里是用了闭包的手法
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 将创建Context的过程对Server模块隐藏，可以更好的将两者解耦
		ctx := NewContext(w, r)
		handlerFunc(ctx)
	})
}

// 注册路由v3
func (s *myHttpServer) Route3(method string, pattern string, handlerFunc MyhandlerFunc) {
	// 增加了访问方法的字段，通过不同的method来访问相同的pattern应该要能够链接到不同的handlerFunc
	// 而http.HandleFunc()很明显不能帮我们实现这个方法，那么就应该是使用自己的Handle来做路由，而不是直接注册Func
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handlerFunc // 成功将method和pattern都作为条件注册func
}

// 注册路由v4
func (s *myHttpServer) Route(method string, pattern string, handlerFunc MyhandlerFunc) {
	// 在实际的构建中，强依赖多层嵌套的结构体来实现某些功能是不可取的，比如这里使用了Server下面的Handler下面的一个map结构来实现Router功能；
	// 理想的情况应该是将功能作为接口方法暴露出来，只需要让Server接口的Router功能调用Handler接口的Router来实现即可
	s.handler.Route(method, pattern, handlerFunc)
}

// 启动ListenAndServe
func (s *myHttpServer) Start(address string) error {
	// Start的时候将自己注册
	// http.Handle("/", s.handler)

	// 在添加filter链机制之后只需要按照Filter链来处理即可
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)

}

// 创建一个新Server
func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBaseOnMap()

	// 想要在Server中添加Filter机制应该将ServerHttp作为Filter执行链的根部添加进来
	var root Filter = handler.ServeHTTP

	for i := len(builders) - 1; i >= 0; i-- {
		builder := builders[i]
		root = builder(root)
	}

	/*
		这里解释一下添加Filter的逻辑，假设Builders只有MetricsFilterBuilder一个，也就是计时filter；
		在执行的时候会调用MetricsFilterBuilder()，在记录下before之后会陷入handler.ServeHTTP(ctx.W, ctx.R)；
		在执行完ServeHTTP之后会跳转回来执行记录after的时间然后再打印出执行时间。
	*/

	return &myHttpServer{
		Name:    name,
		handler: handler,
		root:    root, // 此时的root的最内层是ServerHTTP逻辑
	}
}
