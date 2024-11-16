package handlers

import (
	"fmt"
	external "github/skrebnevf/hs_code/pkg/extertal"
	"log"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func CommandHandlers(b *telebot.Bot, db *supabase.Client) {
	b.Handle("/start", func(ctx telebot.Context) error {
		fmt.Println(ctx.Message().Chat.ID)
		resp, err := external.GetTariffNumber("85177900")
		if err != nil {
			log.Println("cannot get tariff number, err: %v", err)
		}
		fmt.Println(resp.Suggestions[0].Value)
		return ctx.Reply(StartMsg)
	})

	b.Handle("/hs", func(ctx telebot.Context) error {
		WaitingForUserMessage[ctx.Message().Sender.ID] = true
		return ctx.Send(WaitingHsCodeMsg)
	})

	b.Handle("/help", func(ctx telebot.Context) error {
		AwaitngForward = true
		return ctx.Send(HelpCommandMsg)
	})
}
