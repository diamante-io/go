package pipeline

import (
	"context"
	"sync"

	"github.com/diamcircle/go/support/errors"
)

// ErrShutdown is an error send to post-processing hook when pipeline has been
// shutdown.
var ErrShutdown = errors.New("Pipeline shutdown")

// BufferedReadWriter implements Reader and Writer and acts
// like a pipe. All writes are queued in a buffered channel and are waiting
// to be consumed.
//
// Used internally by Pipeline but also helpful for testing.
type BufferedReadWriter struct {
	initOnce sync.Once

	context context.Context

	// readEntriesMutex protects readEntries variable
	readEntriesMutex sync.Mutex
	readEntries      int

	// writeCloseMutex protects from writing to a closed buffer
	// and wroteEntries variable
	writeCloseMutex sync.Mutex
	wroteEntries    int

	// closeOnce protects from closing buffer twice
	closeOnce sync.Once
	buffer    chan interface{}
	closed    bool
}

type multiWriter struct {
	writers []Writer

	mutex        sync.Mutex
	closeAfter   int
	wroteEntries int
}

type Pipeline struct {
	root *PipelineNode

	preProcessingHooks  []func(context.Context) (context.Context, error)
	postProcessingHooks []func(context.Context, error) error

	// mutex protects internal fields that may be modified from
	// multiple go routines.
	mutex      sync.Mutex
	running    bool
	shutDown   bool
	cancelled  bool
	cancelFunc context.CancelFunc
}

// PipelineInterface is an interface that defines common pipeline methods
// in structs that embed Pipeline.
type PipelineInterface interface {
	SetRoot(rootProcessor *PipelineNode)
	// AddPreProcessingHook adds a pre-processing hook function. Returned
	// context.Context will be passed to the processors. If error is returned
	// pipeline will not start processing data.
	AddPreProcessingHook(hook func(context.Context) (context.Context, error))
	AddPostProcessingHook(hook func(context.Context, error) error)
	Shutdown()
	PrintStatus()
}

var _ PipelineInterface = &Pipeline{}

type PipelineNode struct {
	// Remember to update reset() method if you ever add a new field to this struct!
	Processor Processor
	Children  []*PipelineNode

	readEntries     int
	readsPerSecond  int
	queuedEntries   int
	wroteEntries    int
	writesPerSecond int
}

// Reader interface placeholder
type Reader interface {
	// GetContext returns context with values of the current reader. Can be
	// helpful to provide data to structs that wrap `Reader`.
	GetContext() context.Context
	// Read should return next entry. If there are no more
	// entries it should return `io.EOF` error.
	Read() (interface{}, error)
	// Close should be called when reading is finished. This is especially
	// helpful when there are still some entries available so reader can stop
	// streaming them.
	Close() error
}

// Writer interface placeholder
type Writer interface {
	// Write is used to pass entry to the next processor. It can return
	// `io.ErrClosedPipe` when the pipe between processors has been closed meaning
	// that next processor does not need more data. In such situation the current
	// processor can terminate as sending more entries to a `Writer`
	// does not make sense (will not be read).
	Write(interface{}) error
	// Close should be called when there are no more entries
	// to write.
	Close() error
}

// Processor defines methods required by the processing pipeline.
type Processor interface {
	// Process is a main method of `Processor`. It receives `Reader`
	// that contains object passed down the pipeline from the previous procesor. Writes to
	// `Writer` will be passed to the next processor. WARNING! `Process`
	// should **always** call `Close()` on `Writer` when no more object will be
	// written and `Close()` on `Reader` when reading is finished.
	// Data required by following processors (like aggregated data) should be saved in
	// `Store`. Read `Store` godoc to understand how to use it.
	// The first argument `ctx` is a context with cancel. Processor should monitor
	// `ctx.Done()` channel and exit when it returns a value. This can happen when
	// pipeline execution is interrupted, ex. due to an error.
	//
	// Given all information above `Process` should always look like this:
	//
	//    func (p *Processor) Process(ctx context.Context, store *pipeline.Store, r Reader, w Writer) error {
	//    	defer r.Close()
	//    	defer w.Close()
	//
	//    	// Some pre code...
	//
	//    	for {
	//    		entry, err := r.Read()
	//    		if err != nil {
	//    			if err == io.EOF {
	//    				break
	//    			} else {
	//    				return errors.Wrap(err, "Error reading from Reader in [ProcessorName]")
	//    			}
	//    		}
	//
	//    		// Process entry...
	//
	//    		// Write to Writer if needed but exit if pipe is closed:
	//    		err = w.Write(entry)
	//    		if err != nil {
	//    			if err == io.ErrClosedPipe {
	//    				// Reader does not need more data
	//    				return nil
	//    			}
	//    			return errors.Wrap(err, "Error writing to Writer in [ProcessorName]")
	//    		}
	//
	//    		// Return errors if needed...
	//
	//    		// Exit when pipeline terminated due to an error in another processor...
	//    		select {
	//    		case <-ctx.Done():
	//    			return nil
	//    		default:
	//    			continue
	//    		}
	//    	}
	//
	//    	// Some post code...
	//
	//    	return nil
	//    }
	Process(context.Context, *Store, Reader, Writer) error
	// Returns processor name. Helpful for errors, debuging and reports.
	Name() string
	// Reset resets internal state of the processor. This is run by the pipeline
	// everytime before the pipeline starts running.
	// It is extremely important to implement this method, otherwise internal
	// state of the processor will be maintained between pipeline runs and may
	// result in an invalid data.
	Reset()
}

// Store allows storing data connected to pipeline execution.
// It exposes `Lock()` and `Unlock()` methods that must be called
// when accessing the `Store` for both `Put` and `Get` calls.
//
// Example (incrementing a number):
// s.Lock()
// v := s.Get("value")
// s.Put("value", v.(int)+1)
// s.Unlock()
type Store struct {
	sync.Mutex
	initOnce sync.Once
	values   map[string]interface{}
}
