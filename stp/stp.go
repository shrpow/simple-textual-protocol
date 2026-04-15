package stp

import (
	"bufio"
	"io"
	"strings"
)

// STPFrame represents a single STP protocol frame.
type STPFrame struct {
	Command string
	Params  map[string][]string
	Body    io.ReadCloser
}

// StreamSTPFrames reads from r and pushes frames into frameChan.
// It sends a frame as soon as the command line is parsed, allowing
// streaming of the Body content.
func StreamSTPFrames(r io.Reader, frameChan chan<- *STPFrame) error {
	defer close(frameChan)
	scanner := bufio.NewScanner(r)
	var (
		current *STPFrame
		pw      *io.PipeWriter
		sent    bool
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Start of a new frame
		if strings.HasPrefix(trimmed, "@") {
			if current != nil && !sent {
				frameChan <- current // Send previous frame if it had no body
			}
			if pw != nil {
				pw.Close()
			}

			pr, nextPw := io.Pipe()
			pw = nextPw
			current = &STPFrame{
				Command: trimmed[1:],
				Params:  make(map[string][]string),
				Body:    pr,
			}
			sent = false
			continue
		}

		if current == nil {
			continue
		}

		// Body streaming starts
		if strings.HasPrefix(line, ">") {
			if !sent {
				frameChan <- current
				sent = true
			}
			content := strings.TrimPrefix(line[1:], " ")
			pw.Write([]byte(content + "\n"))
			continue
		}

		// Key-value parameters
		if trimmed != "" {
			key, val, _ := strings.Cut(trimmed, " ")
			current.Params[key] = append(current.Params[key], val)
		}
	}

	// Finalize the last frame
	if current != nil && !sent {
		frameChan <- current
	}

	if pw != nil {
		pw.Close()
	}

	return scanner.Err()
}
