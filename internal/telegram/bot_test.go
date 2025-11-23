package telegram

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Test isAddressedToBot with simple mention
func TestIsAddressedToBot_Mention(t *testing.T) {
	b := &Bot{username: "MyBot"}
	msg := &tgbotapi.Message{Text: "Hello @MyBot"}
	if !b.isAddressedToBot(msg) {
		t.Fatalf("expected message to be addressed to bot")
	}
}

func TestIsAddressedToBot_CaseInsensitive(t *testing.T) {
	b := &Bot{username: "MyBot"}
	msg := &tgbotapi.Message{Text: "hello @mybot"}
	if !b.isAddressedToBot(msg) {
		t.Fatalf("expected case-insensitive mention to be addressed")
	}
}

func TestIsAddressedToBot_EntitiesMention(t *testing.T) {
	b := &Bot{username: "MyBot"}
	text := "Hello @MyBot"
	entities := []tgbotapi.MessageEntity{{Type: "mention", Offset: 6, Length: 6}}
	msg := &tgbotapi.Message{Text: text, Entities: &entities}
	if !b.isAddressedToBot(msg) {
		t.Fatalf("expected entity mention to be addressed")
	}
}

func TestIsAddressedToBot_BotCommandWithUsername(t *testing.T) {
	b := &Bot{username: "MyBot"}
	text := "/start@MyBot"
	entities := []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: len(text)}}
	msg := &tgbotapi.Message{Text: text, Entities: &entities}
	if !b.isAddressedToBot(msg) {
		t.Fatalf("expected bot_command with username to be addressed")
	}
}

func TestIsAddressedToBot_NotAddressed(t *testing.T) {
	b := &Bot{username: "MyBot"}
	msg := &tgbotapi.Message{Text: "hello everyone"}
	if b.isAddressedToBot(msg) {
		t.Fatalf("expected message NOT to be addressed to bot")
	}
}

// Stub client to capture Send calls
type stubClient struct {
	sentMessages []tgbotapi.Chattable
}

func (s *stubClient) GetMe() (tgbotapi.User, error) { return tgbotapi.User{UserName: "MyBot"}, nil }
func (s *stubClient) GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel { return nil }
func (s *stubClient) Request(tgbotapi.Chattable) (tgbotapi.APIResponse, error) { return tgbotapi.APIResponse{}, nil }
func (s *stubClient) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	s.sentMessages = append(s.sentMessages, c)
	return tgbotapi.Message{}, nil
}

func TestHandleUpdate_WritesStdoutAndReplies(t *testing.T) {
	// prepare bot with mocked client
	stub := &stubClient{}
	b := &Bot{client: stub, username: "MyBot"}

	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	text := "Hello @MyBot"
	msg := &tgbotapi.Message{Text: text, Chat: &tgbotapi.Chat{ID: 12345, Type: "group"}, MessageID: 10, From: &tgbotapi.User{ID: 111}}
	entities := []tgbotapi.MessageEntity{{Type: "mention", Offset: 6, Length: 6}}
	msg.Entities = &entities

	update := tgbotapi.Update{Message: msg}
	if err := b.handleUpdate(update); err != nil {
		t.Fatalf("handleUpdate returned error: %v", err)
	}

	// read stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	out := buf.String()
	if out == "" {
		t.Fatalf("expected stdout to contain message, got empty")
	}

	// check that Send was called
	if len(stub.sentMessages) == 0 {
		t.Fatalf("expected Send to be called")
	}

	// verify that the sent message text is "Ок"
	if msgCfg, ok := stub.sentMessages[0].(tgbotapi.MessageConfig); ok {
		if msgCfg.Text != "Ок" {
			t.Fatalf("expected reply text 'Ок', got '%s'", msgCfg.Text)
		}
	} else {
		t.Fatalf("sent message is not MessageConfig")
	}
}
