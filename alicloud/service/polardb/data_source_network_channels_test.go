package polardb_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The data source is read-only: use an existing source cluster and channel
// rather than creating billed clusters or altering replication during this test.
func TestAccAliCloudPolarDBNetworkChannels_basic(t *testing.T) {
	clusterID := os.Getenv("POLARDB_NETWORK_CHANNEL_TEST_CLUSTER_ID")
	secondClusterID := os.Getenv("POLARDB_NETWORK_CHANNEL_TEST_SECOND_CLUSTER_ID")
	channelName := os.Getenv("POLARDB_NETWORK_CHANNEL_TEST_CHANNEL_NAME")
	vpcID := os.Getenv("POLARDB_NETWORK_CHANNEL_TEST_VPC_ID")
	resourceGroupID := os.Getenv("POLARDB_NETWORK_CHANNEL_TEST_RESOURCE_GROUP_ID")
	if clusterID == "" || secondClusterID == "" || channelName == "" || vpcID == "" || resourceGroupID == "" {
		t.Skip("set POLARDB_NETWORK_CHANNEL_TEST_CLUSTER_ID, POLARDB_NETWORK_CHANNEL_TEST_SECOND_CLUSTER_ID, POLARDB_NETWORK_CHANNEL_TEST_CHANNEL_NAME, POLARDB_NETWORK_CHANNEL_TEST_VPC_ID and POLARDB_NETWORK_CHANNEL_TEST_RESOURCE_GROUP_ID to existing network channel fixtures")
	}
	if clusterID == secondClusterID {
		t.Fatal("the cluster-switch step requires two distinct existing clusters")
	}
	id := "data.alicloud_polardb_network_channels.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(id, "tf-acc-network-channels", func(string) string { return "" })
	// This test only reads a data source; there is no managed object to destroy.
	// lintignore: AT001
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": clusterID,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(id, "db_cluster_id", clusterID),
					resource.TestCheckResourceAttrSet(id, "channels.0.channel_name"),
					resource.TestCheckResourceAttrSet(id, "names.0"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":     clusterID,
					"channel_name":      channelName,
					"vpc_id":            vpcID,
					"resource_group_id": resourceGroupID,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(id, "channels.#", "1"),
					resource.TestCheckResourceAttr(id, "channels.0.channel_name", channelName),
					resource.TestCheckResourceAttr(id, "channels.0.db_cluster_id", clusterID),
					resource.TestCheckResourceAttr(id, "names.#", "1"),
					resource.TestCheckResourceAttr(id, "names.0", channelName),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id": secondClusterID,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(id, "db_cluster_id", secondClusterID),
					resource.TestCheckResourceAttrSet(id, "channels.#"),
					resource.TestCheckResourceAttrSet(id, "names.#"),
				),
			},
		},
	})
}
