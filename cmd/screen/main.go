package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/RomanVassilchenko/adalqarau-assignment-3"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "screen" {
		fmt.Fprintln(os.Stderr, "usage: screen -input <file> -output <file> -window-days 30")
		os.Exit(2)
	}
	f := flag.NewFlagSet("screen", flag.ExitOnError)
	in := f.String("input", "", "input CSV")
	dest := f.String("output", "", "output CSV")
	window := f.Int("window-days", 30, "inclusive window in days")
	_ = f.Parse(os.Args[2:])
	if *in == "" || *dest == "" || *window < 0 {
		fmt.Fprintln(os.Stderr, "input, output and nonnegative window-days are required")
		os.Exit(2)
	}
	fi, err := os.Open(*in)
	if err != nil {
		fail(err)
	}
	defer fi.Close()
	contracts, err := screening.Read(fi)
	if err != nil {
		fail(err)
	}
	fo, err := os.Create(*dest)
	if err != nil {
		fail(err)
	}
	defer fo.Close()
	w := csv.NewWriter(fo)
	_ = w.Write([]string{"base_contract_id", "repeated_contract_id", "customer_bin", "days_between", "rule_version"})
	for _, m := range screening.Screen(contracts, *window) {
		_ = w.Write([]string{m.Base, m.Repeated, m.Customer, fmt.Sprint(m.Days), screening.RuleVersion})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		fail(err)
	}
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
