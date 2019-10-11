package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"rpclearn/rpcx/data_struct"
	"rpclearn/rpcx/middle"
	"rpclearn/rpcx/tracer_plugin"
)

var (
	addr  = flag.String("addr", "127.0.0.1:8568", "server address")
	addr2 = flag.String("addr2", "127.0.0.1:8569", "server address")
)

func InitTracer() {
	// 加载tracer
	tracer, err := tracer_plugin.GetTracer("rpcx_demo_client")
	if err != nil {
		logrus.Warnf("init tracer failed:err_msg = %v", err)
	} else {
		opentracing.SetGlobalTracer(tracer)
	}
}

// 一对一，服务器地址硬编码或写在配置文件中
func main() {
	flag.Parse()
	//InitTracer()

	//d := client.NewPeer2PeerDiscovery("tcp@"+*addr,"")
	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	// zipkin上报
	reporter := zipkinhttp.NewReporter("http://49.235.1.29:9411/api/v2/spans")
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint("rpcx_demo_client", "127.0.0.1")
	if err != nil {
		logrus.Errorf("zipkin.NewEndPoint new endpoint error,err_msg = %v", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Errorf("zipkin.NewTracer create tracer error,err_msg = %v", err)
	}
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(tracer)

	span, ctx, err := middle.GenerateSpanWithContext(context.Background(), "client-arith-mul")
	if err != nil {
		fmt.Printf("failed to generate span with context,err_msg = %v\n", err)
		return
	}
	defer span.Finish()

	args := data_struct.ArithReq{
		A: 10,
		B: 20,
	}
	reply := &data_struct.ArithResp{}

	fmt.Println("client")

	logrus.Infof("ctx:%v", ctx)
	err = xclient.Call(ctx, "Mul", args, reply)
	if err != nil {
		fmt.Printf("failed to call :%v\n", err)
	}

	fmt.Printf("%v\n", ctx)

	fmt.Printf("%d * %d = %d\n", args.A, args.B, reply.C)
}
