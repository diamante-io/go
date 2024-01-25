package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	aurora "github.com/diamcircle/go/services/aurora/internal"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/log"
)

var recordMetricsCmd = &cobra.Command{
	Use:   "record-metrics",
	Short: "records `/metrics` on admin port for debuging purposes",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := aurora.ApplyFlags(config, flags, aurora.ApplyOptions{}); err != nil {
			return err
		}

		const (
			timeFormat            = "2006-01-02-15-04-05"
			scrapeIntervalSeconds = 15
			scrapesCount          = (60 / scrapeIntervalSeconds) * 10 // remember about rounding if change is required
		)

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		outputFileName := fmt.Sprintf("./metrics-%s.zip", time.Now().Format(timeFormat))
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return err
		}

		w := zip.NewWriter(outputFile)
		defer w.Close()

		for i := 1; i <= scrapesCount; i++ {
			log.Infof(
				"Getting metrics %d/%d... ETA: %s",
				i,
				scrapesCount,
				time.Duration(time.Duration(scrapeIntervalSeconds*(scrapesCount-i))*time.Second),
			)

			metricsResponse, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", config.AdminPort))
			if err != nil {
				return errors.Wrap(err, "Error fetching metrics. Is admin server running?")
			}

			if metricsResponse.StatusCode != http.StatusOK {
				return errors.Errorf("Invalid status code: %d. Is admin server running?", metricsResponse.StatusCode)
			}

			metricsFile, err := w.Create(time.Now().Format(timeFormat))
			if err != nil {
				return err
			}

			if _, err = io.Copy(metricsFile, metricsResponse.Body); err != nil {
				return errors.Wrap(err, "Error reading response body. Is admin server running?")
			}

			// Flush to keep memory usage log and save at least some records in case of errors later.
			err = w.Flush()
			if err != nil {
				return err
			}

			if i < scrapesCount {
				time.Sleep(scrapeIntervalSeconds * time.Second)
			}
		}

		log.Infof("Metrics recorded to %s!", outputFileName)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(recordMetricsCmd)
}
