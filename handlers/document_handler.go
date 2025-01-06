package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/halushko/kino-cat-core-go/nats_helper"
	"gopkg.in/telebot.v3"
	"log"
	"net/http"
	"os"
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

type getFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}

const TgBotApiGetFile = "https://api.telegram.org/bot%s/getFile?file_id=%s"
const TgBotApiDownload = "https://api.telegram.org/file/bot%s/%s"

func HandleDocuments(bot *telebot.Bot) {
	bot.Handle(telebot.OnDocument, func(c telebot.Context) error {
		userId := c.Chat().ID
		document := c.Message().Document
		mimeType := document.MIME

		log.Printf("[HandleDocuments] Отримано файл: %s", document.FileName)

		if !checkMimeType(mimeType) {
			return c.Send("Будь-ласка, відправте файл формату, що підтримується ботом")
		}

		fileId := document.FileID
		fileName := document.FileName
		fileSize := document.FileSize
		fileUrl, err := getFilePathFromTelegram(os.Getenv("BOT_TOKEN"), fileId)
		if err != nil {
			log.Printf("[HandleDocuments] Помилка отримання шляху до файлу: %v", err)
			return err
		}

		log.Printf("[HandleDocuments] userId:%d, uploadedFileId:%s, fileName:%s, size%d, mime:%s", userId, fileId, fileName, fileSize, mimeType)

		if err := nats_helper.PublishFileInfoMessage("TELEGRAM_INPUT_FILE_QUEUE", userId, fileId, fileName, fileSize, mimeType, fileUrl); err != nil {
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

func getFilePathFromTelegram(botToken, fileId string) (string, error) {
	url := fmt.Sprintf(TgBotApiGetFile, botToken, fileId)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("[getFilePathFromTelegram] Помилка виконання запиту getFile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[getFilePathFromTelegram] Неправильний статус відповіді getFile: %s", resp.Status)
	}

	var result getFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("[getFilePathFromTelegram] Помилка розбору відповіді getFile: %w", err)
	}

	if !result.Ok {
		return "", fmt.Errorf("[getFilePathFromTelegram] getFile повернув помилку")
	}

	fileURL := fmt.Sprintf(TgBotApiDownload, botToken, result.Result.FilePath)
	return fileURL, nil
}

func checkMimeType(mimeType string) bool {
	if mimeType != "application/x-bittorrent" {
		log.Printf("[checkMimeType] Невідомий MIME-тип: %s", mimeType)
		return false
	}
	return true
}
