package alicloud

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ecs_image_pipeline",
		&resource.Sweeper{
			Name: "alicloud_ecs_image_pipeline",
			F:    testSweepEcsImagePipeline,
		})
}

func testSweepEcsImagePipeline(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeImagePipelines"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId
	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.ImagePipeline.ImagePipelineSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ImagePipeline.ImagePipelineSet", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs Image Pipeline: %s", item["Name"].(string))
				continue
			}
			action := "DeleteImagePipeline"
			request := map[string]interface{}{
				"ImagePipelineId": item["ImagePipelineId"],
				"RegionId":        aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs Image Pipeline (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs Image Pipeline success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudECSImagePipeline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type": "IMAGE",
					"base_image":      "${data.alicloud_images.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type": "IMAGE",
						"base_image":      CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_content"},
			},
		},
	})
}
func TestAccAliCloudECSImagePipeline_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", defaultRegionToTest, rand)
	prepareEcsImagePipelineFixture(t)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type":            "OSS",
					"name":                       "${var.name}",
					"description":                "${var.name}",
					"system_disk_size":           "40",
					"vswitch_id":                 "${alicloud_vswitch.default.id}",
					"image_name":                 "${var.name}",
					"delete_instance_on_failure": "true",
					"internet_max_bandwidth_out": "20",
					"to_region_id":               []string{"cn-qingdao", "cn-zhangjiakou"},
					"add_account":                []string{"${data.alicloud_account.default.id}"},
					"instance_type":              "${data.alicloud_instance_types.default.ids.0}",
					"build_content":              "RUN yum update -y",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags":                       map[string]interface{}{"Created": "TF", "Foo": "Bar"},
					"test_content":               "RUN echo tf-testacc-imagepipeline",
					"image_family":               "tf-testacc-image-family",
					"repair_mode":                "Standard",
					"nvme_support":               "auto",
					"advanced_options": []map[string]interface{}{
						{
							"retain_cloud_assistant": true,
							"image_name_suffix":      "disable",
						},
					},
					"import_image_options": []map[string]interface{}{
						{
							"architecture":          "x86_64",
							"platform":              "Ubuntu",
							"boot_mode":             "UEFI",
							"license_type":          "Auto",
							"os_type":               "linux",
							"retain_imported_image": true,
							"features": []map[string]interface{}{
								{
									"nvme_support": "supported",
								},
							},
							"disk_device_mappings": []map[string]interface{}{
								{
									"oss_bucket":      "${alicloud_oss_bucket.default.id}",
									"oss_object":      "${alicloud_oss_bucket_object.image.id}",
									"format":          "QCOW2",
									"disk_image_size": 40,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"base_image_type":            "OSS",
						"description":                name,
						"system_disk_size":           "40",
						"vswitch_id":                 CHECKSET,
						"image_name":                 name,
						"delete_instance_on_failure": "true",
						"internet_max_bandwidth_out": "20",
						"to_region_id.#":             "2",
						"add_account.#":              "1",
						"instance_type":              CHECKSET,
						"resource_group_id":          CHECKSET,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.Foo":                   "Bar",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_content"},
			},
		},
	})
}

// collection-order: convergent reorder coverage for the multi-member TypeList
// import_image_options.disk_device_mappings added by this change.
func TestAccAliCloudECSImagePipeline_diskDeviceMappingsOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_image_pipeline.default"
	ra := resourceAttrInit(resourceId, AlicloudECSImagePipelineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsImagePipeline")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsimagepipeline%d", defaultRegionToTest, rand)
	prepareEcsImagePipelineFixture(t)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSImagePipelineBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type":  "OSS",
					"repair_mode":      "Standard",
					"name":             "${var.name}",
					"description":      "${var.name}",
					"system_disk_size": "40",
					"vswitch_id":       "${alicloud_vswitch.default.id}",
					"image_name":       "${var.name}",
					"instance_type":    "${data.alicloud_instance_types.default.ids.0}",
					"import_image_options": []map[string]interface{}{
						{
							"architecture": "x86_64",
							"platform":     "Ubuntu",
							"boot_mode":    "UEFI",
							"license_type": "Auto",
							"os_type":      "linux",
							"disk_device_mappings": []map[string]interface{}{
								{
									"oss_bucket":      "${alicloud_oss_bucket.default.id}",
									"oss_object":      "${alicloud_oss_bucket_object.disk1.id}",
									"format":          "QCOW2",
									"disk_image_size": 40,
								},
								{
									"oss_bucket":      "${alicloud_oss_bucket.default.id}",
									"oss_object":      "${alicloud_oss_bucket_object.disk2.id}",
									"format":          "QCOW2",
									"disk_image_size": 40,
								},
							},
						},
					},
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"base_image_type":  "OSS",
					"repair_mode":      "Standard",
					"name":             "${var.name}",
					"description":      "${var.name}",
					"system_disk_size": "40",
					"vswitch_id":       "${alicloud_vswitch.default.id}",
					"image_name":       "${var.name}",
					"instance_type":    "${data.alicloud_instance_types.default.ids.0}",
					"import_image_options": []map[string]interface{}{
						{
							"architecture": "x86_64",
							"platform":     "Ubuntu",
							"boot_mode":    "UEFI",
							"license_type": "Auto",
							"os_type":      "linux",
							"disk_device_mappings": []map[string]interface{}{
								{
									"oss_bucket":      "${alicloud_oss_bucket.default.id}",
									"oss_object":      "${alicloud_oss_bucket_object.disk2.id}",
									"format":          "QCOW2",
									"disk_image_size": 40,
								},
								{
									"oss_bucket":      "${alicloud_oss_bucket.default.id}",
									"oss_object":      "${alicloud_oss_bucket_object.disk1.id}",
									"format":          "QCOW2",
									"disk_image_size": 40,
								},
							},
						},
					},
				}),
			},
		},
	})
}

