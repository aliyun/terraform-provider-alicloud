package ots_test

import (
	"fmt"
	"testing"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOtsTunnel_basic(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   "BaseAndStream",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": acctest.CHECKSET,
						"table_name":    acctest.CHECKSET,
						"tunnel_name":   acctest.CHECKSET,
						"tunnel_type":   "BaseAndStream",
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
	  instance_type = "Capacity"
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
	`, name)
}

var otsTunnelBasicMap = map[string]string{}
