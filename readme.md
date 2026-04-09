# Simple Textual Protocol

A high-performance, human-readable framing protocol for LLM-to-Tool communication. 
Designed for reliable data exchange without the overhead of JSON escaping.

## Why?

*  **Zero-Escaping**: The body starts where headers end.
*  **Streaming-Ready**: Execute payloads the moment `@headers` transition to the body.
*  **Token-Efficient**: ~20% more compact than standard JSON tool-calling.
*  **Human-Centric**: Pure text, grep-friendly, perfectly readable in any Markdown editor.

## Example of a Frame

```text
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
```

**Parsed as**:

*  **Session**: TOKEN
*  **Headers**: tool: v1Python, tag: [item1, item2]
*  **Body**: raw script content (146 bytes)
    
## Logic

*  **Trigger**: !expression signals the start of a frame.
*  **Headers**: @key value pairs define the execution context.
*  **Payload**: First line without a leading @ switches the parser to raw mode.
*  **Termination**: !expression:end:TOKEN closes the frame and validates the sequence.
