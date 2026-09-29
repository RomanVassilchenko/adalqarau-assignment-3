package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/RomanVassilchenko/adalqarau-assignment-3"
)

func main() {
	f := flag.NewFlagSet("screen", flag.ExitOnError)
	in := f.String("input", "", "input CSV")
	dest := f.String("output", "", "output CSV")
	window := f.Int("window-days", 30, "inclusive window in days")
	if err := f.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	if f.NArg() != 0 || *in == "" || *dest == "" || *window < 0 {
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
	if existing, err := os.Stat(*dest); err == nil {
		inputInfo, err := fi.Stat()
		if err != nil {
			fail(err)
		}
		if os.SameFile(inputInfo, existing) {
			fail(fmt.Errorf("input and output refer to the same file"))
		}
	}
	fo, err := os.Create(*dest)
	if err != nil {
		fail(err)
	}
	w := csv.NewWriter(fo)
	_ = w.Write([]string{"base_contract_id", "repeated_contract_id", "customer_bin", "days_between", "rule_version"})
	for _, m := range screening.Screen(contracts, *window) {
		_ = w.Write([]string{m.Base, m.Repeated, m.Customer, fmt.Sprint(m.Days), screening.RuleVersion})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		_ = fo.Close()
		fail(err)
	}
	if err := fo.Close(); err != nil {
		fail(err)
	}
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
