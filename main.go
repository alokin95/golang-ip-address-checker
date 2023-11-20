package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const lastKnownIpFile string = "last_known_ip_address.txt"

func main() {
	lastKnownIp := getLastKnownIp()
	publicIp := getPublicIp()

	if lastKnownIp != publicIp {
		updateLastKnownIp(publicIp)
		sendTelegramMessage("IP changed to: " + publicIp)
	}
}

func updateLastKnownIp(ip string) {
	ipToWrite := []byte(ip)

	lastIPFile, err := os.Create(lastKnownIpFile)

	checkError(err)

	defer lastIPFile.Close()

	lastIPFile.Write(ipToWrite)
}

func getPublicIp() string {
	resp, err := http.Get("https://api.ipify.org")

	checkError(err)

	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)

	checkError(err)

	return string(ip)
}

func getLastKnownIp() string {
	lastIPFile := lastKnownIpFile
	lastIPBytes, err := os.ReadFile(lastIPFile)

	checkError(err)

	lastIP := string(lastIPBytes)

	return lastIP
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

	checkError(err)

	return os.Getenv(key)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
