package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Log = zerolog.New(os.Stderr).With().Timestamp().Logger()
}
