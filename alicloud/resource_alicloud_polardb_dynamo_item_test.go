package alicloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var AliCloudPolarDBDynamoItemMap0 = map[string]string{
	"endpoint":      CHECKSET,
	"db_cluster_id": CHECKSET,
	"item":          CHECKSET,
}

func TestAccAliCloudPolarDBDynamoItem_basic(t *testing.T) {
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBDynamoItem-%s", rand)
	resourceId := "alicloud_polardb_dynamo_item.default"
	ra := resourceAttrInit(resourceId, AliCloudPolarDBDynamoItemMap0)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBDynamoItemConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPolarDBDynamoItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint":      "http://${alicloud_polardb_endpoint_address.dynamo_public.connection_string}:5432",
					"db_cluster_id": "${alicloud_polardb_cluster.default.id}",
					"account_name":  "${alicloud_polardb_account.dynamo.account_name}",
					"account_auth":  "${alicloud_polardb_account.dynamo.dynamodb_auth_password}",
					"table_name":    "${alicloud_polardb_dynamo_table.default.table_name}",
					"hash_key":      "pk",
					"range_key":     "sk",
					"item":          `{\"pk\": {\"S\": \"test-item-1\"}, \"sk\": {\"S\": \"row1\"}, \"name\": {\"S\": \"Test Item\"}, \"count\": {\"N\": \"42\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_name": name,
						"hash_key":   "pk",
						"range_key":  "sk",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"item": `{\"pk\": {\"S\": \"test-item-1\"}, \"sk\": {\"S\": \"row1\"}, \"name\": {\"S\": \"Updated Test Item\"}, \"count\": {\"N\": \"100\"}, \"active\": {\"BOOL\": true}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"table_name": name,
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

func testAccCheckPolarDBDynamoItemDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_dynamo_item" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 4)
		if err != nil {
			parts, err = ParseResourceId(rs.Primary.ID, 3)
			if err != nil {
				return WrapError(err)
			}
		}
		dbClusterId := parts[0]
		tableName := parts[1]
		endpoint := rs.Primary.Attributes["endpoint"]
		hashKey := rs.Primary.Attributes["hash_key"]
		rangeKey := rs.Primary.Attributes["range_key"]
		accountName := rs.Primary.Attributes["account_name"]
		accountAuth := rs.Primary.Attributes["account_auth"]

		// When the whole cluster is already destroyed the item is gone with it,
		// and its endpoint is no longer resolvable, so only probe the item when
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

		key := map[string]types.AttributeValue{
			hashKey: &types.AttributeValueMemberS{Value: parts[2]},
		}
		if rangeKey != "" && len(parts) > 3 {
			key[rangeKey] = &types.AttributeValueMemberS{Value: parts[3]}
		}

		output, err := dynamoClient.GetItem(context.Background(), &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       key,
		})
		if err == nil && len(output.Item) > 0 {
			return fmt.Errorf("DynamoDB item still exists in table %s", tableName)
		}
		if err != nil && !isDynamoNotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

// resourcePolarDBDynamoItemConfigDependence generates shared infra for item tests:
// VPC, VSwitch, PolarDB Cluster (DynamoDB enabled), DynamoDB Account, Endpoint + Public Address, DynamoDB Table
func resourcePolarDBDynamoItemConfigDependence(name string) string {
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

resource "alicloud_polardb_dynamo_table" "default" {
  endpoint      = "http://${alicloud_polardb_endpoint_address.dynamo_public.connection_string}:5432"
  db_cluster_id = alicloud_polardb_cluster.default.id
  account_name  = alicloud_polardb_account.dynamo.account_name
  account_auth  = alicloud_polardb_account.dynamo.dynamodb_auth_password
  table_name    = var.name
  hash_key      = "pk"
  range_key     = "sk"
  billing_mode  = "PAY_PER_REQUEST"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }
}
`, name)
}
