package helper

import (
	"encoding/xml"
	"strings"
	"testing"
)

// Tests for the OpenAPI param constructors and the XML map marshaler
// moved from the root package (alicloud/common.go); test names are
// unchanged from alicloud/common_test.go.

func TestUnitCommonRpcParam(t *testing.T) {
	param := RpcParam("POST", "2020-01-01", "CreateInstance")

	if *param.Method != "POST" {
		t.Error("Method not set correctly")
	}
	if *param.Version != "2020-01-01" {
		t.Error("Version not set correctly")
	}
	if *param.Action != "CreateInstance" {
		t.Error("Action not set correctly")
	}
	if *param.Style != "RPC" {
		t.Error("Style should be RPC")
	}
}

func TestUnitCommonRoaParam(t *testing.T) {
	param := RoaParam("GET", "2020-01-01", "ListInstances", "/instances")

	if *param.Method != "GET" {
		t.Error("Method not set correctly")
	}
	if *param.Version != "2020-01-01" {
		t.Error("Version not set correctly")
	}
	if *param.Action != "ListInstances" {
		t.Error("Action not set correctly")
	}
	if *param.Pathname != "/instances" {
		t.Error("Pathname not set correctly")
	}
	if *param.Style != "ROA" {
		t.Error("Style should be ROA")
	}
}

func TestUnitCommonXmlParam(t *testing.T) {
	param := XmlParam("PUT", "2020-01-01", "UpdateConfig", "/config")

	if *param.Method != "PUT" {
		t.Error("Method not set correctly")
	}
	if *param.ReqBodyType != "xml" {
		t.Error("ReqBodyType should be xml")
	}
	if *param.BodyType != "xml" {
		t.Error("BodyType should be xml")
	}
}

func TestUnitCommonJsonXmlParam(t *testing.T) {
	param := JsonXmlParam("POST", "2020-01-01", "CreateResource", "/resource")

	if *param.ReqBodyType != "json" {
		t.Error("ReqBodyType should be json")
	}
	if *param.BodyType != "xml" {
		t.Error("BodyType should be xml")
	}
}

func TestUnitCommonXmlJsonParam(t *testing.T) {
	param := XmlJsonParam("POST", "2020-01-01", "CreateResource", "/resource")

	if *param.ReqBodyType != "xml" {
		t.Error("ReqBodyType should be xml")
	}
	if *param.BodyType != "json" {
		t.Error("BodyType should be json")
	}
}

func TestUnitCommonMyMapMarshalXML(t *testing.T) {
	// Test empty map - should return nil error for empty map
	emptyMap := MyMap{}
	var emptyResult strings.Builder
	emptyEncoder := xml.NewEncoder(&emptyResult)
	err := emptyMap.MarshalXML(emptyEncoder, xml.StartElement{Name: xml.Name{Local: "root"}})
	if err != nil {
		t.Errorf("MarshalXML failed for empty map: %v", err)
	}

	// Test non-empty map
	m := MyMap{
		"key1": "value1",
		"key2": "value2",
	}
	var result strings.Builder
	encoder := xml.NewEncoder(&result)
	err = m.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "root"}})
	if err != nil {
		t.Errorf("MarshalXML failed: %v", err)
	}
	encoder.Flush()
	xmlStr := result.String()
	if !strings.Contains(xmlStr, "key1") || !strings.Contains(xmlStr, "value1") {
		t.Errorf("MarshalXML output doesn't contain expected keys: %s", xmlStr)
	}
}
