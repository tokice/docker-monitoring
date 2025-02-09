package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const backendURL = "http://backend:8080/ping-results" // API бэкенда

func getContainerIPs() ([]string, error) {
	out, err := exec.Command("sh", "-c", "docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -q)").Output()
	if err != nil {
		return nil, err
	}
	ips := strings.Fields(string(out))
	return ips, nil
}

func ping(ip string) bool {
	out, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
	if err != nil {
		log.Printf("❌ Не удалось пинговать %s: %v", ip, err)
		return false
	}
	log.Printf("✅ Успешный пинг %s: %s", ip, strings.TrimSpace(string(out)))
	return true
}

func sendPingResult(ip string, success bool) {
	status := "failed"
	if success {
		status = "success"
	}
	_, err := http.Post(backendURL+"?ip="+ip+"&status="+status, "application/json", nil)
	if err != nil {
		log.Printf("❌ Ошибка отправки данных о пинге %s: %v", ip, err)
	}
}

func main() {
	fmt.Println("🚀 Pinger-сервис запущен!")

	for {
		ips, err := getContainerIPs()
		if err != nil {
			log.Printf("❌ Ошибка получения IP контейнеров: %v", err)
		} else {
			for _, ip := range ips {
				success := ping(ip)
				sendPingResult(ip, success)
			}
		}
		time.Sleep(10 * time.Second) // Интервал между проверками
	}
}
