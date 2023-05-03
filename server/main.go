package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/klevtcov/gRPC/tripple/server/server"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server.UnimplementedMultiplierServer
}

func (s *GRPCServer) Myltiply(ctx context.Context, req *server.MyltiplyRequest) (*server.MyltiplyResponse, error) {
	return &server.MyltiplyResponse{
		Result: req.GetX() * req.GetX(),
	}, nil
}

func main() {
	fmt.Println("server app started")

	// Регистрируем gRPC сервер
	s := grpc.NewServer()
	srv := &GRPCServer{}
	server.RegisterMultiplierServer(s, srv)

	// Инициализируем сервер и порт
	lis, err := net.Listen("tcp", ":8088")
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

	fmt.Println("server app stopped")
}
