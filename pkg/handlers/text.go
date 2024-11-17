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
		if WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] {
			org := ctx.Message().Text
			org = strings.TrimSpace(org)

			err := database.WriteUser(ctx, db, org)
			if err != nil {
				log.Println(err)
			}

			WaitingForOrganizationInfoMsg[ctx.Message().Sender.ID] = false
			return ctx.Send(OrgMsg)
		}

		if WaitingForUserMessage[ctx.Message().Sender.ID] {
			text := ctx.Message().Text
			text = strings.TrimSpace(text)
			text = strings.ReplaceAll(text, ".", "")
			hs, err := database.GetHsCode(ctx, db, text)
			if err != nil {
				log.Println(err)
			}

			if len(hs) == 0 {
				resp, err := external.GetTariffNumber(text)
				if err != nil {
					log.Println(err)
					WaitingForUserMessage[ctx.Message().Sender.ID] = false
					return ctx.Send("Sorry external api or database falls down")
				}
				if resp.Total == 0 {
					WaitingForUserMessage[ctx.Message().Sender.ID] = false
					return ctx.Send("Sorry this code does not exist in US and EU HS code database")
				}
				if resp.Query == "" {
					WaitingForUserMessage[ctx.Message().Sender.ID] = false
					return ctx.Send("Sorry this code does not exist in US and EU HS code database")
				}

				msg := ctx.Message()
				code := resp.Query
				desc := clearDescription(resp.Suggestions[0].Value, code)
				category := resp.Query[0:4]
				parentClass := resp.Query[0:2]

				err = database.WriteNewCode(db, code, desc, parentClass, category)
				if err != nil {
					log.Println(err)
				}

				ForwardedMsg, err = b.Send(&telebot.Chat{ID: ChatID}, msg.Text+" need to add decription and other options")
				if err != nil {
					log.Println(err)
				}

				ru, err := database.GetRussianSunctionList(ctx, db, category)
				WaitingForUserMessage[ctx.Message().Sender.ID] = false
				return ctx.Reply(fmt.Sprintf("Entered code: %s\n\nCode discription: %s\n\nInclude in Russian sunction list from: %s\n\nInformation get from: %s\n\n We will soon update this code for dangerous class and more information",
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
			if err != nil {
				log.Println(err)
			}

			WaitingForUserMessage[ctx.Message().Sender.ID] = false
			return ctx.Reply(fmt.Sprintf("Entered code: %s\n\nCode discription: %s\n\nDangerous class: %v\n\nInclude in Russian sunction list from: %s\n\nRelate category: %s\n\nCategory description: %s",
				hs[0].Code,
				hs[0].Description,
				hs[0].ParentCategory.DangerousClass,
				ru.From,
				hs[0].ParentCategory.Category,
				strings.ToLower(hs[0].ParentCategory.Description)))
		}

		if AwaitngForward[ctx.Message().Sender.ID] {
			msg := ctx.Message()

			var err error
			ForwardedMsg, err = b.Forward(&telebot.Chat{ID: ChatID}, msg)
			if err != nil {
				log.Println(err)
				AwaitngForward[ctx.Message().Sender.ID] = false
				return ctx.Reply(CannotForwardedMsg)
			}

			AwaitngForward[ctx.Message().Sender.ID] = false
			return ctx.Reply(CompletlyForwardedMsg)

		}
		return ctx.Reply(BaseMsg)
	})
}
