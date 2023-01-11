package main

import (
	"testing"
	"time"
)

func TestFormattedDate(t *testing.T) {

	tests := []struct {
		name      string
		inputTime time.Time
		required  string
	}{{
		name:      "UTC",
		inputTime: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		required:  "01 Jan 2020 10:00:00",
	},
		{
			name:      "Empty_Date",
			inputTime: time.Time{},
			required:  "",
		},
		{
			name:      "CET",
			inputTime: time.Date(2020, 1, 1, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			required:  "01 Jan 2020 09:00:00",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			formatDate := formattedDate(test.inputTime)

			if formatDate != test.required {
				t.Errorf("want %q got %q", test.required, formatDate)
			}
		})
	}
}
