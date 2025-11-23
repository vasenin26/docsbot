package docs

import (
	"context"
	"testing"
)

func TestStubDocumentService_Success(t *testing.T) {
	svc := NewStubDocumentService()
	res, err := svc.ProcessMessage(context.Background(), "chat123", "hello world")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res != "ок" {
		t.Fatalf("expected 'ок', got %q", res)
	}
}

func TestStubDocumentService_EmptyInputs(t *testing.T) {
	svc := NewStubDocumentService()
	res, err := svc.ProcessMessage(context.Background(), "", "")
	if err != nil {
		t.Fatalf("expected nil error for empty inputs, got %v", err)
	}
	if res != "ок" {
		t.Fatalf("expected 'ок' for empty inputs, got %q", res)
	}
}
