package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "docsbot_requests_total", Help: "Total number of processed requests"},
		[]string{"path", "method", "code"},
	)

	TelegramMessagesReceived = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "telegram_messages_received_total", Help: "Total number of telegram group messages received (text)"},
	)
	TelegramMessagesAddressed = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "telegram_messages_addressed_total", Help: "Total number of messages addressed to the bot"},
	)
	TelegramMessagesProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "telegram_messages_processed_total", Help: "Total number of messages successfully processed (reply sent)"},
	)
	TelegramMessagesErrors = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "telegram_messages_errors_total", Help: "Total number of errors while processing telegram messages"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(TelegramMessagesReceived)
	prometheus.MustRegister(TelegramMessagesAddressed)
	prometheus.MustRegister(TelegramMessagesProcessed)
	prometheus.MustRegister(TelegramMessagesErrors)
}
