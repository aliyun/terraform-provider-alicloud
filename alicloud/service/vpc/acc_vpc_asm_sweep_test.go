package vpc_test

// VPC-owned ACK/ASM sweep cluster. This group of helpers and unit tests
// (vpcASMAPI type, nativeVpcASMAPI, nativeVpcASMDelete, trackVpcACKOwnership,
// asmExactCode, asmList, asmScope, asmClusterScope, vpcASMSweeper, asmUnitMesh,
// asmUnitAPI and the TestUnitVpcASM* cases) lived at the tail of the root
// alicloud package's resource_alicloud_service_mesh_service_mesh_test.go at
// baseline. They form the VPC sweeper's ACK/ASM ownership-and-cleanup adapter
// and are consumed by the VPC sweep tests in this package
// (sweepVpcsWithCallbacks -> trackVpcACKOwnership / nativeVpcASMAPI /
// vpcASMSweeper). The service_mesh resource tests do not reference them, so
// the cluster is relocated here to its natural home alongside the VPC sweep
// tests, keeping both product packages self-contained without duplication.
// Bodies are transcribed verbatim from the baseline source; the only change is
// the owning package/file.

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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

func asmUnitMesh(clusters interface{}) map[string]interface{} {
	return map[string]interface{}{"ServiceMeshInfo": map[string]interface{}{"ServiceMeshId": "mesh", "RegionId": "region"}, "Spec": map[string]interface{}{"Network": map[string]interface{}{"VpcId": "target"}}, "Clusters": clusters}
}
func asmUnitAPI(mesh map[string]interface{}) (vpcASMAPI, *int) {
	deletes := 0
	gone := false
	now := time.Unix(0, 0)
	api := vpcASMAPI{region: "region", now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }, maxPolls: 3, releaseTimeout: time.Minute}
	api.list = func() (map[string]interface{}, error) {
		rows := []interface{}{}
		if !gone {
			rows = append(rows, mesh)
		}
		return map[string]interface{}{"ServiceMeshes": rows}, nil
	}
	api.detail = func(string) (map[string]interface{}, error) {
		if gone {
			return nil, &tea.SDKError{Code: tea.String("ServiceMesh.NotFound")}
		}
		return map[string]interface{}{"ServiceMesh": mesh}, nil
	}
	api.clusterDetail = func(id string) (map[string]interface{}, error) {
		return map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}, nil
	}
	api.remove = func(string) (map[string]interface{}, error) {
		deletes++
		gone = true
		return map[string]interface{}{}, nil
	}
	return api, &deletes
}
func TestUnitVpcASMScopedOwnerAndExactAsyncConfirmation(t *testing.T) {
	for _, knownGone := range []bool{false, true} {
		t.Run(fmt.Sprint(knownGone), func(t *testing.T) {
			api, n := asmUnitAPI(asmUnitMesh([]interface{}{"cluster"}))
			if knownGone {
				api.clusterWasScoped = func(id string) bool { return id == "cluster" }
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
			}
			if e := vpcASMSweeper("region", "target", api).F("region"); e != nil || *n != 1 {
				t.Fatalf("delete=%d err=%v", *n, e)
			}
		})
	}
}
func TestUnitVpcASMSharedUnknownAndOpaqueErrorsNeverDelete(t *testing.T) {
	t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", "true")
	for _, mode := range []string{"shared", "region", "identity", "missingClusters", "malformedClusters", "permission", "opaqueNotFound", "regionalNotFound", "freshVpc", "freshError"} {
		t.Run(mode, func(t *testing.T) {
			mesh := asmUnitMesh([]interface{}{"cluster"})
			api, n := asmUnitAPI(mesh)
			switch mode {
			case "shared", "region", "identity":
				api.clusterDetail = func(id string) (map[string]interface{}, error) {
					d := map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}
					if mode == "shared" {
						d["vpc_id"] = "protected"
					}
					if mode == "region" {
						d["region_id"] = "other"
					}
					if mode == "identity" {
						d["cluster_id"] = "other"
					}
					return d, nil
				}
			case "missingClusters":
				delete(mesh, "Clusters")
			case "malformedClusters":
				mesh["Clusters"] = []interface{}{1}
			case "regionalNotFound":
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
			case "permission", "opaqueNotFound":
				api.clusterDetail = func(string) (map[string]interface{}, error) {
					return nil, &tea.SDKError{Code: tea.String("Forbidden"), Message: tea.String("cluster not found or inaccessible")}
				}
			case "freshVpc":
				api.detail = func(string) (map[string]interface{}, error) {
					d := asmUnitMesh([]interface{}{})
					d["Spec"].(map[string]interface{})["Network"].(map[string]interface{})["VpcId"] = "protected"
					return map[string]interface{}{"ServiceMesh": d}, nil
				}
			case "freshError":
				api.detail = func(string) (map[string]interface{}, error) { return nil, fmt.Errorf("MiranaErr.Common") }
			}
			if e := vpcASMSweeper("region", "target", api).F("region"); e == nil || *n != 0 {
				t.Fatalf("delete=%d err=%v", *n, e)
			}
		})
	}
}
func TestUnitVpcASMCompleteListAndBoundedRelease(t *testing.T) {
	for _, mode := range []string{"missingList", "pagination", "acceptedStillPresent", "opaqueGone", "listStillPresent", "emptyClusters"} {
		t.Run(mode, func(t *testing.T) {
			api, n := asmUnitAPI(asmUnitMesh([]interface{}{}))
			switch mode {
			case "missingList":
				api.list = func() (map[string]interface{}, error) { return map[string]interface{}{}, nil }
			case "pagination":
				orig := api.list
				api.list = func() (map[string]interface{}, error) { d, e := orig(); d["NextToken"] = "more"; return d, e }
			case "acceptedStillPresent":
				api.remove = func(string) (map[string]interface{}, error) { *n++; return map[string]interface{}{}, nil }
			case "opaqueGone":
				api.remove = func(string) (map[string]interface{}, error) {
					*n++
					api.detail = func(string) (map[string]interface{}, error) { return nil, fmt.Errorf("not found or inaccessible") }
					return map[string]interface{}{}, nil
				}
			case "listStillPresent":
				api.remove = func(string) (map[string]interface{}, error) {
					*n++
					api.detail = func(string) (map[string]interface{}, error) {
						return nil, &tea.SDKError{Code: tea.String("ServiceMesh.NotFound")}
					}
					return map[string]interface{}{}, nil
				}
			}
			e := vpcASMSweeper("region", "target", api).F("region")
			if mode == "emptyClusters" {
				if e != nil || *n != 1 {
					t.Fatalf("delete=%d err=%v", *n, e)
				}
			} else if e == nil {
				t.Fatal("invalid/incomplete release accepted")
			}
		})
	}
}

