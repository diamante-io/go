// Copyright 2016 diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"path"

	log "github.com/sirupsen/logrus"
)

func makeTicker(onTick func(uint)) chan bool {
	tick := make(chan bool)
	go func() {
		var k uint = 0
		for range tick {
			k++
			if k&0xfff == 0 {
				onTick(k)
			}
		}
	}()
	return tick
}

func bufReadCloser(in io.ReadCloser) io.ReadCloser {
	return struct {
		io.Reader
		io.Closer
	}{bufio.NewReader(in), in}
}

func copyPath(src *Archive, dst *Archive, pth string, opts *CommandOptions) error {
	if opts.DryRun {
		log.Printf("dryrun skipping " + pth)
		return nil
	}
	exists, err := dst.backend.Exists(pth)
	if err != nil {
		return err
	}
	if exists && !opts.Force {
		log.Printf("skipping existing " + pth)
		return nil
	}
	rdr, err := src.backend.GetFile(pth)
	if err != nil {
		return err
	}
	defer rdr.Close()
	err = dst.backend.PutFile(pth, bufReadCloser(rdr))
	return err
}

func Categories() []string {
	return []string{"history", "ledger", "transactions", "results", "scp"}
}

func categoryExt(n string) string {
	if n == "history" {
		return "json"
	} else {
		return "xdr.gz"
	}
}

func categoryRequired(n string) bool {
	return n != "scp"
}

func CategoryCheckpointPath(cat string, chk uint32) string {
	ext := categoryExt(cat)
	pre := CheckpointPrefix(chk).Path()
	return path.Join(cat, pre, fmt.Sprintf("%s-%8.8x.%s", cat, chk, ext))
}

func BucketPath(bucket Hash) string {
	pre := HashPrefix(bucket)
	return path.Join("bucket", pre.Path(), fmt.Sprintf("bucket-%s.xdr.gz", bucket))
}

// Make a goroutine that unconditionally pulls an error channel into
// (unbounded) local memory, and feeds it to a downstream consumer. This is
// slightly hacky, but the alternatives are to either send {val,err} pairs
// along each "primary" channel, or else risk a primary channel producer
// stalling because nobody's draining its error channel yet. I (vaguely,
// currently) prefer this approach, though time will tell if it has a
// good taste later.
//
// Code here modeled on github.com/eapache/channels/infinite_channel.go
func makeErrorPump(in chan error) chan error {
	buf := make([]error, 0)
	var next error
	ret := make(chan error)
	go func() {
		var out chan error
		for in != nil || out != nil {
			select {
			case err, ok := <-in:
				if ok {
					buf = append(buf, err)
				} else {
					in = nil
				}
			case out <- next:
				buf = buf[1:]
			}
			if len(buf) > 0 {
				out = ret
				next = buf[0]
			} else {
				out = nil
				next = nil
			}
		}
		close(ret)
	}()
	return ret
}

func noteError(e error) uint32 {
	if e != nil {
		log.Errorf("Error: " + e.Error())
		return 1
	}
	return 0
}

func drainErrors(errs chan error) uint32 {
	var count uint32
	for e := range errs {
		count += noteError(e)
	}
	return count
}

func logReq(r *http.Request) {
	if r == nil {
		return
	}
	logFields := log.Fields{"method": r.Method, "url": r.URL.String()}
	log.WithFields(logFields).Trace("http: Req")
}

func logResp(r *http.Response) {
	if r == nil || r.Request == nil {
		return
	}
	logFields := log.Fields{"method": r.Request.Method, "status": r.Status, "url": r.Request.URL.String()}
	if r.StatusCode >= 200 && r.StatusCode < 400 {
		log.WithFields(logFields).Trace("http: OK")
	} else {
		log.WithFields(logFields).Warn("http: Bad")
	}
}
