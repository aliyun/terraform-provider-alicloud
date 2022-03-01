package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

const (
	DataSourceOSS        DataSourceType = "AliyunOSS"
	DataSourceBSS        DataSourceType = "AliyunBSS"
	DataSourceMaxCompute DataSourceType = "AliyunMaxCompute"
	DataSourceJDBC       DataSourceType = "JDBC"
	DataSourceKafka      DataSourceType = "Kafka"
	DataSourceCMS        DataSourceType = "AliyunCloudMonitor"
	DataSourceGeneral    DataSourceType = "General"

	OSSDataFormatTypeLine          OSSDataFormatType = "Line"
	OSSDataFormatTypeMultiline     OSSDataFormatType = "Multiline"
	OSSDataFormatTypeJSON          OSSDataFormatType = "JSON"
	OSSDataFormatTypeParquet       OSSDataFormatType = "Parquet"
	OSSDataFormatTypeDelimitedText OSSDataFormatType = "DelimitedText"

	KafkaValueTypeText KafkaValueType = "Text"
	KafkaValueTypeJSON KafkaValueType = "JSON"

	KafkaPositionGroupOffsets KafkaPosition = "GROUP_OFFSETS"
	KafkaPositionEarliest     KafkaPosition = "EARLIEST"
	KafkaPositionLatest       KafkaPosition = "LATEST"
	KafkaPositionTimeStamp    KafkaPosition = "TIMESTAMP"
)

type (
	BaseJob struct {
		Name           string  `json:"name"`
		DisplayName    string  `json:"displayName,omitempty"`
		Description    string  `json:"description,omitempty"`
		Type           JobType `json:"type"`
		Recyclable     bool    `json:"recyclable"`
		CreateTime     int64   `json:"createTime"`
		LastModifyTime int64   `json:"lastModifyTime"`
	}

	ScheduledJob struct {
		BaseJob
		Status     string    `json:"status"`
		Schedule   *Schedule `json:"schedule"`
		ScheduleId string    `json:"scheduleId"`
	}

	Ingestion struct {
		ScheduledJob
		IngestionConfiguration *IngestionConfiguration `json:"configuration"`
	}

	IngestionConfiguration struct {
		Version          string      `json:"version"`
		LogStore         string      `json:"logstore"`
		NumberOfInstance int32       `json:"numberOfInstance"`
		DataSource       interface{} `json:"source"`
	}

	DataSourceType string

	DataSource struct {
		DataSourceType DataSourceType `json:"type"`
	}

	// >>> ingestion oss source
	OSSDataFormatType string

	AliyunOSSSource struct {
		DataSource
		Bucket                  string      `json:"bucket"`
		Endpoint                string      `json:"endpoint"`
		RoleArn                 string      `json:"roleARN"`
		Prefix                  string      `json:"prefix,omitempty"`
		Pattern                 string      `json:"pattern,omitempty"`
		CompressionCodec        string      `json:"compressionCodec,omitempty"`
		Encoding                string      `json:"encoding,omitempty"`
		Format                  interface{} `json:"format,omitempty"`
		RestoreObjectEnable     bool        `json:"restoreObjectEnable"`
		LastModifyTimeAsLogTime bool        `json:"lastModifyTimeAsLogTime"`
	}

	OSSDataFormat struct {
		Type       OSSDataFormatType `json:"type"`
		TimeFormat string            `json:"timeFormat"`
		TimeZone   string            `json:"timeZone"`
	}

	LineFormat struct {
		OSSDataFormat
		TimePattern string `json:"timePattern"`
	}

	MultiLineFormat struct {
		LineFormat
		MaxLines     int64  `json:"maxLines,omitempty"`
		Negate       bool   `json:"negate"`
		Match        string `json:"match"`
		Pattern      string `json:"pattern"`
		FlushPattern string `json:"flushPattern"`
	}

	StructureDataFormat struct {
		OSSDataFormat
		TimeField string `json:"timeField"`
	}

	JSONFormat struct {
		StructureDataFormat
		SkipInvalidRows bool `json:"skipInvalidRows"`
	}

	ParquetFormat struct {
		StructureDataFormat
	}

	DelimitedTextFormat struct {
		StructureDataFormat
		FieldNames       []string `json:"fieldNames"`
		FieldDelimiter   string   `json:"fieldDelimiter"`
		QuoteChar        string   `json:"quoteChar"`
		EscapeChar       string   `json:"escapeChar"`
		SkipLeadingRows  int64    `json:"skipLeadingRows"`
		MaxLines         int64    `json:"maxLines"`
		FirstRowAsHeader bool     `json:"firstRowAsHeader"`
	}

	// ingestion maxcompute source >>>
	AliyunMaxComputeSource struct {
		DataSource
		AccessKeyID     string `json:"accessKeyID"`
		AccessKeySecret string `json:"accessKeySecret"`
		Endpoint        string `json:"endpoint"`
		TunnelEndpoint  string `json:"tunnelEndpoint,omitempty"`
		Project         string `json:"project"`
		Table           string `json:"table"`
		PartitionSpec   string `json:"partitionSpec"`
		TimeField       string `json:"timeField"`
		TimeFormat      string `json:"timeFormat"`
		TimeZone        string `json:"timeZone"`
	}

	// ingestion cloud monitor source
	AliyunCloudMonitorSource struct {
		DataSource
		AccessKeyID     string   `json:"accessKeyID"`
		AccessKeySecret string   `json:"accessKeySecret"`
		StartTime       int64    `json:"startTime"`
		Namespaces      []string `json:"namespaces"`
		OutputType      string   `json:"outputType"`
		DelayTime       int64    `json:"delayTime"`
	}

	// ingestion kafka source
	KafkaValueType string
	KafkaPosition  string
	KafkaSource    struct {
		DataSource
		Topics           string            `json:"topics"`
		BootStrapServers string            `json:"bootstrapServers"`
		ValueType        KafkaValueType    `json:"valueType"`
		FromPosition     KafkaPosition     `json:"fromPosition"`
		FromTimeStamp    int64             `json:"fromTimestamp"`
		ToTimeStamp      int64             `json:"toTimestamp"`
		TimeField        string            `json:"timeField"`
		TimePattern      string            `json:"timePattern"`
		TimeFormat       string            `json:"timeFormat"`
		TimeZone         string            `json:"timeZone"`
		AdditionalProps  map[string]string `json:"additionalProps"`
	}

	// ingestion JDBC source
	AliyunBssSource struct {
		DataSource
		RoleArn      string `json:"roleARN"`
		HistoryMonth int64  `json:"historyMonth"`
	}

	// ingestion general source
	IngestionGeneralSource struct {
		DataSource
		Fields map[string]interface{}
	}
)

