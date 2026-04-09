package main

import (
	"bufio"
	"fmt"
	"strings"
)

type STPFrame struct {
	Headers map[string][]string
	Body    string
	Session string
}

func (f *STPFrame) String() string {
	return fmt.Sprintf(
		"STPFrame{Session: %q, Headers: %v, BodyLen: %d}",
		f.Session,
		f.Headers,
		len(f.Body),
	)
}

func ParseSTP(input string) (*STPFrame, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	frame := &STPFrame{Headers: make(map[string][]string)}
	var bodyBuilder strings.Builder

	inHeader, inBody := false, false

	for scanner.Scan() {
		line := scanner.Text()

		// 1. Detect Frame Start
		if !inHeader && !inBody && strings.HasPrefix(line, "!expression") {
			inHeader = true
			continue
		}

		// 2. Parse Headers (Key-Value pairs)
		if inHeader {
			if len(line) > 0 && line[0] == '@' {
				// Split strictly into 2 parts: "@key" and "value..."
				parts := strings.SplitN(line[1:], " ", 2)
				key := parts[0]
				val := ""
				if len(parts) > 1 {
					val = parts[1]
				}
				frame.Headers[key] = append(frame.Headers[key], val)
				continue
			}
			// Transition to Body if line doesn't start with @
			inHeader = false
			inBody = true
		}

		// 3. Accumulate Body until End Token
		if inBody {
			// Check length first to avoid panic, then compare the exact prefix
			if len(line) >= 16 && line[:16] == "!expression:end:" {
				frame.Session = line[16:]
				break
			}

			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
		}
	}

	frame.Body = bodyBuilder.String()
	return frame, scanner.Err()
}

func main() {
	rawInput := `
!expression
@tool v1Python
@tag item1
@tag item2

@decorator
def some_func() -> int:
    """
    Everything here is a raw text.
    Symbols like @ or ! are preserved as-is.
    """
    return 0

!expression:end:TOKEN
	`

	frame, _ := ParseSTP(rawInput)
	fmt.Printf("Parsed: %+v\n", frame)
	fmt.Print(frame.Body)
}
