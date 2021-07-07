package datahub

import (
    "encoding/json"
    "fmt"
)

type Field struct {
    Name      string    `json:"name"`
    Type      FieldType `json:"type"`
    AllowNull bool      `json:"notnull"`
}

// RecordSchema
type RecordSchema struct {
    Fields []Field `json:"fields"`
}

// NewRecordSchema create a new record schema for tuple record
func NewRecordSchema() *RecordSchema {
    return &RecordSchema{
        Fields: make([]Field, 0, 0),
    }
}

func NewRecordSchemaFromJson(SchemaJson string) (recordSchema *RecordSchema, err error) {
    recordSchema = &RecordSchema{}
    if err = json.Unmarshal([]byte(SchemaJson), recordSchema); err != nil {
        return
    }
    for _, v := range recordSchema.Fields {
        if !validateFieldType(v.Type) {
            panic(fmt.Sprintf("field type %q illegal", v.Type))
        }
    }
    return
}

func (rs *RecordSchema) String() string {
    for idx := range rs.Fields {
        rs.Fields[idx].AllowNull = !rs.Fields[idx].AllowNull
    }
    byts, _ := json.Marshal(rs)
    return string(byts)
}

// AddField add a field
func (rs *RecordSchema) AddField(f Field) *RecordSchema {
    if !validateFieldType(f.Type) {
        panic(fmt.Sprintf("field type %q illegal", f.Type))
    }
    for _, v := range rs.Fields {
        if v.Name == f.Name {
            panic(fmt.Sprintf("field %q duplicated", f.Name))
        }
    }
    rs.Fields = append(rs.Fields, f)
    return rs
}

// GetFieldIndex get index of given field
func (rs *RecordSchema) GetFieldIndex(fname string) int {
    for idx, v := range rs.Fields {
        if fname == v.Name {
            return idx
        }
    }
    return -1
}

// Size get record schema fields size
func (rs *RecordSchema) Size() int {
    return len(rs.Fields)
}
