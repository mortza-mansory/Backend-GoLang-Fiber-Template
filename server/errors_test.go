package server

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

func TestMapErrorAppError(t *testing.T) {
	err := shared.NewError(shared.CodeNotFound, "user not found")

	he := MapError(err, false)
	if he == nil {
		t.Fatal("expected an HTTPError")
	}
	if he.Status != fiber.StatusNotFound {
		t.Fatalf("expected status 404, got %d", he.Status)
	}
	if he.Code != shared.CodeNotFound {
		t.Fatalf("expected code NOT_FOUND, got %s", he.Code)
	}
}

func TestMapErrorUnknownInternalHiddenInProduction(t *testing.T) {
	he := MapError(errors.New("database connection refused"), true)
	if he == nil {
		t.Fatal("expected an HTTPError")
	}
	if he.Status != fiber.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", he.Status)
	}
	// The internal message must not leak.
	if he.Message == "database connection refused" {
		t.Fatal("internal error message leaked in production")
	}
}

func TestMapErrorNil(t *testing.T) {
	if MapError(nil, false) != nil {
		t.Fatal("expected nil for nil error")
	}
}

func TestStatusForCode(t *testing.T) {
	cases := map[shared.ErrorCode]int{
		shared.CodeBadRequest:    400,
		shared.CodeUnauthorized:  401,
		shared.CodeForbidden:     403,
		shared.CodeNotFound:      404,
		shared.CodeConflict:      409,
		shared.CodeUnprocessable: 422,
		shared.CodeInternal:      500,
	}
	for code, want := range cases {
		if got := StatusForCode(code); got != want {
			t.Errorf("StatusForCode(%s) = %d, want %d", code, got, want)
		}
	}
}
