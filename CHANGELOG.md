# Changelog

All notable changes to this CLI are documented here. Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/).

## [0.1.0] - 2026-05-06

Initial public release.

- Cross-platform binaries: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64).
- Subcommands: `transcribe`, `translate`, `tts`, `embed`, `summarize`, `sentiment`, `moderate`, `detect-lang`, `rerank`, `version`.
- Sync and async (`--callback-url`) modes for `transcribe`.
- File / URL / stdin input for transcription.
- JSON, plain text, SRT, and WebVTT output for transcription.
- `install.sh` one-liner installer for macOS/Linux.
- Built with [github.com/ambie-ai/ambie-go](https://github.com/ambie-ai/ambie-go).
