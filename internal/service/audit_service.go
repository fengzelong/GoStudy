package service

import (
	"context"

	"GoStudy/internal/audit"
)

// AuditPage 是审计记录列表及其分页信息。
type AuditPage struct {
	PageMeta
	Items []audit.Entry `json:"items"`
}

// AuditService 负责审计记录的读取和分页，不承担记录行为。
type AuditService struct {
	entries audit.Reader
}

func NewAuditService(entries audit.Reader) *AuditService {
	return &AuditService{entries: entries}
}

// List 返回分页后的审计记录。
func (s *AuditService) List(ctx context.Context, input PageInput) (AuditPage, error) {
	entries, err := s.entries.List(ctx)
	if err != nil {
		return AuditPage{}, err
	}
	meta, start, end, err := normalizePage(input, len(entries))
	if err != nil {
		return AuditPage{}, err
	}
	return AuditPage{PageMeta: meta, Items: entries[start:end]}, nil
}
