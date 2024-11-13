package handlers

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func CommandHandlers(b *telebot.Bot, db *supabase.Client) {
	b.Handle("/start", func(ctx telebot.Context) error {
		fmt.Println(ctx.Message().Chat.ID)
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
