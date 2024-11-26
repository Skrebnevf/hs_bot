package handlers

import (
	"github/skrebnevf/hs_code/pkg/database"
	"log"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func CommandHandlers(b *telebot.Bot, db *supabase.Client) {
	b.Handle("/start", func(ctx telebot.Context) error {
		WaitingForUserMessage[ctx.Message().Sender.ID] = false
		AwaitngForward[ctx.Message().Sender.ID] = false
		resp, err := database.GetUser(ctx, db)
		if err != nil {
			log.Println(err)
		}

		if len(resp) == 0 {
			WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] = true
			return ctx.Reply(StartMsgWithOrgMsg)
		}
		return ctx.Reply(StartMsg)
	})

	b.Handle("/hs", func(ctx telebot.Context) error {
		WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] = false
		AwaitngForward[ctx.Message().Sender.ID] = false
		WaitingForUserMessage[ctx.Message().Sender.ID] = true
		return ctx.Send(WaitingHsCodeMsg)
	})

	b.Handle("/help", func(ctx telebot.Context) error {
		WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] = false
		WaitingForUserMessage[ctx.Message().Sender.ID] = false
		AwaitngForward[ctx.Message().Sender.ID] = true
		return ctx.Send(HelpCommandMsg)
	})

	b.Handle("/updateinfo", func(ctx telebot.Context) error {
		if ctx.Message().Chat.ID != ChatID {
			return nil
		} else {
			WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] = false
			WaitingForUserMessage[ctx.Message().Sender.ID] = false
			AwaitngForward[ctx.Message().Sender.ID] = false
			WaitingForMessage[ctx.Message().Sender.ID] = true
			return ctx.Send("Type info")
		}
	})
}
