package utility

import (
	"bytes"
	"encoding/gob"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
)

// ParseIntFromString : returns int64 from string, returns 0 for invalid string
func ParseIntFromString(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println("error parsing data")
		return 0
	}
	return i
}

func CovertObjToBytes(obj interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(obj)
	return buf.Bytes()
}

func ConvertPlanObjToBytes(obj interface{}) ([]byte, error) {
	bt, err := kafka.Marshal(obj)
	return bt, err
}
