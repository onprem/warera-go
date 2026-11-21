package warera

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeAndDateJSON(t *testing.T) {
	var record struct {
		At     Time `json:"at"`
		Empty  Time `json:"empty"`
		Null   Time `json:"null"`
		Day    Date `json:"day"`
		NoDay  Date `json:"noDay"`
		Absent Time `json:"absent"`
	}
	input := `{"at":"2026-11-21T16:37:02.123Z","empty":"","null":null,"day":"2026-11-20","noDay":""}`
	if err := json.Unmarshal([]byte(input), &record); err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 11, 21, 16, 37, 2, 123e6, time.UTC); !record.At.Equal(want) {
		t.Errorf("At = %v, want %v", record.At, want)
	}
	if !record.Empty.IsZero() || !record.Null.IsZero() || !record.NoDay.IsZero() || !record.Absent.IsZero() {
		t.Errorf("empty, null, and absent values = %v %v %v %v, want zero", record.Empty, record.Null, record.NoDay, record.Absent)
	}
	if record.At.String() != "2026-11-21T16:37:02Z" || record.Empty.String() != "" {
		t.Errorf("String = %q, %q; want RFC 3339 and empty for zero", record.At, record.Empty)
	}
	if record.Day.String() != "2026-11-20" {
		t.Errorf("Day = %s, want 2026-11-20", record.Day)
	}

	out, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"at":"2026-11-21T16:37:02.123Z","empty":null,"null":null,"day":"2026-11-20","noDay":null,"absent":null}`
	if string(out) != want {
		t.Errorf("Marshal = %s, want %s", out, want)
	}

	for _, bad := range []string{`{"at":"yesterday"}`, `{"at":5}`, `{"day":"2026-11-21T00:00:00Z"}`} {
		if err := json.Unmarshal([]byte(bad), &record); err == nil {
			t.Errorf("Unmarshal(%s) succeeded, want an error", bad)
		}
	}
}
