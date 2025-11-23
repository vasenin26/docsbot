package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vasenin26/docsbot/internal/metrics"
)

// TelegramClient минимальный интерфейс для мокирования внешнего клиента.
type TelegramClient interface {
	GetMe() (tgbotapi.User, error)
	GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Request(tgbotapi.Chattable) (tgbotapi.APIResponse, error)
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

// Bot инкапсулирует логику работы с Telegram API.
type Bot struct {
	client   TelegramClient
t	username string
	logger   *log.Logger
	mu       sync.Mutex
}

// NewBot создаёт новый Bot; token обязателен.
func NewBot(token string) (*Bot, error) {
	if token == "" {
		return nil, errors.New("telegram token is empty")
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = false

	me, err := api.GetMe()
	if err != nil {
		return nil, fmt.Errorf("getMe failed: %w", err)
	}

	b := &Bot{
		client:   api,
		username: me.UserName,
		logger:   log.New(os.Stderr, "telegram: ", log.LstdFlags),
	}
	return b, nil
}

// Start запускает цикл получения обновлений (long polling).
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.client.GetUpdatesChan(u)

	b.logger.Printf("started, username=%s", b.username)

	for {
		select {
		case <-ctx.Done():
			b.logger.Printf("stopping: %v", ctx.Err())
			return nil
		case update, ok := <-updates:
			if !ok {
				b.logger.Printf("updates channel closed")
				return nil
			}
			// Non-blocking handling --- spawn goroutine per update to keep loop responsive
			go func(up tgbotapi.Update) {
				if err := b.handleUpdate(up); err != nil {
					metrics.TelegramMessagesErrors.Inc()
					b.logger.Printf("handleUpdate error: %v", err)
				}
			}(update)
		}
	}
}

// handleUpdate обрабатывает отдельный Update.
func (b *Bot) handleUpdate(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	msg := update.Message

	// Only text messages
	if msg.Text == "" {
		return nil
	}

	// Only group / supergroup
	if !(msg.Chat.IsGroup() || msg.Chat.IsSuperGroup()) {
		return nil
	}

	metrics.TelegramMessagesReceived.Inc()

	if !b.isAddressedToBot(msg) {
		return nil
	}

	metrics.TelegramMessagesAddressed.Inc()

	// Write text to stdout as required
	fmt.Fprintf(os.Stdout, "telegram: received message from chat %d: %s\n", msg.Chat.ID, msg.Text)

	// Reply with "Ок"
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Ок")
	reply.ReplyToMessageID = msg.MessageID

	if _, err := b.client.Send(reply); err != nil {
		metrics.TelegramMessagesErrors.Inc()
		return err
	}

	metrics.TelegramMessagesProcessed.Inc()

	b.logger.Printf("processed addressed message chat=%d user=%d snippet=%.64s", msg.Chat.ID, msg.From.ID, msg.Text)
	return nil
}

// isAddressedToBot определяет, адресовано ли сообщение боту.
func (b *Bot) isAddressedToBot(msg *tgbotapi.Message) bool {
	if msg == nil {
		return false
	}
	textLower := strings.ToLower(msg.Text)
	botMention := "@" + strings.ToLower(b.username)

	// 1) Простая проверка наличия @username в тексте
	if strings.Contains(textLower, botMention) {
		return true
	}

	// 2) Проверяем entities: mention или bot_command
	if msg.Entities != nil {
		for _, e := range *msg.Entities {
			if e.Type == "mention" {
				start := e.Offset
				end := e.Offset + e.Length
				if start >= 0 && end <= len(msg.Text) {
					ent := strings.ToLower(msg.Text[start:end])
					if ent == botMention {
						return true
					}
				}
			}
			if e.Type == "bot_command" {
				start := e.Offset
				end := e.Offset + e.Length
				if start >= 0 && end <= len(msg.Text) {
					ent := msg.Text[start:end]
					// если содержит @
					if strings.Contains(ent, "@") {
						if strings.HasSuffix(strings.ToLower(ent), botMention) {
							return true
						}
					} else {
						// команда без @ может быть личной, игнорируем для групп
					}
				}
			}
		}
	}

	return false
}
