// This file is auto-generated, don't edit it. Thanks.
/**
 * Read data from a readable stream, and parse it by JSON format
 * @param stream the readable stream
 * @return the parsed result
 */
package client

import (
	"bytes"
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/pierrec/lz4"
	"io"
	"strconv"
)

func ReadAndUncompressBlock(stream io.Reader, compressType *string, bodyRawSize *string) (_result io.Reader, _err error) {
	if *compressType != "lz4" {
		return nil, fmt.Errorf("unsupported compress type %s", *compressType)
	}
	rawSize, err := strconv.ParseInt(*bodyRawSize, 10, 64)
	if err != nil {
		return nil, err
	}
	out := make([]byte, rawSize)
	if rawSize != 0 {
		body, _ := util.ReadAsBytes(stream)
		len, err := lz4.UncompressBlock(body, out)
		if err != nil || int64(len) != rawSize {
			return nil, err
		}
	}
	return bytes.NewReader(out), nil
}
