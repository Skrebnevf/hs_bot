package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

type HSCategory struct {
	UUID           uuid.UUID `json:"uuid"`
	Category       string    `json:"category"`
	Description    string    `json:"description"`
	ParentClass    string    `json:"parent_class"`
	DangerousClass string    `json:"dangerous_class"`
}

type HSCode struct {
	UUID           uuid.UUID  `json:"uuid"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	ParentClass    string     `json:"parent_class"`
	ParentCategory HSCategory `json:"parent_category"`
}

func GetHsCode(c telebot.Context, db *supabase.Client, code string) (HSCode, error) {
	resp, _, err := db.From("hs_code").
		Select("*, parent_category(*)", "exact", false).
		Eq("code", code).
		Execute()
	if err != nil {
		return HSCode{}, fmt.Errorf("cannot get hs code, error: %v", err)
	}

	var data []HSCode
	json.Unmarshal(resp, &data)

	if len(data) == 0 {
		return HSCode{
				Code:        "Code not available or does not exist",
				Description: "Description is empty",
				ParentCategory: HSCategory{
					Description:    "Description is empty",
					Category:       "Category not available or does not exist",
					DangerousClass: "Dangerous class not available or does not exist",
				},
			},
			nil
	}

	if data[0].ParentCategory.DangerousClass == "" {
		data[0].ParentCategory.DangerousClass = "Does not have a danger class"
	}

	return data[0], nil
}

type RuSanctionList struct {
	UUID     uuid.UUID `json:"uuid"`
	From     string    `json:"from"`
	Category string    `json:"category"`
}

func GetRussianSunctionList(c telebot.Context, db *supabase.Client, code string) (RuSanctionList, error) {
	resp, _, err := db.From("ru_sanctions").
		Select("*", "exact", false).
		Eq("category", code).
		Execute()
	if err != nil {
		return RuSanctionList{}, fmt.Errorf("cannot get hs code, error: %v", err)
	}

	var data []RuSanctionList
	if err := json.Unmarshal(resp, &data); err != nil {
		log.Println(err)
	}

	if len(data) == 0 {
		return RuSanctionList{
				From: "Not included in the sanctions list",
			},
			nil
	}

	fmt.Println(data[0])

	return data[0], nil
}
