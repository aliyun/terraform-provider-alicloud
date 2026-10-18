package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudApiGatewaySignature covers the alicloud_api_gateway_signature
// resource lifecycle end to end: create, update the name, update the key/secret,
// revert to the original values, and import. Assertions check the real attribute
// values rather than ${var.xxx} placeholders, which the testing framework compares
// literally and would not resolve.
func TestAccAliCloudApiGatewaySignature(t *testing.T) {
	r := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccSignature_%d", r)
	key := "tf_testAccKey"
	secret := "tf_testAccSecret"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewaySignatureConfig(name, key, secret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_signature.default", "id"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_name", name),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_key", key),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_secret", secret),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_signature.default", "created_time"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_signature.default", "modified_time"),
				),
			},
			{
				Config: testAccApiGatewaySignatureConfig(name+"_u", key, secret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_name", name+"_u"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_key", key),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_secret", secret),
				),
			},
			{
				Config: testAccApiGatewaySignatureConfig(name+"_u", key+"_u", secret+"_u"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_name", name+"_u"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_key", key+"_u"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_secret", secret+"_u"),
				),
			},
			{
				Config: testAccApiGatewaySignatureConfig(name, key, secret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_name", name),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_key", key),
					resource.TestCheckResourceAttr("alicloud_api_gateway_signature.default", "signature_secret", secret),
				),
			},
			{
				ResourceName:      "alicloud_api_gateway_signature.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiGatewaySignatureConfig(name, key, secret string) string {
	return fmt.Sprintf(`
resource "alicloud_api_gateway_signature" "default" {
  signature_name   = "%[1]s"
  signature_key    = "%[2]s"
  signature_secret = "%[3]s"
}
`, name, key, secret)
}
