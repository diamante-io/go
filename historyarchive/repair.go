// Copyright 2016 Diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Repair repairs a destination archive based on a source archive, it assumes that the source and destination have the
// same checkpoint ledger frequency
func Repair(src *Archive, dst *Archive, opts *CommandOptions) error {
	state, e := dst.GetRootHAS()
	if e != nil {
		return e
	}
	opts.Range = opts.Range.clamp(state.Range(), src.checkpointManager)

	log.Printf("Starting scan for repair")
	var errs uint32
	errs += noteError(dst.ScanCheckpoints(opts))

	log.Printf("Examining checkpoint files for gaps")
	missingCheckpointFiles := dst.CheckCheckpointFilesMissing(opts)

	repairedHistory := false
	for cat, missing := range missingCheckpointFiles {
		for _, chk := range missing {
			pth := CategoryCheckpointPath(cat, chk)
			exists, err := src.backend.Exists(pth)
			if err != nil {
				return err
			}
			if !categoryRequired(cat) && !exists {
				log.Warnf("Skipping nonexistent, optional %s file %s", cat, pth)
				continue
			}
			log.Printf("Repairing %s", pth)
			errs += noteError(copyPath(src, dst, pth, opts))
			if cat == "history" {
				repairedHistory = true
			}
		}
	}

	if repairedHistory {
		log.Printf("Re-running checkpoing-file scan, for bucket repair")
		dst.ClearCachedInfo()
		errs += noteError(dst.ScanCheckpoints(opts))
	}

	errs += noteError(dst.ScanBuckets(opts))

	log.Printf("Examining buckets referenced by checkpoints")
	missingBuckets := dst.CheckBucketsMissing()

	for bkt := range missingBuckets {
		pth := BucketPath(bkt)
		log.Printf("Repairing %s", pth)
		errs += noteError(copyPath(src, dst, pth, opts))
	}

	if errs != 0 {
		return fmt.Errorf("%d errors while repairing", errs)
	}
	return nil
}
