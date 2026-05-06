# ambie-cli

The official AMBIE command-line tool — speech-to-text, translation, TTS, embeddings, sentiment, summarization, content moderation, and language detection from your terminal or scripts.

[![Release](https://img.shields.io/github/v/release/ambie-ai/ambie-cli)](https://github.com/ambie-ai/ambie-cli/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/ambie-ai/ambie-cli.svg)](https://pkg.go.dev/github.com/ambie-ai/ambie-cli)

## Install

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/ambie-ai/ambie-cli/main/install.sh | sh
```

The installer downloads the right binary for your platform from the [latest release](https://github.com/ambie-ai/ambie-cli/releases/latest) and places it at `/usr/local/bin/ambie`. To install elsewhere, set `AMBIE_INSTALL_DIR`.

### Manual download

Pick the binary for your platform from [the releases page](https://github.com/ambie-ai/ambie-cli/releases/latest), unpack, and put `ambie` somewhere on your `PATH`.

```bash
# Example for Linux amd64:
curl -fsSLo ambie.tar.gz https://github.com/ambie-ai/ambie-cli/releases/latest/download/ambie_<VERSION>_linux_amd64.tar.gz
tar -xzf ambie.tar.gz
sudo mv ambie /usr/local/bin/
```

### Windows

Download the `.zip` for `windows_amd64` from the [releases page](https://github.com/ambie-ai/ambie-cli/releases/latest), extract, and add the folder to your `PATH`.

### Build from source

```bash
go install github.com/ambie-ai/ambie-cli/cmd/ambie@latest
```

## Authenticate

Get an API key at [https://ambie.ai/signup](https://ambie.ai/signup), then set:

```bash
export AMBIE_API_KEY=amb_live_...
```

Or pass `--api-key` on each invocation.

## Quickstart

```bash
# Transcribe a meeting recording, with diarization and an LLM summary
ambie transcribe meeting.mp3 --diarize --summarize

# Translate text to Spanish
ambie translate "Hello, world!" --target es

# Generate embeddings for a list of phrases
ambie embed "alpha" "beta" "gamma"

# Synthesize speech to a wav file
ambie tts "The quick brown fox" --encoding wav --out out.wav.json

# Summarize a URL
ambie summarize --url https://en.wikipedia.org/wiki/AMBIE_(disambiguation)

# Detect language
ambie detect-lang "El zorro marrón rápido salta sobre el perro perezoso."

# Sentiment of a single string or a batch
ambie sentiment "I love this product" "This is awful"

# Moderate text against safety hazards
ambie moderate "Some user-generated content"

# Rerank documents by relevance to a query
ambie rerank "How do I cook pasta?" -d "Boil water and salt it." -d "Cars run on gasoline."
```

## Output

By default `ambie` prints indented JSON to stdout. For the transcribe command you can request raw text or subtitle formats:

```bash
ambie transcribe meeting.mp3 --format text
ambie transcribe meeting.mp3 --format srt --out meeting.srt
```

For any command you can redirect to a file with `--out FILE`.

## Async mode

Pass `--callback-url` to transcribe (and any other endpoint when API support is available) to get the result delivered to your webhook:

```bash
ambie transcribe long-recording.mp3 --callback-url https://example.com/webhooks/ambie
```

The CLI prints the `request_id` and `poll_url` immediately and exits.

## Configuration

| Flag             | Env var          | Description                                  |
| ---------------- | ---------------- | -------------------------------------------- |
| `--api-key`      | `AMBIE_API_KEY`  | API key (required)                           |
| `--base-url`     | —                | Override default `https://ambie.ai`          |
| `--out FILE`     | —                | Write output to FILE instead of stdout       |

## License

MIT © AMBIE
