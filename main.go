package main

import (
	"log"
)

/*
	main.go作为测试使用
	为什么需要web框架？在最基础的http模块的基础上通常情况下会需要更上层的抽象(Server)，这个上层的抽象是作为对某个端口监听的实体存在的
	一个Server作为一个功能的接口来使用，在功能上需要实现一种内聚，在逻辑上是独立的。
	就像是一个中学生需要补课，最理想的情况是报名一个一对一的9科全包干
*/

type signUpReq struct {
	Name       string
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

type commomResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

/* 在每个业务场景中让用户重复操作R&W并不便捷，可以将两者抽象为一次请求
func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed:%v", err)
		return
	}
	// 将body中的数据转化成传进来的obj
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Fprintf(w, "deserialized failed:%v", err)
		return
	}

	resp := commomResponse{
		Data: 123,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, "response failed:%v", err)
		return
	}
	fmt.Fprintf(w, string(respJson))
} */

/* // 将读写都抽象为方法可以大量减少重复化操作，也减少了出错的可能
func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}

	// 在这个版本中可以发现需要用户直接操作w&r，也就是说需要手动创建contex，并不易于代码的稳定
	// 也就是说一个成熟的框架不应该将大量结构体都暴露出来，而是应该在接口方法的调用中隐式的创建需要的部分
	ctx := Context{
		W: w,
		R: r,
	}

	err := ctx.ReadJson(req)
	if err != nil {
		fmt.Fprintf(w, "failed err is:%v", err)
		return
	}

	resp := commomResponse{
		Data: 123,
	}
	err = ctx.WriteOk(resp)
	if err != nil {
		log.Printf("写入响应失败:%v", err)
		return
	}
} */

// 通过隐藏创建context的细节实现了对代码的保护
// 到这里就实现了对Server和Context的简单抽象，即Server是就是服务的路由注册和服务开启，Context就是对一次连接中数据来往
func SignUp(ctx *Context) {
	req := &signUpReq{}

	err := ctx.ReadJson(req)
	if err != nil {
		ctx.WriteBad(err)
		return
	}

	resp := commomResponse{
		Data: 123,
	}
	err = ctx.WriteOk(resp)
	if err != nil {
		log.Printf("写入响应失败:%v", err)
		return
	}
}
func main() {
	/* 	最基本的HttpHandle
		解析下面的代码可以发现一个webServer其实最重要的就只有两个工作，一个是注册路由，一个是启动LAS；
		那么就可以将Server抽象为Router和Start两个方法

	   	http.HandleFunc("/", SignUP)
	   	log.Fatal(http.ListenAndServe(":8080", nil)) */

	server := NewHttpServer("test_Server")
	server.Route("Get", "/", SignUp)
	server.Start(":8080")

}
