package search

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search/model"
	"github.com/golang/protobuf/proto"
)

type FunctionsScoreQuery struct {
	Query       Query            `json:"-"`
	Functions   []*ScoreFunction `json:"-"`
	ScoreMode   *ScoreMode       `json:"-"`
	CombineMode *CombineMode     `json:"-"`
	MinScore    *float32         `json:"-"`
	MaxScore    *float32         `json:"-"`
	// for json marshal and unmarshal
	QueryAlias queryAlias `json:"Query"`
}

func NewFunctionsScoreQuery() *FunctionsScoreQuery {
	return &FunctionsScoreQuery{}
}

func (q *FunctionsScoreQuery) SetQuery(query Query) *FunctionsScoreQuery {
	q.Query = query
	return q
}

func (q *FunctionsScoreQuery) SetFunctions(functions []*ScoreFunction) *FunctionsScoreQuery {
	q.Functions = functions
	return q
}

func (q *FunctionsScoreQuery) AddFunction(function *ScoreFunction) *FunctionsScoreQuery {
	q.Functions = append(q.Functions, function)
	return q
}

func (q *FunctionsScoreQuery) SetScoreMode(mode ScoreMode) *FunctionsScoreQuery {
	q.ScoreMode = &mode
	return q
}

func (q *FunctionsScoreQuery) SetCombineMode(mode CombineMode) *FunctionsScoreQuery {
	q.CombineMode = &mode
	return q
}

func (q *FunctionsScoreQuery) SetMinScore(minScore float32) *FunctionsScoreQuery {
	q.MinScore = &minScore
	return q
}

func (q *FunctionsScoreQuery) SetMaxScore(maxScore float32) *FunctionsScoreQuery {
	q.MaxScore = &maxScore
	return q
}

func (q *FunctionsScoreQuery) MarshalJSON() ([]byte, error) {
	type FunctionsScoreQueryAlias FunctionsScoreQuery
	bqAlias := FunctionsScoreQueryAlias(*q)
	if bqAlias.Query != nil {
		bqAlias.QueryAlias = queryAlias{
			Name:  q.Query.Type().String(),
			Query: q.Query,
		}
	}

	data, err := json.Marshal(bqAlias)
	return data, err
}

func (q *FunctionsScoreQuery) UnmarshalJSON(data []byte) (err error) {
	type FunctionsScoreQueryAlias FunctionsScoreQuery
	bqAlias := &FunctionsScoreQueryAlias{}
	err = json.Unmarshal(data, bqAlias)
	if err != nil {
		return err
	}

	q.Query = bqAlias.QueryAlias.Query
	return
}

func (q *FunctionsScoreQuery) Type() QueryType {
	return QueryType_FunctionsScoreQuery
}

func (q *FunctionsScoreQuery) Serialize() ([]byte, error) {
	if q.Query == nil || q.Functions == nil {
		return nil, errors.New("FunctionsScoreQuery: Query or Functions is nil")
	}
	pb := &otsprotocol.FunctionsScoreQuery{}
	pbQ, err := q.Query.ProtoBuffer()
	if err != nil {
		return nil, err
	}
	pb.Query = pbQ
	for _, sf := range q.Functions {
		f, err := sf.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pb.Functions = append(pb.Functions, f)
	}
	if q.ScoreMode != nil {
		pb.SocreMode = q.ScoreMode.ProtoBuffer()
	}
	if q.CombineMode != nil {
		pb.CombineMode = q.CombineMode.ProtoBuffer()
	}
	if q.MinScore != nil {
		pb.MinScore = q.MinScore
	}
	if q.MaxScore != nil {
		pb.MaxScore = q.MaxScore
	}
	data, err := proto.Marshal(pb)
	return data, err
}

func (q *FunctionsScoreQuery) ProtoBuffer() (*otsprotocol.Query, error) {
	return BuildPBForQuery(q)
}

type ScoreFunction struct {
	FieldValueFactorFunction *FieldValueFactorFunction
	DecayFunction            *DecayFunction
	RandomFunction           *RandomFunction
	Weight                   *float32
	Filter                   Query
}

func NewScoreFunction() *ScoreFunction {
	return &ScoreFunction{}
}

