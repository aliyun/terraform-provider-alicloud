package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	sdkots "github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Reuse the native attachment resource's APIs, never DeleteInstance: an OTS
// instance can be shared by several VPCs and contains persistent table data.
type vpcOtsAttachmentAPI struct {
	listVpc      func(string, int, int) (*sdkots.ListVpcInfoByVpcResponse, error)
	listInstance func(string, int, int) (*sdkots.ListVpcInfoByInstanceResponse, error)
	unbind       func(string, string) error
	pause        func(time.Duration)
	attempts     int
}

func nativeVpcOtsAttachmentSweeper(region, vpcID string, client *connectivity.AliyunClient) *resource.Sweeper {
	api := vpcOtsAttachmentAPI{
		listVpc: func(id string, page, size int) (*sdkots.ListVpcInfoByVpcResponse, error) {
			request := sdkots.CreateListVpcInfoByVpcRequest()
			request.RegionId = region
			request.Method = "GET"
			request.VpcId = id
			request.PageNum = requests.NewInteger(page)
			request.PageSize = requests.NewInteger(size)
			raw, err := client.WithOtsClient(func(c *sdkots.Client) (interface{}, error) { return c.ListVpcInfoByVpc(request) })
			if err != nil {
				return nil, err
			}
			response, ok := raw.(*sdkots.ListVpcInfoByVpcResponse)
			if !ok {
				return nil, fmt.Errorf("invalid OTS VPC response")
			}
			return response, nil
		},
		listInstance: func(name string, page, size int) (*sdkots.ListVpcInfoByInstanceResponse, error) {
			request := sdkots.CreateListVpcInfoByInstanceRequest()
			request.RegionId = region
			request.Method = "GET"
			request.InstanceName = name
			request.PageNum = requests.NewInteger(page)
			request.PageSize = requests.NewInteger(size)
			raw, err := client.WithOtsClient(func(c *sdkots.Client) (interface{}, error) { return c.ListVpcInfoByInstance(request) })
			if err != nil {
				return nil, err
			}
			response, ok := raw.(*sdkots.ListVpcInfoByInstanceResponse)
			if !ok {
				return nil, fmt.Errorf("invalid OTS instance attachment response")
			}
			return response, nil
		},
		unbind: func(name, binding string) error {
			request := sdkots.CreateUnbindInstance2VpcRequest()
			request.RegionId = region
			request.RegionNo = region
			request.InstanceName = name
			request.InstanceVpcName = binding
			_, err := client.WithOtsClient(func(c *sdkots.Client) (interface{}, error) { return c.UnbindInstance2Vpc(request) })
			return err
		}, pause: time.Sleep, attempts: 120,
	}
	return vpcOtsAttachmentSweeper(region, vpcID, api)
}

