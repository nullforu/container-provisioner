package config

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseCPUMilli(s string) (int64, error) {
	v := strings.TrimSpace(s)
	if v == "" {
		return 0, fmt.Errorf("empty value")
	}

	if before, ok := strings.CutSuffix(v, "m"); ok {
		n, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid milli cpu")
		}

		if n <= 0 {
			return 0, fmt.Errorf("cpu must be positive")
		}

		return n, nil
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid cpu")
	}

	if f <= 0 {
		return 0, fmt.Errorf("cpu must be positive")
	}

	return int64(math.Round(f * 1000)), nil
}

func ParseBytes(s string) (int64, error) {
	v := strings.TrimSpace(strings.ToUpper(s))
	if v == "" {
		return 0, fmt.Errorf("empty value")
	}

	units := map[string]int64{
		"KI": 1024,
		"MI": 1024 * 1024,
		"GI": 1024 * 1024 * 1024,
		"TI": 1024 * 1024 * 1024 * 1024,
		"K":  1000,
		"M":  1000 * 1000,
		"G":  1000 * 1000 * 1000,
	}

	for unit, scale := range units {
		if before, ok := strings.CutSuffix(v, unit); ok {
			n, err := strconv.ParseFloat(before, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid memory")
			}

			if n <= 0 {
				return 0, fmt.Errorf("memory must be positive")
			}

			return int64(math.Round(n * float64(scale))), nil
		}
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory")
	}

	if n <= 0 {
		return 0, fmt.Errorf("memory must be positive")
	}

	return n, nil
}