func (f *ScoreFunction) SetFieldValueFactorFunction(function *FieldValueFactorFunction) *ScoreFunction {
	f.FieldValueFactorFunction = function
	return f
}

func (f *ScoreFunction) SetDecayFunction(function *DecayFunction) *ScoreFunction {
	f.DecayFunction = function
	return f
}

func (f *ScoreFunction) SetRandomFunction(function *RandomFunction) *ScoreFunction {
	f.RandomFunction = function
	return f
}

func (f *ScoreFunction) SetWeight(weight float32) *ScoreFunction {
	f.Weight = &weight
	return f
}

func (f *ScoreFunction) SetFilter(filter Query) *ScoreFunction {
	f.Filter = filter
	return f
}

func (f *ScoreFunction) ProtoBuffer() (*otsprotocol.Function, error) {
	pb := &otsprotocol.Function{}
	if f.FieldValueFactorFunction != nil {
		pbFVF, err := f.FieldValueFactorFunction.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pb.FieldValueFactor = pbFVF
	}
	if f.RandomFunction != nil {
		pbR, err := f.RandomFunction.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pb.Random = pbR
	}
	if f.DecayFunction != nil {
		pbD, err := f.DecayFunction.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pb.Decay = pbD
	}
	if f.Weight != nil {
		pb.Weight = f.Weight
	}
	if f.Filter != nil {
		pbF, err := f.Filter.ProtoBuffer()
		if err != nil {
			return nil, err
		}
		pb.Filter = pbF
	}
	return pb, nil
}

type FieldValueFactorFunction struct {
	FieldName *string
	Factor    *float32
	Modifier  *FunctionModifier
	Missing   *float64
}

func NewFieldValueFactorFunction() *FieldValueFactorFunction {
	return &FieldValueFactorFunction{}
}

func (f *FieldValueFactorFunction) SetFieldName(name string) *FieldValueFactorFunction {
	f.FieldName = &name
	return f
}

func (f *FieldValueFactorFunction) SetFactor(factor float32) *FieldValueFactorFunction {
	f.Factor = &factor
	return f
}

func (f *FieldValueFactorFunction) SetFunctionModifier(modifier FunctionModifier) *FieldValueFactorFunction {
	f.Modifier = &modifier
	return f
}

func (f *FieldValueFactorFunction) SetMissing(missing float64) *FieldValueFactorFunction {
	f.Missing = &missing
	return f
}

func (f *FieldValueFactorFunction) ProtoBuffer() (*otsprotocol.FieldValueFactorFunction, error) {
	pb := &otsprotocol.FieldValueFactorFunction{}
	if f.FieldName != nil {
		pb.FieldName = f.FieldName
	}
	if f.Factor != nil {
		pb.Factor = f.Factor
	}
	if f.Modifier != nil {
		pb.Modifier = f.Modifier.ProtoBuffer()
	}
	if f.Missing != nil {
		pb.Missing = f.Missing
	}
	return pb, nil
}

type FunctionModifier int32

const (
	// NONE 不做额外运算
	NONE FunctionModifier = iota
	// LOG 取10为底对数运算
	LOG
	// LOG1P 对真数加1后取10为底对数，防止真数为0
	LOG1P
	// LOG2P 对真数加2后取10为底对数，防止真数为0
	LOG2P
	// LN 取e为底对数运算
	LN
	// LN1P 对真数加1后取e为底对数，防止真数为0
	LN1P
	// LN2P 对真数加2后取e为底对数，防止真数为0
	LN2P
	// SQUARE 平方运算
	SQUARE
	// SQRT 开方运算
	SQRT
	// RECIPROCAL 倒数运算
	RECIPROCAL
)

func (m FunctionModifier) Enum() *FunctionModifier {
	p := new(FunctionModifier)
	*p = m
	return p
}

