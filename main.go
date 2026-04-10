package main

import (
	"bufio"
	"fmt"
	"strings"
)

type STPFrame struct {
	Session string
	Headers map[string][]string
	Body    string
}

func (f *STPFrame) String() string {
	return fmt.Sprintf(
		"STPFrame{Session: %q, Headers: %v, BodyLen: %d}",
		f.Session,
		f.Headers,
		len(f.Body),
	)
}

func ParseSTP(input, sessionID string) (*STPFrame, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	frame := &STPFrame{
		Session: sessionID,
		Headers: make(map[string][]string),
	}
	var bodyBuilder strings.Builder

	endMarker := "!expression:end:" + sessionID
	inHeader, inBody, foundEnd := false, false, false

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
			if line == endMarker {
				foundEnd = true
				break
			}

			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
		}
	}

	if !foundEnd {
		return nil, fmt.Errorf("incomplete frame: missing end marker %q", endMarker)
	}

	frame.Body = bodyBuilder.String()
	return frame, scanner.Err()
}

func main() {
	rawInput := `
!expression
@tool python
@tag 1
@tag 2

def create_squares():
    squares = []
    for i in range(1, 6):
        squares.append(i * i)
    print(squares)

create_squares()

!expression:end:

!expression
@tool test

!expression:end:aXyYYZ
	`

	frame, err := ParseSTP(rawInput, "aXyYYZ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed: %+v\n", frame)
	fmt.Print(frame.Body)
}
