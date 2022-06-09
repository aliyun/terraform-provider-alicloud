package protocol

import "fmt"

type RecordSequenceInfo struct {
	Epoch     int32
	Timestamp int64
	RowIndex  int32
}

func (this *RecordSequenceInfo) String() string {
	return fmt.Sprintf(
		"{\"Epoch\":%d, \"Timestamp\": %d, \"RowIndex\": %d}",
		this.Epoch,
		this.Timestamp,
		this.RowIndex)
}