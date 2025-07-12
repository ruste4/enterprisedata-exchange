package entity

import "time"

type ExchangeNode struct {
	ID           int       `json:"id" db:"id"`
	Description  string    `json:"c_description" db:"c_description"`
	NodeCode     string    `json:"node_code" db:"node_code"`
	ThisNodeCode string    `json:"this_node_code" db:"this_node_code"`
	Prefix       string    `json:"prefix" db:"prefix"`
	ThisPrefix   string    `json:"this_prefix" db:"this_prefix"`
	State        string    `json:"c_state" db:"c_state"`
	User         string    `json:"user" db:"user"`
	Pass         string    `json:"pass" db:"pass"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type MainExchangeParameters struct {
	CorrespondentNodeCode     string `json:"CorrespondentNodeCode"`
	DestinationInfobasePrefix string `json:"DestinationInfobasePrefix"`
	SourceInfobasePrefix      string `json:"SourceInfobasePrefix"`
	NodeCode                  string `json:"NodeCode"`
	ThisInfobaseDescription   string `json:"ThisInfobaseDescription"`
}

type CreateExchangeNodeDto struct {
	MainExchangeParameters MainExchangeParameters `json:"MainExchangeParameters"`
}
