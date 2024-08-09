package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/golang/protobuf/proto"
	"strings"
)

type HighlightEncoder int32

const (
	PlainMode HighlightEncoder = 1
	HtmlMode  HighlightEncoder = 2
)

func (encoder HighlightEncoder) Enum() *HighlightEncoder {
	p := new(HighlightEncoder)
	*p = encoder
	return p
}

func (encoder HighlightEncoder) String() string {
	switch encoder {
	case PlainMode:
		return "plain_mode"
	case HtmlMode:
		return "html_mode"
	default:
		return fmt.Sprintf("%#v", encoder)
	}
}

func (encoder *HighlightEncoder) ProtoBuffer() (*otsprotocol.HighlightEncoder, error) {
	if encoder == nil {
		return otsprotocol.HighlightEncoder_PLAIN_MODE.Enum(), nil
	}

	switch *encoder {
	case PlainMode:
		return otsprotocol.HighlightEncoder_PLAIN_MODE.Enum(), nil
	case HtmlMode:
		return otsprotocol.HighlightEncoder_HTML_MODE.Enum(), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown highlight encoder type: %s", encoder.String()))
	}
}

func (encoder *HighlightEncoder) UnmarshalJSON(data []byte) (err error) {
	var encoderStr string
	err = json.Unmarshal(data, &encoderStr)
	if err != nil {
		return
	}

	*encoder, err = ToHighlightEncoder(encoderStr)
	return
}

func (encoder *HighlightEncoder) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(encoder.String())
	return
}

func ToHighlightEncoder(encoderStr string) (HighlightEncoder, error) {
	switch strings.ToUpper(encoderStr) {
	case "PLAIN_MODE":
		return PlainMode, nil
	case "HTML_MODE":
		return HtmlMode, nil
	default:
		return PlainMode, errors.New(fmt.Sprintf("unknown highlight encoder type string: %#v", encoderStr))
	}
}

type HighlightFragmentOrder int32

const (
	TextSequence HighlightFragmentOrder = 1
	Score        HighlightFragmentOrder = 2
)

func (order HighlightFragmentOrder) Enum() *HighlightFragmentOrder {
	p := new(HighlightFragmentOrder)
	*p = order
	return p
}

func (order HighlightFragmentOrder) String() string {
	switch order {
	case TextSequence:
		return "text_sequence"
	case Score:
		return "score"
	default:
		return fmt.Sprintf("%#v", order)
	}
}

func (order *HighlightFragmentOrder) ProtoBuffer() (*otsprotocol.HighlightFragmentOrder, error) {
	if order == nil {
		return otsprotocol.HighlightFragmentOrder_TEXT_SEQUENCE.Enum(), nil
	}

	switch *order {
	case TextSequence:
		return otsprotocol.HighlightFragmentOrder_TEXT_SEQUENCE.Enum(), nil
	case Score:
		return otsprotocol.HighlightFragmentOrder_SCORE.Enum(), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown highlight fragment order: %s", order.String()))
	}
}

func (order *HighlightFragmentOrder) UnmarshalJSON(data []byte) (err error) {
	var orderStr string
	err = json.Unmarshal(data, &orderStr)
	if err != nil {
		return
	}

	*order, err = ToHighlightFragmentOrder(orderStr)
	return
}

func (order *HighlightFragmentOrder) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(order.String())
	return
}

func ToHighlightFragmentOrder(orderStr string) (HighlightFragmentOrder, error) {
	switch strings.ToUpper(orderStr) {
	case "TEXT_SEQUENCE":
		return TextSequence, nil
	case "SCORE":
		return Score, nil
	default:
		return TextSequence, errors.New(fmt.Sprintf("unknown highlight fragment order string: %s", orderStr))
	}
}

type Highlight struct {
	HighlightEncoder         *HighlightEncoder
	FieldHighlightParameters map[string]*HighlightParameter
}

func NewHighlight() *Highlight {
	return &Highlight{
		HighlightEncoder:         PlainMode.Enum(),
		FieldHighlightParameters: map[string]*HighlightParameter{},
	}
}

func (highlight *Highlight) SetHighlightEncoder(encoder HighlightEncoder) *Highlight {
	highlight.HighlightEncoder = encoder.Enum()
	return highlight
}

func (highlight *Highlight) AddFieldHighlightParameter(fieldName string, param *HighlightParameter) *Highlight {
	if highlight.FieldHighlightParameters == nil {
		highlight.FieldHighlightParameters = map[string]*HighlightParameter{}
	}
	highlight.FieldHighlightParameters[fieldName] = param
	return highlight
}

func (highlight *Highlight) SetFieldHighlightParameters(fieldHighlightParameters map[string]*HighlightParameter) *Highlight {
	highlight.FieldHighlightParameters = fieldHighlightParameters
	return highlight
}

