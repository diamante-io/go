package main

import (
	"context"
	"fmt"

	"github.com/diamcircle/go/ingest/ledgerbackend"
)

// This little app helped testing CaptiveDiamcircleCore.runFromParams on a living
// Diamcircle-Core. Adding it to the repo because it can be useful in a future if
// Diamcircle-Core behaviour changes again.
// To make it work, run standalone network (RUN_STANDALONE=false to allow outside
// connections) and update paths below.
func main() {
	// check(1) // err expected, cannot stream in captive core
	checkLedgers := []uint32{2, 3, 62, 63, 64, 65, 126, 127, 128}
	for _, ledger := range checkLedgers {
		ok := check(ledger)
		if !ok {
			panic("test failed error")
		}
	}
}

func check(ledger uint32) bool {
	c, err := ledgerbackend.NewCaptive(
		ledgerbackend.CaptiveCoreConfig{
			BinaryPath:         "diamcircle-core",
			NetworkPassphrase:  "Standalone Network ; February 2017",
			HistoryArchiveURLs: []string{"http://localhost:1570"},
		},
	)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctx := context.Background()
	err = c.PrepareRange(ctx, ledgerbackend.UnboundedRange(ledger))
	if err != nil {
		fmt.Println(err)
		return false
	}

	meta, err := c.GetLedger(ctx, ledger)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if meta.LedgerSequence() != ledger {
		fmt.Println("wrong ledger", meta.LedgerSequence())
		return false
	}

	fmt.Println(ledger, "ok")
	return true
}
