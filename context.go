package main

import (
	"encoding/json"
	"io"
	"net/http"
)

/*
	是否需要将Context抽象成一个接口?
	作为一个成熟的框架，应该要将所有的功能都抽象成接口，将所有的方法都变成对于接口之间的交互；
	但是在这里为了简单处理将Context设计为了一个结构体
*/

// 将W&R抽象为一个整体存在
type Context struct {
	W http.ResponseWriter
	R *http.Request
}

// 从请求中去读取数据
// 从body中读取，并处理json(反序列化)
func (c *Context) ReadJson(obj interface{}) error {
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// 将body中的数据转化成传进来的obj
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

// 写入response
func (c *Context) WriteJson(code int, resp interface{}) error {
	c.W.WriteHeader(code)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJson)
	if err != nil {
		return err
	}
	return nil
}

// 从WriteJson之上还可以将常见的响应做包装
func (c *Context) WriteOk(resp interface{}) error {
	return c.WriteJson(http.StatusOK, resp)
}
func (c *Context) WriteBad(resp interface{}) error {
	return c.WriteJson(http.StatusBadRequest, resp)
}

// 创建Context的方法
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}
