package service

import (
	"crypto/md5"
	"encoding/base64"
	"hash"
	"io"
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
)

type complexReader struct {
	reader io.Reader
	crc    *digest
	md5    hash.Hash
	refer  map[string]*string
}

func ComplexReader(reader io.Reader, refer map[string]*string) io.ReadCloser {
	return &complexReader{
		reader: reader,
		crc:    NewCRC(crcTable(), 0),
		md5:    md5.New(),
		refer:  refer,
	}
}

func (c *complexReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)

	if err != nil && err == io.EOF {
		c.refer["md5"] = tea.String(base64.StdEncoding.EncodeToString(c.md5.Sum(nil)))
		c.refer["crc"] = tea.String(strconv.FormatUint(c.crc.Sum64(), 10))
	}

	if n > 0 {
		// CRC
		if c.crc != nil {
			if n, err := c.crc.Write(p[:n]); err != nil {
				return n, err
			}
		}

		// MD5
		if c.md5 != nil {
			io.WriteString(c.md5, string(p[:n]))
		}
	}

	return
}

func (c *complexReader) Close() error {
	if rc, ok := c.reader.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}
