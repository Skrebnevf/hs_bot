package handlers

import (
	"fmt"
	"github/skrebnevf/hs_code/pkg/database"
	external "github/skrebnevf/hs_code/pkg/extertal"
	"log"
	"regexp"
	"strings"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func clearDescription(input, code string) string {
	re := regexp.MustCompile(`<[^>]*>`)

	clearHtml := re.ReplaceAllString(input, "")
	clearCode := strings.ReplaceAll(clearHtml, code, "")
	return clearCode
}

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

			if len(hs) == 0 {
				resp, err := external.GetTariffNumber(text)
				if err != nil {
					log.Printf("cannot connect to tariffnumber api, err: %v", err)
					return ctx.Send("Sorry external api or database falls down")
				}
				if resp.Total == 0 {
					return ctx.Send("Sorry this code does not exist in US and EU HS code database")
				}

				category := resp.Query[1:5]

				ru, err := database.GetRussianSunctionList(ctx, db, category)

				code := resp.Query
				desc := clearDescription(resp.Suggestions[0].Value, code)
				parentClass := resp.Query[1:3]

				err = database.WriteNewCode(db, code, desc, parentClass, category)
				if err != nil {
					log.Println(err)
				}

				WaitingForUserMessage[ctx.Message().Sender.ID] = false

				return ctx.Reply(fmt.Sprintf("<b>Entered code:</b> %s\n\n<b>Code discription:</b> %s\n\n<b>Include in Russian sunction list from:</b> %s\n\n <b>Information get from:</b> %s\n\n We will soon update this code for dangerous class and more information",
					code,
					desc,
					ru.From,
					resp.Suggestions[0].Data,
				))
			}

			if hs[0].ParentCategory.DangerousClass == "" {
				hs[0].ParentCategory.DangerousClass = "Does not have a danger class"
			}

			ru, err := database.GetRussianSunctionList(ctx, db, hs[0].ParentCategory.Category)

			WaitingForUserMessage[ctx.Message().Sender.ID] = false
			return ctx.Reply(fmt.Sprintf("<b>Entered code:</b> %s\n\n<b>Code discription:</b> %s\n\n<b>Dangerous class:</b> %v\n\n<b>Include in Russian sunction list from:</b> %s\n\n<b>Relate category:</b> %s\n\n<b>Category description:</b> %s",
				hs[0].Code,
				hs[0].Description,
				hs[0].ParentCategory.DangerousClass,
				ru.From,
				hs[0].ParentCategory.Category,
				strings.ToLower(hs[0].ParentCategory.Description)))
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
