package ots

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	otsUnitInstanceName = "terraform-ex"
	otsUnitTableName    = "terraform_example"
)

// otsUnitWire captures the serialized OTS requests and replays canned responses so the stream
// contract can be asserted at the protocol payload level without a cloud account.
type otsUnitWire struct {
	mu      sync.Mutex
	bodies  map[string][][]byte
	respond func(path string) []byte
}

func (w *otsUnitWire) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	w.mu.Lock()
	w.bodies[req.URL.Path] = append(w.bodies[req.URL.Path], body)
	w.mu.Unlock()

	var payload []byte
	if w.respond != nil {
		payload = w.respond(req.URL.Path)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader(payload)),
		Request:    req,
	}, nil
}

func (w *otsUnitWire) last(t *testing.T, path string) []byte {
	t.Helper()
	w.mu.Lock()
	defer w.mu.Unlock()
	bodies := w.bodies[path]
	if len(bodies) == 0 {
		t.Fatalf("no OTS request captured for %s", path)
	}
	return bodies[len(bodies)-1]
}

func newOtsUnitClient(wire *otsUnitWire) *tablestore.TableStoreClient {
	config := tablestore.NewDefaultTableStoreConfig()
	config.Transport = wire
	config.RetryTimes = 0
	config.MaxRetryTime = 0
	return tablestore.NewClientWithConfig("http://ots-unit.invalid", otsUnitInstanceName, "unit-ak", "unit-sk", "", config)
}

func otsUnitTableMeta() *tablestore.TableMeta {
	meta := new(tablestore.TableMeta)
	meta.TableName = otsUnitTableName
	meta.AddPrimaryKeyColumn("pk1", tablestore.PrimaryKeyType_INTEGER)
	return meta
}

func otsUnitTableOptions() *otsprotocol.TableOptions {
	return &otsprotocol.TableOptions{
		TimeToLive:                proto.Int32(-1),
		MaxVersions:               proto.Int32(1),
		DeviationCellVersionInSec: proto.Int64(86400),
	}
}

func otsUnitReservedThroughputDetails() *otsprotocol.ReservedThroughputDetails {
	return &otsprotocol.ReservedThroughputDetails{
		CapacityUnit:     &otsprotocol.CapacityUnit{Read: proto.Int32(0), Write: proto.Int32(0)},
		LastIncreaseTime: proto.Int64(0),
	}
}

func otsUnitDescribeResponse(enableStream bool, expirationTime int32) *otsprotocol.DescribeTableResponse {
	primaryKeyType := otsprotocol.PrimaryKeyType_INTEGER
	resp := &otsprotocol.DescribeTableResponse{
		TableMeta: &otsprotocol.TableMeta{
			TableName:  proto.String(otsUnitTableName),
			PrimaryKey: []*otsprotocol.PrimaryKeySchema{{Name: proto.String("pk1"), Type: &primaryKeyType}},
		},
		ReservedThroughputDetails: otsUnitReservedThroughputDetails(),
		TableOptions:              otsUnitTableOptions(),
		TableStatus:               otsprotocol.TableStatus_ACTIVE.Enum(),
		StreamDetails:             &otsprotocol.StreamDetails{EnableStream: proto.Bool(enableStream)},
	}
	if enableStream {
		resp.StreamDetails.StreamId = proto.String("unit-stream-id")
		resp.StreamDetails.ExpirationTime = proto.Int32(expirationTime)
		resp.StreamDetails.LastEnableTime = proto.Int64(1700000000)
	}
	return resp
}

func otsUnitUpdateResponse(enableStream bool, expirationTime int32) *otsprotocol.UpdateTableResponse {
	details := &otsprotocol.StreamDetails{EnableStream: proto.Bool(enableStream)}
	if enableStream {
		details.StreamId = proto.String("unit-stream-id")
		details.ExpirationTime = proto.Int32(expirationTime)
		details.LastEnableTime = proto.Int64(1700000000)
	}
	return &otsprotocol.UpdateTableResponse{
		ReservedThroughputDetails: otsUnitReservedThroughputDetails(),
		TableOptions:              otsUnitTableOptions(),
		StreamDetails:             details,
	}
}

