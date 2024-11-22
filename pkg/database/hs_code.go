package database

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
	"gopkg.in/telebot.v4"
)

type HSCategory struct {
	Category       string `json:"category"`
	Description    string `json:"description"`
	ParentClass    string `json:"parent_class"`
	DangerousClass string `json:"dangerous_class"`
}

type HSCode struct {
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	ParentClass    string     `json:"parent_class"`
	ParentCategory HSCategory `json:"parent_category"`
}

type NewHSCode struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	ParentClass    string `json:"parent_class"`
	ParentCategory string `json:"parent_category"`
}

func WriteNewCode(db *supabase.Client, code, description, parentClass, parentCategory string) error {
	insert := NewHSCode{
		Code:           code,
		Description:    description,
		ParentClass:    parentClass,
		ParentCategory: parentCategory,
	}

	_, _, err := db.From("hs_code").
		Insert(insert, true, "uuid", "representation", "exact").
		Execute()
	if err != nil {
		return fmt.Errorf("cannot write new code, err: %v", err)
	}
	return nil
}

func GetHsCode(c telebot.Context, db *supabase.Client, code string) ([]HSCode, error) {
	resp, _, err := db.From("hs_code").
		Select("*, parent_category(*)", "exact", false).
		Eq("code", code).
		Execute()
	if err != nil {
		return []HSCode{}, fmt.Errorf("cannot get hs code, error: %v", err)
	}

	var data []HSCode
	json.Unmarshal(resp, &data)

	return data, nil
}
