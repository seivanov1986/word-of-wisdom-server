package vo

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	value string
}

func ParsePayload(value string) (Payload, error) {
	payload := Payload{value}

	if value == "" {
		return payload, nil
	}

	if !json.Valid([]byte(value)) {
		fmt.Println(value)
		return payload, fmt.Errorf("payload error")
	}

	return payload, nil
}

func (p Payload) String() string {
	return p.value
}
