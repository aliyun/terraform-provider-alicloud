package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) { resource.TestMain(m) }

var sharedClientForRegion = sweep.SharedClientForRegion
