package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
)

func main() {
	socketPath := "/tmp/dns_server.sock"

	// Подключаемся к Unix-сокету
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("Failed to connect to Unix socket: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Создаём DNS-запрос
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)
	msg.RecursionDesired = true

	// Упаковываем запрос в бинарный формат
	data, err := msg.Pack()
	if err != nil {
		fmt.Printf("Failed to pack DNS message: %v\n", err)
		os.Exit(1)
	}

	// Отправляем запрос через сокет
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Failed to write to Unix socket: %v\n", err)
		os.Exit(1)
	}

	// Читаем ответ от сервера
	buf := make([]byte, 65535)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read from Unix socket: %v\n", err)
		os.Exit(1)
	}

	// Распаковываем ответ в структуру DNS
	resp := new(dns.Msg)
	if err := resp.Unpack(buf[:n]); err != nil {
		fmt.Printf("Failed to unpack DNS response: %v\n", err)
		os.Exit(1)
	}

	// Выводим результат
	fmt.Println("Received DNS response:")
	fmt.Println(resp)
}
