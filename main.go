package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	text, err := getPublicIp()

	if err != nil {
		sendTelegramMessage(err.Error())
	}

	sendTelegramMessage(text)
}

func getPublicIp() (string, error) {
	resp, err := http.Get("https://api.ipify.org")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)

	return string(ip), err
}

func sendTelegramMessage(message string) {
	var botToken string = getEnvVariable("TELEGRAM_BOT_API_TOKEN")
	var chatId string = getEnvVariable("TELEGRAM_CHAT_ID")

	endpoint := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	data := url.Values{}
	data.Set("chat_id", chatId)
	data.Set("text", message)

	_, err := http.PostForm(endpoint, data)
	if err != nil {
		log.Fatal("Error sending message: ", err)
	}
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