func (m *FunctionModifier) ProtoBuffer() *otsprotocol.FunctionModifier {
	if m == nil {
		return nil
	}
	switch *m {
	case NONE:
		return otsprotocol.FunctionModifier_FM_NONE.Enum()
	case LOG:
		return otsprotocol.FunctionModifier_FM_LOG.Enum()
	case LOG1P:
		return otsprotocol.FunctionModifier_FM_LOG1P.Enum()
	case LOG2P:
		return otsprotocol.FunctionModifier_FM_LOG2P.Enum()
	case LN:
		return otsprotocol.FunctionModifier_FM_LN.Enum()
	case LN1P:
		return otsprotocol.FunctionModifier_FM_LN1P.Enum()
	case LN2P:
		return otsprotocol.FunctionModifier_FM_LN2P.Enum()
	case SQUARE:
		return otsprotocol.FunctionModifier_FM_SQUARE.Enum()
	case RECIPROCAL:
		return otsprotocol.FunctionModifier_FM_RECIPROCAL.Enum()
	case SQRT:
		return otsprotocol.FunctionModifier_FM_SQRT.Enum()
	default:
		return nil
	}
}

type DecayFunction struct {
	FieldName      *string
	ParamType	   ParamType	`json:"ParamType"`
	DecayParam     DecayParam	`json:"DecayParam"`
	MathFunction   *MathFunction
	Decay          *float64
	MultiValueMode *MultiValueMode
}

func NewDecayFunction() *DecayFunction {
	return &DecayFunction{}
}

func (f *DecayFunction) SetFieldName(name string) *DecayFunction {
	f.FieldName = &name
	return f
}

func (f *DecayFunction) SetDecayParam(param DecayParam) *DecayFunction {
	f.DecayParam = param
	f.ParamType = param.GetType()
	return f
}

func (f *DecayFunction) SetMathFunction(function MathFunction) *DecayFunction {
	f.MathFunction = &function
	return f
}

func (f *DecayFunction) SetDecay(decay float64) *DecayFunction {
	f.Decay = &decay
	return f
}

func (f *DecayFunction) SetMultiValueMode(mode MultiValueMode) *DecayFunction {
	f.MultiValueMode = &mode
	return f
}

func (f *DecayFunction) ProtoBuffer() (*otsprotocol.DecayFunction, error) {
	pb := &otsprotocol.DecayFunction{}
	if f.FieldName != nil {
		pb.FieldName = f.FieldName
	}
	if f.MathFunction != nil {
		pb.MathFunction = f.MathFunction.ProtoBuffer()
	}
	if f.DecayParam != nil {
		ParamType := f.DecayParam.GetType()
		pb.ParamType = ParamType.ProtoBuffer()
		switch ParamType {
		case PT_DATE:
			param := f.DecayParam.(*DecayFuncDateParam)
			pbParam, err := param.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pb.Param, err = proto.Marshal(pbParam)
			if err != nil {
				return nil, err
			}
			break
		case PT_GEO:
			param := f.DecayParam.(*DecayFuncGeoParam)
			pbParam, err := param.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pb.Param, err = proto.Marshal(pbParam)
			if err != nil {
				return nil, err
			}
			break
		case PT_NUMERIC:
			param := f.DecayParam.(*DecayFuncNumericParam)
			pbParam, err := param.ProtoBuffer()
			if err != nil {
				return nil, err
			}
			pb.Param, err = proto.Marshal(pbParam)
			if err != nil {
				return nil, err
			}
			break
		default:
			return nil, errors.New("param type is illegal")
		}
	}
	if f.Decay != nil {
		pb.Decay = f.Decay
	}
	if f.MultiValueMode != nil {
		pb.MultiValueMode = f.MultiValueMode.ProtoBuffer()
	}
	return pb, nil
}

type DecayParam interface {
	GetType() ParamType
}

type ParamType int32

const (
	PT_DATE ParamType = iota
	PT_GEO
	PT_NUMERIC
)

func (t ParamType) Enum() *ParamType {
	p := new(ParamType)
	*p = t
	return p
}

func (t *ParamType) ProtoBuffer() *otsprotocol.DecayFuncParamType {
	if t == nil {
		return nil
	}
	switch *t {
	case PT_DATE:
		return otsprotocol.DecayFuncParamType_DF_DATE_PARAM.Enum()
	case PT_GEO:
		return otsprotocol.DecayFuncParamType_DF_GEO_PARAM.Enum()
	case PT_NUMERIC:
		return otsprotocol.DecayFuncParamType_DF_NUMERIC_PARAM.Enum()
	default:
		return nil
	}
}

