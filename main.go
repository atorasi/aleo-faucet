package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/atorasi/aleo-faucet/discord"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("File '%s' not found", filename))
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func main() {
	log.Println("\nSay hello to braindead Aleo founder")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(color.RedString("File '.env' not found"))
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	sleepFrom, sleepTo := os.Getenv("SLEEP_FROM"), os.Getenv("SLEEP_TO")

	proxy := os.Getenv("PROXY")
	useProxy := false
	if proxy == "1" {
		useProxy = true
	}
	delayMin, errMin := strconv.Atoi(sleepFrom)
	delayMax, errMax := strconv.Atoi(sleepTo)
	if errMin != nil || errMax != nil {
		log.Fatal(color.RedString("Error converting strings to numbers."))
	}

	tokenList := readFile("tokens.txt")
	addrList := readFile("addresses.txt")
	proxyList := readFile("addresses.txt")
	if len(addrList) > len(tokenList) {
		log.Fatal(color.RedString("Tokens must be <= addresses"))
	}

	log.Print(color.MagentaString("Uploaded %d tokens and %d addresses", len(tokenList), len(addrList)))
	for i, address := range addrList {

		client := discord.NewClient(i, tokenList[i], useProxy)
		client.SetProxy(proxyList)

		if err := client.SendMessage("1153398971109220432", "/sendcredits "+address); err != nil {
			log.Print(color.RedString("Acc.%v Error sending message: %v.", i+1, err))
		} else {
			log.Print(color.GreenString("Acc.%v Successfully sent message.", i+1))
		}

		delay := delayMin + rng.Intn(delayMax-delayMin)
		log.Print(color.MagentaString("Acc.%v Sleeping %d seconds before next account.", i+1, delay))
		time.Sleep(time.Second * time.Duration(delay))
	}

	log.Println("omg.")
}
