package service

import "fmt"

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// PageInput 表示列表查询的页码和每页数量，零值使用默认分页参数。
type PageInput struct {
	Page     int
	PageSize int
}

// PageMeta 是列表响应中通用的分页元数据。
type PageMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func normalizePage(input PageInput, total int) (PageMeta, int, int, error) {
	page := input.Page
	pageSize := input.PageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if page < 1 || pageSize < 1 || pageSize > maxPageSize {
		return PageMeta{}, 0, 0, fmt.Errorf("%w: page must be positive and page_size must be between 1 and %d", ErrInvalidInput, maxPageSize)
	}

	start := 0
	if page-1 > total/pageSize {
		start = total
	} else {
		start = (page - 1) * pageSize
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return PageMeta{Page: page, PageSize: pageSize, Total: total}, start, end, nil
}
