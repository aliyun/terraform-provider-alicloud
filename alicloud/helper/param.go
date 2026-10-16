package helper

import (
	"encoding/xml"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// OpenAPI param constructors and the XML map marshaler moved from the
// root package (alicloud/common.go); package alicloud re-exports them
// as var/type aliases so the flat layout keeps compiling unchanged.

func RpcParam(method, version, action string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Pathname:    tea.String("/"),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
}
func RoaParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
}

func XmlParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("xml"),
		BodyType:    tea.String("xml"),
	}
}

func JsonXmlParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("xml"),
	}
}

func XmlJsonParam(method, version, action, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String(path),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		ReqBodyType: tea.String("xml"),
		BodyType:    tea.String("json"),
	}
}

type MyMap map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

func (m MyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}
