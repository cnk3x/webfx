package strs

import (
	"fmt"
)

func HumanBytes[T Num](i T) string {
	n := float64(i)
	const unit = float64(1024)
	if n < unit {
		return fmt.Sprintf("%.0f B", n)
	}
	div, exp := unit, 0
	for n := n / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", n/div, "KMGTPE"[exp])
}
