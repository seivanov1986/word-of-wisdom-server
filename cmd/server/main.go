package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/cache"
	"github.com/seivanov1986/word-of-wisdom-server/internal/server"
)

func main() {
	ctx := context.Background()
	zerosCountStr := os.Getenv("ZEROS_COUNT")

	zerosCount, _ := strconv.ParseInt(zerosCountStr, 10, 64)

	err := server.New(
		os.Getenv(server.ENVServerAddress),
		int(zerosCount),
		cache.New(time.Now()),
	).Start(ctx)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
