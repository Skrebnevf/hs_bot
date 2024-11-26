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
	tablesClass := []string{
		"ru_sanction_class",
		"bel_sanction_class",
	}

	tablesCategory := []string{
		"ru_sanction_category",
		"bel_sanction_category",
		"iran_sanction_category",
	}

	tablesCode := []string{
		"ru_sanctions_code",
		"bel_sanction_code",
		"iran_sanction_code",
	}

	checkSanction := func(description string, value string, table string, getSanction func() (interface{}, error)) error {
		sanctions, err := getSanction()
		if err != nil {
			log.Println(err)
			return ctx.Send(fmt.Sprintf("Sorry, we cannot get %s info", description))
		}

		switch s := sanctions.(type) {
		case []database.SanctionClassList:
			if len(s) > 0 {
				switch table {
				case "ru_sanction_class":
					return ctx.Send(fmt.Sprintf("For Russia sanction:\nFrom: %s\nClass: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Class, s[0].Ban, s[0].LastUpdate, s[0].Source))
				case "bel_sanction_class":
					return ctx.Send(fmt.Sprintf("For Belarus sanction:\nFrom: %s\nClass: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Class, s[0].Ban, s[0].LastUpdate, s[0].Source))
				}
			}
		case []database.SanctionCategoryList:
			if len(s) > 0 {
				switch table {
				case "ru_sanction_category":
					return ctx.Send(fmt.Sprintf("For Russia sanction:\nFrom: %s\nCategory: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Category, s[0].Ban, s[0].LastUpdate, s[0].Source))
				case "bel_sanction_category":
					return ctx.Send(fmt.Sprintf("For Belarus sanction:\nFrom: %s\nCategory: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Category, s[0].Ban, s[0].LastUpdate, s[0].Source))
				case "iran_sanction_category":
					return ctx.Send(fmt.Sprintf("For Iran sanction:\nFrom: %s\nCategory: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Category, s[0].Ban, s[0].LastUpdate, s[0].Source))
				}
			}
		case []database.SanctionCodeList:
			if len(s) > 0 {
				switch table {
				case "ru_sanctions_code":
					return ctx.Send(fmt.Sprintf("For Russia sanction:\nFrom: %s\nCode: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Code, s[0].Ban, s[0].LastUpdate, s[0].Source))
				case "bel_sanction_code":
					return ctx.Send(fmt.Sprintf("For Belarus sanction:\nFrom: %s\nCode: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Code, s[0].Ban, s[0].LastUpdate, s[0].Source))
				case "iran_sanction_code":
					return ctx.Send(fmt.Sprintf("For Iran sanction:\nFrom: %s\nCode: %s\nBan: %s\nLast Update: %s\nSource: %s", s[0].From, s[0].Code, s[0].Ban, s[0].LastUpdate, s[0].Source))
				}
			}
		}

		return nil
	}

	for _, table := range tablesClass {
		if err := checkSanction("class", class, table, func() (interface{}, error) {
			return database.GetSanctionByClass(db, table, class)
		}); err != nil {
			return err
		}
	}

	for _, table := range tablesCategory {
		if err := checkSanction("category", category, table, func() (interface{}, error) {
			return database.GetSanctionByCategory(db, table, category)
		}); err != nil {
			return err
		}
	}

	for _, table := range tablesCode {
		if err := checkSanction("code", code, table, func() (interface{}, error) {
			return database.GetSanctionByCode(db, table, code)
		}); err != nil {
			return err
		}
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

	if err := ctx.Send(fmt.Sprintf("Entered code: %s\nDescription: %s\nCategory: %s\nDangerous Class: %s\n\nIn the next one(s) message(s) will be provided information about the presence of the code in the sanctions lists.\n\nIf there is no message, then the information about the sanctions is not in our database.\n\nYou can always check up-to-date information about EU sanctions at https://eur-lex.europa.eu/homepage.html?locale=en.", text, desc, category, dangerousClass)); err != nil {
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

		case WaitingForMessage[userID]:
			msg := ctx.Message().Text
			WaitingForMessage[userID] = false
			users, err := database.GetUsersID(db)
			if err != nil {
				log.Println(err)
				return ctx.Reply("Sorry DB have error")
			}

			for _, user := range users {
				_, err := ctx.Bot().Send(&telebot.Chat{ID: user.ID}, msg)
				if err != nil {
					log.Println(err)
					return ctx.Reply("Cannot send update info message")
				}
			}

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
