package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	hlog "github.com/diamcircle/go/support/log"
)

// PanicIfError is an utility function that panics if err != nil
func PanicIfError(e error) {
	if e != nil {
		panic(e)
	}
}

// WriteJSONToFile atomically writes a json []byte dump to <filename>
// It ensures atomicity by first creating a tmp file (filename.tmp), writing
// the contents to it, then renaming it to the originally specified filename.
func WriteJSONToFile(jsonBytes []byte, filename string) (numBytes int, err error) {
	tmp := fmt.Sprintf("%s.tmp", filename)
	f, err := os.Create(tmp)
	PanicIfError(err)
	defer f.Close()

	numBytes, err = f.Write(jsonBytes)
	if err != nil {
		return
	}

	err = f.Sync()
	if err != nil {
		return
	}

	err = os.Rename(tmp, filename)
	return
}

// SliceDiff returns the elements in `a` that aren't in `b`.
func SliceDiff(a, b []string) (diff []string) {
	bmap := map[string]bool{}
	for _, x := range b {
		bmap[x] = true
	}
	for _, x := range a {
		if _, ok := bmap[x]; !ok {
			diff = append(diff, x)
		}
	}
	return
}

// GetAssetString returns a string representation of an asset
func GetAssetString(assetType string, code string, issuer string) string {
	if assetType == "native" {
		return "native"
	}
	return fmt.Sprintf("%s:%s", code, issuer)
}

// TimeToTimestamp converts a time.Time into a Unix epoch
func TimeToUnixEpoch(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// TimeToRFC3339 converts a time.Time to a string in RFC3339 format
func TimeToRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// CalcSpread calculates the spread stats for the given bidMax and askMin orderbook values
func CalcSpread(bidMax float64, askMin float64) (spread float64, midPoint float64) {
	if askMin == 0 || bidMax == 0 {
		return 0, 0
	}
	spread = (askMin - bidMax) / askMin
	midPoint = bidMax + spread/2.0
	return
}

// Retry retries running a function that returns an error numRetries times, multiplying
// the sleep time by a factor of 2 each time it retries.
func Retry(numRetries int, delay time.Duration, logger *hlog.Entry, f func() error) error {
	if err := f(); err != nil {
		if numRetries--; numRetries > 0 {
			jitter := time.Duration(rand.Int63n(int64(delay)))
			delay = delay + jitter/2

			logger.Infof("Backing off for %.3f seconds before retrying", delay.Seconds())

			time.Sleep(delay)
			return Retry(numRetries, 2*delay, logger, f)
		}
		return err
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