var AlicloudECSImagePipelineMap0 = map[string]string{}

func AlicloudECSImagePipelineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
data "alicloud_account" "default" {}
`, name)
}

// ecsImagePipelineCloudImageURL is the Canonical Ubuntu Minimal 22.04 (jammy)
// release-20260818 amd64 cloud image (QCOW2, UEFI/GPT, bootable). The image is
// fully downloaded and its SHA256 is verified against the hardcoded trust anchor
// below (file integrity). Ubuntu's GPG-signed SHA256SUMS is NOT fetched or
// verified on the runner, so authenticity beyond the hardcoded anchor remains
// unverified; the hardcoded SHA256 is the only integrity guarantee asserted here.
// The downloaded qcow2 is uploaded to OSS as-is and imported with Format=QCOW2,
// which matches the real file format and is a server-accepted enum
// (CreateImagePipeline DiskDeviceMappings.Format ∈ {RAW, VHD, QCOW2}); no
// qemu-img conversion is performed, so the fixture does not depend on qemu-img
// being installed on the ACC runner. The prior RAW-conversion path is removed;
// RAW-format import is no longer covered by this fixture (gap).
const (
	ecsImagePipelineCloudImageURL    = "https://cloud-images.ubuntu.com/minimal/releases/jammy/release-20260818/ubuntu-22.04-minimal-cloudimg-amd64.img"
	ecsImagePipelineCloudImageSHA256 = "bc866170b04b0e2aedaa7baabd9baee6199d8c296dd174e1c11af43906aa28eb"
)

// ecsImagePipelineFixturePath holds the SHA256-verified Canonical Ubuntu qcow2
// downloaded by prepareEcsImagePipelineFixture into t.TempDir(), so the ~310MB
// qcow2 never enters the bundle, tool package or checkpoint. The same verified
// qcow2 backs all three private OSS object keys (image/disk1/disk2): identical
// bytes across keys only prove the CreateImagePipeline request/config ordering of
// disk_device_mappings, not that different system contents were imported.
var ecsImagePipelineFixturePath string

// prepareEcsImagePipelineFixture materializes a real, importable disk image in
// t.TempDir(): download the Ubuntu cloud image (qcow2) and verify its SHA256
// against the hardcoded anchor. The qcow2 is used directly with Format=QCOW2 —
// no qemu-img check/convert — so the fixture runs on ACC runners without
// qemu-img. When the download is unavailable on the runner, it t.Skip's honestly
// rather than silently falling back to a placeholder — the fixture and ImportImage
// are then UNVERIFIED, not PASS. The actual cloud import
// (StartImagePipelineExecution) is a separate paid scope and is never claimed as
// verified by this fixture. RAW-format import is not covered (the prior
// qemu-img convert-to-raw path was removed to eliminate the qemu-img dependency);
// see the const comment above.
func prepareEcsImagePipelineFixture(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	qcow2 := filepath.Join(dir, "ubuntu-minimal.qcow2")
	if err := downloadAndVerifySHA256(ecsImagePipelineCloudImageURL, qcow2, ecsImagePipelineCloudImageSHA256); err != nil {
		t.Skipf("Ubuntu cloud image download/SHA256 verify failed; qcow2 fixture + import UNVERIFIED: %v", err)
	}
	ecsImagePipelineFixturePath = qcow2
}

// downloadAndVerifySHA256 downloads imageURL to dst and fails unless the SHA256
// matches wantHex.
func downloadAndVerifySHA256(imageURL, dst, wantHex string) error {
	client := &http.Client{Timeout: 30 * time.Minute}
	resp, err := client.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", imageURL, resp.StatusCode)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		return err
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != wantHex {
		return fmt.Errorf("sha256 mismatch for %s: got %s want %s", imageURL, got, wantHex)
	}
	return nil
}

// AlicloudECSImagePipelineBasicDependence1 is a self-contained fixture for the
// import_image_options / disk_device_mappings tests: it dynamically provisions a
// VPC, a vSwitch and a private OSS bucket whose objects reference the real qcow2
// produced by prepareEcsImagePipelineFixture (source, not a content placeholder),
// so basic1 and diskDeviceMappingsOrder no longer depend on a pre-existing
// default-NODELETING VPC or on persistent OSS objects. ImportImageOptions
// parameters are consistent with the imported Ubuntu image: Platform=Ubuntu,
// Architecture=x86_64, BootMode=UEFI, Format=QCOW2 (matches the real file
// format), system/disk size 40GiB >= the ~310MB qcow2 source. NOTE: this
// fixture proves the CreateImagePipeline template CRUD
// (create/read/import-state) and the disk_device_mappings TypeList ordering; the
// actual cloud import (StartImagePipelineExecution) is a separate paid scope and is
// UNVERIFIED. The full RPC mapping and ordering are covered by the no-cloud unit
// tests TestUnitAlicloudEcsImagePipelineExpandCreateRequest and
// TestUnitAlicloudEcsImagePipelineFlattenRead. The bucket uses force_destroy for
// cleanup; failure-path object deletion must be confirmed by real cleanup receipts
// rather than assumed.
func AlicloudECSImagePipelineBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name      = var.name
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  image_id = data.alicloud_images.default.ids.0
}
data "alicloud_account" "default" {}
resource "alicloud_oss_bucket" "default" {
  bucket        = var.name
  force_destroy = true
}
resource "alicloud_oss_bucket_object" "image" {
  bucket = alicloud_oss_bucket.default.id
  key    = "image.qcow2"
  source = "%s"
}
resource "alicloud_oss_bucket_object" "disk1" {
  bucket = alicloud_oss_bucket.default.id
  key    = "disk1.qcow2"
  source = "%s"
}
resource "alicloud_oss_bucket_object" "disk2" {
  bucket = alicloud_oss_bucket.default.id
  key    = "disk2.qcow2"
  source = "%s"
}
`, name, ecsImagePipelineFixturePath, ecsImagePipelineFixturePath, ecsImagePipelineFixturePath)
}

