package utils

import (
	"fmt"
	"strings"
)

func FormatRupiah(amount float64) string {
	res := fmt.Sprintf("Rp %.0f", amount)
	res = strings.ReplaceAll(res, " ", "")

	parts := strings.Split(fmt.Sprintf("%.0f", amount), "")
	var result []string
	count := 0
	for i := len(parts) - 1; i >= 0; i-- {
		result = append([]string{parts[i]}, result...)
		count++
		if count%3 == 0 && i != 0 {
			result = append([]string{"."}, result...)
		}
	}

	return "Rp " + strings.Join(result, "")
}
