package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudKeyPair_basic(t *testing.T) {
	var keypair ecs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_key_pair.basic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKeyPairConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(
						"alicloud_key_pair.basic", &keypair),
				),
			},
		},
	})

}

func TestAccAlicloudKeyPair_prefix(t *testing.T) {
	var keypair ecs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_key_pair.prefix",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKeyPairConfigPrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(
						"alicloud_key_pair.prefix", &keypair),
					testAccCheckKeyPairHasPrefix(
						"alicloud_key_pair.prefix", &keypair, "terraform-test-key-pair-prefix"),
				),
			},
		},
	})

}

func TestAccAlicloudKeyPair_publicKey(t *testing.T) {
	var keypair ecs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_key_pair.publickey",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKeyPairConfigPublicKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists(
						"alicloud_key_pair.publickey", &keypair),
					testAccCheckKeyPairHasPrefix(
						"alicloud_key_pair.publickey", &keypair, resource.UniqueIdPrefix),
				),
			},
		},
	})

}

func testAccCheckKeyPairExists(n string, keypair *ecs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeKeyPair(rs.Primary.ID)

		log.Printf("[WARN] disk ids %#v", rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Finding Key Pair %#v got an error: %#v.", rs.Primary.ID, err)
		}
		*keypair = response
		return nil
	}
}

func testAccCheckKeyPairHasPrefix(n string, keypair *ecs.KeyPair, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Key Pair ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeKeyPair(rs.Primary.ID)

		log.Printf("[WARN] disk ids %#v", rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Finding Key Pair prefix %#v got an error: %#v.", rs.Primary.ID, err)
		}

		if strings.HasPrefix(response.KeyPairName, prefix) {
			*keypair = response
		}
		return nil
	}
}

func testAccCheckKeyPairDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_key_pair" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeKeyPair(rs.Primary.ID)
		os.Remove(rs.Primary.Attributes["key_file"])

		if err != nil {
			// Verify the error is what we want
			if NotFoundError(err) || IsExceptedError(err, KeyPairNotFound) {
				continue
			}
			return err
		}
		if response.KeyPairName != "" {
			return fmt.Errorf("Error Key Pair still exist")
		}
	}

	return nil
}

const testAccKeyPairConfig = `
resource "alicloud_key_pair" "basic" {
	key_name = "terraform-test-key-pair"
}
`
const testAccKeyPairConfigPrefix = `
resource "alicloud_key_pair" "prefix" {
	key_name_prefix = "terraform-test-key-pair-prefix"
}
`

const testAccKeyPairConfigPublicKey = `
resource "alicloud_key_pair" "publickey" {
  	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`
