package models

import (
		"github.com/go-telegram/bot/models"
		"fmt"
		"log"
)

type MemoryStorage struct{
	stor map[string][]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{ stor: make(map[string][]string)}
}

func (ms *MemoryStorage) InsertUpdate(update *models.Update) {
	key := UpdateToKey(update)
	log.Print(key)
}

func UpdateToKey(update *models.Update) string {
	query_id := update.CallbackQuery.ID
	user_id := update.CallbackQuery.From.ID
	chat_id := update.CallbackQuery.Message.Message.Chat.ID
	return fmt.Sprintf("%s_%d_%d", query_id, user_id, chat_id)
}