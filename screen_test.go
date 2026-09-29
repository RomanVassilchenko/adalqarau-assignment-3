package screening

import (
	"strings"
	"testing"
	"time"
)

const header = "contract_id,customer_bin,contract_date,item_name,unit,quantity_units,amount_kzt\n"

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		csv  string
		want bool
	}{
		{"valid and normalized", header + "A,111,2025-01-01, Paper ,PACK,2,100\n", false},
		{"headers", "contract_id,bad\n", true},
		{"duplicate id", header + "A,111,2025-01-01,paper,pack,2,100\nA,111,2025-01-02,paper,pack,2,100\n", true},
		{"empty customer", header + "A, ,2025-01-01,paper,pack,2,100\n", true},
		{"empty id", header + ",111,2025-01-01,paper,pack,2,100\n", true},
		{"empty item", header + "A,111,2025-01-01, ,pack,2,100\n", true},
		{"empty unit", header + "A,111,2025-01-01,paper, ,2,100\n", true},
		{"field count", header + "A,111,2025-01-01,paper,pack,2\n", true},
		{"invalid date", header + "A,111,2025-02-30,paper,pack,2,100\n", true},
		{"invalid quantity", header + "A,111,2025-01-01,paper,pack,-1,100\n", true},
		{"invalid amount", header + "A,111,2025-01-01,paper,pack,2,1.5\n", true},
		{"quantity overflow", header + "A,111,2025-01-01,paper,pack,9223372036854775808,100\n", true},
		{"amount overflow", header + "A,111,2025-01-01,paper,pack,2,9223372036854775808\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := Read(strings.NewReader(tt.csv))
			if (err != nil) != tt.want {
				t.Fatalf("Read() error = %v, want error %v", err, tt.want)
			}
			if !tt.want && (len(cs) != 1 || cs[0].ItemName != "paper" || cs[0].Unit != "pack") {
				t.Fatalf("Read() = %#v", cs)
			}
		})
	}
}

func TestScreen(t *testing.T) {
	date := func(s string) time.Time { d, _ := time.Parse("2006-01-02", s); return d }
	cs := []Contract{
		{ID: "later", CustomerBIN: "111", Date: date("2025-01-31"), ItemName: "paper", Unit: "pack", Quantity: 2, Amount: 100},
		{ID: "earlier", CustomerBIN: "111", Date: date("2025-01-01"), ItemName: "paper", Unit: "pack", Quantity: 2, Amount: 100},
		{ID: "edge", CustomerBIN: "111", Date: date("2025-01-31"), ItemName: "paper", Unit: "pack", Quantity: 2, Amount: 100},
		{ID: "outside", CustomerBIN: "111", Date: date("2025-02-01"), ItemName: "paper", Unit: "pack", Quantity: 2, Amount: 100},
		{ID: "same-day", CustomerBIN: "111", Date: date("2025-01-01"), ItemName: "paper", Unit: "pack", Quantity: 2, Amount: 100},
	}
	tests := []struct {
		window int
		want   []Match
	}{
		{30, []Match{{"earlier", "later", "111", 30}, {"later", "edge", "111", 0}, {"later", "outside", "111", 1}, {"same-day", "later", "111", 30}, {"earlier", "edge", "111", 30}, {"earlier", "same-day", "111", 0}, {"edge", "outside", "111", 1}, {"same-day", "edge", "111", 30}}},
		{0, []Match{{"later", "edge", "111", 0}, {"earlier", "same-day", "111", 0}}},
		{-1, nil},
	}
	for _, tt := range tests {
		got := Screen(cs, tt.window)
		if len(got) != len(tt.want) {
			t.Fatalf("Screen(window=%d) = %#v, want %#v", tt.window, got, tt.want)
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("Screen(window=%d)[%d] = %#v, want %#v", tt.window, i, got[i], tt.want[i])
			}
		}
	}
	for _, field := range []string{"CustomerBIN", "ItemName", "Unit", "Quantity", "Amount"} {
		t.Run("mismatch "+field, func(t *testing.T) {
			a := cs[0]
			b := a
			b.ID = "different"
			switch field {
			case "CustomerBIN":
				b.CustomerBIN = "222"
			case "ItemName":
				b.ItemName = "pen"
			case "Unit":
				b.Unit = "box"
			case "Quantity":
				b.Quantity++
			case "Amount":
				b.Amount++
			}
			if got := Screen([]Contract{a, b}, 30); len(got) != 0 {
				t.Fatalf("mismatched %s produced %#v", field, got)
			}
		})
	}
}

func TestScreenFullDateRange(t *testing.T) {
	a := Contract{ID: "old", Date: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)}
	b := Contract{ID: "new", Date: time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)}
	if got := Screen([]Contract{b, a}, int(^uint(0)>>1)); len(got) != 1 || got[0].Base != "old" {
		t.Fatalf("full date range: %#v", got)
	}
}
