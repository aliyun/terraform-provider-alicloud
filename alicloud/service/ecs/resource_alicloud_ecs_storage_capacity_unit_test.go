package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The white-box TestUnitAlicloudECSStorageCapacityUnit (which calls unexported
// resourceEcsStorageCapacityUnit Create/Read/Update, the package-ecs aliases
// GetNotFoundErrorFromString/GetNotFoundMessage and the unit-test engine
// helpers) has been split into resource_alicloud_ecs_storage_capacity_unit_unit_test.go
// (package ecs); only the external acceptance tests remain here.
func TestAccAliCloudECSStorageCapacityUnit_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_storage_capacity_unit.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsStorageCapacityUnitMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsStorageCapacityUnit")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsstoragecapacityunit%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsStorageCapacityUnitBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                name,
					"capacity":                   "20",
					"period":                     "1",
					"period_unit":                "Month",
					"storage_capacity_unit_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                CHECKSET,
						"capacity":                   "20",
						"storage_capacity_unit_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "period"},
			},
		},
	})
}

var AlicloudEcsStorageCapacityUnitMap0 = map[string]string{}

func AlicloudEcsStorageCapacityUnitBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAliCloudECSStorageCapacityUnit_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_storage_capacity_unit.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsStorageCapacityUnitMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsStorageCapacityUnit")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsstoragecapacityunit%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEcsStorageCapacityUnitBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time":                 "2022-01-01T12:05:05Z",
					"storage_capacity_unit_name": name,
					"capacity":                   "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_capacity_unit_name": CHECKSET,
						"capacity":                   "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_capacity_unit_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_capacity_unit_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                name,
					"storage_capacity_unit_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                name,
						"storage_capacity_unit_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "period"},
			},
		},
	})
}

var AlicloudEcsStorageCapacityUnitMap1 = map[string]string{}

func AlicloudEcsStorageCapacityUnitBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