func vpcOtsAttachmentSweeper(region, vpcID string, api vpcOtsAttachmentAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_ots_instance_attachment", F: func(actual string) error {
		if region == "" || vpcID == "" || actual != region || api.listVpc == nil || api.listInstance == nil || api.unbind == nil || api.pause == nil || api.attempts < 1 {
			return fmt.Errorf("OTS attachment requires verified VPC/region and bounded native API")
		}
		candidates, err := listVpcOtsAttachments(region, vpcID, "", api)
		if err != nil {
			return err
		}
		var failures []string
		for _, candidate := range candidates {
			rows, err := listVpcOtsAttachments(region, vpcID, candidate.InstanceName, api)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			present := false
			for _, row := range rows {
				if row.InstanceVpcName == candidate.InstanceVpcName {
					if row.VpcId != vpcID {
						err = fmt.Errorf("OTS attachment ownership changed; refusing unbind")
						break
					}
					present = true
				}
			}
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if !present {
				continue
			}
			if err = api.unbind(candidate.InstanceName, candidate.InstanceVpcName); err != nil {
				failures = append(failures, fmt.Sprintf("OTS attachment unbind: %v", err))
				continue
			}
			gone := false
			for attempt := 0; attempt < api.attempts; attempt++ {
				rows, err = listVpcOtsAttachments(region, vpcID, candidate.InstanceName, api)
				if err != nil {
					break
				}
				present = false
				for _, row := range rows {
					if row.InstanceVpcName == candidate.InstanceVpcName {
						if row.VpcId != vpcID {
							err = fmt.Errorf("OTS attachment changed VPC during wait")
							break
						}
						present = true
					}
				}
				if err != nil {
					break
				}
				if !present {
					gone = true
					break
				}
				if attempt+1 < api.attempts {
					api.pause(5 * time.Second)
				}
			}
			if err != nil {
				failures = append(failures, err.Error())
			} else if !gone {
				failures = append(failures, "OTS attachment still present after bounded unbind wait")
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("OTS VPC attachment cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}

// By-VPC responses provide VpcId in the response context, whereas by-instance
// responses provide InstanceName there. Validate context before filling only
// that contextual field, preserving every other attachment's independent VPC.
func listVpcOtsAttachments(region, vpcID, instance string, api vpcOtsAttachmentAPI) ([]sdkots.VpcInfo, error) {
	var all []sdkots.VpcInfo
	seen := map[string]bool{}
	const size = 100
	total := -1
	for page := 1; page <= 10000; page++ {
		var rows []sdkots.VpcInfo
		var count, returnedPage, returnedSize int
		if instance == "" {
			response, err := api.listVpc(vpcID, page, size)
			if err != nil {
				return nil, err
			}
			if response == nil || response.VpcId != vpcID {
				return nil, fmt.Errorf("OTS VPC inventory ownership unknown or changed")
			}
			rows = response.VpcInfos.VpcInfo
			count = response.TotalCount
			returnedPage = response.PageNum
			returnedSize = response.PageSize
		} else {
			response, err := api.listInstance(instance, page, size)
			if err != nil {
				return nil, err
			}
			if response == nil || response.InstanceName != instance {
				return nil, fmt.Errorf("OTS instance inventory ownership unknown or changed")
			}
			rows = response.VpcInfos.VpcInfo
			count = response.TotalCount
			returnedPage = response.PageNum
			returnedSize = response.PageSize
		}
		if returnedPage != page || returnedSize != size || count < 0 || (total >= 0 && count != total) {
			return nil, fmt.Errorf("OTS attachment pagination metadata inconsistent")
		}
		total = count
		for _, row := range rows {
			if instance == "" {
				if row.VpcId != "" && row.VpcId != vpcID {
					return nil, fmt.Errorf("OTS VPC inventory crossed scope")
				}
				row.VpcId = vpcID
			} else {
				if row.InstanceName != "" && row.InstanceName != instance {
					return nil, fmt.Errorf("OTS instance inventory crossed scope")
				}
				row.InstanceName = instance
			}
			if row.RegionNo != region || row.VpcId == "" || row.InstanceName == "" || row.InstanceVpcName == "" {
				return nil, fmt.Errorf("OTS attachment identity unknown; refusing unbind")
			}
			key := row.InstanceName + "/" + row.InstanceVpcName
			if seen[key] {
				return nil, fmt.Errorf("OTS attachment pagination repeated identity")
			}
			seen[key] = true
			all = append(all, row)
		}
		if len(all) > total {
			return nil, fmt.Errorf("OTS attachment count exceeded inventory")
		}
		if len(all) == total {
			return all, nil
		}
		if len(rows) != size {
			return nil, fmt.Errorf("OTS attachment inventory truncated")
		}
	}
	return nil, fmt.Errorf("OTS attachment pagination exceeded bound")
}

func vpcOtsTestAPI(t *testing.T, scenario string) (vpcOtsAttachmentAPI, *int, *int) {
	t.Helper()
	deletes, reads := 0, 0
	target := sdkots.VpcInfo{InstanceName: "instance", InstanceVpcName: "target-link", RegionNo: "cn-shanghai"}
	protected := sdkots.VpcInfo{InstanceVpcName: "protected-link", RegionNo: "cn-shanghai", VpcId: "vpc-protected"}
	api := vpcOtsAttachmentAPI{pause: func(time.Duration) {}, attempts: 2}
	api.listVpc = func(vpc string, page, size int) (*sdkots.ListVpcInfoByVpcResponse, error) {
		if vpc != "vpc-target" || page != 1 || size != 100 {
			t.Fatalf("unexpected inventory scope")
		}
		r := &sdkots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 1}
		r.VpcInfos.VpcInfo = []sdkots.VpcInfo{target}
		return r, nil
	}
	api.listInstance = func(instance string, page, size int) (*sdkots.ListVpcInfoByInstanceResponse, error) {
		reads++
		if instance != "instance" {
			t.Fatal("wrong shared instance")
		}
		if scenario == "read-error" || (scenario == "wait-read-error" && deletes > 0) {
			return nil, fmt.Errorf("read denied")
		}
		row := target
		row.InstanceName = ""
		row.VpcId = "vpc-target"
		switch scenario {
		case "wrong-vpc":
			row.VpcId = "vpc-protected"
		case "unknown-region":
			row.RegionNo = ""
		case "unknown-vpc":
			row.VpcId = ""
		case "unknown-binding":
			row.InstanceVpcName = ""
		}
		if scenario == "wait-scope-change" && deletes > 0 {
			row.VpcId = "vpc-protected"
		}
		r := &sdkots.ListVpcInfoByInstanceResponse{InstanceName: instance, PageNum: page, PageSize: size}
		r.VpcInfos.VpcInfo = []sdkots.VpcInfo{protected}
		if deletes == 0 || scenario == "timeout" || scenario == "wait-scope-change" || (scenario == "async" && reads == 2) {
			r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, row)
		}
		r.TotalCount = len(r.VpcInfos.VpcInfo)
		if scenario == "wrong-context" {
			r.InstanceName = "other"
		}
		return r, nil
	}
	api.unbind = func(instance, binding string) error {
		if instance != "instance" || binding != "target-link" {
			t.Fatal("shared/protected attachment touched")
		}
		deletes++
		if scenario == "write-error" {
			return fmt.Errorf("unbind denied")
		}
		return nil
	}
	return api, &deletes, &reads
}

func TestUnitVpcSweepOtsAttachmentSharedInstanceExactUnbind(t *testing.T) {
	api, deletes, reads := vpcOtsTestAPI(t, "async")
	pauses := 0
	api.pause = func(time.Duration) { pauses++ }
	if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err != nil {
		t.Fatal(err)
	}
	if *deletes != 1 || *reads != 3 || pauses != 1 {
		t.Fatalf("unbind=%d fresh reads=%d waits=%d", *deletes, *reads, pauses)
	}
}

func TestUnitVpcSweepOtsAttachmentUnknownOwnershipNeverUnbinds(t *testing.T) {
	for _, scenario := range []string{"wrong-vpc", "unknown-region", "unknown-vpc", "unknown-binding", "wrong-context", "read-error"} {
		t.Run(scenario, func(t *testing.T) {
			api, deletes, _ := vpcOtsTestAPI(t, scenario)
			if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err == nil || *deletes != 0 {
				t.Fatalf("unsafe scope err=%v unbind=%d", err, *deletes)
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentWaitAndWriteFailures(t *testing.T) {
	for _, scenario := range []string{"timeout", "write-error", "wait-read-error", "wait-scope-change"} {
		t.Run(scenario, func(t *testing.T) {
			api, deletes, _ := vpcOtsTestAPI(t, scenario)
			if err := vpcOtsAttachmentSweeper("cn-shanghai", "vpc-target", api).F("cn-shanghai"); err == nil || *deletes != 1 {
				t.Fatalf("release failure hidden err=%v unbind=%d", err, *deletes)
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentSDKContextMapping(t *testing.T) {
	var byVpc sdkots.ListVpcInfoByVpcResponse
	if err := json.Unmarshal([]byte(`{"VpcId":"vpc-target","TotalCount":1,"PageNum":1,"PageSize":100,"VpcInfos":{"VpcInfo":[{"InstanceName":"instance","InstanceVpcName":"link","RegionNo":"cn-shanghai"}]}}`), &byVpc); err != nil {
		t.Fatal(err)
	}
	var byInstance sdkots.ListVpcInfoByInstanceResponse
	if err := json.Unmarshal([]byte(`{"InstanceName":"instance","TotalCount":1,"PageNum":1,"PageSize":100,"VpcInfos":{"VpcInfo":[{"VpcId":"vpc-target","InstanceVpcName":"link","RegionNo":"cn-shanghai"}]}}`), &byInstance); err != nil {
		t.Fatal(err)
	}
	api := vpcOtsAttachmentAPI{listVpc: func(string, int, int) (*sdkots.ListVpcInfoByVpcResponse, error) { return &byVpc, nil }, listInstance: func(string, int, int) (*sdkots.ListVpcInfoByInstanceResponse, error) { return &byInstance, nil }}
	for _, instance := range []string{"", "instance"} {
		rows, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", instance, api)
		if err != nil || len(rows) != 1 || rows[0].InstanceName != "instance" || rows[0].VpcId != "vpc-target" {
			t.Fatalf("native contextual identity lost rows=%v err=%v", rows, err)
		}
	}
}

func TestUnitVpcSweepOtsAttachmentPaginationRejectsIncomplete(t *testing.T) {
	for _, scenario := range []string{"truncated", "wrong-page", "wrong-context", "cross-vpc", "repeated"} {
		t.Run(scenario, func(t *testing.T) {
			api := vpcOtsAttachmentAPI{listVpc: func(vpc string, page, size int) (*sdkots.ListVpcInfoByVpcResponse, error) {
				r := &sdkots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 2}
				row := sdkots.VpcInfo{InstanceName: "instance", InstanceVpcName: "link", RegionNo: "cn-shanghai"}
				r.VpcInfos.VpcInfo = []sdkots.VpcInfo{row}
				switch scenario {
				case "wrong-page":
					r.PageNum = 0
				case "wrong-context":
					r.VpcId = "vpc-protected"
				case "cross-vpc":
					r.VpcInfos.VpcInfo[0].VpcId = "vpc-protected"
				case "repeated":
					r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, row)
				}
				return r, nil
			}}
			if _, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", "", api); err == nil {
				t.Fatal("unverified pagination accepted")
			}
		})
	}
}

func TestUnitVpcSweepOtsAttachmentCompletePagination(t *testing.T) {
	calls := 0
	api := vpcOtsAttachmentAPI{listVpc: func(vpc string, page, size int) (*sdkots.ListVpcInfoByVpcResponse, error) {
		calls++
		r := &sdkots.ListVpcInfoByVpcResponse{VpcId: vpc, PageNum: page, PageSize: size, TotalCount: 101}
		count := 100
		if page == 2 {
			count = 1
		}
		for i := 0; i < count; i++ {
			r.VpcInfos.VpcInfo = append(r.VpcInfos.VpcInfo, sdkots.VpcInfo{InstanceName: "instance", InstanceVpcName: fmt.Sprintf("page%d-item%d", page, i), RegionNo: "cn-shanghai"})
		}
		return r, nil
	}}
	rows, err := listVpcOtsAttachments("cn-shanghai", "vpc-target", "", api)
	if err != nil || len(rows) != 101 || calls != 2 {
		t.Fatalf("inventory incomplete rows=%d calls=%d err=%v", len(rows), calls, err)
	}
}
