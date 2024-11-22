package database

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

type RuSanctionCodeList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Code       string `json:"code"`
	Ban        string `json:"ban"`
}

func GetSanctionByCode(db *supabase.Client, code string) ([]RuSanctionCodeList, error) {
	resp, _, err := db.From("ru_sanctions_code").
		Select("*", "exact", false).
		Eq("code", code).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get hs code, error: %v", err)
	}

	var data []RuSanctionCodeList
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse code data from db, error: %v", err)
	}

	return data, nil
}

type RuSanctionCategoryList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Category   string `json:"category"`
	Ban        string `json:"ban"`
}

func GetSanctionByCategory(db *supabase.Client, category string) ([]RuSanctionCategoryList, error) {
	resp, _, err := db.From("ru_sanction_category").
		Select("*", "exact", false).
		Eq("category", category).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get sanction category, error: %v", err)
	}

	var data []RuSanctionCategoryList
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse category data from db, error: %v", err)
	}

	return data, nil
}

type RuSanctionClassList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Class      string `json:"class"`
	Ban        string `json:"ban"`
}

func GetSanctionByClass(db *supabase.Client, class string) ([]RuSanctionClassList, error) {
	resp, _, err := db.From("ru_sanction_class").
		Select("*", "exact", false).
		Eq("class", class).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get sanction class list, error: %v", err)
	}

	var data []RuSanctionClassList
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse class data from db, error: %v", err)
	}

	return data, nil
}
