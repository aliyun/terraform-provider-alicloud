package realtimecompute_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test RealtimeCompute UdfArtifact. >>> Resource test cases, handwritten.
//
// The ververica UdfArtifact API registers a Java jar or Python file stored in OSS
// with a namespace; the service derives udfClasses from the uploaded artifact. A
// Python artifact is used here so the test can ship a real, plain-text UDF file as
// an OSS object instead of a binary jar fixture; the bucket is the one bound to the
// VVP instance so the Flink service can read the artifact back. The composite id is
// <workspace>:<namespace>:<udfArtifactName>, all three ForceNew identity parts.

func TestAccAliCloudRealtimeComputeUdfArtifact_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_realtime_compute_udf_artifact.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRealtimeComputeUdfArtifactMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return realtimecompute.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRealtimeComputeUdfArtifact")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRealtimeComputeUdfArtifactBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":     "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":     "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":          name,
					"artifact_type": "ARTIFACT_TYPE_PYTHON",
					"jar_url":       "oss://${alicloud_oss_bucket.default.id}/udf/UpperUdf.py",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"artifact_type": "ARTIFACT_TYPE_PYTHON",
						"workspace":     acctest.CHECKSET,
						"namespace":     acctest.CHECKSET,
						"jar_url":       acctest.CHECKSET,
					}),
				),
			},
			// jar_url is mutable in place: pointing the artifact at a second UDF file reaches
			// UpdateUdfArtifact (PUT) and re-derives udfClasses from the new file.
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace":     "${alicloud_realtime_compute_vvp_instance.default.resource_id}",
					"namespace":     "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default",
					"name":          name,
					"artifact_type": "ARTIFACT_TYPE_PYTHON",
					"jar_url":       "oss://${alicloud_oss_bucket.default.id}/udf/LowerUdf.py",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"artifact_type": "ARTIFACT_TYPE_PYTHON",
						"workspace":     acctest.CHECKSET,
						"namespace":     acctest.CHECKSET,
						"jar_url":       acctest.CHECKSET,
						"modified_at":   acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// AlicloudRealtimeComputeUdfArtifactMap0 holds the server-derived Computed fields
// that are checked across every step. creator/created_at are populated on create;
// modified_at is additionally asserted after the update step.
var AlicloudRealtimeComputeUdfArtifactMap0 = map[string]string{
	"creator":    acctest.CHECKSET,
	"created_at": acctest.CHECKSET,
}

func AlicloudRealtimeComputeUdfArtifactBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_oss_bucket_object" "udf_v1" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "udf/UpperUdf.py"
  content = <<-EOF
class UpperUdf:
    def eval(self, s):
        if s is None:
            return None
        return str(s).upper()
  EOF
}

resource "alicloud_oss_bucket_object" "udf_v2" {
  bucket  = alicloud_oss_bucket.default.id
  key     = "udf/LowerUdf.py"
  content = <<-EOF
class LowerUdf:
    def eval(self, s):
        if s is None:
            return None
        return str(s).lower()
  EOF
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = ["${alicloud_vswitch.default.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.default.zone_id
}
`, name)
}

// Test RealtimeCompute UdfArtifact. <<< Resource test cases.
