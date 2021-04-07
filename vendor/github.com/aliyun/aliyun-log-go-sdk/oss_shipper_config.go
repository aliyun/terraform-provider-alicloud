package sls

import (
	"encoding/json"
)

const (
	OSSShipperType = "oss"
)

type Shipper struct {
	ShipperName            string          `json:"shipperName"`
	TargetType             string          `json:"targetType"`
	RawTargetConfiguration json.RawMessage `json:"targetConfiguration"`

	TargetConfiguration interface{} `json:"-"`
}

type shipperAlias Shipper

type shipperDisplay struct {
	ShipperName         string      `json:"shipperName"`
	TargetType          string      `json:"targetType"`
	TargetConfiguration interface{} `json:"targetConfiguration"`
}

type OssStoreageCsvDetail struct {
	Delimiter      string   `json:"delemiter"`
	Header         bool     `json:"header"`
	LineFeed       string   `json:"lineFeed"`
	Columns        []string `json:"columns"`
	NullIdentifier string   `json:"nullIdentfifier"`
	Quote          string   `json:"quote"`
}

type ParquetConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OssStoreageParquet struct {
	Columns []ParquetConfig `json:"columns"`
}
type OSSShipperConfig struct {
	OssBucket      string         `json:"ossBucket"`
	OssPrefix      string         `json:"ossPrefix"`
	RoleArn        string         `json:"roleArn"`
	BufferInterval int            `json:"bufferInterval"`
	BufferSize     int            `json:"bufferSize"`
	CompressType   string         `json:"compressType"`
	PathFormat     string         `json:"pathFormat"`
	Format         string         `json:"format"`
	Storage        ShipperStorage `json:"storage"`
}

type ShipperStorage struct {
	Format string      `json:"format"`
	Detail interface{} `json:"detail"`
}

type OssStorageJsonDetail struct {
	EnableTag bool `json:"enableTag"`
}
