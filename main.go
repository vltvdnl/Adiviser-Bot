package main

import (
	"flag"
	"log"

	"github.com/vltvdnl/Adviser-Bot/clients/telegram"
)

func main() {
	tgClient := telegram.New("api.telegram.org", MustToken())

}

func MustToken() string {
	token := flag.String("token-bot-token", "", "token for usage to tg bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
