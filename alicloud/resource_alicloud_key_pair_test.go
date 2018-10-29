package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_key_pair", &resource.Sweeper{
		Name: "alicloud_key_pair",
		F:    testSweepKeyPairs,
	})
}

func testSweepKeyPairs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
		"terraform-test-",
	}

	var pairs []ecs.KeyPair
	req := ecs.CreateDescribeKeyPairsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeKeyPairs(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Key Pairs: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeKeyPairsResponse)
		if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
			break
		}
		pairs = append(pairs, resp.KeyPairs.KeyPair...)

		if len(resp.KeyPairs.KeyPair) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range pairs {
		name := v.KeyPairName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Key Pair: %s", name)
			continue
		}
		log.Printf("[INFO] Deleting Key Pair: %s", name)
		req := ecs.CreateDeleteKeyPairsRequest()
		req.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), name))
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Key Pair (%s): %s", name, err)
		}
	}
	return nil
}

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
						"alicloud_key_pair.prefix", &keypair, "tf-testAccKeyPairConfigPrefix"),
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
						"alicloud_key_pair.publickey", &keypair, "tf-testAccKeyPairConfigPublicKey"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeKeyPair(rs.Primary.ID)

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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeKeyPair(rs.Primary.ID)

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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}

		response, err := ecsService.DescribeKeyPair(rs.Primary.ID)
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
	key_name = "tf-testAccKeyPairConfig"
}
`
const testAccKeyPairConfigPrefix = `
resource "alicloud_key_pair" "prefix" {
	key_name_prefix = "tf-testAccKeyPairConfigPrefix"
}
`

const testAccKeyPairConfigPublicKey = `
resource "alicloud_key_pair" "publickey" {
	key_name = "tf-testAccKeyPairConfigPublicKey"
  	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`
