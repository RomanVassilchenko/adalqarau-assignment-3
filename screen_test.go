package screening

import (
	"strings"
	"testing"
)

const csvData = `contract_id,customer_bin,contract_date,item_name,unit,quantity_units,amount_kzt
A,111,2025-01-01, Paper ,pack,2,100
B,111,2025-01-31,paper,PACK,2,100
C,222,2025-01-31,paper,pack,2,100
D,111,2025-02-02,paper,pack,2,100
E,111,2025-01-15,pen,pack,2,100
`

func TestScreenBoundaryAndRules(t *testing.T) {
	cs, err := Read(strings.NewReader(csvData))
	if err != nil {
		t.Fatal(err)
	}
	m := Screen(cs, 30)
	if len(m) != 2 || m[0].Base != "A" || m[0].Repeated != "B" || m[0].Days != 30 || m[1].Base != "B" || m[1].Repeated != "D" {
		t.Fatalf("matches: %#v", m)
	}
}

func TestInvalidCSV(t *testing.T) {
	if _, err := Read(strings.NewReader("contract_id,bad\nA,1\n")); err == nil {
		t.Fatal("expected header error")
	}
	if _, err := Read(strings.NewReader(strings.Replace(csvData, "E,111", "A,111", 1))); err == nil {
		t.Fatal("expected duplicate error")
	}
}
