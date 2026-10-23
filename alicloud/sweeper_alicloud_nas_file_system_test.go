package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The alicloud_nas_file_system sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 and plugin-testing
// helper/resource packages cannot both register -sweep in one binary. The
// resource implementation itself lives in alicloud/service/nas.
func init() {
	resource.AddTestSweepers("alicloud_nas_file_system",
		&resource.Sweeper{
			Name: "alicloud_nas_file_system",
			F:    testSweepNasFileSystem,
		})
}

func testSweepNasFileSystem(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeFileSystems"
	request := make(map[string]interface{})
	request["RegionId"] = client.Region
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	ids := make([]string, 0)
	for {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Error retrieving filesystem: %s", err)
		}
		resp, err := jsonpath.Get("$.FileSystems.FileSystem", response)
		if err != nil {
			log.Println("Get $.FileSystems.FileSystem failed. err:", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			description, _ := item["Description"].(string)
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping FileSystem: %s (%s)", description, item["FileSystemId"])
					continue
				}
			}
			// 删除 fileSystem 时需要先删除其挂载关系
			if v, ok := item["MountTargets"].(map[string]interface{})["MountTarget"].([]interface{}); ok && len(v) > 0 {
				log.Printf("[INFO] Delete mount targets with filesystem: %v", item["FileSystemId"])
				for _, domain := range v {
					domainInfo := domain.(map[string]interface{})
					request := map[string]interface{}{
						"FileSystemId":      item["FileSystemId"],
						"MountTargetDomain": domainInfo["MountTargetDomain"],
					}
					action := "DeleteMountTarget"
					response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
					if err != nil {
						log.Printf("[ERROR] Error delete mount target: %v with filesystem: %v err: %v", domainInfo["MountTargetDomain"], item["FileSystemId"], err)
					}
				}
			}
			ids = append(ids, item["FileSystemId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, filesystemId := range ids {
		request := map[string]interface{}{
			"FileSystemId": filesystemId,
		}
		action := "DeleteFileSystem"
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Error delete filesystem: %s err: %v", filesystemId, err)
		}
	}
	return nil
}
