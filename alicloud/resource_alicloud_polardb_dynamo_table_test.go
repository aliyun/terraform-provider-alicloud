package alicloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var AliCloudPolarDBDynamoTableMap0 = map[string]string{
	"endpoint":      CHECKSET,
	"db_cluster_id": CHECKSET,
}

func TestAccAliCloudPolarDBDynamoTable_basic(t *testing.T) {
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBDynamoTable%s", rand)
	resourceId := "alicloud_polardb_dynamo_table.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDBDynamoTableMap0)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDynamoTableConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPolarDBDynamoTableDestroy,
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_name":  name,
						"hash_key":    "pk",
						"range_key":   "sk",
						"attribute.#": "4",
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
					"ttl": []map[string]interface{}{
						{
							"enabled":        "true",
							"attribute_name": "expire_at",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"billing_mode":   "PROVISIONED",
						"read_capacity":  "10",
						"write_capacity": "10",
						"attribute.#":    "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// The DynamoDB-compatible DescribeTable does not return
				// BillingModeSummary / ProvisionedThroughput (neither on the table
				// nor on its GSIs), so billing_mode and the capacities cannot be
				// rebuilt from the API on import. The GSI set is hashed on the
				// index name only (schema.HashString("idx_gsi_v2") = 1579187990),
				// so the flat state keys of its capacity fields are stable and
				// can be ignored precisely.
				ImportStateVerifyIgnore: []string{"endpoint", "account_name", "account_auth", "billing_mode", "read_capacity", "write_capacity", "global_secondary_index.1579187990.read_capacity", "global_secondary_index.1579187990.write_capacity"},
			},
		},
	})
}

func testAccCheckPolarDBDynamoTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_dynamo_table" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		dbClusterId := parts[0]
		tableName := parts[1]
		endpoint := rs.Primary.Attributes["endpoint"]
		accountName := rs.Primary.Attributes["account_name"]
		accountAuth := rs.Primary.Attributes["account_auth"]

		// When the whole cluster is already destroyed the table is gone with it,
		// and its endpoint is no longer resolvable, so only probe the table when
		// the cluster still exists.
		if _, err := polarDBService.DescribePolarDBClusterAttribute(dbClusterId); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		dynamoClient, err := client.NewPolarDBDynamoClient(endpoint, accountName, accountAuth)
		if err != nil {
			return WrapError(err)
		}

		_, err = dynamoClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err == nil {
			return fmt.Errorf("DynamoDB table %s still exists", tableName)
		}
		if !isDynamoNotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

// resourcePolarDBDynamoTableConfigDependence generates the shared infrastructure config:
// VPC, VSwitch, PolarDB Cluster (DynamoDB enabled), DynamoDB Account, DynamoDB Endpoint + Public Address
func resourcePolarDBDynamoTableConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  db_type       = "PostgreSQL"
  db_version    = "16"
  pay_type      = "PostPaid"
  db_node_class = "polar.pg.x4.medium"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type                    = "PostgreSQL"
  db_version                 = "16"
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
