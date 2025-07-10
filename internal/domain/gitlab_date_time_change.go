package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DateTimeChange описывает {previous, current} с GitLab-овским форматом.
type DateTimeChange struct {
	Previous *time.Time `json:"previous"`
	Current  *time.Time `json:"current"`
}

// layout для строки вида "2006-01-02 15:04:05 MST"
const gitlabTimeLayout = "2006-01-02 15:04:05 MST"

// UnmarshalJSON разбирает JSON типа
//
//	{ "previous": null|"2025-07-09 19:11:00 UTC", "current": … }
//
// в *time.Time.
func (d *DateTimeChange) UnmarshalJSON(data []byte) error {
	// временная оболочка
	var aux struct {
		Previous *string `json:"previous"`
		Current  *string `json:"current"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parse := func(src *string) (*time.Time, error) {
		if src == nil {
			return nil, nil
		}
		// убираем лишние кавычки/пробелы
		s := strings.TrimSpace(*src)
		t, err := time.Parse(gitlabTimeLayout, s)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q: %w", s, err)
		}
		return &t, nil
	}

	var err error
	if d.Previous, err = parse(aux.Previous); err != nil {
		return err
	}
	if d.Current, err = parse(aux.Current); err != nil {
		return err
	}
	return nil
}
