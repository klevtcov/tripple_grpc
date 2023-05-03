package main

import (
	"context"
	"flag"
	"fmt"

	"log"
	"strconv"

	"github.com/klevtcov/gRPC/tripple/client2/api"
	"google.golang.org/grpc"
)

func main() {

	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatalf("work only with two arguments")
	}

	x, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	y, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	// Запуск клиента на подключение к серверу
	conn, err := grpc.Dial(":8099", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewSquareAdderClient(conn)
	res, err := c.SquareAdd(context.Background(), &api.SquareAddRequest{X: int32(x), Y: int32(y)})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.GetResult())

}