func TestUnitVpcASMNativeForceOnlyForFullCleanup(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprint(full), func(t *testing.T) {
			t.Setenv("ALICLOUD_SWEEP_ALL_RESOURCES", fmt.Sprint(full))
			calls := 0
			expected := fmt.Errorf("native ASM delete rejected")
			_, err := nativeVpcASMDelete(func(product, version, action string, query, request map[string]interface{}, auto bool) (map[string]interface{}, error) {
				calls++
				force, present := request["Force"]
				fields := 1
				if full {
					fields++
				}
				if product != "servicemesh" || version != "2020-01-11" || action != "DeleteServiceMesh" || query != nil || request["ServiceMeshId"] != "mesh-target" || present != full || (full && force != true) || len(request) != fields {
					t.Fatalf("wrong native full/ordinary ASM delete contract: %s %s %s %v", product, version, action, request)
				}
				return nil, expected
			}, "mesh-target")
			if err != expected || calls != 1 {
				t.Fatalf("native deletion result was lost: err=%v calls=%d", err, calls)
			}
		})
	}
}

func TestUnitVpcASMACKOwnershipProofIsLocalAndExact(t *testing.T) {
	for _, scenario := range []string{"owned", "other-vpc", "other-region", "other-id", "read-error", "missing-detail"} {
		t.Run(scenario, func(t *testing.T) {
			api := vpcACKAPI{detail: func(id string) (map[string]interface{}, error) {
				row := map[string]interface{}{"cluster_id": id, "region_id": "region", "vpc_id": "target"}
				switch scenario {
				case "other-vpc":
					row["vpc_id"] = "protected"
				case "other-region":
					row["region_id"] = "reserved-region"
				case "other-id":
					row["cluster_id"] = "different-cluster"
				case "read-error":
					return nil, &tea.SDKError{Code: tea.String("ErrorClusterNotFound")}
				}
				return row, nil
			}}
			if scenario == "missing-detail" {
				api.detail = nil
			}
			tracked, proof := trackVpcACKOwnership("region", "target", api)
			if proof("cluster") {
				t.Fatal("ownership existed before a successful exact read")
			}
			if tracked.detail != nil {
				_, _ = tracked.detail("cluster")
			}
			if proof("cluster") != (scenario == "owned") || proof("another-cluster") {
				t.Fatal("unknown, protected or cross-region identity counted as owned")
			}
			_, freshProof := trackVpcACKOwnership("region", "target", api)
			if freshProof("cluster") {
				t.Fatal("ownership leaked into another invocation")
			}
		})
	}
}
