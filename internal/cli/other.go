package cli

import (
	"errors"
	"strings"

	ambie "github.com/ambie-ai/ambie-go"
	"github.com/spf13/cobra"
)

func translateCmd() *cobra.Command {
	var sourceLang, targetLang, formality, context string
	cmd := &cobra.Command{
		Use:   "translate TEXT",
		Short: "Translate text",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			r, err := c.Translate(ctx(), ambie.TranslateOptions{
				Text:       strings.Join(args, " "),
				TargetLang: targetLang,
				SourceLang: sourceLang,
				Formality:  formality,
				Context:    context,
			})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	cmd.Flags().StringVarP(&targetLang, "target", "t", "en", "target language code (e.g. es, fr, ja)")
	cmd.Flags().StringVarP(&sourceLang, "source", "s", "", "source language (auto-detect if omitted)")
	cmd.Flags().StringVar(&formality, "formality", "", "formal, informal, or neutral")
	cmd.Flags().StringVar(&context, "context", "", "domain hint to improve quality")
	return cmd
}

func ttsCmd() *cobra.Command {
	var voice, engine, language, encoding string
	var speed float64
	cmd := &cobra.Command{
		Use:   "tts TEXT",
		Short: "Synthesize speech from text",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			r, err := c.TTS(ctx(), ambie.TtsOptions{
				Text:     strings.Join(args, " "),
				Voice:    voice,
				Engine:   engine,
				Language: language,
				Encoding: encoding,
				Speed:    speed,
			})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	cmd.Flags().StringVar(&voice, "voice", "", "voice ID (engine-specific)")
	cmd.Flags().StringVar(&engine, "engine", "", "aura or melo")
	cmd.Flags().StringVar(&language, "language", "", "language code")
	cmd.Flags().StringVar(&encoding, "encoding", "", "mp3, wav, ogg, or linear16")
	cmd.Flags().Float64Var(&speed, "speed", 0, "playback speed multiplier")
	return cmd
}

func embedCmd() *cobra.Command {
	var model string
	cmd := &cobra.Command{
		Use:   "embed TEXT [TEXT ...]",
		Short: "Generate text embeddings",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			var text any = args[0]
			if len(args) > 1 {
				text = args
			}
			r, err := c.Embeddings(ctx(), ambie.EmbeddingsOptions{Text: text, Model: model})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "embedding model (default: BGE base)")
	return cmd
}

func summarizeCmd() *cobra.Command {
	var url string
	var maxLen int
	cmd := &cobra.Command{
		Use:   "summarize [TEXT]",
		Short: "Summarize text or a URL",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			opts := ambie.SummarizeOptions{URL: url, MaxLength: maxLen}
			if len(args) == 1 {
				opts.Text = args[0]
			} else if url == "" {
				return errors.New("summarize requires TEXT or --url")
			}
			r, err := c.Summarize(ctx(), opts)
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "URL to summarize instead of inline TEXT")
	cmd.Flags().IntVar(&maxLen, "max-length", 0, "maximum summary length in tokens")
	return cmd
}

func sentimentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sentiment TEXT [TEXT ...]",
		Short: "Analyze sentiment (POSITIVE/NEGATIVE)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			var text any = args[0]
			if len(args) > 1 {
				text = args
			}
			r, err := c.Sentiment(ctx(), ambie.SentimentOptions{Text: text})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	return cmd
}

func moderateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moderate TEXT",
		Short: "Screen content for safety hazards",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			r, err := c.Moderate(ctx(), ambie.ModerateOptions{Text: strings.Join(args, " ")})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	return cmd
}

func detectLangCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detect-lang TEXT",
		Short: "Detect the language of TEXT",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newClient()
			if err != nil {
				return err
			}
			r, err := c.DetectLanguage(ctx(), ambie.DetectLangOptions{Text: strings.Join(args, " ")})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	return cmd
}

func rerankCmd() *cobra.Command {
	var topK int
	cmd := &cobra.Command{
		Use:   "rerank QUERY -d DOC [-d DOC ...]",
		Short: "Rerank documents by relevance to QUERY",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			docs, _ := cmd.Flags().GetStringArray("doc")
			if len(docs) == 0 {
				return errors.New("rerank requires at least one --doc")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			r, err := c.Rerank(ctx(), ambie.RerankOptions{
				Query:     args[0],
				Documents: docs,
				TopK:      topK,
			})
			if err != nil {
				return err
			}
			return emit(r, "")
		},
	}
	cmd.Flags().StringArrayP("doc", "d", nil, "document to rerank (repeat for multiple)")
	cmd.Flags().IntVar(&topK, "top-k", 0, "return only the top K results")
	return cmd
}
