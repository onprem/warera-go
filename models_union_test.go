package warera

import (
	"encoding/json"
	"testing"
)

func TestStringOrStringsJSONForms(t *testing.T) {
	tests := []struct {
		name  string
		value *StringOrStrings
		want  string
	}{
		{"single string", SingleString("wage"), `"wage"`},
		{"string list", StringList("wage", "trading"), `["wage","trading"]`},
		{"empty string list", StringList(), `[]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("marshal union: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("marshal union = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestStringOrStringsRejectsAmbiguousValue(t *testing.T) {
	single := "wage"
	_, err := json.Marshal(StringOrStrings{Single: &single, Multiple: []string{}})
	if err == nil {
		t.Fatal("expected an error when both variants are set")
	}
}

func TestStringOrStringsRejectsEmptyValue(t *testing.T) {
	if _, err := json.Marshal(StringOrStrings{}); err == nil {
		t.Fatal("expected an error when neither variant is set")
	}
}

func TestTransactionFilterMarshalsStringOrList(t *testing.T) {
	for _, tt := range []struct {
		value *StringOrStrings
		want  string
	}{
		{SingleString("wage"), `"wage"`},
		{StringList("wage", "trading"), `["wage","trading"]`},
	} {
		data, err := json.Marshal(TransactionFilter{TransactionType: tt.value})
		if err != nil {
			t.Fatalf("marshal transaction filter: %v", err)
		}
		var got map[string]json.RawMessage
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("decode marshaled filter: %v", err)
		}
		if string(got["transactionType"]) != tt.want {
			t.Fatalf("transactionType = %s, want %s", got["transactionType"], tt.want)
		}
	}
}
