package ros_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ros"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudROSChangeSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ros.NewRosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRosChangeSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRosChangeSetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// ROS DeleteChangeSet is a soft delete: GetChangeSet keeps returning
		// the change set with Status=DELETE_COMPLETE rather than NotFound, so the
		// generic CheckResourceDestroy (which only treats NotFound as destroyed)
		// would always report "not destroyed". Mirror the ROS Stack test's
		// CheckDestroy and accept DELETE_COMPLETE as a successfully destroyed
		// change set.
		CheckDestroy: func(state *terraform.State) error {
			rosServiceClient := ros.NewRosService(acctest.Provider.Meta().(*connectivity.AliyunClient))
			for _, rs := range state.RootModule().Resources {
				if "alicloud_ros_change_set" != rs.Type {
					continue
				}
				changeSetId := rs.Primary.ID
				object, err := rosServiceClient.DescribeRosChangeSet(changeSetId)
				if err != nil {
					if errs.NotFoundError(err) {
						continue
					}
					return errs.WrapError(err)
				}
				if status, ok := object["Status"].(string); ok && status == "DELETE_COMPLETE" {
					continue
				}
				return fmt.Errorf("ROS ChangeSet %s was not destroyed (status: %v)", changeSetId, object["Status"])
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
					"stack_name":      name + "stack",
					"change_set_type": "CREATE",
					"description":     "Test From Terraform",
					"template_body":   `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack",
						"change_set_type": "CREATE",
						"description":     "Test From Terraform",
						"template_body":   acctest.CHECKSET,
					}),
				),
			},
			{
				// Removing the required change_set_name verifies the plan
				// fails before reaching the API. The final step below
				// restores change_set_name so the post-test destroy renders
				// a valid config; the persistent attributeMap would otherwise
				// keep it deleted after REMOVEKEY, and the import step's
				// config is not used for deferred destroy.
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": acctest.REMOVEKEY,
				}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
				}),
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"notification_urls", "ram_role_name", "replacement_option", "stack_policy_body", "stack_policy_during_update_body", "stack_policy_during_update_url", "stack_policy_url", "template_url", "use_previous_parameters"},
			},
			{
				// Restoring change_set_name as the final non-import step
				// ensures the post-test destroy renders a valid config
				// containing the required change_set_name argument. The
				// import step above does not supply the config used for
				// deferred destroy, so without this step the destroy would
				// reuse the REMOVEKEY state and fail schema validation.
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
				}),
			},
		},
	})
}

var AlicloudRosChangeSetMap = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudRosChangeSetBasicDependence(name string) string {
	return ""
}
