package util

import (
	"github.com/cheggaaa/pb/v3"
)

var format = pb.New(0).Set(pb.Bytes, true).Set(pb.SIBytesPrefix, true)

func Bytes(i int64) string {
	return format.Format(i)
}

func BytesSec(i int64) string {
	return Bytes(i) + "/s"
}
