package main

import (
	"context"
	"os"
	"strconv"

	"github.com/seivanov1986/word-of-wisdom-server/internal/client"
	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/logger"
)

func main() {
	ctx := context.Background()
	log := logger.New()

	hashCacheMaxIterationsStr := os.Getenv("HASH_CASH_MAX_ITERATIONS")
	hashCacheMaxIterations, err := strconv.ParseInt(hashCacheMaxIterationsStr, 10, 64)
	if err != nil {
		log.Fatal("zerosCount error:", err)
	}

	err = client.New(
		os.Getenv(client.ENVClientAddress),
		log,
		int(hashCacheMaxIterations),
	).Start(ctx)
	if err != nil {
		log.Fatal("client error:", err)
	}
}
