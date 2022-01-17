package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOssBucketReplication(t *testing.T) {
	var v oss.GetBucketReplicationResult

	resourceId := "alicloud_oss_bucket_replication.default"
	ra := resourceAttrInit(resourceId, ossBucketReplicationMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-replication-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketReplicationDependence)
	preHashcode := strconv.Itoa(prefixSetHash(map[string]interface{}{
		"prefixes": []string{
			"test/",
			"src/",
		},
	}))
	desHashcode := strconv.Itoa(destinationHash(map[string]interface{}{
		"bucket":        name + "-t",
		"location":      "oss-cn-beijing",
		"transfer_type": string(oss.TransferInternal),
	}))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replication_rule": []map[string]interface{}{
						{
							"action":                        string(oss.ReplicationActionALL),
							"historical_object_replication": string(oss.HistoricalEnabled),
							"prefix_set": []map[string]interface{}{
								{
									"prefixes": []string{
										"test/",
										"src/",
									},
								},
							},
							"destination": []map[string]string{
								{
									"bucket":        name + "-t",
									"location":      "oss-cn-beijing", // todo 这个也改成动态的
									"transfer_type": string(oss.TransferInternal),
								},
							},
							//"source_selection_criteria": []map[string]interface{}{
							//	{
							//		"sse_kms_encrypted_objects": []map[string]string{
							//			{
							//				"status": string(oss.SourceSSEKMSEnabled),
							//			},
							//		},
							//	},
							//},
							//"encryption_configuration": []map[string]string{
							//	{
							//		"replica_kms_key_id": "1a8d780d-0d34-49e5-ab45-7b9b5fcca954",
							//	},
							//},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replication_rule.#":                                               "1",
						"replication_rule.0.action":                                        string(oss.ReplicationActionALL),
						"replication_rule.0.historical_object_replication":                 string(oss.HistoricalEnabled),
						"replication_rule.0.prefix_set." + preHashcode + ".prefixes.#":     "2",
						"replication_rule.0.destination." + desHashcode + ".bucket":        name + "-t",
						"replication_rule.0.destination." + desHashcode + ".location":      "oss-cn-beijing",
						"replication_rule.0.destination." + desHashcode + ".transfer_type": string(oss.TransferInternal),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replication_rule": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replication_rule.#":                                               "0",
						"replication_rule.0.action":                                        REMOVEKEY,
						"replication_rule.0.historical_object_replication":                 REMOVEKEY,
						"replication_rule.0.prefix_set." + preHashcode + ".prefixes.#":     "0",
						"replication_rule.0.destination." + desHashcode + ".bucket":        REMOVEKEY,
						"replication_rule.0.destination." + desHashcode + ".location":      REMOVEKEY,
						"replication_rule.0.destination." + desHashcode + ".transfer_type": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourceOssBucketReplicationDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket_replication" "dst"{
	bucket = "%s-t"
}
`, name)
}

var ossBucketReplicationMap = map[string]string{
	"replication_rule.#": "0",
}
