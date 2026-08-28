# Two reproducible API/CLI failures, found 2026-08-28

Found while transcribing 21 radio commercials for damngoodpepper.com. Both
were worked around by calling the API directly with curl, which succeeded on
every one of the 21 files.

## 1. `--engine deepgram` returns 502

```
curl -H "Authorization: Bearer $AMBIE_API_KEY" \
  -F "audio=@spot.mp3" -F "engine=deepgram" -F "language=th" \
  https://ambie.ai/api/v1/transcribe
→ HTTP 502, body: "error code: 502"
```

`engine=whisper` on the same file and the same key returns HTTP 200 with a
full result. deepgram is the CLI's DEFAULT engine, so the default path is the
broken one.

The 502 body is not JSON, which is what produces the CLI error below when the
default engine is used.

## 2. The CLI fails where the identical curl succeeds

```
ambie transcribe spot.mp3 --engine whisper --language en --format text
→ invalid character 'I' looking for beginning of value
```

That is a Go JSON decode error on a non-JSON body beginning with `I`
(most likely `Internal server error`). It reproduces with `--engine whisper`
too, so it is NOT only the deepgram 502 above:

| call | result |
|---|---|
| `curl` + `engine=whisper` | 200, full transcript |
| `ambie transcribe --engine whisper` | `invalid character 'I'` |
| `ambie transcribe` (default deepgram) | `invalid character 'I'` |

Same key, same file, same host. So there are two defects, not one.

Worth checking in `ambie-go/client.go`: `Transcribe` reads and re-wraps the
request body after `http.NewRequestWithContext` has already consumed the
`*bytes.Buffer`:

```go
bodyBytes, _ := io.ReadAll(req.Body)
req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
```

The multipart field name (`audio`) and the URL (`BaseURL + /api/v1/transcribe`)
both match the working curl, so the difference is in how the request is sent,
not what it is addressed to.

Two smaller notes:

- `ambie --version` is not a flag (`unknown flag: --version`).
- A non-JSON error body should surface as the HTTP status and the body text,
  not as a JSON parse error. The current message tells the caller nothing
  about a 502.

## What worked

`POST https://ambie.ai/api/v1/transcribe` with `engine=whisper` and a
`language` hint transcribed all 21 spots across 15 languages, with usable
`transcript_confidence`, `duration_seconds` and word timings. Quality caveats
(hallucinated tails on Chinese/Japanese, Thai unusable at 0.53-0.72) are
whisper-model behaviour, not platform bugs.
