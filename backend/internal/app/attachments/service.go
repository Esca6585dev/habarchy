// Package attachments stores small binary blobs (user avatars and chat
// images) in Postgres and serves them back. Images only; capped size.
package attachments

import (
	"bytes"
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres"
	"github.com/Esca6585dev/habarchy/backend/internal/adapters/postgres/sqlcgen"
	"github.com/Esca6585dev/habarchy/backend/internal/domain"
	"github.com/Esca6585dev/habarchy/backend/pkg/ids"
)

// MaxImageBytes caps a stored image.
const MaxImageBytes = 5 << 20 // 5 MiB

var allowed = map[string]bool{"image/png": true, "image/jpeg": true, "image/gif": true, "image/webp": true}

// Service stores and reads attachments.
type Service struct {
	db *postgres.DB
}

// New creates the service.
func New(db *postgres.DB) *Service { return &Service{db: db} }

// Meta is an attachment without its bytes.
type Meta struct {
	ID          uuid.UUID `json:"id"`
	ContentType string    `json:"content_type"`
	SizeBytes   int32     `json:"size_bytes"`
	Filename    string    `json:"filename"`
}

// SaveImage validates (by magic bytes) and stores an image, returning its id.
func (s *Service) SaveImage(ctx context.Context, uploadedBy uuid.UUID, filename string, data []byte) (*Meta, error) {
	if len(data) == 0 {
		return nil, domain.ErrValidation.WithMessage("empty file")
	}
	if len(data) > MaxImageBytes {
		return nil, domain.ErrValidation.WithMessage("image is larger than 5 MB")
	}
	ct := http.DetectContentType(data)
	if isWebP(data) {
		ct = "image/webp"
	}
	if !allowed[ct] {
		return nil, domain.ErrValidation.WithMessage("file is not a supported image (png, jpeg, gif or webp)")
	}
	id := ids.New()
	row, err := s.db.Queries.CreateAttachment(ctx, sqlcgen.CreateAttachmentParams{
		ID: id, UploadedBy: &uploadedBy, ContentType: ct, SizeBytes: int32(len(data)), Filename: filename, Data: data, //nolint:gosec // capped above
	})
	if err != nil {
		return nil, err
	}
	return &Meta{ID: row.ID, ContentType: row.ContentType, SizeBytes: row.SizeBytes, Filename: row.Filename}, nil
}

// Attachment carries bytes for serving.
type Attachment struct {
	ContentType string
	Data        []byte
}

// Get returns an attachment's bytes.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Attachment, error) {
	a, err := s.db.Queries.GetAttachment(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound.WithMessage("attachment not found")
	}
	return &Attachment{ContentType: a.ContentType, Data: a.Data}, nil
}

func isWebP(b []byte) bool {
	return len(b) >= 12 && bytes.Equal(b[0:4], []byte("RIFF")) && bytes.Equal(b[8:12], []byte("WEBP"))
}
