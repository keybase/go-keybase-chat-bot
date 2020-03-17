package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func runningtime(s string) (string, time.Time) {
	fmt.Println("Start: ", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println("End: ", s, "took", endTime.Sub(startTime))
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}
	// send
	c := 200
	to := "TESTACCOUNT" ///
	txids := make([]string, 0, c)
	fmt.Printf("....start sending....\n")
	payingLabel, payingTime := runningtime("timepaying")
	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res, err := kbc.SendWalletTx(to, "3", nil, nil, nil, nil)
			if err != nil {
				fmt.Printf("--ERROR: i=%d, err=%+v\n", i, err.Error())
			} else {
				fmt.Printf("SUCCESS: %d, %+v\n", i, res.TxID)
				txids = append(txids, string(res.TxID))
			}
		}(i)
	}
	wg.Wait()
	track(payingLabel, payingTime)
	fmt.Printf("....finished sending.... len=%d\n\n", len(txids))

	// cancel
	successCount := 0
	fmt.Printf("....start canceling....\n")
	cancelingLabel, cancelingTime := runningtime("timecanceling")
	for i, txID := range txids {
		wg.Add(1)
		go func(i int, txID string) {
			defer wg.Done()
			res, err := kbc.CancelWalletTx(txID)
			if err != nil {
				fmt.Printf("--ERROR: i=%d, err=%+v\n", i, err.Error())
			} else {
				fmt.Printf("SUCCESS: %d, %+v, %+v\n", i, txID, res.ClaimStellarID)
				successCount++
			}
		}(i, txID)
	}
	wg.Wait()
	track(cancelingLabel, cancelingTime)
	fmt.Printf("....finished sending.... successCount=%d\n\n", successCount)
}
