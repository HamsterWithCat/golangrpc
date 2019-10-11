package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"rpclearn/rpcx/data_struct"
)

var (
	consulAddr = flag.String("consulAddr", "49.235.1.29:8500", "consul address")
	basePath   = flag.String("base", "/rpcx_test/Arith", "prefix path")
)

func ConsulClient() {
	flag.Parse()

	d := client.NewConsulDiscovery(*basePath, "", []string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &data_struct.ArithReq{
		A: 10,
		B: 20,
	}

	reply := &data_struct.ArithResp{}

	err := xclient.Call(context.Background(), "Add", args, reply)
	if err != nil {
		logrus.Errorf("call Arith Mul method failed,err_msg = %v", err)
	}

	logrus.Infof("%d + %d = %d", args.A, args.B, reply.C)
}

func main() {
	ConsulClient()
}
