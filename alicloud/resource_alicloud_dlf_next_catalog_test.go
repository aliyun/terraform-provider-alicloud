package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudDlfNextCatalog_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dlf_next_catalog.default"
	rc := acctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-catalog-%d", rc)

	testAccConfig := resourceTestAccConfigFunc(resourceId, catalogName, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": catalogName,
					"type": "PAIMON",
					"options": map[string]interface{}{
						"key1": "value1",
						"key2": "value2",
					},
					"is_shared": false,
					"share_id":  "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "name", catalogName),
					resource.TestCheckResourceAttr(resourceId, "type", "PAIMON"),
					resource.TestCheckResourceAttr(resourceId, "options.key1", "value1"),
					resource.TestCheckResourceAttr(resourceId, "options.key2", "value2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": map[string]interface{}{
						"key1": "updated1",
						"key3": "value3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.key1", "updated1"),
					resource.TestCheckResourceAttr(resourceId, "options.key3", "value3"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": CLEARMAP,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
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

func TestAccAliCloudDlfNextCatalog_iceberg(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dlf_next_catalog.default"
	rc := acctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-catalog-iceberg-%d", rc)

	testAccConfig := resourceTestAccConfigFunc(resourceId, catalogName, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": catalogName,
					"type": "ICEBERG",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "type", "ICEBERG"),
				),
			},
		},
	})
}

func TestAccAliCloudDlfNextCatalogDataSource_basic(t *testing.T) {
	rc := acctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-ds-catalog-%d", rc)

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_dlf_next_catalogs.default", catalogName, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"catalog_name_pattern": catalogName,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_dlf_next_catalogs.default", "catalogs.#"),
				),
			},
		},
	})
}

func testAccAlicloudDlfNextCatalogExists(resourceId string, v *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceId)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DlfNext Catalog ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		dlfNextServiceV2 := DlfNextServiceV2{client}

		object, err := dlfNextServiceV2.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error describing DlfNext Catalog: %s", err)
		}

		*v = object
		return nil
	}
}

func testAccCheckAlicloudDlfNextCatalogDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dlf_next_catalog" {
			continue
		}

		dlfNextServiceV2 := DlfNextServiceV2{client}
		_, err := dlfNextServiceV2.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error checking DlfNext Catalog destroy: %s", err)
		}
		return fmt.Errorf("DlfNext Catalog still exists: %s", rs.Primary.ID)
	}
	return nil
}
