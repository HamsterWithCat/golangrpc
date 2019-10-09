package service

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	. "rpclearn/rpcx/data_struct"
)

type Arith int

func (t *Arith)Mul(ctx context.Context,args *ArithReq,reply *ArithResp) error {
	/*fmt.Printf("ctx:%v\n",ctx)
	span,ctx,_ := middle.GenerateSpanWithContext(ctx,"Arith-Mul")
	span.SetTag("http.method","Mul")
	span.SetTag("span.kind","server")
	span.LogKV("step","rpcx")
	defer span.Finish()*/
	span := opentracing.SpanFromContext(ctx)
	if nil == span {
		return errors.New("span is nil")
	}
	span.SetTag("http.method","Mul")
	span.SetTag("span.kind","server")
	span.LogKV("step","rpcx")
	defer span.Finish()

	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Add(ctx context.Context,args *ArithReq,reply *ArithResp) error {
	reply.C = args.A +  args.B
	return nil
}

func (t *Arith) Say(ctx context.Context,args *string,reply *string) error {
	*reply = "hello " + *args
	return nil
}