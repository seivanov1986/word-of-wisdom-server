package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/seivanov1986/word-of-wisdom-server/internal/helpers/hash_cache_data"
	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/cache"
	"github.com/seivanov1986/word-of-wisdom-server/internal/pkg/logger"
	"github.com/seivanov1986/word-of-wisdom-server/internal/proto"
	"github.com/seivanov1986/word-of-wisdom-server/internal/vo"
)

const (
	ENVServerAddress = "SERVER_ADDRESS"
)

type server struct {
	address          string
	zerosCount       int
	hashcashDuration int
	cache            cache.Cache
	log              logger.Logger
}

func New(address string, log logger.Logger, zerosCount int, cache cache.Cache) *server {
	return &server{address: address, log: log, zerosCount: zerosCount, cache: cache}
}

func (c *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", c.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}
		go c.handle(ctx, conn)
	}
}

func (c *server) handle(ctx context.Context, conn net.Conn) {
	c.log.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			c.log.Println("err read connection:", err)
			return
		}

		msg, err := c.process(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			c.log.Println("err process request:", err)
			return
		}

		if msg != nil {
			err := c.send(proto.MessageToBytes(*msg), conn)
			if err != nil {
				c.log.Println("err send message:", err)
			}
		}
	}
}

func (c *server) process(ctx context.Context, msgStr string, clientInfo string) (*proto.Message, error) {
	msg, err := proto.Parse(msgStr)
	if err != nil {
		return nil, err
	}

	switch msg.Header {
	case vo.QuitHeader:
		return nil, errors.New("client requests to close connection")
	case vo.RequestChallengeHeader:
		randValue, err := c.generateRandomValue()
		if err != nil {
			return nil, err
		}

		hashData, err := c.makeHashCacheData(clientInfo, randValue)
		if err != nil {
			return nil, err
		}

		payload, err := vo.ParsePayload(*hashData)
		if err != nil {
			return nil, err
		}

		msg := proto.Message{
			Header:  vo.ResponseChallengeHeader,
			Payload: payload,
		}
		return &msg, nil
	case vo.RequestResourceHeader:
		hashcash, err := c.getHashCacheData(*msg)
		if err != nil {
			return nil, err
		}

		err = c.validateHashCacheData(*hashcash, clientInfo)
		if err != nil {
			return nil, err
		}

		c.log.Println("client %s succesfully computed hashcash %s", clientInfo, msg.Payload)

		payload, err := vo.ParsePayload(string(`{"message": "test string"}`))
		if err != nil {
			return nil, err
		}

		msg := proto.Message{
			Header:  vo.ResponseResourceHeader,
			Payload: payload,
		}
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func (c *server) validateHashCacheData(hashcash hash_cache_data.HashcashData, clientInfo string) error {
	if hashcash.Resource != clientInfo {
		return fmt.Errorf("invalid hashcash resource")
	}

	randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
	if err != nil {
		return fmt.Errorf("err decode rand: %w", err)
	}
	randValue, err := strconv.Atoi(string(randValueBytes))
	if err != nil {
		return fmt.Errorf("err decode rand: %w", err)
	}

	exists, err := c.cache.Get(randValue)
	if err != nil {
		return fmt.Errorf("err get rand from cache: %w", err)
	}
	if !exists {
		return fmt.Errorf("challenge expired or not sent")
	}

	if time.Now().Unix()-hashcash.Date > int64(c.hashcashDuration) {
		return fmt.Errorf("challenge expired")
	}

	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}
	_, err = hashcash.ComputeHashcash(maxIter)
	if err != nil {
		return fmt.Errorf("invalid hashcash")
	}

	c.cache.Delete(randValue)

	return nil
}

func (c *server) getHashCacheData(msg proto.Message) (*hash_cache_data.HashcashData, error) {
	var hashcash hash_cache_data.HashcashData
	err := json.Unmarshal([]byte(msg.Payload.String()), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
	}

	return &hashcash, nil
}

func (c *server) makeHashCacheData(clientInfo string, randValue int) (*string, error) {
	hashcash := hash_cache_data.HashcashData{
		Version:    1,
		ZerosCount: c.zerosCount,
		Date:       time.Now().Unix(),
		Resource:   clientInfo,
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
		Counter:    0,
	}

	hashcashMarshaled, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashcash: %v", err)
	}

	hashCacheData := string(hashcashMarshaled)

	return &hashCacheData, nil
}

func (c *server) generateRandomValue() (int, error) {
	randValue := rand.Intn(100000)
	err := c.cache.Add(randValue, int64(c.hashcashDuration))
	return randValue, err
}

func (c *server) send(msg []byte, conn io.Writer) error {
	c.log.Println(string(msg))

	_, err := conn.Write(msg)
	return err
}
