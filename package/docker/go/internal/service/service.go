package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-backend/internal"
)

type service struct {
	Entries  map[string]internal.Entry
	Versions map[string][]string
}

func New(entries map[string]internal.Entry, versions map[string][]string) *service {
	var Service service
	Service.Entries = entries
	Service.Versions = versions

	return &Service
}

func (s *service) CreateOrUpdateEntry(ctx context.Context, entry internal.Entry) (internal.Entry, error) {
	if entry.Id == "" {
		entry.Id = uuid.New().String()
	}

	s.Entries[entry.Id] = entry
	s.Versions[entry.Key] = append(s.Versions[entry.Key], entry.Id)

	return entry, nil
}

func (s *service) GetEntries(ctx context.Context) (map[string]internal.Entry, error) {
	return s.Entries, nil
}

func (s *service) GetEntry(ctx context.Context, id string) (internal.Entry, error) {
	if val, err := s.Entries[id]; err {
		return val, nil
	} else {
		return internal.Entry{}, errors.New("Entry doesn't exists.")
	}
}

func (s *service) GetKey(ctx context.Context, key string) ([]internal.Entry, error) {
	if val, err := s.Versions[key]; err {
		var entries []internal.Entry = make([]internal.Entry, 0, 1)

		for _, val := range val {
			entries = append(entries, s.Entries[val])
		}

		return entries, nil
	} else {
		return []internal.Entry{}, errors.New("Key doesn't exists.")
	}
}

func (s *service) DeleteEntry(ctx context.Context, id string) error {
	s.Versions[s.Entries[id].Key] = remove(s.Versions[s.Entries[id].Key], id)
	delete(s.Entries, id)
	return nil
}