type DecayFuncDateParam struct {
	OriginLong   *int64
	OriginString *string
	Scale        *model.DateTimeValue
	Offset       *model.DateTimeValue
}

func NewDecayFuncDateParam() *DecayFuncDateParam {
	return &DecayFuncDateParam{}
}

func (p *DecayFuncDateParam) SetOriginLong(origin int64) *DecayFuncDateParam {
	p.OriginLong = &origin
	return p
}

func (p *DecayFuncDateParam) SetOriginString(origin string) *DecayFuncDateParam {
	p.OriginString = &origin
	return p
}

func (p *DecayFuncDateParam) SetScale(value *model.DateTimeValue) *DecayFuncDateParam {
	p.Scale = value
	return p
}

func (p *DecayFuncDateParam) SetOffset(value *model.DateTimeValue) *DecayFuncDateParam {
	p.Offset = value
	return p
}

func (p *DecayFuncDateParam) GetType() ParamType {
	return PT_DATE
}

func (p *DecayFuncDateParam) ProtoBuffer() (*otsprotocol.DecayFuncDateParam, error) {
	pb := &otsprotocol.DecayFuncDateParam{}
	if p.OriginLong != nil {
		pb.OriginLong = p.OriginLong
	}
	if p.OriginString != nil {
		pb.OriginString = p.OriginString
	}
	if p.Scale != nil {
		pb.Scale = p.Scale.ProtoBuffer()
	}
	if p.Offset != nil {
		pb.Offset = p.Offset.ProtoBuffer()
	}
	return pb, nil
}

type DecayFuncGeoParam struct {
	Origin *string
	Scale  *float64
	Offset *float64
}

func NewDecayFuncGeoParam() *DecayFuncGeoParam {
	return &DecayFuncGeoParam{}
}

func (p *DecayFuncGeoParam) SetOrigin(origin string) *DecayFuncGeoParam {
	p.Origin = &origin
	return p
}

func (p *DecayFuncGeoParam) SetScale(value float64) *DecayFuncGeoParam {
	p.Scale = &value
	return p
}

func (p *DecayFuncGeoParam) SetOffset(value float64) *DecayFuncGeoParam {
	p.Offset = &value
	return p
}

func (p *DecayFuncGeoParam) GetType() ParamType {
	return PT_GEO
}

func (p *DecayFuncGeoParam) ProtoBuffer() (*otsprotocol.DecayFuncGeoParam, error) {
	pb := &otsprotocol.DecayFuncGeoParam{}
	if p.Origin != nil {
		pb.Origin = p.Origin
	}
	if p.Scale != nil {
		pb.Scale = p.Scale
	}
	if p.Offset != nil {
		pb.Offset = p.Offset
	}
	return pb, nil
}

type DecayFuncNumericParam struct {
	Origin *float64
	Scale  *float64
	Offset *float64
}

func NewDecayFuncNumericParam() *DecayFuncNumericParam {
	return &DecayFuncNumericParam{}
}

func (p *DecayFuncNumericParam) SetOrigin(origin float64) *DecayFuncNumericParam {
	p.Origin = &origin
	return p
}

func (p *DecayFuncNumericParam) SetScale(value float64) *DecayFuncNumericParam {
	p.Scale = &value
	return p
}

func (p *DecayFuncNumericParam) SetOffset(value float64) *DecayFuncNumericParam {
	p.Offset = &value
	return p
}

func (p *DecayFuncNumericParam) GetType() ParamType {
	return PT_NUMERIC
}

func (p *DecayFuncNumericParam) ProtoBuffer() (*otsprotocol.DecayFuncNumericParam, error) {
	pb := &otsprotocol.DecayFuncNumericParam{}
	if p.Origin != nil {
		pb.Origin = p.Origin
	}
	if p.Scale != nil {
		pb.Scale = p.Scale
	}
	if p.Offset != nil {
		pb.Offset = p.Offset
	}
	return pb, nil
}

type MathFunction int32

