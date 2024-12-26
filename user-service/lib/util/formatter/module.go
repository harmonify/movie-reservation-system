package formatter

import "go.uber.org/fx"

var (
	FormatterModule = fx.Module(
		"formatter-util",
		fx.Provide(NewFormatterUtil),
	)
)