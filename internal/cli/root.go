// Package cli implements the AMBIE command-line interface.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	ambie "github.com/ambie-ai/ambie-go"
	"github.com/spf13/cobra"
)

// Version is set via -ldflags at release time.
var Version = "dev"

// flags shared across subcommands
var (
	flagAPIKey  string
	flagBaseURL string
	flagOut     string
)

// Root returns the root cobra command.
func Root() *cobra.Command {
	root := &cobra.Command{
		Use:           "ambie",
		Short:         "AMBIE — speech-to-text, translation, embeddings, and more.",
		Long:          "Official command-line interface for AMBIE (https://ambie.ai). Get an API key at https://ambie.ai/signup.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "AMBIE API key (or set AMBIE_API_KEY)")
	root.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "override AMBIE base URL")
	root.PersistentFlags().StringVarP(&flagOut, "out", "o", "", "write output to FILE instead of stdout")

	root.AddCommand(
		versionCmd(),
		transcribeCmd(),
		translateCmd(),
		ttsCmd(),
		embedCmd(),
		summarizeCmd(),
		sentimentCmd(),
		moderateCmd(),
		detectLangCmd(),
		rerankCmd(),
	)

	return root
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ambie version", Version)
		},
	}
}

func newClient() (*ambie.Client, error) {
	key := flagAPIKey
	if key == "" {
		key = os.Getenv("AMBIE_API_KEY")
	}
	if key == "" {
		return nil, errors.New("AMBIE API key required: set AMBIE_API_KEY or pass --api-key. Get one at https://ambie.ai/signup")
	}
	var opts []ambie.Option
	if flagBaseURL != "" {
		opts = append(opts, ambie.WithBaseURL(flagBaseURL))
	}
	return ambie.New(key, opts...)
}

func ctx() context.Context {
	return context.Background()
}

// emit writes either JSON-marshaled v or s (when v is nil) to --out file or stdout.
func emit(v any, raw string) error {
	var w io.Writer = os.Stdout
	if flagOut != "" {
		f, err := os.Create(flagOut)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	if raw != "" {
		_, err := fmt.Fprint(w, raw)
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
