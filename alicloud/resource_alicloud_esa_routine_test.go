package alicloud

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA Routine. >>> Resource test cases, automatically generated.
// Case resource_Routine_new_test
func TestAccAliCloudESARoutineresource_Routine_new_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineresource_Routine_new_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutine")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARoutine%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineresource_Routine_new_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test-routine2",
					"name":        "test-routine2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARoutineresource_Routine_new_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARoutineresource_Routine_new_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ESA Routine. <<< Resource test cases, automatically generated.

// esaRoutineWriteCodeFixture writes a small edge routine JS file to a temp path and
// returns the absolute path. It is used to exercise the local filename upload flow.
func esaRoutineWriteCodeFixture(t *testing.T, name, body string) string {
	dir := t.TempDir()
	path := filepath.Join(dir, name+".js")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("failed to write ESA routine code fixture: %v", err)
	}
	return path
}

var AliCloudESARoutineCodeMap = map[string]string{
	"code_checksum":       CHECKSET,
	"latest_code_version": CHECKSET,
	"create_time":         CHECKSET,
}

// TestAccAliCloudESARoutine_code exercises the hand-written code lifecycle: creating a
// routine with a local code file (staging upload + commit), then updating the code to a
// new committed version by changing the local file content.
func TestAccAliCloudESARoutine_code(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineCodeMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutine")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%sesaroutine%d", defaultRegionToTest, rand)

	codeV1 := esaRoutineWriteCodeFixture(t, name+"-v1", "addEventListener('fetch', e => e.respondWith(new Response('v1')))")
	codeV2 := esaRoutineWriteCodeFixture(t, name+"-v2", "addEventListener('fetch', e => e.respondWith(new Response('v2 updated')))")

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineresource_Routine_new_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"description":      "tf-testacc esa routine code v1",
					"filename":         codeV1,
					"code_description": "code version v1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"description":      "tf-testacc esa routine code v1",
						"filename":         codeV1,
						"code_description": "code version v1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             name,
					"description":      "tf-testacc esa routine code v1",
					"filename":         codeV2,
					"code_description": "code version v2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filename":         codeV2,
						"code_description": "code version v2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filename", "code_description", "code_checksum", "latest_code_version"},
			},
		},
	})
}
