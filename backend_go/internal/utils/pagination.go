package utils

import "github.com/gofiber/fiber/v3"

const (
	// DefaultPage adalah halaman default jika tidak dispesifikasikan.
	DefaultPage = 1
	// DefaultLimit adalah jumlah item per halaman default.
	DefaultLimit = 10
	// MaxLimit adalah batas maksimum item per halaman untuk mencegah query berlebihan.
	MaxLimit = 100
)

// ParsePagination membaca parameter pagination (page & limit) dari query string
// dan menerapkan nilai default serta validasi batas.
func ParsePagination(ctx fiber.Ctx) (page, limit int) {
	page = fiber.Query[int](ctx, "page", DefaultPage)
	limit = fiber.Query[int](ctx, "limit", DefaultLimit)

	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return
}
