package response

// PaginationMeta berisi metadata pagination untuk respons API yang dipaginasi.
type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}

// NewPaginationMeta membuat PaginationMeta dengan kalkulasi total halaman otomatis.
func NewPaginationMeta(page, limit int, total int64) PaginationMeta {
	pages := (total + int64(limit) - 1) / int64(limit)
	return PaginationMeta{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: pages,
	}
}

// Paginated membuat respons sukses dengan data dan metadata pagination.
// Contoh: response.Paginated("students", students, page, limit, total)
func Paginated(key string, data interface{}, page, limit int, total int64) WebResponse {
	return Success(map[string]interface{}{
		key:          data,
		"pagination": NewPaginationMeta(page, limit, total),
	})
}
