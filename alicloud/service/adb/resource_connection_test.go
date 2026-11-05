// Package adb_test.
package adb_test

import (
	"fmt"
	"regexp"
	"testing"

	adbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/adb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudADBConnectionConfig(t *testing.T) {
	var v *adbsdk.Address
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	rand2 := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccAdbConnection%s", rand)
	prefix := fmt.Sprintf("tf-testacc%s", rand)
	prefix2 := fmt.Sprintf("tf-testacc%s", rand2)
	var basicMap = map[string]string{
		"db_cluster_id":            acctest.CHECKSET,
		"connection_prefix":        acctest.CHECKSET,
		"connection_string_prefix": acctest.CHECKSET,
		"connection_string":        acctest.CHECKSET,
		"ip_address":               acctest.CHECKSET,
		"port":                     acctest.CHECKSET,
	}
	resourceId := "alicloud_adb_connection.default"
	ra := acctest.ResourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return adb.NewAdbService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbConnection")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAdbConnectionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		},

		// module name

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": "${alicloud_adb_db_cluster.cluster.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": prefix,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix":        prefix,
						"connection_string_prefix": prefix,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix":        acctest.REMOVEKEY,
					"connection_string_prefix": prefix2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_prefix":        prefix2,
						"connection_string_prefix": prefix2,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_prefix": prefix,
				}),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func resourceAdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "alicloud_adb_db_cluster" "cluster" {
	db_cluster_category = "MixedStorage"
	mode = "flexible"
	compute_resource = "8Core32GB"
	vswitch_id              = local.vswitch_id
	description             = "${var.name}"
	}`, AdbCommonTestCase, name)
}