func (highlight *Highlight) ProtoBuffer() (*otsprotocol.Highlight, error) {
	pbHighlight := &otsprotocol.Highlight{}

	if highlight == nil {
		return pbHighlight, nil
	}

	if pbEncoder, err := highlight.HighlightEncoder.ProtoBuffer(); err != nil {
		return nil, err
	} else {
		pbHighlight.HighlightEncoder = pbEncoder
	}

	if highlight.FieldHighlightParameters != nil && len(highlight.FieldHighlightParameters) != 0 {
		for fieldName, param := range highlight.FieldHighlightParameters {
			if pbHighlightParameter, err := param.ProtoBuffer(); err != nil {
				return nil, err
			} else {
				pbHighlightParameter.FieldName = proto.String(fieldName)
				pbHighlight.HighlightParameters = append(pbHighlight.HighlightParameters, pbHighlightParameter)
			}
		}
	}

	return pbHighlight, nil
}

func ToHighlight(pbHighlight *otsprotocol.Highlight) (*Highlight, error) {
	highlight := &Highlight{
		HighlightEncoder:         PlainMode.Enum(),
		FieldHighlightParameters: map[string]*HighlightParameter{},
	}

	if pbHighlight.HighlightEncoder != nil {
		if highlightEncoder, err := ToHighlightEncoder(pbHighlight.HighlightEncoder.String()); err != nil {
			return highlight, err
		} else {
			highlight.HighlightEncoder = &highlightEncoder
		}
	}

	if pbHighlight.HighlightParameters == nil {
		return highlight, nil
	}

	for _, pbHighlightParameter := range pbHighlight.GetHighlightParameters() {
		if highlightParameter, err := ToHighlightParameter(pbHighlightParameter); err != nil {
			return highlight, err
		} else {
			highlight.FieldHighlightParameters[pbHighlightParameter.GetFieldName()] = highlightParameter
		}
	}

	return highlight, nil
}

type HighlightParameter struct {
	NumberOfFragments      *int32
	FragmentSize           *int32
	PreTag                 *string
	PostTag                *string
	HighlightFragmentOrder *HighlightFragmentOrder
}

func NewHighlightParameter() *HighlightParameter {
	return &HighlightParameter{
		HighlightFragmentOrder: TextSequence.Enum(),
	}
}

func (param *HighlightParameter) SetNumberOfFragments(numberOfFragments int32) *HighlightParameter {
	param.NumberOfFragments = &numberOfFragments
	return param
}

func (param *HighlightParameter) SetFragmentSize(fragmentSize int32) *HighlightParameter {
	param.FragmentSize = &fragmentSize
	return param
}

func (param *HighlightParameter) SetPreTag(preTag string) *HighlightParameter {
	param.PreTag = &preTag
	return param
}

func (param *HighlightParameter) SetPostTag(postTag string) *HighlightParameter {
	param.PostTag = &postTag
	return param
}

func (param *HighlightParameter) SetHighlightFragmentOrder(order HighlightFragmentOrder) *HighlightParameter {
	param.HighlightFragmentOrder = &order
	return param
}

func (param *HighlightParameter) ProtoBuffer() (*otsprotocol.HighlightParameter, error) {
	pbHighlightParameter := &otsprotocol.HighlightParameter{}
	if param == nil {
		return pbHighlightParameter, nil
	}

	if param.PreTag != nil {
		pbHighlightParameter.PreTag = proto.String(*param.PreTag)
	}

	if param.PostTag != nil {
		pbHighlightParameter.PostTag = proto.String(*param.PostTag)
	}

	if param.FragmentSize != nil {
		pbHighlightParameter.FragmentSize = proto.Int32(*param.FragmentSize)
	}

	if param.NumberOfFragments != nil {
		pbHighlightParameter.NumberOfFragments = proto.Int32(*param.NumberOfFragments)
	}

	if pbFragmentOrder, err := param.HighlightFragmentOrder.ProtoBuffer(); err != nil {
		return nil, err
	} else {
		pbHighlightParameter.FragmentsOrder = pbFragmentOrder
	}

	return pbHighlightParameter, nil
}

func ToHighlightParameter(pbParameter *otsprotocol.HighlightParameter) (param *HighlightParameter, err error) {
	param = &HighlightParameter{
		HighlightFragmentOrder: TextSequence.Enum(),
	}

	if pbParameter.PreTag != nil {
		param.PreTag = pbParameter.PreTag
	}

	if pbParameter.PostTag != nil {
		param.PostTag = pbParameter.PostTag
	}

	if pbParameter.NumberOfFragments != nil {
		param.NumberOfFragments = pbParameter.NumberOfFragments
	}

	if pbParameter.FragmentsOrder != nil {
		param.FragmentSize = pbParameter.FragmentSize
	}

	if pbParameter.FragmentsOrder != nil {
		var highlightFragmentOrder HighlightFragmentOrder
		if highlightFragmentOrder, err = ToHighlightFragmentOrder(pbParameter.FragmentsOrder.String()); err != nil {
			return
		}
		param.HighlightFragmentOrder = &highlightFragmentOrder
	}

	return
}
