package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc_demo/pb" // Импортируем сгенерированный пакет

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Подключаемся к серверу
	// Используем insecure, так как у нас нет SSL-сертификатов для локального теста
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	// 2. Создаем "клиента" для нашего сервиса
	client := pb.NewCommandServiceClient(conn)

	// Контекст с таймаутом (чтобы не ждать вечно, если сервер упал)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// --- ОТПРАВЛЯЕМ КОМАНДЫ ---

	// Команда 1: Время
	fmt.Println("\n--- Отправляю команду GetTime ---")
	timeResp, err := client.GetTime(ctx, &pb.TimeRequest{})
	if err != nil {
		log.Fatalf("Ошибка GetTime: %v", err)
	}
	fmt.Printf("✅ Ответ от сервера: Сейчас время %s\n\n", timeResp.CurrentTime)

	// Команда 2: Калькулятор (Сложение)
	fmt.Println("--- Отправляю команду Calculate (Сложение) ---")
	calcResp, err := client.Calculate(ctx, &pb.CalcRequest{
		A:  15.5,
		B:  4.5,
		Op: pb.Operation_ADD,
	})
	if err != nil {
		log.Fatalf("Ошибка Calculate: %v", err)
	}
	fmt.Printf("✅ Ответ от сервера: 15.5 + 4.5 = %.2f\n\n", calcResp.Result)

	// Команда 3: Калькулятор (Умножение)
	fmt.Println("--- Отправляю команду Calculate (Умножение) ---")
	calcResp2, err := client.Calculate(ctx, &pb.CalcRequest{
		A:  10,
		B:  5,
		Op: pb.Operation_MULTIPLY,
	})
	if err != nil {
		log.Fatalf("Ошибка Calculate: %v", err)
	}
	fmt.Printf("✅ Ответ от сервера: 10 * 5 = %.2f\n\n", calcResp2.Result)

	// Команда 4: Текст
	fmt.Println("--- Отправляю команду ProcessText ---")
	textResp, err := client.ProcessText(ctx, &pb.TextRequest{
		Text: "привет, grpc!",
	})
	if err != nil {
		log.Fatalf("Ошибка ProcessText: %v", err)
	}
	fmt.Printf("✅ Ответ от сервера: %s\n", textResp.ProcessedText)
}