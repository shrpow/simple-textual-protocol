# Simple Textual Protocol (STP)

A high-performance framing protocol for LLM-to-Tool communication.
Designed for raw data exchange without JSON escaping.

## Why

- **Zero-Escaping**: Raw payload after the `>` marker.
- **Streaming-Ready**: Execute payloads the moment `@headers` transition to the body.
- **Token-Optimized**: Often 10-30% cheaper than JSON for code-heavy payloads.
- **Human-Centric**: Pure text, grep-friendly, perfectly readable in any Markdown editor.
- **Linear**: O(1) parsing.

## Example

```text
@a16c2f35
param value
someList item1
someList item2
> @decorator
> def some_func():
>     return "xd"
```

## Specification

- **Trigger**: @ at the start of a line starts a new frame.
- **Metadata**: key value pairs (one per line) before the body.
- **Body**: Lines starting with `>` are raw data. Strips `>` and one optional space.
- **Termination**: A new @ at the start of a line or EOF (End of File) terminates the current body.
