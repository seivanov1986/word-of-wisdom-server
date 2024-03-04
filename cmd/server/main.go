package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/cache"
	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/logger"
	"github.com/seivanov1986/word-of-wisdom-server/internal/server"
)

func main() {
	ctx := context.Background()
	zerosCountStr := os.Getenv("ZEROS_COUNT")
	log := logger.New()

	zerosCount, err := strconv.ParseInt(zerosCountStr, 10, 64)
	if err != nil {
		log.Fatal("zerosCount error:", err)
	}

	err = server.New(
		os.Getenv(server.ENVServerAddress),
		log,
		int(zerosCount),
		cache.New(time.Now()),
	).Start(ctx)
	if err != nil {
		log.Fatal("server error:", err)
	}
}
