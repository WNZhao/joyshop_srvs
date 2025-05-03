package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"joyshop_srvs/user_srv/handler"
	"joyshop_srvs/user_srv/proto"
	"net"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	PORT := flag.Int("port", 50051, "端口号")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *PORT)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
