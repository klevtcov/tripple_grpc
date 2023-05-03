package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/klevtcov/gRPC/tripple/client1/api"
	"github.com/klevtcov/gRPC/tripple/client1/server"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server.UnimplementedSquareAdderServer
}

func (s *GRPCServer) SquareAdd(ctx context.Context, req *server.SquareAddRequest) (*server.SquareAddResponse, error) {
	return &server.SquareAddResponse{
		Result: multiply(req.GetX()) + multiply(req.GetY()),
	}, nil
}

// Возводим в квадрат с помощью запроса на сервере
func multiply(x int32) int32 {
	// Запуск клиента на подключение к серверу
	conn, err := grpc.Dial(":8088", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return 0
	}

	c := api.NewMultiplierClient(conn)
	res, err := c.Myltiply(context.Background(), &api.MyltiplyRequest{X: int32(x)})
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return res.GetResult()
}

func main() {
	fmt.Println("client 1 started")

	// Регистрируем gRPC сервер
	s := grpc.NewServer()
	srv := &GRPCServer{}
	server.RegisterSquareAdderServer(s, srv)

	// Инициализируем сервер и порт
	lis, err := net.Listen("tcp", ":8099")
	if err != nil {
		log.Fatalf("cannot start tcp server: %s", err)
	}

	fmt.Println("listener started")

	// Добавляем обработчик к прослушиванию порта
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("cant start add server: %s", err)
		}
	}()

	fmt.Println("grpc handler started")

	// Ожидаем сигнала на закрытие приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("client1 app stopped")

}
