package main

import (
	"flag"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"rpclearn/rpcx/server/service"
	"rpclearn/rpcx/tracer_plugin"
)

var (
	addr  = flag.String("addr", "localhost:8568", "server address")
	addr2 = flag.String("addr2", "localhost:8569", "server address")

	consulAddr = flag.String("consulAddr", "49.235.1.29:8500", "consul address")
	basePath   = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	InitTracer()

	// zipkin上报
	reporter := zipkinhttp.NewReporter("http://49.235.1.29:9411/api/v2/spans")
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("rpcx_demo_server", "127.0.0.1")
	if err != nil {
		logrus.Errorf("zipkin.NewEndPoint new endpoint error,err_msg = %v", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Errorf("zipkin.NewTracer create tracer error,err_msg = %v", err)
	}
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)

	//一对一
	// singleServer()

	//启动多个服务
	multiplyServer()

	// consul 服务发现
	//consulServer()
}

func InitTracer() {
	// 加载tracer
	tracer, err := tracer_plugin.GetTracer("rpcx_demo_server")
	if err != nil {
		logrus.Warnf("init tracer failed:err_msg = %v", err)
	} else {
		opentracing.SetGlobalTracer(tracer)
	}
}

//一对一  没有配置中心
func singleServer() {
	s := server.NewServer()
	s.RegisterName("Arith", new(service.Arith), "")
	go s.Serve("tcp", *addr)
	select {}
}

//启动多个服务
func multiplyServer() {
	go createServer(*addr)
	go createServer(*addr2)

	select {}
}

/*
//consul服务发现
func consulServer() {
	consultPlugin := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		ConsulServers:  []string{*consulAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := consultPlugin.Start()
	if err != nil {
		logrus.Errorf("start consul plugin failed,err_msg = %v", err)
	}

	s := server.NewServer()
	s.Plugins.Add(consultPlugin)
	s.RegisterName("Arith", new(service.Arith), "")
	s.Serve("tcp", *addr)
}*/

//
func createServer(addr string) {
	s := server.NewServer()
	s.RegisterName("Arith", new(service.Arith), "")
	s.Serve("tcp", addr)
}
