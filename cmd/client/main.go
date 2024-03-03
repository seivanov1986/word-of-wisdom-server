package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/seivanov1986/word-of-wisdom-server/internal/client"
)

func main() {
	ctx := context.Background()

	hashCacheMaxIterationsStr := os.Getenv("HASH_CASH_MAX_ITERATIONS")
	hashCacheMaxIterations, _ := strconv.ParseInt(hashCacheMaxIterationsStr, 10, 64)

	err := client.New(
		os.Getenv(client.ENVClientAddress),
		int(hashCacheMaxIterations),
	).Start(ctx)
	if err != nil {
		fmt.Println("client error:", err)
	}
}
