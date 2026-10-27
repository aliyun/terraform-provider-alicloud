package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnitMongodbShardingAuditFilterStorageValues pins the edition-specific retention mapping against
// the Dds 2015-12-01 API definition: DescribeMongoDBLogConfig reports the Standard edition retention in
// TtlForStandard (the value ModifyAuditPolicy's StoragePeriod writes for that edition) and the
// V2_Standard cold / hot storage days in TtlForV2Standard / HotTtlForV2Standard (StoragePeriod /
// HotStoragePeriod). The reported edition decides which pair is authoritative; when it is neither of
// the two supported editions, the declared edition decides; an omitted field stays nil so the caller
// keeps the state value.
func TestUnitMongodbShardingAuditFilterStorageValues(t *testing.T) {
	cases := []struct {
		name     string
		raw      map[string]interface{}
		reported string
		declared string
		storage  interface{}
		hot      interface{}
	}{
		{
			name:     "v2 standard reads cold and hot storage",
			raw:      map[string]interface{}{"TtlForStandard": 30, "TtlForV2Standard": 180, "HotTtlForV2Standard": 5},
			reported: "V2_Standard",
			declared: "V2_Standard",
			storage:  180,
			hot:      5,
		},
		{
			name:     "standard reads retention only",
			raw:      map[string]interface{}{"TtlForStandard": 180, "TtlForV2Standard": 30},
			reported: "Standard",
			declared: "Standard",
			storage:  180,
			hot:      nil,
		},
		{
			name:     "reported edition wins over the declared one",
			raw:      map[string]interface{}{"TtlForStandard": 30, "TtlForV2Standard": 180, "HotTtlForV2Standard": 5},
			reported: "Standard",
			declared: "V2_Standard",
			storage:  30,
			hot:      nil,
		},
		{
			name:     "disabled reports trial so the declared edition decides",
			raw:      map[string]interface{}{"TtlForStandard": 30, "TtlForV2Standard": 180, "HotTtlForV2Standard": 5},
			reported: "Trial",
			declared: "V2_Standard",
			storage:  180,
			hot:      5,
		},
		{
			name:     "unknown edition without declaration falls back to standard",
			raw:      map[string]interface{}{"TtlForStandard": 30},
			reported: "Trial",
			declared: "",
			storage:  30,
			hot:      nil,
		},
		{
			name:     "omitted fields keep the state value",
			raw:      map[string]interface{}{},
			reported: "Trial",
			declared: "V2_Standard",
			storage:  nil,
			hot:      nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			storage, hot := mongodbShardingAuditFilterStorageValues(c.raw, c.reported, c.declared)
			assert.Equal(t, c.storage, storage)
			assert.Equal(t, c.hot, hot)
		})
	}
}
