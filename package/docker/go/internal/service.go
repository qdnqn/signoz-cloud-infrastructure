package internal

import (
	"context"
)

type Entry struct {
	Id    string `json:"id,omitempty"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Service interface {
	CreateOrUpdateEntry(context.Context, Entry) (Entry, error)
	GetEntries(context.Context) (map[string]Entry, error)
	GetEntry(context.Context, string) (Entry, error)
	GetKey(context.Context, string) ([]Entry, error)
	DeleteEntry(context.Context, string) error
}
