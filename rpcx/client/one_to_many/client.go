package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/prometheus/common/log"
	"github.com/smallnest/rpcx/client"
	"rpclearn/rpcx/data_struct"
	"rpclearn/rpcx/middle"
)
var (
	addr = flag.String("addr","127.0.0.1:8568","server address")
	addr2 = flag.String("addr2","127.0.0.1:8569","server address")
)

func init(){

}

// 一对一，服务器地址硬编码或写在配置文件中
func main(){
	flag.Parse()

	//d := client.NewPeer2PeerDiscovery("tcp@"+*addr,"")
	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key:*addr},{Key:*addr2}})
	xclient := client.NewXClient("Arith",client.Failtry,client.RandomSelect,d,client.DefaultOption)
	defer xclient.Close()

	/**************************/
	// zipkin上报
	reporter := zipkinhttp.NewReporter("http://49.235.1.29:9411/api/v1/spans")
	defer reporter.Close()

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
    /**************************/

	span,ctx,err := middle.GenerateSpanWithContext(context.Background(),"start point")
	if err != nil {
		fmt.Printf("failed to generate span with context,err_msg = %v\n",err)
		return
	}
	span.SetTag("span.kind","client")
	span.Finish()
	/*span := tracer.StartSpan("oneToManyClient-Main",ext.SpanKindRPCClient)
	ctx := opentracing.ContextWithSpan(context.Background(),span)*/

	args := data_struct.ArithReq{
		A:10,
		B:20,
	}
	reply := &data_struct.ArithResp{}


	err = xclient.Call(ctx,"Mul",args,reply)
	if err != nil {
		fmt.Printf("failed to call :%v\n",err)
	}

	fmt.Printf("%v\n",ctx)


	fmt.Printf("%d * %d = %d\n",args.A,args.B,reply.C)
}
