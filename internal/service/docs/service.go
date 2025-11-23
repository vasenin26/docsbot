package docs

import "context"

// DocumentService описывает контракт для обработки сообщений документа/чата.
//
// Implementations may perform I/O or long-running operations; the context
// parameter allows cancellation and deadline propagation.
type DocumentService interface {
	// ProcessMessage обрабатывает сообщение для заданного chatID.
	// Возвращает строковый результат и ошибку.
	ProcessMessage(ctx context.Context, chatID string, message string) (string, error)
}

// StubDocumentService — простая заглушка реализации DocumentService.
// Всегда возвращает "ок" без ошибок.
type StubDocumentService struct{}

// NewStubDocumentService возвращает новую заглушечную реализацию DocumentService.
func NewStubDocumentService() DocumentService {
	return StubDocumentService{}
}

// ProcessMessage реализует DocumentService и всегда возвращает "ок" и nil.
func (StubDocumentService) ProcessMessage(ctx context.Context, chatID string, message string) (string, error) {
	return "ок", nil
}
