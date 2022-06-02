package say_hello

import (
	"context"
	"gocalcharger/api/service"
	"google.golang.org/grpc"
	"time"
)

func SayHello(conn *grpc.ClientConn, name string) (*service.HelloReply, error) {
	c := service.NewGocalChargerServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := c.SayHello(ctx, &service.HelloRequest{Name: name})
	return r, err
}
