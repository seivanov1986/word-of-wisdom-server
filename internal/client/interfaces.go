package client

import (
	"context"
)

type Client interface {
	Start(ctx context.Context) error
}
