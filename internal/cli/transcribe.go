package cli

import (
	"errors"
	"os"
	"strings"

	ambie "github.com/ambie-ai/ambie-go"
	"github.com/spf13/cobra"
)

func transcribeCmd() *cobra.Command {
	var (
		url, engine, language, format, callbackURL, clientID         string
		translate, diarize, summarize, keyPhrases, actions, chapters bool
	)

	cmd := &cobra.Command{
		Use:   "transcribe [FILE]",
		Short: "Transcribe an audio file or URL",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			opts := ambie.TranscribeOptions{
				URL:         url,
				Engine:      engine,
				Language:    language,
				Format:      format,
				Translate:   translate,
				CallbackURL: callbackURL,
				ClientID:    clientID,
				Diarize:     diarize,
				Summarize:   summarize,
				KeyPhrases:  keyPhrases,
				ActionItems: actions,
				Chapters:    chapters,
			}
			if len(args) == 1 {
				path := args[0]
				audio, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				opts.Audio = audio
				opts.AudioName = baseName(path)
			} else if url == "" {
				return errors.New("transcribe requires a FILE argument or --url")
			}
			if callbackURL != "" {
				accepted, err := c.TranscribeAsync(ctx(), opts)
				if err != nil {
					return err
				}
				return emit(accepted, "")
			}
			r, err := c.Transcribe(ctx(), opts)
			if err != nil {
				return err
			}
			if format != "" && format != "json" {
				return emit(nil, r.Text)
			}
			return emit(r, "")
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "remote audio URL (alternative to FILE)")
	cmd.Flags().StringVar(&engine, "engine", "deepgram", "deepgram or whisper")
	cmd.Flags().StringVar(&language, "language", "", "BCP-47 hint (e.g. en, es, fr-CA)")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "json, text, srt, or vtt")
	cmd.Flags().BoolVar(&translate, "translate", false, "translate to English (whisper only)")
	cmd.Flags().StringVar(&callbackURL, "callback-url", "", "webhook URL for async delivery")
	cmd.Flags().StringVar(&clientID, "client-id", "", "your correlation ID")
	cmd.Flags().BoolVar(&diarize, "diarize", false, "label speakers (deepgram only)")
	cmd.Flags().BoolVar(&summarize, "summarize", false, "generate LLM summary")
	cmd.Flags().BoolVar(&keyPhrases, "key-phrases", false, "extract key phrases")
	cmd.Flags().BoolVar(&actions, "action-items", false, "extract action items")
	cmd.Flags().BoolVar(&chapters, "chapters", false, "split into timestamped chapters")
	return cmd
}

func baseName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}
