// Code in this file is a root-package copy of the VPC ASM sweep helpers that
// originally lived in resource_alicloud_service_mesh_service_mesh_test.go. The
// servicemesh tests were relocated to alicloud/service/servicemesh as part of
// the whole-product sub-package migration; the root package VPC sweeper chain
// in resource_alicloud_vpc_test.go still references these helpers, so a copy is
// retained here to keep the root test package compiling independently. The
// canonical copy now lives in alicloud/service/servicemesh.

package alicloud

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type vpcASMAPI struct {
	region           string
	list             func() (map[string]interface{}, error)
	detail           func(string) (map[string]interface{}, error)
	clusterDetail    func(string) (map[string]interface{}, error)
	clusterWasScoped func(string) bool
	remove           func(string) (map[string]interface{}, error)
	now              func() time.Time
	sleep            func(time.Duration)
	maxPolls         int
	releaseTimeout   time.Duration
}

func nativeVpcASMAPI(client *connectivity.AliyunClient, region string) vpcASMAPI {
	ack := nativeVpcACKAPI(client, region)
	return vpcASMAPI{region: region,
		list: func() (map[string]interface{}, error) {
			return client.RpcGet("servicemesh", "2020-01-11", "DescribeServiceMeshes", map[string]interface{}{}, nil)
		},
		detail: func(id string) (map[string]interface{}, error) {
			return client.RpcPost("servicemesh", "2020-01-11", "DescribeServiceMeshDetail", nil, map[string]interface{}{"ServiceMeshId": id}, true)
		},
		clusterDetail: ack.detail,
		// Scope and associated-cluster checks run before this delete adapter.
		remove: func(id string) (map[string]interface{}, error) {
			return nativeVpcASMDelete(client.RpcPost, id)
		},
		now: time.Now, sleep: time.Sleep, maxPolls: 180, releaseTimeout: 30 * time.Minute}
}

// Force is the native ASM parameter, used only by explicitly authorized full
// cleanup after VPC/cluster ownership verification. Ordinary deletion is unchanged.
func nativeVpcASMDelete(rpc vpcSweepRPC, id string) (map[string]interface{}, error) {
	request := map[string]interface{}{"ServiceMeshId": id}
	if sweepAll() {
		request["Force"] = true
	}
	return rpc("servicemesh", "2020-01-11", "DeleteServiceMesh", nil, request, true)
}

// A regional not-found alone cannot prove that a cluster in another region
// is gone. Retain only positive, exact ownership observed by this invocation's
// ACK sweeper so its completed deletion can safely precede ASM cleanup.
func trackVpcACKOwnership(region, vpcID string, api vpcACKAPI) (vpcACKAPI, func(string) bool) {
	var mu sync.RWMutex
	verified := map[string]bool{}
	proof := func(id string) bool {
		mu.RLock()
		defer mu.RUnlock()
		return verified[id]
	}
	original := api.detail
	if original == nil {
		return api, proof
	}
	api.detail = func(id string) (map[string]interface{}, error) {
		row, err := original(id)
		if err == nil && region != "" && vpcID != "" && id != "" && row["cluster_id"] == id && row["region_id"] == region && row["vpc_id"] == vpcID {
			mu.Lock()
			verified[id] = true
			mu.Unlock()
		}
		return row, err
	}
	return api, proof
}

