package warera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Time is a timestamp from the API, such as a record's CreatedAt. It embeds
// [time.Time], so its methods are available directly. The zero value means the
// API sent no time (an empty string, null, or a missing field); check it with
// IsZero.
type Time struct {
	time.Time
}

// UnmarshalJSON decodes an RFC 3339 timestamp. Empty strings and null decode
// to the zero Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	if isEmptyJSONString(data) {
		t.Time = time.Time{}
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("warera: decode time: %w", err)
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return fmt.Errorf("warera: decode time: %w", err)
	}
	t.Time = parsed
	return nil
}

// MarshalJSON encodes the time in RFC 3339 format, or null for the zero Time.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// String returns the time in RFC 3339 format, or an empty string for the zero
// Time.
func (t Time) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// Date is a calendar day from the API, sent as YYYY-MM-DD, such as a day in a
// price history. It embeds a [time.Time] at midnight UTC. The zero value
// means the API sent no date.
type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

// UnmarshalJSON decodes a YYYY-MM-DD date. Empty strings and null decode to the
// zero Date.
func (d *Date) UnmarshalJSON(data []byte) error {
	if isEmptyJSONString(data) {
		d.Time = time.Time{}
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("warera: decode date: %w", err)
	}
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return fmt.Errorf("warera: decode date: %w", err)
	}
	d.Time = parsed
	return nil
}

// MarshalJSON encodes the date as YYYY-MM-DD, or null for the zero Date.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// String returns the date as YYYY-MM-DD.
func (d Date) String() string {
	return d.Format(dateLayout)
}

func isEmptyJSONString(data []byte) bool {
	data = bytes.TrimSpace(data)
	return bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`))
}