func (c *Client) CreateIngestion(project string, ingestion *Ingestion) error {
	body, err := json.Marshal(ingestion)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}
	uri := "/jobs"
	r, err := c.request(project, "POST", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) UpdateIngestion(project string, ingestion *Ingestion) error {
	body, err := json.Marshal(ingestion)
	if err != nil {
		return NewClientError(err)
	}
	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + ingestion.Name
	r, err := c.request(project, "PUT", uri, h, body)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) GetIngestion(project string, name string) (*Ingestion, error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + name
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	buf, _ := ioutil.ReadAll(r.Body)
	ingestion := &Ingestion{}
	if err = json.Unmarshal(buf, ingestion); err != nil {
		err = NewClientError(err)
	}
	return ingestion, err
}

func (c *Client) ListIngestion(project, logstore, name, displayName string, offset, size int) (ingestions []*Ingestion, total, count int, error error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	v := url.Values{}
	v.Add("logstore", logstore)
	v.Add("jobName", name)
	if displayName != "" {
		v.Add("displayName", displayName)
	}
	v.Add("jobType", "Ingestion")
	v.Add("offset", fmt.Sprintf("%d", offset))
	v.Add("size", fmt.Sprintf("%d", size))
	uri := "/jobs?" + v.Encode()
	r, err := c.request(project, "GET", uri, h, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	defer r.Body.Close()
	type ingestionList struct {
		Total   int          `json:"total"`
		Count   int          `json:"count"`
		Results []*Ingestion `json:"results"`
	}
	buf, _ := ioutil.ReadAll(r.Body)
	is := &ingestionList{}
	if err = json.Unmarshal(buf, is); err != nil {
		err = NewClientError(err)
	}
	return is.Results, is.Total, is.Count, err
}

func (c *Client) DeleteIngestion(project string, name string) error {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Content-Type":      "application/json",
	}
	uri := "/jobs/" + name
	r, err := c.request(project, "DELETE", uri, h, nil)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}