func asmExactCode(err error, code string) bool {
	for {
		if e, ok := err.(*ComplexError); ok {
			err = e.Cause
			continue
		}
		break
	}
	switch e := err.(type) {
	case *tea.SDKError:
		return tea.StringValue(e.Code) == code
	case *sdkerrors.ServerError:
		return e.ErrorCode() == code
	}
	return false
}
func asmList(api vpcASMAPI) ([]map[string]interface{}, error) {
	r, e := api.list()
	if e != nil {
		return nil, e
	}
	for _, k := range []string{"NextToken", "nextToken", "HasMore", "PageNumber", "PageSize"} {
		if v, ok := r[k]; ok && v != nil {
			return nil, fmt.Errorf("ASM unexpected pagination %s", k)
		}
	}
	raw, ok := r["ServiceMeshes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("ASM missing complete ServiceMeshes")
	}
	rows := []map[string]interface{}{}
	seen := map[string]bool{}
	for _, v := range raw {
		row, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ASM malformed list")
		}
		info, ok := row["ServiceMeshInfo"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ASM missing identity")
		}
		id, ok := info["ServiceMeshId"].(string)
		if !ok || id == "" || seen[id] {
			return nil, fmt.Errorf("ASM invalid/duplicate identity")
		}
		seen[id] = true
		rows = append(rows, row)
	}
	return rows, nil
}
func asmScope(row map[string]interface{}, region, vpc, id string) (bool, error) {
	info, ok := row["ServiceMeshInfo"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing info")
	}
	if info["ServiceMeshId"] != id {
		return false, fmt.Errorf("ASM identity changed")
	}
	spec, ok := row["Spec"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing spec")
	}
	network, ok := spec["Network"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ASM missing network")
	}
	owner, ok := network["VpcId"].(string)
	if !ok || owner == "" {
		return false, fmt.Errorf("ASM unknown VPC")
	}
	if owner != vpc {
		return false, nil
	}
	if info["RegionId"] != region {
		return false, fmt.Errorf("ASM region mismatch")
	}
	return true, nil
}
func asmClusterScope(row map[string]interface{}, api vpcASMAPI, region, vpc string) error {
	raw, ok := row["Clusters"].([]interface{})
	if !ok {
		return fmt.Errorf("ASM complete associated cluster list missing")
	}
	seen := map[string]bool{}
	for _, v := range raw {
		id, ok := v.(string)
		if !ok || id == "" || seen[id] {
			return fmt.Errorf("ASM malformed associated cluster")
		}
		seen[id] = true
		d, e := api.clusterDetail(id)
		if e != nil {
			if asmExactCode(e, "ErrorClusterNotFound") && api.clusterWasScoped != nil && api.clusterWasScoped(id) {
				continue
			}
			return fmt.Errorf("ASM associated ACK detail: %w", e)
		}
		if d["cluster_id"] != id || d["region_id"] != region || d["vpc_id"] != vpc {
			return fmt.Errorf("ASM shared/unknown associated ACK ownership; retained")
		}
	}
	return nil
}
func vpcASMSweeper(region, vpcID string, api vpcASMAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_service_mesh_service_mesh", F: func(actual string) error {
		if actual != region || region == "" || vpcID == "" || api.region != region || api.list == nil || api.detail == nil || api.clusterDetail == nil || api.remove == nil || api.now == nil || api.sleep == nil || api.maxPolls < 1 || api.releaseTimeout <= 0 {
			return fmt.Errorf("ASM requires exact scope and bounded API")
		}
		rows, e := asmList(api)
		if e != nil {
			return e
		}
		var failures []string
		for _, row := range rows {
			id := row["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"].(string)
			selected, e := asmScope(row, region, vpcID, id)
			if e != nil {
				failures = append(failures, e.Error())
				continue
			}
			if !selected {
				continue
			}
			r, e := api.detail(id)
			if e != nil {
				failures = append(failures, e.Error())
				continue
			}
			detail, ok := r["ServiceMesh"].(map[string]interface{})
			if !ok {
				failures = append(failures, "ASM missing fresh detail")
				continue
			}
			selected, e = asmScope(detail, region, vpcID, id)
			if e != nil || !selected {
				failures = append(failures, "ASM fresh ownership changed")
				continue
			}
			if e = asmClusterScope(detail, api, region, vpcID); e != nil {
				failures = append(failures, e.Error())
				continue
			}
			if _, e = api.remove(id); e != nil {
				failures = append(failures, e.Error())
				continue
			}
			deadline := api.now().Add(api.releaseTimeout)
			gone := false
			for poll := 0; poll < api.maxPolls && api.now().Before(deadline); poll++ {
				_, e = api.detail(id)
				if e != nil {
					if !asmExactCode(e, "ServiceMesh.NotFound") {
						break
					}
					var current []map[string]interface{}
					current, e = asmList(api)
					if e != nil {
						break
					}
					found := false
					for _, x := range current {
						if x["ServiceMeshInfo"].(map[string]interface{})["ServiceMeshId"] == id {
							found = true
						}
					}
					if !found {
						gone = true
						break
					}
				}
				api.sleep(10 * time.Second)
			}
			if !gone {
				failures = append(failures, fmt.Sprintf("ASM release not confirmed: %v", e))
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("ASM owner cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}
