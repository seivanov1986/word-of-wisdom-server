package vo

import (
	"fmt"
	"strconv"
)

const (
	Quit = iota
	RequestChallenge
	ResponseChallenge
	RequestResource
	ResponseResource
)

var (
	QuitHeader              = Header{Quit}
	RequestChallengeHeader  = Header{RequestChallenge}
	ResponseChallengeHeader = Header{ResponseChallenge}
	RequestResourceHeader   = Header{RequestResource}
	ResponseResourceHeader  = Header{ResponseResource}
)

type Header struct {
	value int64
}

func ParseHeader(value int64) (Header, error) {
	header := Header{}

	switch value {
	case Quit:
		header = QuitHeader
	case RequestChallenge:
		header = RequestChallengeHeader
	case ResponseChallenge:
		header = ResponseChallengeHeader
	case RequestResource:
		header = RequestResourceHeader
	case ResponseResource:
		header = ResponseResourceHeader
	default:
		return header, fmt.Errorf("header not found")
	}

	return header, nil
}

func (h Header) String() string {
	return strconv.FormatInt(h.value, 10)
}
