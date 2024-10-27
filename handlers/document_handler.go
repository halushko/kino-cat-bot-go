package handlers

import (
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"gopkg.in/telebot.v3"
	"log"
)

type TorrentFile struct {
	ChatID   int64  `json:"chat_id"`
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	Text     string `json:"text"`
	Caption  string `json:"caption"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

func HandleDocuments(bot *telebot.Bot) {
	bot.Handle(telebot.OnDocument, func(c telebot.Context) error {
		userId := c.Chat().ID
		document := c.Message().Document
		mimeType := document.MIME

		log.Printf("[HandleDocuments] Отримано файл: %s", document.FileName)

		if document.MIME != "application/x-bittorrent" {
			return c.Send("[HandleDocuments] Будь-ласка, відправте .torrent файл.")
		}

		fileID := document.FileID
		fileName := document.FileName
		fileSize := document.FileSize

		log.Printf("[HandleDocuments] userId:%d, uploadedFileId:%s, fileName:%s, size%d, mime:%s", userId, fileID, fileName, fileSize, mimeType)

		if err := nats_helper.PublishFileInfoMessage("TELEGRAM_INPUT_FILE_QUEUE", userId, fileID, fileName, fileSize, mimeType); err != nil {
			log.Printf("[HandleDocuments] Error:%s", err)
			return err
		}

		if err := c.Send("Файл " + document.FileName + " додано до обробки"); err != nil {
			log.Printf("[HandleDocuments] Error:%s", err)
			return err
		}
		return nil
	})
}
