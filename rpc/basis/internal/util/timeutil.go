package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
