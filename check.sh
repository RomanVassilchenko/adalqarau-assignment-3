#!/bin/sh
set -eu

test -z "$(gofmt -l .)"
go vet ./...
go test ./...
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
go build -o "$tmp/screen" ./cmd/screen
"$tmp/screen" -input testdata/contracts.csv -output "$tmp/out.csv" -window-days 30
printf '%s\n' 'base_contract_id,repeated_contract_id,customer_bin,days_between,rule_version' 'A,B,111,30,repeated-need-v2' 'B,D,111,2,repeated-need-v2' > "$tmp/want.csv"
cmp "$tmp/want.csv" "$tmp/out.csv"
"$tmp/screen" -input testdata/contracts.csv -output "$tmp/out2.csv" -window-days 30
cmp "$tmp/out.csv" "$tmp/out2.csv"
cp testdata/contracts.csv "$tmp/same.csv"
if "$tmp/screen" -input "$tmp/same.csv" -output "$tmp/same.csv" 2>/dev/null; then
	echo 'CLI accepted the same input and output file' >&2
	exit 1
fi
cmp testdata/contracts.csv "$tmp/same.csv"
printf '%s\n' 'contract_id,customer_bin,contract_date,item_name,unit,quantity_units,amount_kzt' 'bad,111,not-a-date,paper,pack,2,100' > "$tmp/invalid.csv"
printf 'keep this output\n' > "$tmp/kept.csv"
if "$tmp/screen" -input "$tmp/invalid.csv" -output "$tmp/kept.csv" 2>/dev/null; then
	echo 'CLI accepted an invalid CSV' >&2
	exit 1
fi
printf 'keep this output\n' > "$tmp/want-kept.csv"
cmp "$tmp/want-kept.csv" "$tmp/kept.csv"
