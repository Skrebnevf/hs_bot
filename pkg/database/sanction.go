package database

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

type SanctionCodeList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Code       string `json:"code"`
	Ban        string `json:"ban"`
}

func GetSanctionByCode(db *supabase.Client, table, code string) ([]SanctionCodeList, error) {
	resp, _, err := db.From(table).
		Select("*", "exact", false).
		Eq("code", code).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get hs code, error: %v", err)
	}

	var data []SanctionCodeList
	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse code data from db, error: %v", err)
	}

	return data, nil
}

type SanctionCategoryList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Category   string `json:"category"`
	Ban        string `json:"ban"`
}

func GetSanctionByCategory(db *supabase.Client, table, category string) ([]SanctionCategoryList, error) {
	resp, _, err := db.From(table).
		Select("*", "exact", false).
		Eq("category", category).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get sanction category, error: %v", err)
	}

	var data []SanctionCategoryList
	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse category data from db, error: %v", err)
	}

	return data, nil
}

type SanctionClassList struct {
	From       string `json:"from"`
	Source     string `json:"source"`
	LastUpdate string `json:"last_update"`
	Class      string `json:"class"`
	Ban        string `json:"ban"`
}

func GetSanctionByClass(db *supabase.Client, table, class string) ([]SanctionClassList, error) {
	resp, _, err := db.From(table).
		Select("*", "exact", false).
		Eq("class", class).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("cannot get sanction class list, error: %v", err)
	}

	var data []SanctionClassList
	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, fmt.Errorf("cannot parse class data from db, error: %v", err)
	}

	return data, nil
}
