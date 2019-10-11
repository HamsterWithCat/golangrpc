package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"rpclearn/rpcx/data_struct"
)

var (
	addr = flag.String("addr", "127.0.0.1:8568", "server address")
)

// 一对一，服务器地址硬编码或写在配置文件中
func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := data_struct.ArithReq{
		A: 10,
		B: 20,
	}
	reply := &data_struct.ArithResp{}

	err := xclient.Call(context.Background(), "Mul", args, reply)

	if err != nil {
		fmt.Printf("failed to call :%v\n", err)
	}
	fmt.Printf("%d * %d = %d\n", args.A, args.B, reply.C)
}
