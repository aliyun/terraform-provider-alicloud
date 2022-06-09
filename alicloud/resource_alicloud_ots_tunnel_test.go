package alicloud

import (
	"fmt"
	"testing"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOtsTunnel_basic(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := resourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsTunnelConfigDependence)

	randTunnelType := string(getRandomTunnelType())
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   randTunnelType,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-" + name,
						"table_name":    name,
						"tunnel_name":   name,
						"tunnel_id":     CHECKSET,
						"tunnel_type":   randTunnelType,
						"tunnel_stage":  "InitBaseDataAndStreamShard",
						"tunnel_rpo":    CHECKSET,
						"expired":       "false",
						"create_time":   CHECKSET,
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

func resourceOtsTunnelConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "%s"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	resource "alicloud_ots_table" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }
	  time_to_live = -1
	  max_version = 1
	}
	`, name, string(OtsCapacity))
}

var otsTunnelBasicMap = map[string]string{}

func getRandomTunnelType() TunnelTypeString {
	r := acctest.RandInt() % 3
	if r == 1 {
		return BaseData
	} else if r == 2 {
		return Stream
	} else {
		return BaseAndStream
	}
}
