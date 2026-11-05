package ddoscoo_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ddoscoo "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddoscoo"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const udpReflectResourceId = "alicloud_ddoscoo_udp_reflect.default"

// TestAccAliCloudDdosCooUdpReflect_basic exercises the full lifecycle that the
// single Config/Describe API pair supports:
//   - create applies the full UdpSports list;
//   - import by instance id reproduces the same state;
//   - update reducing the list drops the removed port server-side (full
//     replacement, not append);
//   - destroy clears the policy via an empty UdpSports list.
func TestAccAliCloudDdosCooUdpReflect_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
		},
		IDRefreshName:            udpReflectResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDdosCooUdpReflectDestroy(udpReflectResourceId),
		Steps: []resource.TestStep{
			{
				Config: testAccDdosCooUdpReflectConfig("[17, 19]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.#", "2"),
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.0", "17"),
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.1", "19"),
					testAccCheckDdosCooUdpReflectFromAPI(udpReflectResourceId, []int{17, 19}),
				),
			},
			{
				ResourceName:      udpReflectResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ConfigUdpReflect applies the full list, so reducing the set must drop
			// port 19 from the server-side policy rather than leave it behind.
			{
				Config: testAccDdosCooUdpReflectConfig("[17]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.#", "1"),
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.0", "17"),
					testAccCheckDdosCooUdpReflectFromAPI(udpReflectResourceId, []int{17}),
				),
			},
			// Expanding the list again confirms the previous reduction was real and
			// that an added port lands server-side.
			{
				Config: testAccDdosCooUdpReflectConfig("[17, 123]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.#", "2"),
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.0", "17"),
					resource.TestCheckResourceAttr(udpReflectResourceId, "udp_sports.1", "123"),
					testAccCheckDdosCooUdpReflectFromAPI(udpReflectResourceId, []int{17, 123}),
				),
			},
		},
	})
}

func testAccDdosCooUdpReflectConfig(ports string) string {
	return fmt.Sprintf(`
data "alicloud_ddoscoo_instances" "default" {}

resource "alicloud_ddoscoo_udp_reflect" "default" {
  instance_id = data.alicloud_ddoscoo_instances.default.ids.0
  udp_sports  = %s
}
`, ports)
}

// testAccCheckDdosCooUdpReflectDestroy asserts that destroying the resource
// cleared the policy: DescribeUdpReflect must report an empty UdpSports list
// (or the owning instance is gone, which is also a clean destroy).
func testAccCheckDdosCooUdpReflectDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		service := ddoscoo.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_ddoscoo_udp_reflect" {
				continue
			}
			object, err := service.DescribeDdoscooUdpReflect(rs.Primary.ID)
			if err != nil {
				return nil
			}
			if ports := ddoscoo.UdpReflectResponsePorts(object["UdpSports"]); len(ports) != 0 {
				return fmt.Errorf("DdosCoo UdpReflect policy %s still has ports %v after destroy", rs.Primary.ID, ports)
			}
		}
		return nil
	}
}

// testAccCheckDdosCooUdpReflectFromAPI reads the policy directly from the API
// and compares it to the expected set, independently of Terraform state, so a
// state/plan mismatch is caught.
func testAccCheckDdosCooUdpReflectFromAPI(resourceId string, expected []int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := ddoscoo.NewService(client).DescribeDdoscooUdpReflect(rs.Primary.ID)
		if err != nil {
			return err
		}
		got := ddoscoo.UdpReflectResponsePorts(object["UdpSports"])
		if len(got) != len(expected) {
			return fmt.Errorf("UdpSports mismatch: got %v, expected %v", got, expected)
		}
		for i, v := range got {
			if int(v) != expected[i] {
				return fmt.Errorf("UdpSports mismatch: got %v, expected %v", got, expected)
			}
		}
		return nil
	}
}
