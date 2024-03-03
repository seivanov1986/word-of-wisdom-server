package proto

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/seivanov1986/word-of-wisdom-server/internal/vo"
)

type Message struct {
	Header  vo.Header
	Payload vo.Payload
}

func MessageToBytes(msg Message) []byte {
	format := fmt.Sprintf("%v|%v", msg.Header.String(), msg.Payload)
	msgStr := fmt.Sprintf("%v\n", format)
	return []byte(msgStr)
}

func Parse(in string) (*Message, error) {
	in = strings.TrimSpace(in)
	var msgType int

	parts := strings.Split(in, "|")
	if len(parts) < 1 || len(parts) > 2 { //only 1 or 2 parts allowed
		return nil, fmt.Errorf("message doesn't match protocol")
	}

	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}

	header, err := vo.ParseHeader(int64(msgType))
	if err != nil {
		return nil, err
	}

	msg := Message{
		Header: header,
	}

	payload, err := vo.ParsePayload(parts[1])
	if err != nil {
		return nil, err
	}

	if len(parts) == 2 {
		msg.Payload = payload
	}
	return &msg, nil
}
