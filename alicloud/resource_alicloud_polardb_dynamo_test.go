package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var AliCloudPolarDBDynamoMap0 = map[string]string{
	"endpoint":      CHECKSET,
	"db_cluster_id": CHECKSET,
}

func TestAccAliCloudPolarDBDynamo_basic(t *testing.T) {
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBDynamo%s", rand)
	resourceId := "alicloud_polardb_dynamo.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDBDynamoMap0)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDynamoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPolarDBDynamoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint":       "http://${alicloud_polardb_endpoint_address.dynamo_public.connection_string}:5432",
					"db_cluster_id":  "${alicloud_polardb_cluster.default.id}",
					"account_name":   "${alicloud_polardb_account.dynamo.account_name}",
					"account_auth":   "${alicloud_polardb_account.dynamo.dynamodb_auth_password}",
					"table_name":     name,
					"hash_key":       "pk",
					"range_key":      "sk",
					"read_capacity":  "5",
					"write_capacity": "5",
					"attribute": []map[string]interface{}{
						{
							"name": "pk",
							"type": "S",
						},
						{
							"name": "sk",
							"type": "S",
						},
						{
							"name": "lsi_sk",
							"type": "N",
						},
						{
							"name": "gsi_pk",
							"type": "S",
						},
					},
					"local_secondary_index": []map[string]interface{}{
						{
							"name":               "idx_lsi",
							"range_key":          "lsi_sk",
							"projection_type":    "INCLUDE",
							"non_key_attributes": []string{"detail"},
						},
					},
					"global_secondary_index": []map[string]interface{}{
						{
							"name":            "idx_gsi",
							"hash_key":        "gsi_pk",
							"range_key":       "sk",
							"projection_type": "KEYS_ONLY",
							"read_capacity":   "5",
							"write_capacity":  "5",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "PolarDBDynamo",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_name":  name,
						"hash_key":    "pk",
						"range_key":   "sk",
						"attribute.#": "4",
						"tags.%":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"billing_mode":   "PROVISIONED",
					"read_capacity":  "10",
					"write_capacity": "10",
					"attribute": []map[string]interface{}{
						{
							"name": "pk",
							"type": "S",
						},
						{
							"name": "sk",
							"type": "S",
						},
						{
							"name": "lsi_sk",
							"type": "N",
						},
						{
							"name": "gsi_pk",
							"type": "S",
						},
						{
							"name": "expire_at",
							"type": "N",
						},
					},
					"global_secondary_index": []map[string]interface{}{
						{
							"name":               "idx_gsi_v2",
							"hash_key":           "expire_at",
							"range_key":          "gsi_pk",
							"projection_type":    "INCLUDE",
							"non_key_attributes": []string{"detail"},
							"read_capacity":      "10",
							"write_capacity":     "10",
						},
					},
					"stream_enabled":   "true",
					"stream_view_type": "NEW_AND_OLD_IMAGES",
					"ttl": []map[string]interface{}{
						{
							"enabled":        "true",
							"attribute_name": "expire_at",
						},
					},
					"point_in_time_recovery": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"server_side_encryption": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "PolarDBDynamo-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"billing_mode":   "PROVISIONED",
						"read_capacity":  "10",
						"write_capacity": "10",
						"tags.%":         "2",
						"tags.Created":   "TF-update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"endpoint", "account_name", "account_auth"},
			},
		},
	})
}

func testAccCheckPolarDBDynamoDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_dynamo" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			log.Printf("[WARN] testAccCheckPolarDBDynamoDestroy: failed to parse resource ID %s: %s", rs.Primary.ID, err)
			continue
		}
		tableName := parts[1]
		endpoint := rs.Primary.Attributes["endpoint"]
		accountName := rs.Primary.Attributes["account_name"]
		accountAuth := rs.Primary.Attributes["account_auth"]

		dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accountName, accountAuth)
		if err != nil {
			log.Printf("[WARN] testAccCheckPolarDBDynamoDestroy: failed to create dynamo client for table %s: %s", tableName, err)
			continue
		}

		_, err = dynamoClient.DescribeTable(&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err == nil {
			return fmt.Errorf("DynamoDB table %s still exists", tableName)
		}
		if !isDynamoNotFoundError(err) {
			log.Printf("[WARN] testAccCheckPolarDBDynamoDestroy: DescribeTable for %s returned non-not-found error (treating as destroyed): %s", tableName, err)
		}
	}
	return nil
}

// resourcePolarDBDynamoConfigDependence generates the shared infrastructure config:
// VPC, VSwitch, PolarDB Cluster (DynamoDB enabled), DynamoDB Account, DynamoDB Endpoint + Public Address
func resourcePolarDBDynamoConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.7.id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type                    = "PostgreSQL"
  db_version                 = "14"
  db_node_class              = "polar.pg.x4.medium"
  pay_type                   = "PostPaid"
  vswitch_id                 = alicloud_vswitch.default.id
  description                = var.name
  enable_dynamodb            = true
  global_security_group_list = [alicloud_polardb_global_security_ip_group.default.id]

  depends_on = [alicloud_polardb_global_security_ip_group.default]
}

resource "alicloud_polardb_global_security_ip_group" "default" {
  global_ip_group_name = "tf_dynamo_whitelist"
  global_ip_list       = "0.0.0.0/0"
}

resource "alicloud_polardb_account" "dynamo" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "tf_dynamo_acc"
  account_password = "TfTestDynamo2026!"
  account_type     = "DynamoDB"
}

resource "alicloud_polardb_endpoint" "dynamo" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  endpoint_type = "DynamoDB"
  read_write_mode = "ReadWrite"

  depends_on = [alicloud_polardb_account.dynamo]
}

resource "alicloud_polardb_endpoint_address" "dynamo_public" {
  db_cluster_id  = alicloud_polardb_cluster.default.id
  db_endpoint_id = alicloud_polardb_endpoint.dynamo.db_endpoint_id
  net_type       = "Public"

  depends_on = [alicloud_polardb_endpoint.dynamo]
}
`, name)
}
