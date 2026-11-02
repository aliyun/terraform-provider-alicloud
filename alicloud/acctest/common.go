package acctest

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest/helpers"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	NOSET             = helpers.NOSET
	CHECKSET          = helpers.CHECKSET
	REMOVEKEY         = helpers.REMOVEKEY
	CLEARLIST         = helpers.CLEARLIST
	CLEARMAP          = helpers.CLEARMAP
	REGEXMATCH        = helpers.REGEXMATCH
	ForceSleep        = helpers.ForceSleep
	INDENTATIONSYMBOL = helpers.INDENTATIONSYMBOL
	CHILDINDEND       = helpers.CHILDINDEND
)

type ResourceAttrMapUpdate = helpers.ResourceAttrMapUpdate
type ResourceTestAccConfigFunc = helpers.ResourceTestAccConfigFunc
type ResourceCheck = helpers.ResourceCheck
type ResourceAttr = helpers.ResourceAttr
type ResourceAttrCheck = helpers.ResourceAttrCheck
type DataSourceTestAccConfig = helpers.DataSourceTestAccConfig

var (
	ResourceAttrInit            = helpers.ResourceAttrInit
	ResourceAttrCheckInit       = helpers.ResourceAttrCheckInit
	ResourceTestAccConfig       = helpers.ResourceTestAccConfig
	DataSourceTestAccConfigFunc = helpers.DataSourceTestAccConfigFunc
	CheckAlicloudDataSourceID   = helpers.CheckAlicloudDataSourceID
)

func ResourceCheckInit(id string, object interface{}, service func() interface{}, attrs ...string) *ResourceCheck {
	return helpers.ResourceCheckInit(Provider, id, object, service, attrs...)
}
func ResourceCheckInitWithDescribeMethod(id string, object interface{}, service func() interface{}, describe string, attrs ...string) *ResourceCheck {
	return helpers.ResourceCheckInitWithDescribeMethod(Provider, id, object, service, describe, attrs...)
}
func NewInstanceDiff(name string, attrs, diff map[string]interface{}, state *sdktf.InstanceState) (*sdktf.InstanceDiff, error) {
	return helpers.NewInstanceDiff(Provider, name, attrs, diff, state)
}

// DataSourceAttr retains the full-provider defaults used by existing tests.
type DataSourceAttr struct {
	ResourceId   string
	ExistMapFunc func(int) map[string]string
	FakeMapFunc  func(int) map[string]string
}

func (d *DataSourceAttr) helper() *helpers.DataSourceAttr {
	return &helpers.DataSourceAttr{ResourceId: d.ResourceId, ExistMapFunc: d.ExistMapFunc, FakeMapFunc: d.FakeMapFunc, ProviderFactories: SDKProviderFactories, PreCheck: PreCheck}
}
func (d *DataSourceAttr) DataSourceTestCheck(t *testing.T, rand int, configs ...DataSourceTestAccConfig) {
	d.helper().DataSourceTestCheck(t, rand, configs...)
}
func (d *DataSourceAttr) DataSourceTestCheckWithPreCheck(t *testing.T, rand int, check func(), configs ...DataSourceTestAccConfig) {
	d.helper().DataSourceTestCheckWithPreCheck(t, rand, check, configs...)
}
