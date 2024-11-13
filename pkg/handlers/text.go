package handlers

import (
	"fmt"
	"github/skrebnevf/hs_code/pkg/database"
	"log"
	"strings"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func TextHandlers(b *telebot.Bot, db *supabase.Client) {
	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		if WaitingForUserMessage[ctx.Message().Sender.ID] {
			text := ctx.Message().Text
			text = strings.TrimSpace(text)
			text = strings.ReplaceAll(text, ".", "")
			hs, err := database.GetHsCode(ctx, db, text)
			if err != nil {
				log.Printf("cannot get hs code, err: %v", err)
			}

			ru, err := database.GetRussianSunctionList(ctx, db, hs.ParentCategory.Category)

			WaitingForUserMessage[ctx.Message().Sender.ID] = false
			return ctx.Reply(fmt.Sprintf("<b>Entered code:</b> %s\n\n<b>Code discription:</b> %s\n\n<b>Dangerous class:</b> %v\n\n<b>Include in Russian sunction list from:</b> %s\n\n<b>Relate category:</b> %s\n\n<b>Category description:</b> %s",
				hs.Code,
				hs.Description,
				hs.ParentCategory.DangerousClass,
				ru.From,
				hs.ParentCategory.Category,
				strings.ToLower(hs.ParentCategory.Description)))
		}

		if AwaitngForward {
			msg := ctx.Message()

			var err error
			ForwardedMsg, err = b.Forward(&telebot.Chat{ID: ChatID}, msg)
			if err != nil {
				log.Printf("cannot forvared message, err: %v", err)
				AwaitngForward = false
				return ctx.Reply(CannotForwardedMsg)
			}

			AwaitngForward = false
			return ctx.Reply(CompletlyForwardedMsg)

		}
		return ctx.Reply(BaseMsg)
	})
}
