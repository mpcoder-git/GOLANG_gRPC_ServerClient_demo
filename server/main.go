package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"grpc_demo/pb" // Импортируем сгенерированный пакет (путь зависит от твоего go.mod)

	"google.golang.org/grpc"
)

// Создаем структуру, которая будет реализовывать наш сервис
type server struct {
	pb.UnimplementedCommandServiceServer // Встраиваем для безопасности (чтобы не упало, если забудешь метод)
}

// Реализуем команду 1: GetTime
func (s *server) GetTime(ctx context.Context, req *pb.TimeRequest) (*pb.TimeResponse, error) {
	fmt.Println("[Сервер] Выполняю команду GetTime...")
	return &pb.TimeResponse{
		CurrentTime: time.Now().Format("15:04:05"),
	}, nil
}

// Реализуем команду 2: Calculate
func (s *server) Calculate(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	fmt.Printf("[Сервер] Выполняю команду Calculate: %v %v %v\n", req.A, req.Op, req.B)
	
	var result float64
	
	// Используем switch вместо if-else (это то, что рекомендует линтер)
	switch req.Op {
	case pb.Operation_ADD:
		result = req.A + req.B
	case pb.Operation_MULTIPLY:
		result = req.A * req.B
	default:
		// Хорошая практика: вернуть понятную ошибку, если пришла неизвестная команда
		return nil, fmt.Errorf("неизвестная операция: %v", req.Op)
	}

	return &pb.CalcResponse{Result: result}, nil
}

// Реализуем команду 3: ProcessText
func (s *server) ProcessText(ctx context.Context, req *pb.TextRequest) (*pb.TextResponse, error) {
	fmt.Printf("[Сервер] Выполняю команду ProcessText: %s\n", req.Text)
	
	return &pb.TextResponse{
		ProcessedText: strings.ToUpper(req.Text),
	}, nil
}

func main() {
	
	// Слушаем порт 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось слушать порт: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()
	// Регистрируем наш сервис в сервере
	pb.RegisterCommandServiceServer(s, &server{})

	fmt.Println("🚀 Сервер запущен на порту 50051...")
	// Запускаем сервер (блокирующий вызов)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}