const (
	GAUSS MathFunction = iota
	EXP
	LINEAR
)

func (m MathFunction) Enum() *MathFunction {
	p := new(MathFunction)
	*p = m
	return p
}

func (f *MathFunction) ProtoBuffer() *otsprotocol.DecayMathFunction {
	if f == nil {
		return nil
	}
	switch *f {
	case GAUSS:
		return otsprotocol.DecayMathFunction_GAUSS.Enum()
	case EXP:
		return otsprotocol.DecayMathFunction_EXP.Enum()
	case LINEAR:
		return otsprotocol.DecayMathFunction_LINEAR.Enum()
	default:
		return nil
	}
}

type MultiValueMode int32

const (
	MVM_MAX MultiValueMode = iota
	MVM_MIN
	MVM_SUM
	MVM_AVG
)

func (m MultiValueMode) Enum() *MultiValueMode {
	p := new(MultiValueMode)
	*p = m
	return p
}

func (m *MultiValueMode) ProtoBuffer() *otsprotocol.MultiValueMode {
	if m == nil {
		return nil
	}
	switch *m {
	case MVM_MAX:
		return otsprotocol.MultiValueMode_MVM_MAX.Enum()
	case MVM_MIN:
		return otsprotocol.MultiValueMode_MVM_MIN.Enum()
	case MVM_SUM:
		return otsprotocol.MultiValueMode_MVM_SUM.Enum()
	case MVM_AVG:
		return otsprotocol.MultiValueMode_MVM_AVG.Enum()
	default:
		return nil
	}
}

type RandomFunction struct {
}

func NewRandomFunction() *RandomFunction {
	return &RandomFunction{}
}

func (f *RandomFunction) ProtoBuffer() (*otsprotocol.RandomScoreFunction, error) {
	pb := &otsprotocol.RandomScoreFunction{}
	return pb, nil
}

type ScoreMode int32

const (
	SM_AVG ScoreMode = iota
	SM_MAX
	SM_SUM
	SM_MIN
	SM_MULTIPLY
	SM_FIRST
)

func (m ScoreMode) Enum() *ScoreMode {
	p := new(ScoreMode)
	*p = m
	return p
}

func (m *ScoreMode) ProtoBuffer() *otsprotocol.FunctionScoreMode {
	if m == nil {
		return nil
	}
	switch *m {
	case SM_AVG:
		return otsprotocol.FunctionScoreMode_FSM_AVG.Enum()
	case SM_MAX:
		return otsprotocol.FunctionScoreMode_FSM_MAX.Enum()
	case SM_SUM:
		return otsprotocol.FunctionScoreMode_FSM_SUM.Enum()
	case SM_MIN:
		return otsprotocol.FunctionScoreMode_FSM_MIN.Enum()
	case SM_MULTIPLY:
		return otsprotocol.FunctionScoreMode_FSM_MULTIPLY.Enum()
	case SM_FIRST:
		return otsprotocol.FunctionScoreMode_FSM_FIRST.Enum()
	default:
		return nil
	}
}

type CombineMode int32

const (
	CM_MULTIPLY CombineMode = iota
	CM_AVG
	CM_MAX
	CM_SUM
	CM_MIN
	CM_REPLACE
)

func (m CombineMode) Enum() *CombineMode {
	p := new(CombineMode)
	*p = m
	return p
}

func (m *CombineMode) ProtoBuffer() *otsprotocol.FunctionCombineMode {
	if m == nil {
		return nil
	}
	switch *m {
	case CM_MULTIPLY:
		return otsprotocol.FunctionCombineMode_FCM_MULTIPLY.Enum()
	case CM_AVG:
		return otsprotocol.FunctionCombineMode_FCM_AVG.Enum()
	case CM_MAX:
		return otsprotocol.FunctionCombineMode_FCM_MAX.Enum()
	case CM_SUM:
		return otsprotocol.FunctionCombineMode_FCM_SUM.Enum()
	case CM_MIN:
		return otsprotocol.FunctionCombineMode_FCM_MIN.Enum()
	case CM_REPLACE:
		return otsprotocol.FunctionCombineMode_FCM_REPLACE.Enum()
	default:
		return nil
	}
}
