package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbListener_importHttp(t *testing.T) {
	resourceName := "alicloud_slb_listener.http"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerHttp,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSlbListener_importHttps(t *testing.T) {
	resourceName := "alicloud_slb_listener.https"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerHttps,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSlbListener_importTcp(t *testing.T) {
	resourceName := "alicloud_slb_listener.tcp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerTcp,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSlbListener_importUdp(t *testing.T) {
	resourceName := "alicloud_slb_listener.udp"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbListenerUdp,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
