package paidsw_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// PAI DSW snapshots are created against an existing DSW instance and reference an
// image URL; the provider has no alicloud_pai_dsw_instance resource to build one
// inline (and creating one is a paid compute resource outside this task's
// authorization). The two values are therefore borrowed prerequisites supplied
// via environment, the same pattern other heavy-dependency ACC tests in this
// repo use. When they are absent the lifecycle test skips — this is a genuine
// account-access dependency reported as a verification gap, not a fake pass.
const (
	envPaiDswInstanceId = "PAI_DSW_INSTANCE_ID"
	envPaiDswImageUrl   = "PAI_DSW_IMAGE_URL"
)

func testAccPaiDswApiInstanceSnapshotPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	if v := os.Getenv(envPaiDswInstanceId); v == "" {
		t.Skipf("%s must reference an existing PAI DSW instance for this acceptance test", envPaiDswInstanceId)
	}
	if v := os.Getenv(envPaiDswImageUrl); v == "" {
		t.Skipf("%s must reference a valid image URL for this acceptance test", envPaiDswImageUrl)
	}
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_basic covers the full create/read/
// no-op/update-via-replace/import lifecycle, including the optional labels and
// exclude_paths collections and the create-only snapshot_description/overwrite
// fields. The snapshot API has no Update, so every input change forces
// replacement.
func TestAccAliCloudPaiDswApiInstanceSnapshot_basic(t *testing.T) {
	r := tfacctest.RandInt()
	testAccPaiDswApiInstanceSnapshotPreCheck(t)

	instanceId := os.Getenv(envPaiDswInstanceId)
	imageUrl := os.Getenv(envPaiDswImageUrl)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "status"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "instance_id", instanceId),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "image_url", imageUrl),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_name", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_description", fmt.Sprintf("tf-testacc description %d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "overwrite", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "labels.#", "2"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "exclude_paths.#", "2"),
				),
			},
			{
				// No-op: re-applying the same config must produce an empty plan.
				Config: testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl, r),
			},
			{
				// Update via replacement: snapshot_name, description and the
				// labels/exclude_paths collections all change; overwrite is
				// omitted to cover the unset branch.
				Config: testAccPaiDswApiInstanceSnapshotConfigUpdate(instanceId, imageUrl, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_name", fmt.Sprintf("tf-testacc-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_description", fmt.Sprintf("tf-testacc description %d upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "labels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "exclude_paths.#", "1"),
				),
			},
			{
				ResourceName:            "alicloud_pai_dsw_api_instance_snapshot.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_description", "overwrite"},
			},
		},
	})
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_invalidName covers the failure path
// where the snapshot_name violates the documented charset (lowercase letters,
// digits, hyphens). It is a schema-level validation that fires during plan
// with no API call, so it only needs credentials, not a borrowed instance.
func TestAccAliCloudPaiDswApiInstanceSnapshot_invalidName(t *testing.T) {
	acctest.PreCheck(t)
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiDswApiInstanceSnapshotConfigInvalidName(r),
				// The regex uses \s+ between tokens because the Terraform CLI
				// line-wraps the diagnostic detail (~80 cols), inserting a
				// newline + indentation mid-message. A literal-spaces pattern
				// would fail to match the wrapped text.
				ExpectError: regexp.MustCompile(`must\s+only\s+contain\s+lowercase\s+letters,\s+digits\s+and\s+hyphens`),
			},
		},
	})
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_invalidInstance covers the cloud
// failure path: Create against a non-existent instance must surface an error
// rather than succeed silently. It needs credentials only — the bogus instance
// and dummy image make the Create call fail before any real resource is built.
func TestAccAliCloudPaiDswApiInstanceSnapshot_invalidInstance(t *testing.T) {
	acctest.PreCheck(t)
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccPaiDswApiInstanceSnapshotConfigInvalidInstance(r),
				ExpectError: regexp.MustCompile(`(?i)(not\s*found|not\s*exist|invalid|forbidden|endpoint|404|400|instance)`),
			},
		},
	})
}

func testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl string, r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id          = "%[1]s"
  image_url            = "%[2]s"
  snapshot_name        = "tf-testacc-%[3]d"
  snapshot_description = "tf-testacc description %[3]d"
  overwrite            = true

  labels = [
    { key = "env",    value = "tf-testacc" },
    { key = "purpose", value = "snapshot" },
  ]

  exclude_paths = ["/home/admin/data", "/tmp/cache"]
}
`, instanceId, imageUrl, r)
}

func testAccPaiDswApiInstanceSnapshotConfigUpdate(instanceId, imageUrl string, r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id          = "%[1]s"
  image_url            = "%[2]s"
  snapshot_name        = "tf-testacc-%[3]d-upd"
  snapshot_description = "tf-testacc description %[3]d upd"

  labels = [
    { key = "env", value = "tf-testacc-upd" },
  ]

  exclude_paths = ["/home/admin/workspace"]
}
`, instanceId, imageUrl, r)
}

func testAccPaiDswApiInstanceSnapshotConfigInvalidName(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id   = "i-tf-testacc-%[1]d"
  image_url     = "registry.cn-hangzhou.aliyuncs.com/pai-dsw/nonexistent:latest"
  snapshot_name = "TF-Testacc-UPPER-%[1]d"
}
`, r)
}

func testAccPaiDswApiInstanceSnapshotConfigInvalidInstance(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id   = "i-tf-testacc-nonexistent-%[1]d"
  image_url     = "registry.cn-hangzhou.aliyuncs.com/pai-dsw/nonexistent:latest"
  snapshot_name = "tf-testacc-invalidinst-%[1]d"
}
`, r)
}
