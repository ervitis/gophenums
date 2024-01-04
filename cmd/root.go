package cmd

import (
	"log/slog"
	"os"

	"github.com/ervitis/gophenums/enum"
	"github.com/urfave/cli/v2"
)

func Execute() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := (&cli.App{
		Name:  "gophenums",
		Usage: "Generate enums with type safe from types",
		Action: func(ctx *cli.Context) error {
			return enum.NewGenerator().Generate()
		},
	}).Run(os.Args); err != nil {
		logger.Error("Error generating enums", slog.Any("error", err))
	}
}
