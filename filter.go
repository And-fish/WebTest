package main

import (
	"fmt"
	"time"
)

type Filter func(*Context)

type FilterBuilder func(next Filter) Filter

func MetricsFilterBuilder(next Filter) Filter {
	return func(ctx *Context) {
		before := time.Now().Nanosecond()
		next(ctx)
		after := time.Now().Nanosecond()
		fmt.Printf("使用了%d纳秒", before-after)
	}
}

/*
	通过一个简单的MetricsFilterBuilder就能理解一个成熟的web框架中是如何解决像middleware之类的功能的；
	本质上是使用闭包的机制将需要实现的功能添加到这条链中，这条链在执行的过程中会在处理一部分逻辑之后，陷入更深的处理逻辑，最后在链的末端跳出
*/
