package screening

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const RuleVersion = "repeated-need-v2"

type Contract struct {
	ID, CustomerBIN, ItemName, Unit string
	Date                            time.Time
	Quantity, Amount                int64
}

type Match struct {
	Base, Repeated, Customer string
	Days                     int
}

var headers = []string{"contract_id", "customer_bin", "contract_date", "item_name", "unit", "quantity_units", "amount_kzt"}

func Read(r io.Reader) ([]Contract, error) {
	c := csv.NewReader(r)
	h, err := c.Read()
	if err != nil {
		return nil, err
	}
	if len(h) != len(headers) {
		return nil, errors.New("invalid CSV headers")
	}
	for i := range h {
		if strings.TrimSpace(h[i]) != headers[i] {
			return nil, errors.New("invalid CSV headers")
		}
	}
	seen := map[string]bool{}
	var out []Contract
	for n := 2; ; n++ {
		rec, err := c.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", n, err)
		}
		if len(rec) != len(headers) {
			return nil, fmt.Errorf("row %d: expected 7 fields", n)
		}
		id := strings.TrimSpace(rec[0])
		if id == "" || seen[id] {
			return nil, fmt.Errorf("row %d: duplicate or empty contract_id", n)
		}
		seen[id] = true
		customer := strings.TrimSpace(rec[1])
		item := strings.ToLower(strings.TrimSpace(rec[3]))
		unit := strings.ToLower(strings.TrimSpace(rec[4]))
		if customer == "" || item == "" || unit == "" {
			return nil, fmt.Errorf("row %d: customer_bin, item_name and unit are required", n)
		}
		date, err := time.Parse("2006-01-02", strings.TrimSpace(rec[2]))
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid contract_date", n)
		}
		q, err := nonNegative(rec[5])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid quantity_units", n)
		}
		a, err := nonNegative(rec[6])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid amount_kzt", n)
		}
		out = append(out, Contract{ID: id, CustomerBIN: customer, Date: date, ItemName: item, Unit: unit, Quantity: q, Amount: a})
	}
	return out, nil
}

func nonNegative(s string) (int64, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil || n < 0 {
		return 0, errors.New("must be a nonnegative integer")
	}
	return n, nil
}

func Screen(cs []Contract, window int) []Match {
	if window < 0 {
		return nil
	}
	var out []Match
	// ponytail: O(n²) fits the small synthetic sets; group and sort by keys if measured input scale needs it.
	for i, base := range cs {
		for j := i + 1; j < len(cs); j++ {
			r := cs[j]
			if base.CustomerBIN != r.CustomerBIN || base.ItemName != r.ItemName || base.Unit != r.Unit || base.Quantity != r.Quantity || base.Amount != r.Amount {
				continue
			}
			// Unix seconds cover the full four-digit CSV year range without time.Duration overflow.
			first, second := base, r
			if second.Date.Before(first.Date) {
				first, second = second, first
			}
			days := int((second.Date.Unix() - first.Date.Unix()) / 86400)
			if days <= window {
				out = append(out, Match{first.ID, second.ID, first.CustomerBIN, days})
			}
		}
	}
	return out
}
