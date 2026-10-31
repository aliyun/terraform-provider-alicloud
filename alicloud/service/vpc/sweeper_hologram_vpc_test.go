package vpc_test

// Hologram sweep cluster for the migrated VPC sweep graph. The baseline root
// package alicloud hosted this cluster (type vpcHologramAPI, nativeVpcHologramAPI,
// vpcHologramSweeper and the hologram* helpers) in sweeper_hologram_vpc_test.go,
// reachable from the root resource_alicloud_vpc_test.go sweep graph because both
// lived in package alicloud. The migration moves the VPC sweep graph
// (sweepVpcDependentResources and its 117 TestUnit sweep tests) into the external
// package vpc_test, which cannot see root-package helpers. The cluster is
// re-hosted here verbatim so the migrated sweep graph's hologram branch
// (vpcHologramSweeper(client.RegionId, id, nativeVpcHologramAPI(client))) and
// its four TestUnitVpcHologram* tests resolve. ComplexError resolves via
// aliases_test.go. The root copy is removed: no remaining root file references
// these symbols (root resource_alicloud_vpc_test.go was migrated).

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type vpcHologramAPI struct {
	region         string
	list           func() (map[string]interface{}, error)
	get            func(string) (map[string]interface{}, error)
	remove         func(string) (map[string]interface{}, error)
	now            func() time.Time
	sleep          func(time.Duration)
	releaseTimeout time.Duration
	maxPolls       int
}

func nativeVpcHologramAPI(client *connectivity.AliyunClient) vpcHologramAPI {
	return vpcHologramAPI{region: client.RegionId,
		// Official ListInstances has a single optional opaque body and a complete
		// InstanceList response: no page/size/next-token fields exist in its schema.
		list: func() (map[string]interface{}, error) {
			return client.RoaPost("Hologram", "2022-06-01", "/api/v1/instances", map[string]*string{}, nil, map[string]interface{}{}, true)
		},
		get: func(id string) (map[string]interface{}, error) {
			return client.RoaGet("Hologram", "2022-06-01", "/api/v1/instances/"+id, map[string]*string{}, nil, map[string]interface{}{"instanceId": id})
		},
		remove: func(id string) (map[string]interface{}, error) {
			return client.RoaPost("Hologram", "2022-06-01", "/api/v1/instances/"+id+"/delete", map[string]*string{}, nil, map[string]interface{}{"instanceId": id, "RegionId": client.RegionId}, true)
		},
		now: time.Now, sleep: time.Sleep, releaseTimeout: 30 * time.Minute, maxPolls: 180,
	}
}

