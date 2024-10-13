package handlers

import (
	"encoding/json"
	"github.com/halushko/kino-cat-core-go"
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
		document := c.Message().Document

		log.Printf("[TorrentFileHandler] Отримано файл: %s", document.FileName)

		if document.MIME != "application/x-bittorrent" {
			return c.Send("[TorrentFileHandler] Будь-ласка, відправте .torrent файл.")
		}

		chatId := c.Chat().ID
		fileID := document.FileID
		fileName := document.FileName
		fileSize := document.FileSize
		mimeType := document.MIME
		messageText := c.Message().Text
		caption := c.Message().Caption

		msg := TorrentFile{
			ChatID:   chatId,
			FileID:   fileID,
			FileName: fileName,
			Text:     messageText,
			Caption:  caption,
			Size:     fileSize,
			MimeType: mimeType,
		}

		log.Printf(
			"[TorrentFileHandler] chatId:%d, uploadedFileId:%s, fileName:%s, message:%s, caption:%s",
			chatId, fileID, fileName, messageText, caption,
		)

		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[TorrentFileHandler] Error:%s", err)
			return err
		}

		if err = kino_cat_core_go.PublishToNATS("TELEGRAM_INPUT_FILE_QUEUE", jsonData); err != nil {
			log.Printf("[TorrentFileHandler] Error:%s", err)
			return err
		}

		if err = c.Send("Файл " + document.FileName + " додано до обробки"); err != nil {
			log.Printf("[TorrentFileHandler] Error:%s", err)
			return err
		}
		return nil
	})
}
