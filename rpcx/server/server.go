package main

import (
	"flag"
	"github.com/smallnest/rpcx/server"
	"rpclearn/rpcx/server/service"
)

var (
	addr = flag.String("addr","localhost:8568","server address")
	addr2 = flag.String("addr2","localhost:8569","server address")

	consulAddr = flag.String("consulAddr","49.235.1.29:8500","consul address")
	basePath = flag.String("base","/rpcx_test","prefix path")
)


func main(){
	flag.Parse()
	//一对一
	// singleServer()

	//启动多个服务
	multiplyServer()

}

//一对一  没有配置中心
func singleServer(){
	s := server.NewServer()
	s.RegisterName("Arith",new(service.Arith),"")
	go s.Serve("tcp",*addr)
	select {
	}
}

//启动多个服务
func multiplyServer(){
	go createServer(*addr)
	go createServer(*addr2)

	select {

	}
}

//consul服务发现

//
func createServer(addr string){
	s := server.NewServer()
	s.RegisterName("Arith",new(service.Arith),"")
	s.Serve("tcp",addr)
}

