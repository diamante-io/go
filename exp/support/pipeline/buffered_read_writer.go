package pipeline

import (
	"context"
	"io"
)

// bufferSize is a size of a buffered channel in BufferedReadWriter.
// This should be big enough to hold a short lag of items in a pipeline
// but small enough to not consume too much memory.
// In pipelines with no slow processors a buffered channel will be empty
// or almost empty most of the time.
const bufferSize = 50000

func (b *BufferedReadWriter) init() {
	b.buffer = make(chan interface{}, bufferSize)
}

func (b *BufferedReadWriter) close() {
	b.writeCloseMutex.Lock()
	defer b.writeCloseMutex.Unlock()

	close(b.buffer)
	b.closed = true
}

func (b *BufferedReadWriter) GetContext() context.Context {
	return b.context
}

func (b *BufferedReadWriter) Read() (interface{}, error) {
	b.initOnce.Do(b.init)

	entry, more := <-b.buffer
	if more {
		b.readEntriesMutex.Lock()
		b.readEntries++
		b.readEntriesMutex.Unlock()
		return entry, nil
	}

	return nil, io.EOF
}

func (b *BufferedReadWriter) Write(entry interface{}) error {
	b.initOnce.Do(b.init)

	b.writeCloseMutex.Lock()
	defer b.writeCloseMutex.Unlock()

	if b.closed {
		return io.ErrClosedPipe
	}

	b.buffer <- entry
	b.wroteEntries++
	return nil
}

func (b *BufferedReadWriter) QueuedEntries() int {
	b.initOnce.Do(b.init)
	return len(b.buffer)
}

// Close can be called in `Writer` and `Reader` context.
//
// In `Reader` it means that no more values will be read so writer can
// stop writing to a buffer (`io.ErrClosedPipe` will be returned for calls to
// `Write()`).
//
// In `Writer` it means that no more values will be written so reader
// should start returning `io.EOF` error after returning all queued values.
func (b *BufferedReadWriter) Close() error {
	b.initOnce.Do(b.init)
	b.closeOnce.Do(b.close)
	return nil
}

var _ Reader = &BufferedReadWriter{}
var _ Writer = &BufferedReadWriter{}
