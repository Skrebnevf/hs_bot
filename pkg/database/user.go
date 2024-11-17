package database

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Organization string `json:"organization"`
}

func WriteUser(ctx telebot.Context, db *supabase.Client, org string) error {
	id := ctx.Sender().ID
	username := ctx.Sender().Username
	name := ctx.Sender().FirstName
	surname := ctx.Sender().LastName
	organization := org

	insert := User{
		ID:           id,
		Username:     username,
		Name:         name,
		Surname:      surname,
		Organization: organization,
	}

	_, _, err := db.From("users").
		Insert(insert, true, "uuid", "representation", "exact").
		Execute()
	if err != nil {
		return fmt.Errorf("cannot write user, err: %v", err)
	}
	return nil
}

func GetUser(ctx telebot.Context, db *supabase.Client) ([]User, error) {
	id := ctx.Sender().ID
	stringID := strconv.FormatInt(int64(id), 10)
	resp, _, err := db.From("users").
		Select("id", "exact", false).
		Eq("id", stringID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get user from db, err: %v", err)
	}

	var u []User
	err = json.Unmarshal(resp, &u)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal user data, err: %v", err)
	}

	return u, nil
}
