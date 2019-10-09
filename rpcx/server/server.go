package main

import (
	"flag"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/smallnest/rpcx/server"
	"log"
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

	zikpinReporter()

	//一对一
	// singleServer()

	//启动多个服务
	multiplyServer()

}

func zikpinReporter(){
	// zipkin上报
	reporter := zipkinhttp.NewReporter("http://49.235.1.29:9411/api/v1/spans")

	endpoint,err := zipkin.NewEndpoint("myService","myservice.mydomain.com:80")
	if err != nil {
		log.Fatalf("unable to create local endpoint:%v\n",err)
	}

	nativeTracer,err := zipkin.NewTracer(reporter,zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer:%+v\n",err)
	}

	tracer := zipkinot.Wrap(nativeTracer)

	opentracing.SetGlobalTracer(tracer)
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

