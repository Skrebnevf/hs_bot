package handlers

import (
	"fmt"
	"github/skrebnevf/hs_code/pkg/database"
	external "github/skrebnevf/hs_code/pkg/extertal"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

func clearDescription(input, code string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	clearHtml := re.ReplaceAllString(input, "")
	return strings.ReplaceAll(clearHtml, code, "")
}

func handleSanctions(ctx telebot.Context, db *supabase.Client, code, category, class string) error {
	checkSanction := func(description string, value string, getSanction func() (interface{}, error)) error {
		sanctions, err := getSanction()
		if err != nil {
			log.Println(err)
			return ctx.Send(fmt.Sprintf("Sorry, we cannot get %s info", description))
		}

		switch s := sanctions.(type) {
		case []database.RuSanctionClassList:
			if len(s) > 0 {
				return ctx.Send(fmt.Sprintf("Sanction:\nFrom: %s\nClass: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Class, s[0].Ban, s[0].LastUpdate, s[0].Source))
			}
		case []database.RuSanctionCategoryList:
			if len(s) > 0 {
				return ctx.Send(fmt.Sprintf("Sanction:\nFrom: %s\nCategory: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Category, s[0].Ban, s[0].LastUpdate, s[0].Source))
			}
		case []database.RuSanctionCodeList:
			if len(s) > 0 {
				return ctx.Send(fmt.Sprintf("Sanction:\nFrom: %s\nCode: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Code, s[0].Ban, s[0].LastUpdate, s[0].Source))
			}
		}

		return nil
	}

	if err := checkSanction("class", class, func() (interface{}, error) {
		return database.GetSanctionByClass(db, class)
	}); err != nil {
		return err
	}

	if err := checkSanction("category", category, func() (interface{}, error) {
		return database.GetSanctionByCategory(db, category)
	}); err != nil {
		return err
	}

	if err := checkSanction("code", code, func() (interface{}, error) {
		return database.GetSanctionByCode(db, code)
	}); err != nil {
		return err
	}

	return nil
}

func handleUserMessage(ctx telebot.Context, db *supabase.Client, text string) error {
	text = strings.TrimSpace(strings.ReplaceAll(text, ".", ""))
	if len(text) < 6 {
		return ctx.Send("Length of HS code should be 6 or more digits")
	}

	hs, err := database.GetHsCode(ctx, db, text)
	if err != nil {
		log.Println(err)
		return ctx.Send("Sorry something went wrong, database is not available")
	}
	if len(hs) == 0 {
		resp, err := external.GetTariffNumber(text)
		if err != nil {
			log.Println(err)
			return ctx.Send("External server error, HS code cannot read from third party API service")
		}
		if resp.Total == 0 || resp.Query == "" {
			return ctx.Send("Sorry, this code does not exist in US and EU HS code database")
		}

		code, desc := resp.Query, clearDescription(resp.Suggestions[0].Value, resp.Query)
		category, class := code[:4], code[:2]

		err = performDBOperation(ctx, func() error {
			return database.WriteNewCode(db, code, desc, class, category)
		}, 40*time.Second, "cannot write code to db")
		if err != nil {
			log.Println(err)
		}

		_, err = ctx.Bot().Send(&telebot.Chat{ID: ChatID}, fmt.Sprintf(text+" need to add decription and other options"))
		if err != nil {
			log.Println("Failed to forward message to ChatID:", err)
		}

		err = ctx.Send("It looks like the code was missing in our database, later we will make a more detailed description and add a danger class.")
		if err != nil {
			return err
		}

		err = ctx.Send(fmt.Sprintf("Entered code: %s\nDescription: %s\nCategory: %s", code, desc, category))
		if err != nil {
			return err
		}
		return handleSanctions(ctx, db, code, category, class)
	}

	code, desc := text[:6], hs[0].Description
	category, class := code[:4], code[:2]
	dangerousClass := hs[0].ParentCategory.DangerousClass

	if dangerousClass == "" {
		dangerousClass = "Without dangerous class"
	}

	if err := ctx.Send(fmt.Sprintf("Entered code: %s\nDescription: %s\nCategory: %s\nDangerous Class: %s", text, desc, category, dangerousClass)); err != nil {
		return err
	}
	return handleSanctions(ctx, db, code, category, class)
}

func performDBOperation(ctx telebot.Context, dbFunc func() error, timeout time.Duration, errorMsg string) error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		errChan <- dbFunc()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			log.Println(errorMsg, err)
			ctx.Bot().Send(&telebot.Chat{ID: ChatID}, fmt.Sprintf("%s: %v", errorMsg, err))
			return err
		}
	case <-time.After(timeout):
		log.Println("database timeout")
		ctx.Bot().Send(&telebot.Chat{ID: ChatID}, "DB Timeout, check supabase ASAP!")
		return fmt.Errorf("database timeout")
	}
	return nil
}

func TextHandlers(b *telebot.Bot, db *supabase.Client) {
	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		userID := ctx.Message().Sender.ID

		switch {
		case WaitingForOrganizationInfoMsg[userID]:
			org := strings.TrimSpace(ctx.Message().Text)
			err := performDBOperation(ctx, func() error {
				return database.WriteUser(ctx, db, org)
			}, 40*time.Second, "cannot write user to db")
			WaitingForOrganizationInfoMsg[userID] = false
			if err == nil {
				return ctx.Send(OrgMsg)
			}

		case WaitingForUserMessage[userID]:
			WaitingForUserMessage[userID] = false
			err := performDBOperation(ctx, func() error {
				return database.WriteMsgLog(ctx, db)
			}, 40*time.Second, "cannot write user log to db")
			if err == nil {
				return handleUserMessage(ctx, db, ctx.Message().Text)
			}

		case AwaitngForward[userID]:
			AwaitngForward[userID] = false
			_, err := b.Forward(&telebot.Chat{ID: ChatID}, ctx.Message())
			if err != nil {
				log.Println(err)
				return ctx.Reply(CannotForwardedMsg)
			}
			return ctx.Reply(CompletlyForwardedMsg)

		default:
			err := performDBOperation(ctx, func() error {
				return database.WriteMsgLog(ctx, db)
			}, 40*time.Second, "cannot write user log to db")
			if err == nil {
				return ctx.Reply(BaseMsg)
			}
		}
		return nil
	})
}
