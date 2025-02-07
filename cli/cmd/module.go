package cmd

import (
	"github.com/harmonify/movie-reservation-system/cli/cmd/kafka"
	"github.com/harmonify/movie-reservation-system/cli/shared"

	"go.uber.org/fx"
)

var (
	CommandModule = fx.Module(
		"command",
		fx.Provide(
			AsCommand(kafka.NewKafkaCmd),
			AsCommand(
				kafka.NewKafkaMigrateUpCmd,
				fx.ParamTags(`group:"migrations"`),
			),
			AsCommand(
				kafka.NewKafkaMigrateDownCmd,
				fx.ParamTags(`group:"migrations"`),
			),
			AsCommand(
				kafka.NewKafkaMigrateNewCmd,
			),
			fx.Annotate(
				NewRootCmd,
				fx.ParamTags(`group:"commands"`),
			),
		),
	)
)

func AsCommand(f any, anns ...fx.Annotation) any {
	finalAnns := []fx.Annotation{
		fx.As(new(shared.CobraCommand)),
		fx.ResultTags(`group:"commands"`),
	}
	if len(anns) > 0 {
		finalAnns = append(finalAnns, anns...)
	}

	return fx.Annotate(
		f,
		finalAnns...,
	)
}
