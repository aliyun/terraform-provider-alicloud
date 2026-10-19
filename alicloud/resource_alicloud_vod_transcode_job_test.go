package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudVodTranscodeJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_transcode_job.default"
	checkoutSupportedRegions(t, true, connectivity.VODSupportRegions)
	ra := resourceAttrInit(resourceId, AliCloudVodTranscodeJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodTranscodeJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svodtranscodejob%d", defaultRegionToTest, rand)
	videoId := os.Getenv("ALICLOUD_VOD_VIDEO_ID")
	if videoId == "" {
		t.Skip("Skipping the test case without ALICLOUD_VOD_VIDEO_ID")
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVodTranscodeJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"video_id":               "${var.video_id}",
					"reference_id":           "${var.name}",
					"height":                 "480",
					"width":                  "640",
					"specified_offset_time":  0,
					"count":                  10,
					"interval":               10,
					"user_data":              "tf-testacc-userdata",
					"specified_offset_times": "[0, 1000, 2000]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"video_id":     videoId,
						"reference_id": name,
						"height":       "480",
						"width":        "640",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"video_id":               "${var.video_id}",
					"reference_id":           "${var.name}",
					"snapshot_template_id":   "${var.snapshot_template_id}",
					"sprite_snapshot_config": "{\"CellWidth\":\"100\",\"CellHeight\":\"100\",\"Columns\":\"10\",\"Lines\":\"10\",\"Padding\":\"10\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"video_id":             videoId,
						"snapshot_template_id": "tf-testacc-snapshot-template",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reference_id", "height", "width", "specified_offset_time", "count", "interval", "sprite_snapshot_config", "snapshot_template_id", "user_data", "specified_offset_times"},
			},
		},
	})
}

var AliCloudVodTranscodeJobMap0 = map[string]string{
	"transcode_job_id": CHECKSET,
	"video_id":         CHECKSET,
	"create_time":      CHECKSET,
	"region_id":        CHECKSET,
}

func AliCloudVodTranscodeJobBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "video_id" {
  default = "%s"
}
variable "snapshot_template_id" {
  default = "tf-testacc-snapshot-template"
}
`, name, os.Getenv("ALICLOUD_VOD_VIDEO_ID"))
}

// lintignore: R001
func TestUnitAccAlicloudVodTranscodeJob(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_vod_transcode_job"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_vod_transcode_job"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"video_id":               "CreateVodTranscodeJobValue",
		"reference_id":           "CreateVodTranscodeJobValue",
		"height":                 "480",
		"width":                  "640",
		"specified_offset_time":  0,
		"count":                  10,
		"interval":               10,
		"sprite_snapshot_config": "CreateVodTranscodeJobValue",
		"snapshot_template_id":   "CreateVodTranscodeJobValue",
		"user_data":              "CreateVodTranscodeJobValue",
		"specified_offset_times": []interface{}{0, 1000, 2000},
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
		"TranscodeTask": map[string]interface{}{
			"CreationTime": "2024-01-01T00:00:00Z",
			"VideoId":      "CreateVodTranscodeJobValue",
			"TranscodeJobInfoList": map[string]interface{}{
				"TranscodeJobInfo": []interface{}{
					map[string]interface{}{
						"TranscodeJobId": "CreateVodTranscodeJobValue",
					},
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"SnapshotJob": map[string]interface{}{
			"JobId": "CreateVodTranscodeJobValue",
		},
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_vod_transcode_job", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVodClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudVodTranscodeJobCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "SubmitSnapshotJob" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudVodTranscodeJobCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vod_transcode_job"].Schema).Data(dInit.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetTranscodeTask" {
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
		err := resourceAlicloudVodTranscodeJobRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudVodTranscodeJobDelete(dExisted, rawClient)
	assert.Nil(t, err)
}
