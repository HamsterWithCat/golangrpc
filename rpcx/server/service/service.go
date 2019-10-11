package service

import (
	"context"
	. "rpclearn/rpcx/data_struct"
	"rpclearn/rpcx/middle"
)

type Arith int

func (t *Arith) Mul(ctx context.Context, args *ArithReq, reply *ArithResp) error {
	span, ctx, _ := middle.GenerateSpanWithContext(ctx, "server-arith-mul")
	span.SetTag("http.method", "Mul")
	span.SetTag("span.kind", "server")
	span.LogKV("step", "rpcx")
	defer span.Finish()

	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Add(ctx context.Context, args *ArithReq, reply *ArithResp) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Say(ctx context.Context, args *string, reply *string) error {
	*reply = "hello " + *args
	return nil
}
