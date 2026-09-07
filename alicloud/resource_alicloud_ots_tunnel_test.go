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
					testAccCheck(getCheckMapWithTunnelType(name, randTunnelType)),
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
		return BaseDataTunnel
	} else if r == 2 {
		return StreamTunnel
	} else {
		return BaseAndStreamTunnel
	}
}

func getCheckMapWithTunnelType(name string, tunnelType string) map[string]string {
	accMap := make(map[string]string)
	accMap["instance_name"] = "tf-" + name
	accMap["table_name"] = name
	accMap["tunnel_name"] = name
	accMap["tunnel_id"] = CHECKSET
	accMap["tunnel_type"] = tunnelType
	if tunnelType == string(StreamTunnel) {
		accMap["tunnel_stage"] = "ProcessStream"
	} else {
		accMap["tunnel_stage"] = "ProcessBaseData"
	}
	accMap["tunnel_rpo"] = "0"
	accMap["expired"] = "false"
	accMap["create_time"] = CHECKSET
	if tunnelType == string(BaseAndStreamTunnel) {
		accMap["channels.#"] = "2"
		accMap["channels.0.channel_id"] = CHECKSET
		accMap["channels.0.channel_type"] = CHECKSET
		accMap["channels.0.channel_status"] = CHECKSET
		accMap["channels.0.client_id"] = ""
		accMap["channels.0.channel_rpo"] = "0"
		accMap["channels.1.channel_id"] = CHECKSET
		accMap["channels.1.channel_type"] = CHECKSET
		accMap["channels.1.channel_status"] = CHECKSET
		accMap["channels.1.client_id"] = ""
		accMap["channels.1.channel_rpo"] = "0"
	} else {
		accMap["channels.#"] = "1"
		accMap["channels.0.channel_id"] = CHECKSET
		accMap["channels.0.channel_type"] = CHECKSET
		accMap["channels.0.channel_status"] = CHECKSET
		accMap["channels.0.client_id"] = ""
		accMap["channels.0.channel_rpo"] = "0"
	}
	return accMap
}