func vpcHologramSweeper(region, vpcID string, api vpcHologramAPI) *resource.Sweeper {
	return &resource.Sweeper{Name: "alicloud_hologram_instance", F: func(actual string) error {
		if region == "" || vpcID == "" || actual != region || api.region != region {
			return fmt.Errorf("Hologram cleanup requires exact VPC and region")
		}
		if api.list == nil || api.get == nil || api.remove == nil || api.now == nil || api.sleep == nil || api.maxPolls < 1 || api.releaseTimeout <= 0 {
			return fmt.Errorf("Hologram cleanup requires bounded injected API")
		}
		rows, err := hologramCompleteInstanceList(api, region)
		if err != nil {
			return err
		}
		var failures []string
		for _, row := range rows {
			selected, err := hologramExclusiveVpc(row, vpcID)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			if !selected {
				continue
			}
			id := row["InstanceId"].(string)
			response, err := api.get(id)
			if err != nil {
				failures = append(failures, fmt.Sprintf("Hologram detail %s: %v", id, err))
				continue
			}
			detail, err := hologramVerifiedDetail(response, region, id)
			if err != nil {
				failures = append(failures, err.Error())
				continue
			}
			selected, err = hologramExclusiveVpc(detail, vpcID)
			if err != nil || !selected {
				failures = append(failures, fmt.Sprintf("Hologram %s complete endpoint ownership changed or shared: %v", id, err))
				continue
			}
			if detail["InstanceChargeType"] != "PostPaid" {
				failures = append(failures, fmt.Sprintf("Hologram %s billing requires manual release; not PostPaid", id))
				continue
			}
			response, err = api.remove(id)
			if err != nil && !hologramExactNotFound(err) {
				failures = append(failures, fmt.Sprintf("Hologram delete %s: %v", id, err))
				continue
			}
			if err == nil {
				if err = hologramSuccessfulResponse(response); err != nil {
					failures = append(failures, err.Error())
					continue
				}
			}
			if err = hologramWaitReleased(api, region, vpcID, id); err != nil {
				failures = append(failures, err.Error())
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("Hologram owner cleanup incomplete: %s", strings.Join(failures, "; "))
		}
		return nil
	}}
}

func hologramSuccessfulResponse(response map[string]interface{}) error {
	success, ok := response["Success"].(bool)
	if !ok || !success {
		return fmt.Errorf("Hologram response did not explicitly succeed")
	}
	if status, exists := response["HttpStatusCode"]; exists && fmt.Sprint(status) != "200" {
		return fmt.Errorf("Hologram unexpected response status")
	}
	return nil
}

func hologramCompleteInstanceList(api vpcHologramAPI, region string) ([]map[string]interface{}, error) {
	response, err := api.list()
	if err != nil {
		return nil, err
	}
	if err = hologramSuccessfulResponse(response); err != nil {
		return nil, err
	}
	// A changed/partial response must not masquerade as the documented full list.
	for _, key := range []string{"NextToken", "nextToken", "HasMore", "hasMore", "PageNumber", "pageNumber", "PageSize", "pageSize"} {
		if value, exists := response[key]; exists && value != nil {
			return nil, fmt.Errorf("Hologram unexpected pagination field %s", key)
		}
	}
	raw, ok := response["InstanceList"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Hologram missing complete InstanceList")
	}
	rows := make([]map[string]interface{}, 0, len(raw))
	seen := map[string]bool{}
	for _, value := range raw {
		row, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Hologram malformed instance list")
		}
		id, ok := row["InstanceId"].(string)
		if !ok || !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`).MatchString(id) || seen[id] {
			return nil, fmt.Errorf("Hologram missing/duplicate/invalid instance identity")
		}
		if row["RegionId"] != region {
			return nil, fmt.Errorf("Hologram list region missing or mismatched")
		}
		seen[id] = true
		rows = append(rows, row)
	}
	return rows, nil
}

func hologramVerifiedDetail(response map[string]interface{}, region, id string) (map[string]interface{}, error) {
	if err := hologramSuccessfulResponse(response); err != nil {
		return nil, err
	}
	detail, ok := response["Instance"].(map[string]interface{})
	if !ok || detail["InstanceId"] != id || detail["RegionId"] != region {
		return nil, fmt.Errorf("Hologram detail identity/region mismatched")
	}
	return detail, nil
}

func hologramExclusiveVpc(instance map[string]interface{}, vpcID string) (bool, error) {
	endpoints, ok := instance["Endpoints"].([]interface{})
	if !ok || len(endpoints) == 0 {
		return false, fmt.Errorf("Hologram complete Endpoints missing")
	}
	selected, foreign := false, false
	for _, raw := range endpoints {
		endpoint, ok := raw.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("Hologram malformed endpoint")
		}
		kind, ok := endpoint["Type"].(string)
		if !ok {
			return false, fmt.Errorf("Hologram endpoint Type missing")
		}
		switch kind {
		case "Intranet", "Internet":
			if v, exists := endpoint["VpcId"]; exists && v != nil && v != "" {
				return false, fmt.Errorf("Hologram non-VPC endpoint has unexpected VPC binding")
			}
		case "VPCSingleTunnel":
			id, ok := endpoint["VpcId"].(string)
			if !ok || id == "" {
				return false, fmt.Errorf("Hologram VPC endpoint missing VpcId")
			}
			if id == vpcID {
				selected = true
			} else {
				foreign = true
			}
		default:
			return false, fmt.Errorf("Hologram unknown endpoint Type; refusing release")
		}
	}
	if selected && foreign {
		return false, fmt.Errorf("Hologram shared VPC endpoints; instance retained")
	}
	return selected, nil
}

func hologramExactNotFound(err error) bool {
	for {
		if wrapped, ok := err.(*ComplexError); ok {
			err = wrapped.Cause
			continue
		}
		break
	}
	switch e := err.(type) {
	case *tea.SDKError:
		return tea.StringValue(e.Code) == "resource not exists failed" || tea.StringValue(e.Message) == "resource not exists failed"
	case *sdkerrors.ServerError:
		return e.ErrorCode() == "resource not exists failed" || e.Message() == "resource not exists failed"
	}
	return false
}

func hologramWaitReleased(api vpcHologramAPI, region, vpcID, id string) error {
	deadline := api.now().Add(api.releaseTimeout)
	for attempt := 0; attempt < api.maxPolls && api.now().Before(deadline); attempt++ {
		response, err := api.get(id)
		if err != nil {
			if !hologramExactNotFound(err) {
				return fmt.Errorf("Hologram release verification: %w", err)
			}
			rows, listErr := hologramCompleteInstanceList(api, region)
			if listErr != nil {
				return listErr
			}
			present := false
			for _, row := range rows {
				if row["InstanceId"] == id {
					present = true
					break
				}
			}
			if !present {
				return nil
			}
		} else {
			detail, verifyErr := hologramVerifiedDetail(response, region, id)
			if verifyErr != nil {
				return verifyErr
			}
			selected, scopeErr := hologramExclusiveVpc(detail, vpcID)
			if scopeErr != nil || !selected {
				return fmt.Errorf("Hologram ownership changed during release verification")
			}
		}
		if attempt+1 < api.maxPolls && api.now().Before(deadline) {
			api.sleep(10 * time.Second)
		}
	}
	return fmt.Errorf("Hologram exact instance still present after bounded release verification")
}

func hologramUnitInstance(vpc, charge string) map[string]interface{} {
	return map[string]interface{}{"InstanceId": "hg-instance", "RegionId": "cn-hangzhou", "InstanceChargeType": charge, "InstanceStatus": "Running", "Endpoints": []interface{}{map[string]interface{}{"Type": "Intranet"}, map[string]interface{}{"Type": "VPCSingleTunnel", "VpcId": vpc}}}
}
func hologramUnitList(row map[string]interface{}) map[string]interface{} {
	rows := []interface{}{}
	if row != nil {
		rows = append(rows, row)
	}
	return map[string]interface{}{"Success": true, "HttpStatusCode": 200, "InstanceList": rows}
}
func hologramUnitDetail(row map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"Success": true, "HttpStatusCode": 200, "Instance": row}
}
func hologramUnitAPI() vpcHologramAPI {
	now := time.Unix(10, 0)
	return vpcHologramAPI{region: "cn-hangzhou", now: func() time.Time { return now }, sleep: func(d time.Duration) { now = now.Add(d) }, releaseTimeout: time.Minute, maxPolls: 3,
		list: func() (map[string]interface{}, error) {
			return hologramUnitList(hologramUnitInstance("target", "PostPaid")), nil
		},
		get: func(string) (map[string]interface{}, error) {
			return hologramUnitDetail(hologramUnitInstance("target", "PostPaid")), nil
		},
		remove: func(string) (map[string]interface{}, error) { return map[string]interface{}{"Success": true}, nil }}
}

func TestUnitVpcHologramExclusivePostpaidOwnerAndAsyncConfirmation(t *testing.T) {
	api := hologramUnitAPI()
	writes, polls, lists := 0, 0, 0
	api.list = func() (map[string]interface{}, error) {
		lists++
		if polls >= 2 {
			return hologramUnitList(nil), nil
		}
		return hologramUnitList(hologramUnitInstance("target", "PostPaid")), nil
	}
	api.get = func(id string) (map[string]interface{}, error) {
		if id != "hg-instance" {
			t.Fatal("wrong exact ID")
		}
		if writes > 0 {
			polls++
			if polls >= 2 {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "resource not exists failed", "message": "resource not exists failed"})
			}
		}
		return hologramUnitDetail(hologramUnitInstance("target", "PostPaid")), nil
	}
	api.remove = func(id string) (map[string]interface{}, error) {
		writes++
		if id != "hg-instance" {
			t.Fatal("wrong delete ID")
		}
		return map[string]interface{}{"Success": true}, nil
	}
	if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err != nil {
		t.Fatal(err)
	}
	if writes != 1 || polls != 2 || lists != 2 {
		t.Fatalf("writes=%d polls=%d lists=%d", writes, polls, lists)
	}
}

func TestUnitVpcHologramSharedVpcAndBillingNeverRelease(t *testing.T) {
	for _, scenario := range []string{"shared", "protected-only", "PrePaid", "unknown-billing", "unknown-type", "unexpected-internet-vpc", "missing-vpc", "missing-endpoints", "changed-detail"} {
		api := hologramUnitAPI()
		writes := 0
		row := hologramUnitInstance("target", "PostPaid")
		switch scenario {
		case "shared":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "VPCSingleTunnel", "VpcId": "protected"})
		case "protected-only":
			row = hologramUnitInstance("protected", "PostPaid")
		case "PrePaid":
			row["InstanceChargeType"] = "PrePaid"
		case "unknown-billing":
			delete(row, "InstanceChargeType")
		case "unknown-type":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "Unknown"})
		case "unexpected-internet-vpc":
			row["Endpoints"] = append(row["Endpoints"].([]interface{}), map[string]interface{}{"Type": "Internet", "VpcId": "protected"})
		case "missing-vpc":
			row["Endpoints"] = []interface{}{map[string]interface{}{"Type": "VPCSingleTunnel"}}
		case "missing-endpoints":
			delete(row, "Endpoints")
		}
		api.list = func() (map[string]interface{}, error) { return hologramUnitList(row), nil }
		api.get = func(string) (map[string]interface{}, error) {
			if scenario == "changed-detail" {
				return hologramUnitDetail(hologramUnitInstance("protected", "PostPaid")), nil
			}
			return hologramUnitDetail(row), nil
		}
		api.remove = func(string) (map[string]interface{}, error) {
			writes++
			return map[string]interface{}{"Success": true}, nil
		}
		err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou")
		if writes != 0 {
			t.Fatal("unsafe release", scenario)
		}
		if scenario != "protected-only" && err == nil {
			t.Fatal("unsafe boundary accepted", scenario)
		}
	}
}

func TestUnitVpcHologramMalformedRegionalListFailsBeforeDeletion(t *testing.T) {
	for _, scenario := range []string{"wrong-region", "missing-region", "duplicate", "missing-list", "success-false", "partial-token", "bad-id"} {
		api := hologramUnitAPI()
		writes := 0
		row := hologramUnitInstance("target", "PostPaid")
		response := hologramUnitList(row)
		switch scenario {
		case "wrong-region":
			row["RegionId"] = "cn-beijing"
		case "missing-region":
			delete(row, "RegionId")
		case "duplicate":
			response["InstanceList"] = []interface{}{row, row}
		case "missing-list":
			delete(response, "InstanceList")
		case "success-false":
			response["Success"] = false
		case "partial-token":
			response["NextToken"] = "not-complete"
		case "bad-id":
			row["InstanceId"] = "../wrong-path"
		}
		api.list = func() (map[string]interface{}, error) { return response, nil }
		api.remove = func(string) (map[string]interface{}, error) { writes++; return nil, nil }
		if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err == nil || writes != 0 {
			t.Fatalf("scenario=%s writes=%d err=%v", scenario, writes, err)
		}
	}
}

func TestUnitVpcHologramOpaqueErrorsAndReleaseWaitAreBounded(t *testing.T) {
	for _, scenario := range []string{"list-denied", "get-denied", "delete-denied", "delete-unknown-success", "poll-denied", "poll-unknown-notfound", "timeout", "notfound-but-listed", "ownership-change"} {
		api := hologramUnitAPI()
		writes, gets := 0, 0
		originalGet := api.get
		api.get = func(id string) (map[string]interface{}, error) {
			gets++
			if scenario == "get-denied" || writes > 0 && scenario == "poll-denied" {
				return nil, fmt.Errorf("permission denied")
			}
			if writes > 0 && scenario == "poll-unknown-notfound" {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "Forbidden", "message": "not found or inaccessible"})
			}
			if writes > 0 && scenario == "notfound-but-listed" {
				return nil, tea.NewSDKError(map[string]interface{}{"code": "resource not exists failed"})
			}
			if writes > 0 && scenario == "ownership-change" {
				return hologramUnitDetail(hologramUnitInstance("protected", "PostPaid")), nil
			}
			return originalGet(id)
		}
		if scenario == "list-denied" {
			api.list = func() (map[string]interface{}, error) { return nil, fmt.Errorf("list denied") }
		}
		api.remove = func(string) (map[string]interface{}, error) {
			writes++
			if scenario == "delete-denied" {
				return nil, fmt.Errorf("delete denied")
			}
			if scenario == "delete-unknown-success" {
				return map[string]interface{}{}, nil
			}
			return map[string]interface{}{"Success": true}, nil
		}
		if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-hangzhou"); err == nil {
			t.Fatal("false success", scenario)
		}
		if writes > 1 || gets > 4 {
			t.Fatalf("unbounded scenario=%s writes=%d gets=%d", scenario, writes, gets)
		}
	}
	api := hologramUnitAPI()
	if err := vpcHologramSweeper("cn-hangzhou", "target", api).F("cn-beijing"); err == nil {
		t.Fatal("region mismatch accepted")
	}
}
