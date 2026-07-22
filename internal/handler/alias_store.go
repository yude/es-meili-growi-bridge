package handler

import (
	"sync"
)

type aliasEntry struct {
	Aliases map[string]interface{} `json:"aliases"`
}

type AliasStore struct {
	mu      sync.RWMutex
	aliases map[string]*aliasEntry
}

func NewAliasStore() *AliasStore {
	return &AliasStore{
		aliases: make(map[string]*aliasEntry),
	}
}

func (s *AliasStore) PutAlias(index, alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.aliases[index]
	if !ok {
		entry = &aliasEntry{Aliases: make(map[string]interface{})}
		s.aliases[index] = entry
	}
	entry.Aliases[alias] = struct{}{}
}

func (s *AliasStore) RemoveAlias(index, alias string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.aliases[index]
	if !ok {
		return
	}
	delete(entry.Aliases, alias)
	if len(entry.Aliases) == 0 {
		delete(s.aliases, index)
	}
}

func (s *AliasStore) ExistsAlias(name string, index string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.aliases[index]
	if !ok {
		return false
	}
	_, exists := entry.Aliases[name]
	return exists
}

func (s *AliasStore) GetIndexAliases(index string) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.aliases[index]
	if !ok {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})
	for alias := range entry.Aliases {
		result[alias] = struct{}{}
	}
	return result
}

func (s *AliasStore) ResolveAlias(alias string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for index, entry := range s.aliases {
		if _, ok := entry.Aliases[alias]; ok {
			return index, true
		}
	}
	return "", false
}

func (s *AliasStore) GetAllAliases() map[string]map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]map[string]interface{})
	for index, entry := range s.aliases {
		aliases := make(map[string]interface{})
		for alias := range entry.Aliases {
			aliases[alias] = struct{}{}
		}
		result[index] = aliases
	}
	return result
}
