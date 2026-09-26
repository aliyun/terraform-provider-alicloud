// Package alicloud. This file is hand-written for PaiDsw TempFileTask.
package alicloud

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPaiDswTempFileTask_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_dsw_temp_file_task.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiDswTempFileTaskMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiDswService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiDswTempFileTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-tftask%d", rand)
	// Build an env-based client via sharedClientForRegion so the PAI DSW
	// Instance fixture can be created BEFORE resource.Test runs. This avoids
	// the chicken-and-egg problem where testAccProvider.Meta() is nil until
	// the provider Configure executes inside resource.Test (PreCheck runs
	// before Configure, so touching Meta() in PreCheck always panics).
	testAccPreCheck(t)
	rawClient, err := sharedClientForRegion(defaultRegionToTest)
	if err != nil {
		t.Fatalf("failed to build shared client for region %s: %s", defaultRegionToTest, err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	instanceId := createPaiDswTempFileTaskInstanceForTest(t, client, name)
	t.Cleanup(func() { deletePaiDswTempFileTaskInstanceForTest(client, instanceId) })

	expiredTime := time.Now().Add(72 * time.Hour).UTC().Format("2006-01-02T15:04:05Z")
	updatedExpiredTime := time.Now().Add(168 * time.Hour).UTC().Format("2006-01-02T15:04:05Z")

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiDswTempFileTaskBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      instanceId,
					"gmt_expired_time": expiredTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":      instanceId,
						"gmt_expired_time": expiredTime,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      instanceId,
					"gmt_expired_time": updatedExpiredTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gmt_expired_time": updatedExpiredTime,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      instanceId,
					"gmt_expired_time": updatedExpiredTime,
				}),
			},
		},
	})
}

// createPaiDswTempFileTaskInstanceForTest creates a PAI DSW Instance fixture.
// It accepts an env-based client built via sharedClientForRegion so it can be
// called from the test body BEFORE resource.Test (PreCheck runs before the
// provider Configure, so testAccProvider.Meta() would be nil there).
func createPaiDswTempFileTaskInstanceForTest(t *testing.T, client *connectivity.AliyunClient, name string) string {
	action := "/api/v2/instances"
	query := make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	body := map[string]interface{}{
		"InstanceName":  name,
		"EcsSpec":       "ecs.c6.large",
		"ImageId":       "image-3a6z2s7hs26zw067yo",
		"Accessibility": "PRIVATE",
	}
	response, err := client.RoaPost("pai-dsw", "2022-01-01", action, query, nil, body, true)
	addDebug(action, response, body)
	if err != nil {
		t.Fatalf("Failed to create PAI DSW instance for test: %s", err)
	}
	if instanceId, ok := response["InstanceId"].(string); ok && instanceId != "" {
		return instanceId
	}
	t.Fatalf("Failed to extract InstanceId from CreateInstance response: %v", response)
	return ""
}

// deletePaiDswTempFileTaskInstanceForTest deletes the Instance fixture.
// It takes the same env-based client built via sharedClientForRegion, so it
// is safe to invoke from a t.Cleanup registered before resource.Test.
func deletePaiDswTempFileTaskInstanceForTest(client *connectivity.AliyunClient, instanceId string) {
	if instanceId == "" {
		return
	}
	action := fmt.Sprintf("/api/v2/instances/%s", instanceId)
	query := make(map[string]*string)
	query["RegionId"] = StringPointer(client.RegionId)
	_, err := client.RoaDelete("pai-dsw", "2022-01-01", action, query, nil, nil, true)
	if err != nil {
		log.Printf("[WARN] Failed to delete PAI DSW instance %s: %s", instanceId, err)
	}
}

var AliCloudPaiDswTempFileTaskMap = map[string]string{
	"instance_id":       CHECKSET,
	"gmt_expired_time":  CHECKSET,
	"temp_file_task_id": CHECKSET,
	"owner_id":          CHECKSET,
	"user_id":           CHECKSET,
	"create_time":       CHECKSET,
	"gmt_modified_time": CHECKSET,
	"region_id":         CHECKSET,
}

func AliCloudPaiDswTempFileTaskBasicDependence(name string) string {
	return ""
}
