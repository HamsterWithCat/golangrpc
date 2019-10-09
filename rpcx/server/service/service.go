package service

import (
	"context"
	"fmt"
	. "rpclearn/rpcx/data_struct"
	"rpclearn/rpcx/middle"
)

type Arith int

func (t *Arith)Mul(ctx context.Context,args *ArithReq,reply *ArithResp) error {
	span,ctx,_ := middle.GenerateSpanWithContext(ctx,"Arith-Mul")
	span.SetTag("method","Mul")
	span.LogKV("step","rpcx")
	fmt.Printf("span:%v\n",span)
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