// lintignore: R001
func TestUnitAlicloudEcsImagePipeline(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"base_image_type":            "CreateEcsImagePipelineValue",
		"name":                       "CreateEcsImagePipelineValue",
		"description":                "CreateEcsImagePipelineValue",
		"system_disk_size":           40,
		"vswitch_id":                 "CreateEcsImagePipelineValue",
		"image_name":                 "CreateEcsImagePipelineValue",
		"delete_instance_on_failure": true,
		"internet_max_bandwidth_out": 20,
		"to_region_id":               []string{"CreateEcsImagePipelineValue"},
		"add_account":                []string{"CreateEcsImagePipelineValue"},
		"instance_type":              "CreateEcsImagePipelineValue",
		"base_image":                 "CreateEcsImagePipelineValue",
		"build_content":              "CreateEcsImagePipelineValue",
		"resource_group_id":          "CreateEcsImagePipelineValue",
		"tags":                       map[string]interface{}{"TestKey": "TestValue"},
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"ImagePipeline": map[string]interface{}{
			"ImagePipelineSet": []interface{}{
				map[string]interface{}{
					"Description":     "CreateEcsImagePipelineValue",
					"BaseImageType":   "CreateEcsImagePipelineValue",
					"ResourceGroupId": "CreateEcsImagePipelineValue",
					"SystemDiskSize":  40,
					"ImagePipelineId": "CreateEcsImagePipelineValue",
					"VSwitchId":       "CreateEcsImagePipelineValue",
					"ToRegionIds": map[string]interface{}{
						"ToRegionId": []interface{}{
							"CreateEcsImagePipelineValue",
						},
					},
					"Name":                    "CreateEcsImagePipelineValue",
					"DeleteInstanceOnFailure": true,
					"ImageName":               "CreateEcsImagePipelineValue",
					"AddAccounts": map[string]interface{}{
						"AddAccount": []interface{}{
							"CreateEcsImagePipelineValue",
						},
					},
					"InternetMaxBandwidthOut": 20,
					"CreationTime":            "CreateEcsImagePipelineValue",
					"InstanceType":            "CreateEcsImagePipelineValue",
					"BaseImage":               "CreateEcsImagePipelineValue",
					"BuildContent":            "CreateEcsImagePipelineValue",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"TagKey":   "TestKey",
								"TagValue": "TestValue",
							},
						},
					},
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_image_pipeline", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsImagePipelineCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateImagePipeline" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						_, _ = successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsImagePipelineCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ecs_image_pipeline", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeImagePipelines" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsImagePipelineRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsImagePipelineDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecs_image_pipeline", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteImagePipeline" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsImagePipelineDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// TestUnitAlicloudEcsImagePipelineExpandCreateRequest is a no-cloud unit test for the
// CreateImagePipeline request mapping. It asserts that every schema field maps to
// the correct RPC parameter key and that ImportImageOptions.DiskDeviceMappings
// preserve the schema list order with 1-based indexing — including a reversed-order
// case that proves the index is derived from position rather than a fixed value.
// It invokes no SDK client and no real API; it cannot validate image import, which
// remains an unverified gap (see the acceptance-test fixture comment above).
func TestUnitAlicloudEcsImagePipelineExpandCreateRequest(t *testing.T) {
	p := Provider().ResourcesMap
	buildRequest := func(diskOrder []interface{}) map[string]interface{} {
		d, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
		assert.Nil(t, d.Set("test_content", "tc"))
		assert.Nil(t, d.Set("image_family", "if"))
		assert.Nil(t, d.Set("repair_mode", "Standard"))
		assert.Nil(t, d.Set("nvme_support", "supported"))
		assert.Nil(t, d.Set("advanced_options", []interface{}{
			map[string]interface{}{
				"retain_cloud_assistant": true,
				"image_name_suffix":      "sfx",
			},
		}))
		assert.Nil(t, d.Set("import_image_options", []interface{}{
			map[string]interface{}{
				"architecture":          "x86_64",
				"platform":              "AliLinux",
				"boot_mode":             "BIOS",
				"license_type":          "BYOL",
				"os_type":               "Linux",
				"retain_imported_image": true,
				"features": []interface{}{
					map[string]interface{}{"nvme_support": "supported"},
				},
				"disk_device_mappings": diskOrder,
			},
		}))
		return expandImagePipelineCreateRequest(d, "cn-hangzhou")
	}

	// Every schema field is mapped to the correct RPC parameter key.
	req := buildRequest([]interface{}{
		map[string]interface{}{"oss_bucket": "bkt", "oss_object": "disk1.raw", "format": "raw", "disk_image_size": 40},
		map[string]interface{}{"oss_bucket": "bkt", "oss_object": "disk2.raw", "format": "raw", "disk_image_size": 40},
	})
	assert.Equal(t, "tc", req["TestContent"])
	assert.Equal(t, "if", req["ImageFamily"])
	assert.Equal(t, "Standard", req["RepairMode"])
	assert.Equal(t, "supported", req["NvmeSupport"])
	assert.Equal(t, "cn-hangzhou", req["RegionId"])
	assert.Equal(t, true, req["AdvancedOptions.RetainCloudAssistant"])
	assert.Equal(t, "sfx", req["AdvancedOptions.ImageNameSuffix"])
	assert.Equal(t, "x86_64", req["ImportImageOptions.Architecture"])
	assert.Equal(t, "AliLinux", req["ImportImageOptions.Platform"])
	assert.Equal(t, "BIOS", req["ImportImageOptions.BootMode"])
	assert.Equal(t, "BYOL", req["ImportImageOptions.LicenseType"])
	assert.Equal(t, "Linux", req["ImportImageOptions.OSType"])
	assert.Equal(t, true, req["ImportImageOptions.RetainImportedImage"])
	assert.Equal(t, "supported", req["ImportImageOptions.Features.NvmeSupport"])
	// disk_device_mappings preserve schema order with 1-based indexing.
	assert.Equal(t, "disk1.raw", req["ImportImageOptions.DiskDeviceMappings.1.OSSObject"])
	assert.Equal(t, "disk2.raw", req["ImportImageOptions.DiskDeviceMappings.2.OSSObject"])
	assert.Equal(t, "bkt", req["ImportImageOptions.DiskDeviceMappings.1.OSSBucket"])
	assert.Equal(t, "raw", req["ImportImageOptions.DiskDeviceMappings.1.Format"])
	assert.Equal(t, 40, req["ImportImageOptions.DiskDeviceMappings.1.DiskImageSize"])
	// No 0-based index leaks into the request map.
	_, hasZero := req["ImportImageOptions.DiskDeviceMappings.0.OSSObject"]
	assert.False(t, hasZero)

	// Reversed order proves the index follows position, not a fixed value.
	reqRev := buildRequest([]interface{}{
		map[string]interface{}{"oss_bucket": "bkt", "oss_object": "disk2.raw", "format": "raw", "disk_image_size": 40},
		map[string]interface{}{"oss_bucket": "bkt", "oss_object": "disk1.raw", "format": "raw", "disk_image_size": 40},
	})
	assert.Equal(t, "disk2.raw", reqRev["ImportImageOptions.DiskDeviceMappings.1.OSSObject"])
	assert.Equal(t, "disk1.raw", reqRev["ImportImageOptions.DiskDeviceMappings.2.OSSObject"])
}

// TestUnitAlicloudEcsImagePipelineFlattenRead is a no-cloud unit test for the Read/Importer
// flatten path (Importer uses ImportStatePassthrough → Read, which calls this same
// flatten). It builds a DescribeImagePipelines response object containing the
// nested AdvancedOptions, ImportImageOptions.Features and
// ImportImageOptions.DiskDeviceMappings blocks and asserts the flattened
// ResourceData preserves values and the disk_device_mappings list order.
func TestUnitAlicloudEcsImagePipelineFlattenRead(t *testing.T) {
	p := Provider().ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ecs_image_pipeline"].Schema).Data(nil, nil)
	object := map[string]interface{}{
		"BaseImage":       "m-xxx",
		"ImagePipelineId": "ip-xxx",
		"TestContent":     "tc",
		"ImageFamily":     "if",
		"RepairMode":      "Standard",
		"NvmeSupport":     "supported",
		"AdvancedOptions": map[string]interface{}{
			"RetainCloudAssistant": true,
			"ImageNameSuffix":      "sfx",
		},
		"ImportImageOptions": map[string]interface{}{
			"Architecture":        "x86_64",
			"Platform":            "AliLinux",
			"BootMode":            "BIOS",
			"LicenseType":         "BYOL",
			"OSType":              "Linux",
			"RetainImportedImage": true,
			"Features":            map[string]interface{}{"NvmeSupport": "supported"},
			"DiskDeviceMappings": map[string]interface{}{
				"DiskDeviceMapping": []interface{}{
					map[string]interface{}{"OSSBucket": "bkt", "OSSObject": "disk1.raw", "Format": "raw", "DiskImageSize": 40},
					map[string]interface{}{"OSSBucket": "bkt", "OSSObject": "disk2.raw", "Format": "raw", "DiskImageSize": 40},
				},
			},
		},
	}
	err := flattenImagePipelineRead(d, object)
	assert.Nil(t, err)

	assert.Equal(t, "tc", d.Get("test_content"))
	assert.Equal(t, "if", d.Get("image_family"))
	assert.Equal(t, "Standard", d.Get("repair_mode"))
	assert.Equal(t, "supported", d.Get("nvme_support"))

	advancedOptions := d.Get("advanced_options").([]interface{})
	assert.Equal(t, 1, len(advancedOptions))
	ao := advancedOptions[0].(map[string]interface{})
	assert.Equal(t, true, ao["retain_cloud_assistant"])
	assert.Equal(t, "sfx", ao["image_name_suffix"])

	assert.Equal(t, "x86_64", d.Get("import_image_options.0.architecture"))
	assert.Equal(t, true, d.Get("import_image_options.0.retain_imported_image"))
	features := d.Get("import_image_options.0.features").([]interface{})
	assert.Equal(t, "supported", features[0].(map[string]interface{})["nvme_support"])

	// disk_device_mappings preserve response order through the flatten.
	mappings := d.Get("import_image_options.0.disk_device_mappings").([]interface{})
	assert.Equal(t, 2, len(mappings))
	assert.Equal(t, "disk1.raw", mappings[0].(map[string]interface{})["oss_object"])
	assert.Equal(t, "disk2.raw", mappings[1].(map[string]interface{})["oss_object"])
	assert.Equal(t, "bkt", mappings[0].(map[string]interface{})["oss_bucket"])
	assert.Equal(t, "raw", mappings[0].(map[string]interface{})["format"])
	assert.Equal(t, 40, mappings[0].(map[string]interface{})["disk_image_size"])
}
