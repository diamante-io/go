package pipeline

import (
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	buffer := &BufferedReadWriter{}
	write := 20
	read := 0

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			_, err := buffer.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}
			read++
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < write; i++ {
			buffer.Write("test")
		}
		buffer.Close()
	}()

	wg.Wait()

	assert.Equal(t, 20, read)
}
