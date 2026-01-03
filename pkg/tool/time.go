package tool

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ToProtoTimestamp Time -> Proto Timestamp
func ToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
