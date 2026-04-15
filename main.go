package main

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/shrpow/simple-textual-protocol/stp"
)

func main() {
	// Mock input stream from LLM
	rawInput := `
@a16c2f35
param value
someList item1
someList item2
> @decorator
> def some_func():
> 	return "this is a first frame"
@b25d3e46
key param
> this is a second frame
`
	frameChan := make(chan *stp.STPFrame)

	// Start the parser in the background
	go func() {
		if err := stp.StreamSTPFrames(strings.NewReader(rawInput), frameChan); err != nil {
			fmt.Printf("Parser error: %v\n", err)
		}
	}()

	var wg sync.WaitGroup

	// Main processing loop
	for frame := range frameChan {
		fmt.Printf("[%s] params: %v\n", frame.Command, frame.Params)

		// IMPORTANT: Read the Body in a goroutine so as not to block the parser thread.
		// This enables concurrent tool execution.
		wg.Add(1)

		go func(f *stp.STPFrame) {
			defer wg.Done()
			defer f.Body.Close()

			// Here you would implement tool logic, e.g., exec.Command or Python interpreter call.

			content, _ := io.ReadAll(f.Body)
			fmt.Printf("[%s] body: \n%s", f.Command, string(content))
		}(frame)
	}

	wg.Wait()
}