func otsUnitStreamSpecs(enableStream bool, expirationTime int) []interface{} {
	spec := map[string]interface{}{"enable_stream": enableStream}
	if enableStream {
		spec["expiration_time"] = expirationTime
	}
	return []interface{}{spec}
}

func marshalOtsUnitResponse(t *testing.T, message proto.Message) []byte {
	t.Helper()
	payload, err := proto.Marshal(message)
	if err != nil {
		t.Fatalf("marshal canned OTS response: %s", err)
	}
	return payload
}

func TestUnitAliCloudOtsTableStreamSpecificationWire(t *testing.T) {
	t.Run("create sends stream target values", func(t *testing.T) {
		wire := &otsUnitWire{bodies: map[string][][]byte{}}
		request := &tablestore.CreateTableRequest{
			TableMeta:          otsUnitTableMeta(),
			TableOption:        &tablestore.TableOption{TimeToAlive: -1, MaxVersion: 1},
			ReservedThroughput: new(tablestore.ReservedThroughput),
			StreamSpec:         expandOtsStreamSpecification(otsUnitStreamSpecs(true, 24)),
		}
		if _, err := newOtsUnitClient(wire).CreateTable(request); err != nil {
			t.Fatalf("CreateTable: %s", err)
		}

		payload := new(otsprotocol.CreateTableRequest)
		if err := proto.Unmarshal(wire.last(t, "/CreateTable"), payload); err != nil {
			t.Fatalf("decode CreateTable request: %s", err)
		}
		if payload.GetStreamSpec() == nil {
			t.Fatalf("CreateTable payload has no stream_spec")
		}
		if !payload.GetStreamSpec().GetEnableStream() {
			t.Fatalf("CreateTable payload enable_stream = false, want true")
		}
		if got := payload.GetStreamSpec().GetExpirationTime(); got != 24 {
			t.Fatalf("CreateTable payload expiration_time = %d, want 24", got)
		}
	})

	t.Run("update sends new retention", func(t *testing.T) {
		wire := &otsUnitWire{bodies: map[string][][]byte{}}
		wire.respond = func(path string) []byte {
			if path == "/UpdateTable" {
				return marshalOtsUnitResponse(t, otsUnitUpdateResponse(true, 168))
			}
			return nil
		}
		request := &tablestore.UpdateTableRequest{
			TableName:  otsUnitTableName,
			StreamSpec: expandOtsStreamSpecification(otsUnitStreamSpecs(true, 168)),
		}
		if _, err := newOtsUnitClient(wire).UpdateTable(request); err != nil {
			t.Fatalf("UpdateTable: %s", err)
		}

		payload := new(otsprotocol.UpdateTableRequest)
		if err := proto.Unmarshal(wire.last(t, "/UpdateTable"), payload); err != nil {
			t.Fatalf("decode UpdateTable request: %s", err)
		}
		if payload.GetStreamSpec() == nil {
			t.Fatalf("UpdateTable payload has no stream_spec")
		}
		if !payload.GetStreamSpec().GetEnableStream() {
			t.Fatalf("UpdateTable payload enable_stream = false, want true")
		}
		if got := payload.GetStreamSpec().GetExpirationTime(); got != 168 {
			t.Fatalf("UpdateTable payload expiration_time = %d, want 168", got)
		}
	})

	t.Run("update disables stream without retention", func(t *testing.T) {
		wire := &otsUnitWire{bodies: map[string][][]byte{}}
		wire.respond = func(path string) []byte {
			if path == "/UpdateTable" {
				return marshalOtsUnitResponse(t, otsUnitUpdateResponse(false, 0))
			}
			return nil
		}
		request := &tablestore.UpdateTableRequest{
			TableName:  otsUnitTableName,
			StreamSpec: expandOtsStreamSpecification(otsUnitStreamSpecs(false, 0)),
		}
		if _, err := newOtsUnitClient(wire).UpdateTable(request); err != nil {
			t.Fatalf("UpdateTable: %s", err)
		}

		payload := new(otsprotocol.UpdateTableRequest)
		if err := proto.Unmarshal(wire.last(t, "/UpdateTable"), payload); err != nil {
			t.Fatalf("decode UpdateTable request: %s", err)
		}
		if payload.GetStreamSpec() == nil {
			t.Fatalf("UpdateTable payload has no stream_spec")
		}
		if payload.GetStreamSpec().GetEnableStream() {
			t.Fatalf("UpdateTable payload enable_stream = true, want false")
		}
		if payload.GetStreamSpec().ExpirationTime != nil {
			t.Fatalf("UpdateTable payload expiration_time = %d, want omitted", payload.GetStreamSpec().GetExpirationTime())
		}
	})

	t.Run("describe readback flows into state", func(t *testing.T) {
		wire := &otsUnitWire{bodies: map[string][][]byte{}}
		wire.respond = func(path string) []byte {
			if path == "/DescribeTable" {
				return marshalOtsUnitResponse(t, otsUnitDescribeResponse(true, 24))
			}
			return nil
		}
		response, err := newOtsUnitClient(wire).DescribeTable(&tablestore.DescribeTableRequest{TableName: otsUnitTableName})
		if err != nil {
			t.Fatalf("DescribeTable: %s", err)
		}
		if response.StreamDetails == nil || !response.StreamDetails.EnableStream {
			t.Fatalf("DescribeTable stream details = %#v, want enabled", response.StreamDetails)
		}
		if got := response.StreamDetails.ExpirationTime; got != 24 {
			t.Fatalf("DescribeTable expiration_time = %d, want 24", got)
		}

		d := schema.TestResourceDataRaw(t, ResourceOtsTable().Schema, map[string]interface{}{
			"instance_name": otsUnitInstanceName,
			"table_name":    otsUnitTableName,
			"primary_key": []interface{}{
				map[string]interface{}{"name": "pk1", "type": "Integer"},
			},
			"time_to_live": -1,
			"max_version":  1,
		})
		setOtsStreamSpecification(d, response.StreamDetails)
		specs, ok := d.Get("stream_specification").([]interface{})
		if !ok || len(specs) != 1 {
			t.Fatalf("stream_specification = %#v, want one block", d.Get("stream_specification"))
		}
		spec := specs[0].(map[string]interface{})
		if enabled, _ := spec["enable_stream"].(bool); !enabled {
			t.Fatalf("state enable_stream = %v, want true", spec["enable_stream"])
		}
		if got, _ := spec["expiration_time"].(int); got != 24 {
			t.Fatalf("state expiration_time = %v, want 24", spec["expiration_time"])
		}
	})

	t.Run("describe readback clears disabled stream", func(t *testing.T) {
		wire := &otsUnitWire{bodies: map[string][][]byte{}}
		wire.respond = func(path string) []byte {
			if path == "/DescribeTable" {
				return marshalOtsUnitResponse(t, otsUnitDescribeResponse(false, 0))
			}
			return nil
		}
		response, err := newOtsUnitClient(wire).DescribeTable(&tablestore.DescribeTableRequest{TableName: otsUnitTableName})
		if err != nil {
			t.Fatalf("DescribeTable: %s", err)
		}

		d := schema.TestResourceDataRaw(t, ResourceOtsTable().Schema, map[string]interface{}{
			"instance_name": otsUnitInstanceName,
			"table_name":    otsUnitTableName,
			"primary_key": []interface{}{
				map[string]interface{}{"name": "pk1", "type": "Integer"},
			},
			"time_to_live": -1,
			"max_version":  1,
		})
		setOtsStreamSpecification(d, response.StreamDetails)
		if specs, _ := d.Get("stream_specification").([]interface{}); len(specs) != 0 {
			t.Fatalf("stream_specification = %#v, want empty", specs)
		}
	})
}
