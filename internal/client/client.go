package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/seivanov1986/word-of-wisdom-server/internal/helpers/hash_cache_data"
	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/logger"
	"github.com/seivanov1986/word-of-wisdom-server/internal/proto"
	"github.com/seivanov1986/word-of-wisdom-server/internal/vo"
)

const (
	ENVClientAddress = "CLIENT_ADDRESS"
)

type client struct {
	address               string
	hashcashMaxIterations int
	log                   logger.Logger
}

func New(address string, log logger.Logger, hashcashMaxIterations int) *client {
	return &client{address: address, log: log, hashcashMaxIterations: hashcashMaxIterations}
}

func (c *client) Start(ctx context.Context) error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		message, err := c.handle(ctx, conn)

		if err != nil {
			c.log.Println("client error handle connection: ", err)
		} else {
			c.log.Println("quote result:", *message)
		}

		time.Sleep(5 * time.Second)
	}
}

func (c *client) handle(ctx context.Context, conn io.ReadWriter) (*string, error) {
	reader := bufio.NewReader(conn)

	err := c.send(proto.MessageToBytes(proto.Message{
		Header: vo.RequestChallengeHeader,
	}), conn)
	if err != nil {
		return nil, err
	}

	msgStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msg, err := proto.Parse(msgStr)
	if err != nil {
		return nil, err
	}

	var hashcash hash_cache_data.HashcashData
	err = json.Unmarshal([]byte(msg.Payload.String()), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err parse hashcash: %w", err)
	}
	c.log.Println("got hashcash:", hashcash)

	hashcash, err = hashcash.ComputeHashcash(c.hashcashMaxIterations)
	if err != nil {
		return nil, fmt.Errorf("err compute hashcash: %w", err)
	}
	c.log.Println("hashcash computed:", hashcash)

	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashcash: %w", err)
	}

	payload, _ := vo.ParsePayload(string(byteData))

	err = c.send(proto.MessageToBytes(proto.Message{
		Header:  vo.RequestResourceHeader,
		Payload: payload,
	}), conn)
	if err != nil {
		return nil, fmt.Errorf("err send request: %w", err)
	}
	c.log.Println("challenge sent to server")

	msgStr, err = reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("err read msg: %w", err)
	}
	msg, err = proto.Parse(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}

	result := msg.Payload.String()

	return &result, nil
}

func (c *client) send(msg []byte, conn io.Writer) error {
	c.log.Println(string(msg))

	_, err := conn.Write(msg)
	return err
}
