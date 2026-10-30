package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudApiGatewayModel_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_api_gateway_model.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudApiGatewayModelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApiGatewayModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-m%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudApiGatewayModelBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":    "${alicloud_api_gateway_group.default.id}",
					"model_name":  name,
					"schema":      `{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}`,
					"description": name,
					"tags": map[string]interface{}{
						"Created": "tf-testacc",
						"Env":     "test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":     CHECKSET,
						"model_name":   name,
						"schema":       "{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}",
						"description":  name,
						"tags.Created": "tf-testacc",
						"tags.Env":     "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schema": `{\"type\":\"object\",\"properties\":{\"name\":{\"maxLength\":10,\"type\":\"string\"}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"schema": "{\"type\":\"object\",\"properties\":{\"name\":{\"maxLength\":10,\"type\":\"string\"}}}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"model_name": name + "-renamed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_name": name + "-renamed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{
						"Created": "tf-testacc",
						"Env":     "prod",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.Created": "tf-testacc",
						"tags.Env":     "prod",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var resourceAlicloudApiGatewayModelMap = map[string]string{
	"group_id": CHECKSET,
}

func resourceAlicloudApiGatewayModelBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_api_gateway_group" "default" {
  		name        = var.name
  		description = var.name
	}
`, name)
}

// lintignore: R001
