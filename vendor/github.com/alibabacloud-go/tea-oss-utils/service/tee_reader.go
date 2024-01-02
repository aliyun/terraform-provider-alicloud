package service

import (
	"io"

	"github.com/alibabacloud-go/tea/utils"
)

type teeReader struct {
	reader        io.Reader
	writer        io.Writer
	listener      utils.ProgressListener
	consumedBytes int64
	totalBytes    int64
	tracker       *utils.ReaderTracker
}

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w.  There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(reader io.Reader, writer io.Writer, totalBytes int64, listener utils.ProgressListener, tracker *utils.ReaderTracker) io.ReadCloser {
	return &teeReader{
		reader:        reader,
		writer:        writer,
		listener:      listener,
		consumedBytes: 0,
		totalBytes:    totalBytes,
		tracker:       tracker,
	}
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)

	// Read encountered error
	if err != nil && err != io.EOF {
		event := utils.NewProgressEvent(utils.TransferFailedEvent, t.consumedBytes, t.totalBytes, 0)
		utils.PublishProgress(t.listener, event)
	}

	if n > 0 {
		t.consumedBytes += int64(n)
		// CRC
		if t.writer != nil {
			if n, err := t.writer.Write(p[:n]); err != nil {
				return n, err
			}
		}
		// Progress
		if t.listener != nil {
			event := utils.NewProgressEvent(utils.TransferDataEvent, t.consumedBytes, t.totalBytes, int64(n))
			utils.PublishProgress(t.listener, event)
		}
		// Track
		if t.tracker != nil {
			t.tracker.CompletedBytes = t.consumedBytes
		}
	}

	return
}

func (t *teeReader) Close() error {
	if rc, ok := t.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}
