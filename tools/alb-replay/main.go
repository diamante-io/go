// alb-replay replays the successful GET requests found in an AWS ALB log file
// see https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	albLogEntryCount               = 29
	albTargetStatusCodeRecordIndex = 9
	albRequestIndex                = 12
)

type NumberedURL struct {
	Number int
	URL    string
}

func isSuccesfulStatusCode(statusCode int) bool {
	// consider all 2XX HTTP errors a success
	return statusCode/100 == 2
}

type ALBLogEntryReader csv.Reader

func newALBLogEntryReader(input io.Reader) *ALBLogEntryReader {
	reader := csv.NewReader(input)
	reader.Comma = ' '
	reader.FieldsPerRecord = albLogEntryCount
	reader.ReuseRecord = true
	return (*ALBLogEntryReader)(reader)
}

func (r *ALBLogEntryReader) GetRequestURI() (string, error) {
	records, err := (*csv.Reader)(r).Read()
	if err != nil {
		return "", err
	}

	statusCodeStr := records[albTargetStatusCodeRecordIndex]
	// discard requests with unknown status code
	if statusCodeStr == "-" {
		return "", nil
	}
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return "", fmt.Errorf("error parsing target status code %q: %v", statusCodeStr, err)
	}

	// discard unsuccesful requests
	if !isSuccesfulStatusCode(statusCode) {
		return "", nil
	}

	reqStr := records[albRequestIndex]
	reqFields := strings.Split(reqStr, " ")
	if len(reqFields) != 3 {
		return "", fmt.Errorf("error parsing request %q: 3 fields exepcted, found %d", reqStr, len(reqFields))
	}
	method := reqFields[0]

	// discard non-get requests
	if method != http.MethodGet {
		return "", nil
	}

	urlStr := reqFields[1]
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("error parsing url %q: %v", urlStr, err)
	}

	return parsed.RequestURI(), nil
}

func parseURLs(ctx context.Context, startFrom int, baseURL string, logReader *ALBLogEntryReader, urlChan chan NumberedURL) {
	counter := 0
	for {
		uri, err := logReader.GetRequestURI()
		if err != nil {
			if err == io.EOF {
				// we are done
				return
			}
			log.Fatal(err.Error())
		}
		if uri == "" {
			// no usable URL found in the current log line
			continue
		}
		counter++
		if counter < startFrom {
			// we haven't yet reached the expected start point
			continue
		}
		url := NumberedURL{
			Number: counter,
			URL:    baseURL + uri,
		}
		select {
		case <-ctx.Done():
			return
		case urlChan <- url:
		}
	}
}

func queryURLs(ctx context.Context, timeout time.Duration, urlChan chan NumberedURL) {
	client := http.Client{
		Timeout: timeout,
	}
	for numURL := range urlChan {
		if ctx.Err() != nil {
			return
		}

		req, err := http.NewRequest(http.MethodGet, numURL.URL, nil)
		if err != nil {
			log.Printf("(%d) unexpected error creating request: %v", numURL.Number, err)
			continue
		}
		req = req.WithContext(ctx)
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			// we don't want to print cancel errors due to a signal interrupt
			if errors.Unwrap(err) != context.Canceled {
				log.Printf("(%d) unexpected request error: %v %q", numURL.Number, errors.Unwrap(err), numURL.URL)
			}
			continue
		}
		resp.Body.Close()
		if !isSuccesfulStatusCode(resp.StatusCode) {
			log.Printf("(%d) unexpected status code: %d %q", numURL.Number, resp.StatusCode, numURL.URL)
			continue
		}
		log.Printf("(%d) %s %s", numURL.Number, time.Since(start), numURL.URL)
	}
}

func main() {
	workers := flag.Int("workers", 1, "How many parallel workers to use")
	startFromURLNum := flag.Int("start-from", 1, "What URL number to start from")
	timeout := flag.Duration("timeout", time.Second*5, "HTTP request timeout")
	flag.Parse()
	if *workers < 1 {
		log.Fatal("--workers parameter must be > 0")
	}
	if *startFromURLNum < 1 {
		log.Fatal("--start-from must be > 0")
	}
	if flag.NArg() != 2 {
		log.Fatalf("usage: %s <aws_log_file> <target_host_base_url>", os.Args[0])
	}

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("error opening file %q: %v", os.Args[1], err)
	}
	baseURL := flag.Args()[1]
	logReader := newALBLogEntryReader(file)
	urlChan := make(chan NumberedURL, *workers)
	var wg sync.WaitGroup

	// setup interrupt cleanup code
	ctx, stopped := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopped()

	// spawn url consumers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			queryURLs(ctx, *timeout, urlChan)
			wg.Done()
		}()
	}

	parseURLs(ctx, *startFromURLNum, baseURL, logReader, urlChan)
	// signal the consumers there won't be more urls
	close(urlChan)
	// wait for to consumers to be done
	wg.Wait()
}
