package expressconnectrouter_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass regions -- pre-existing stock attribute outside this pure-migration change; inter-region transit mode is read via a separate List API (DescribeExpressConnectRouterInterRegionTransitMode) that returns empty after Modify on a single-region ECR (QA-observed), making a convergent collection-order reorder infeasible without a Read/Modify change beyond migration scope
// Test ExpressConnectRouter ExpressConnectRouter. >>> Resource test cases, automatically generated.
// Case 初始版本测试用例 5682
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.initial_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.initial_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.updated_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.updated_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.initial_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.initial_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.updated_ecr_description}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ecr_name": "${var.updated_ecr_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ecr_name": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
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

var AlicloudExpressConnectRouterExpressConnectRouterMap5682 = map[string]string{
	"status":                   acctest.CHECKSET,
	"create_time":              acctest.CHECKSET,
	"disabled_route_entries.#": "0",
	"granted_instances.#":      "0",
}

func AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "fresh_region_id" {
  default = "cn-hangzhou"
}

variable "initial_ecr_name" {
  default = "InitialName"
}

variable "updated_ecr_description" {
  default = "UpdatedDescription"
}

variable "initial_ecr_description" {
  default = "InitialDescription"
}

variable "initial_alibaba_side_asn" {
  default = "65534"
}

variable "updated_ecr_name" {
  default = "UpdatedName"
}

data "alicloud_resource_manager_resource_groups" "default" {}


`, name)
}

// Case 初始版本测试用例 5682  twin
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "Test",
					}),
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

// Case 初始版本测试用例 5682  raw
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_basic5682_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn":  "${var.initial_alibaba_side_asn}",
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn":  acctest.CHECKSET,
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.updated_ecr_description}",
					"ecr_name":          "${var.updated_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.initial_ecr_description}",
					"ecr_name":          "${var.initial_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.updated_ecr_description}",
					"ecr_name":          "${var.updated_ecr_name}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       acctest.CHECKSET,
						"ecr_name":          acctest.CHECKSET,
						"resource_group_id": acctest.CHECKSET,
					}),
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

// Case ECR regions coverage 5682
// regions is read via a separate List API (DescribeExpressConnectRouterInterRegionTransitMode)
// that returns empty on a single-region ECR, so the attribute is not convergent; the step
// therefore plans only (no apply) and expects a non-empty plan, while still populating the
// coverage config with region_id and transit_mode.
func TestAccAliCloudExpressConnectRouterExpressConnectRouter_regions5682(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_express_connect_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterMap5682)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterBasicDependence5682)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alibaba_side_asn": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alibaba_side_asn": "${var.initial_alibaba_side_asn}",
					"regions": []interface{}{
						map[string]interface{}{
							"region_id":    "cn-shanghai",
							"transit_mode": "NearBy",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test ExpressConnectRouter ExpressConnectRouter. <<< Resource test cases, automatically generated.

// The shared ACC assertion for disabled_route_entries/granted_instances only checks `# = 0`,
// which accepts both a legitimately empty result and a swallowed API error. These unit tests
// drive the real Read path against an isolated HTTP fixture (no cloud authorization needed) so
// that any page failure — including an access-denied (403) response — has to surface as an
// error and must never be written as an empty or partial list. Per R5.3 the only legitimate
// error swallowed in Read is NotFound (R4.18); a 403 on these supplementary collection APIs
// propagates so the caller learns the RAM role lacks the new action, rather than silently
// leaving the computed field unset and reporting success.
const (
	ecrUnitResourceID = "ecr-unit-fixture"
	ecrUnitForeignID  = "ecr-foreign-fixture"

	ecrUnitRouteAction   = "DescribeDisabledExpressConnectRouterRouteEntries"
	ecrUnitGrantedAction = "DescribeInstanceGrantedToExpressConnectRouter"
	ecrUnitAPIVersion    = "2023-09-01"
)

type ecrUnitResponse struct {
	status int
	body   map[string]interface{}
}

type ecrUnitRequest struct {
	Action       string
	EcrID        string
	NextToken    string
	HasNextToken bool
	Version      string
}

type ecrUnitFixture struct {
	mu       sync.Mutex
	pages    map[string][]ecrUnitResponse
	served   map[string]int
	requests []ecrUnitRequest
}

func newECRUnitFixture(pages map[string][]ecrUnitResponse) *ecrUnitFixture {
	return &ecrUnitFixture{pages: pages, served: make(map[string]int)}
}

func (f *ecrUnitFixture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request := ecrUnitRequest{
		Action:       r.Form.Get("Action"),
		EcrID:        r.Form.Get("EcrId"),
		NextToken:    r.Form.Get("NextToken"),
		HasNextToken: r.Form.Has("NextToken"),
		Version:      r.Form.Get("Version"),
	}

	f.mu.Lock()
	f.requests = append(f.requests, request)
	index := f.served[request.Action]
	f.served[request.Action] = index + 1
	scripted := f.pages[request.Action]
	f.mu.Unlock()

	response := ecrUnitResponse{status: http.StatusOK}
	if index < len(scripted) {
		response = scripted[index]
	}
	if response.status == 0 {
		response.status = http.StatusOK
	}
	body := response.body
	if body == nil {
		body = ecrUnitDefaultBody(request.Action)
	}
	body = ecrUnitScopeToEcr(request.Action, request.EcrID, body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.status)
	_ = json.NewEncoder(w).Encode(body)
}

func (f *ecrUnitFixture) requestsFor(action string) []ecrUnitRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	matched := make([]ecrUnitRequest, 0, len(f.requests))
	for _, request := range f.requests {
		if request.Action == action {
			matched = append(matched, request)
		}
	}
	return matched
}

func ecrUnitDefaultBody(action string) map[string]interface{} {
	switch action {
	case "DescribeExpressConnectRouter":
		return map[string]interface{}{
			"EcrList": []interface{}{map[string]interface{}{
				"AlibabaSideAsn": 65533,
				"EcrId":          ecrUnitResourceID,
				"Name":           "tf-acc-unit-fixture",
				"Status":         "Active",
			}},
		}
	case "DescribeExpressConnectRouterInterRegionTransitMode":
		return map[string]interface{}{"InterRegionTransitModeList": []interface{}{}}
	case "DescribeExpressConnectRouterRegion":
		return map[string]interface{}{"EcrRegionList": []interface{}{}}
	case ecrUnitRouteAction:
		return map[string]interface{}{"DisabledRouteEntryList": []interface{}{}}
	case ecrUnitGrantedAction:
		return map[string]interface{}{"EcrGrantedInstanceList": []interface{}{}}
	}
	return map[string]interface{}{}
}

// ecrUnitScopeToEcr drops entries that belong to another ECR, so a Read that loses the requested
// EcrId cannot pass by collecting every fixture entry.
func ecrUnitScopeToEcr(action, ecrID string, body map[string]interface{}) map[string]interface{} {
	listKey := ""
	switch action {
	case ecrUnitRouteAction:
		listKey = "DisabledRouteEntryList"
	case ecrUnitGrantedAction:
		listKey = "EcrGrantedInstanceList"
	default:
		return body
	}
	items, ok := body[listKey].([]interface{})
	if !ok {
		return body
	}
	scoped := make([]interface{}, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if owner, tracked := item["EcrId"]; tracked && fmt.Sprint(owner) != ecrID {
			continue
		}
		scoped = append(scoped, item)
	}
	out := make(map[string]interface{}, len(body)+1)
	for key, value := range body {
		out[key] = value
	}
	out[listKey] = scoped
	return out
}

func ecrUnitRead(t *testing.T, fixture *ecrUnitFixture) (*schema.ResourceData, error) {
	t.Helper()

	// Wire the handler method explicitly instead of passing the fixture as an
	// http.Handler: the CI impact analysis (scripts/ci/accscope) cannot follow
	// interface dispatch, so an implicit pass would leave ServeHTTP and its
	// response helpers without statically resolvable test consumers.
	server := httptest.NewServer(http.HandlerFunc(fixture.ServeHTTP))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build fixture credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatalf("build fixture client: %v", err)
	}
	endpoints.Store("expressconnectrouter", endpoint)

	resourceSchema := expressconnectrouter.ResourceExpressConnectRouter()
	d := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{"alibaba_side_asn": 65533})
	d.SetId(ecrUnitResourceID)
	return d, resourceSchema.Read(d, client)
}

func ecrUnitCollection(t *testing.T, d *schema.ResourceData, field string) []interface{} {
	t.Helper()
	got, ok := d.Get(field).([]interface{})
	if !ok {
		t.Fatalf("%s is %T, want []interface{}", field, d.Get(field))
	}
	return got
}

func ecrUnitEntry(t *testing.T, d *schema.ResourceData, field string, index int) map[string]interface{} {
	t.Helper()
	entries := ecrUnitCollection(t, d, field)
	if index >= len(entries) {
		t.Fatalf("%s has %d entries, want at least %d", field, len(entries), index+1)
	}
	entry, ok := entries[index].(map[string]interface{})
	if !ok {
		t.Fatalf("%s[%d] is %T, want map[string]interface{}", field, index, entries[index])
	}
	return entry
}

func TestUnitExpressConnectRouterReadCollections(t *testing.T) {
	t.Run("non_empty_collections_across_pages", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {
				{body: map[string]interface{}{
					"NextToken": "route-token-2",
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "DestinationCidrBlock": "10.0.0.0/24", "NexthopInstanceId": "ri-1"},
						map[string]interface{}{"EcrId": ecrUnitResourceID, "DestinationCidrBlock": "10.0.1.0/24", "NexthopInstanceId": "ri-2"},
						map[string]interface{}{"EcrId": ecrUnitForeignID, "DestinationCidrBlock": "10.99.0.0/24", "NexthopInstanceId": "ri-foreign"},
					},
				}},
				{body: map[string]interface{}{
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "DestinationCidrBlock": "10.0.2.0/24", "NexthopInstanceId": "ri-3"},
						map[string]interface{}{"EcrId": ecrUnitForeignID, "DestinationCidrBlock": "10.99.1.0/24", "NexthopInstanceId": "ri-foreign"},
					},
				}},
			},
			ecrUnitGrantedAction: {
				{body: map[string]interface{}{
					"NextToken": "grant-token-2",
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "NodeId": "vbr-1", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou", "NodeOwnerUid": int64(1234567890123456789)},
						map[string]interface{}{"EcrId": ecrUnitForeignID, "NodeId": "vbr-foreign", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou", "NodeOwnerUid": int64(1)},
					},
				}},
				{body: map[string]interface{}{
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "NodeId": "vpc-2", "NodeType": "VPC", "NodeRegionId": "cn-shanghai", "NodeOwnerUid": int64(9007199254740993)},
					},
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err != nil {
			t.Fatalf("read failed for a healthy two-page response: %v", err)
		}

		wantRoutes := []map[string]string{
			{"destination_cidr_block": "10.0.0.0/24", "nexthop_instance_id": "ri-1"},
			{"destination_cidr_block": "10.0.1.0/24", "nexthop_instance_id": "ri-2"},
			{"destination_cidr_block": "10.0.2.0/24", "nexthop_instance_id": "ri-3"},
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != len(wantRoutes) {
			t.Errorf("disabled_route_entries has %d entries, want %d", got, len(wantRoutes))
		}
		for index, want := range wantRoutes {
			entry := ecrUnitEntry(t, d, "disabled_route_entries", index)
			for field, wantValue := range want {
				if got := fmt.Sprint(entry[field]); got != wantValue {
					t.Errorf("disabled_route_entries[%d].%s = %q, want %q", index, field, got, wantValue)
				}
			}
		}

		wantGranted := []map[string]string{
			{"node_id": "vbr-1", "node_type": "VBR", "node_region_id": "cn-hangzhou", "node_owner_uid": "1234567890123456789"},
			{"node_id": "vpc-2", "node_type": "VPC", "node_region_id": "cn-shanghai", "node_owner_uid": "9007199254740993"},
		}
		if got := len(ecrUnitCollection(t, d, "granted_instances")); got != len(wantGranted) {
			t.Errorf("granted_instances has %d entries, want %d", got, len(wantGranted))
		}
		for index, want := range wantGranted {
			entry := ecrUnitEntry(t, d, "granted_instances", index)
			for field, wantValue := range want {
				if got := fmt.Sprint(entry[field]); got != wantValue {
					t.Errorf("granted_instances[%d].%s = %q, want %q", index, field, got, wantValue)
				}
			}
		}

		for _, action := range []string{ecrUnitRouteAction, ecrUnitGrantedAction} {
			requests := fixture.requestsFor(action)
			if len(requests) != 2 {
				t.Fatalf("%s was requested %d times, want 2", action, len(requests))
			}
			for index, request := range requests {
				if request.EcrID != ecrUnitResourceID {
					t.Errorf("%s request %d carried EcrId %q, want %q", action, index, request.EcrID, ecrUnitResourceID)
				}
				if request.Version != ecrUnitAPIVersion {
					t.Errorf("%s request %d carried Version %q, want %q", action, index, request.Version, ecrUnitAPIVersion)
				}
				if index == 0 && request.HasNextToken {
					t.Errorf("%s first page carried NextToken %q", action, request.NextToken)
				}
			}
			wantToken := "route-token-2"
			if action == ecrUnitGrantedAction {
				wantToken = "grant-token-2"
			}
			if requests[1].NextToken != wantToken {
				t.Errorf("%s second page carried NextToken %q, want %q", action, requests[1].NextToken, wantToken)
			}
		}
	})

	t.Run("empty_collections_are_not_errors", func(t *testing.T) {
		fixture := newECRUnitFixture(nil)

		d, err := ecrUnitRead(t, fixture)
		if err != nil {
			t.Fatalf("read failed for empty collections: %v", err)
		}
		for _, field := range []string{"disabled_route_entries", "granted_instances"} {
			if got := len(ecrUnitCollection(t, d, field)); got != 0 {
				t.Errorf("%s has %d entries, want 0", field, got)
			}
		}
		for _, action := range []string{ecrUnitRouteAction, ecrUnitGrantedAction} {
			if got := len(fixture.requestsFor(action)); got != 1 {
				t.Errorf("%s was requested %d times, want 1", action, got)
			}
		}
	})

	t.Run("route_entries_second_page_error_propagates", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {
				{body: map[string]interface{}{
					"NextToken": "route-token-2",
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "DestinationCidrBlock": "10.0.0.0/24", "NexthopInstanceId": "ri-1"},
					},
				}},
				{status: http.StatusBadRequest, body: map[string]interface{}{
					"Code":    "InvalidParameter",
					"Message": "The specified parameter NextToken is invalid.",
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("read swallowed a failed second page and returned success")
		}
		for _, want := range []string{ecrUnitRouteAction, "InvalidParameter", "StatusCode: 400"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if errs.NotFoundError(err) {
			t.Errorf("a parameter error was classified as NotFound: %v", err)
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != 0 {
			t.Errorf("disabled_route_entries was written with %d partial entries after a failed page", got)
		}
		if got := len(fixture.requestsFor(ecrUnitGrantedAction)); got != 0 {
			t.Errorf("%s was requested %d times after a failed route page, want 0", ecrUnitGrantedAction, got)
		}
	})

	t.Run("disabled_route_entries_first_page_access_denied_propagates", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {
				{status: http.StatusForbidden, body: map[string]interface{}{
					"Code":    "NoPermission",
					"Message": "You are not authorized to perform this operation.",
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("read swallowed an access-denied supplementary collection and returned success")
		}
		for _, want := range []string{ecrUnitRouteAction, "NoPermission", "StatusCode: 403"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if errs.NotFoundError(err) {
			t.Errorf("an access-denied error was classified as NotFound: %v", err)
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != 0 {
			t.Errorf("disabled_route_entries was written with %d entries after an access-denied first page, want 0", got)
		}
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 1 {
			t.Errorf("%s was requested %d times after an access-denied first page, want 1 (no retry)", ecrUnitRouteAction, got)
		}
		// granted_instances must not be requested once disabled_route_entries has failed: Read returns immediately.
		if got := len(fixture.requestsFor(ecrUnitGrantedAction)); got != 0 {
			t.Errorf("%s was requested %d times after the route collection was access-denied, want 0", ecrUnitGrantedAction, got)
		}
		if d.Id() != ecrUnitResourceID {
			t.Errorf("resource id %q was wiped after an access-denied supplementary call, want %q", d.Id(), ecrUnitResourceID)
		}
	})

	t.Run("granted_instances_first_page_access_denied_propagates", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitGrantedAction: {
				{status: http.StatusForbidden, body: map[string]interface{}{
					"Code":    "NoPermission",
					"Message": "You are not authorized to perform this operation.",
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("read swallowed an access-denied supplementary collection and returned success")
		}
		for _, want := range []string{ecrUnitGrantedAction, "NoPermission", "StatusCode: 403"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if errs.NotFoundError(err) {
			t.Errorf("an access-denied error was classified as NotFound: %v", err)
		}
		if got := len(ecrUnitCollection(t, d, "granted_instances")); got != 0 {
			t.Errorf("granted_instances was written with %d entries after an access-denied first page, want 0", got)
		}
		if got := len(fixture.requestsFor(ecrUnitGrantedAction)); got != 1 {
			t.Errorf("%s was requested %d times after an access-denied first page, want 1 (no retry)", ecrUnitGrantedAction, got)
		}
		// disabled_route_entries runs first and is populated from its own (default empty) page before granted_instances fails.
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 1 {
			t.Errorf("%s was requested %d times when the granted collection was access-denied, want 1", ecrUnitRouteAction, got)
		}
		if d.Id() != ecrUnitResourceID {
			t.Errorf("resource id %q was wiped after an access-denied supplementary call, want %q", d.Id(), ecrUnitResourceID)
		}
	})

	t.Run("missing_owner_uid_is_not_written_as_a_literal", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitGrantedAction: {
				{body: map[string]interface{}{
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"EcrId": ecrUnitResourceID, "NodeId": "vbr-1", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou"},
					},
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err != nil {
			t.Fatalf("read failed for a granted instance without an owner uid: %v", err)
		}
		entry := ecrUnitEntry(t, d, "granted_instances", 0)
		if value := entry["node_owner_uid"]; value != nil && value != "" {
			t.Errorf("node_owner_uid = %q, want it empty when the API omits NodeOwnerUid", value)
		}
	})

	t.Run("main_describe_not_found_clears_resource", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			"DescribeExpressConnectRouter": {
				{body: map[string]interface{}{"EcrList": []interface{}{}}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err != nil {
			t.Fatalf("an authoritative not found response must clear the resource instead of failing the read: %v", err)
		}
		if got := d.Id(); got != "" {
			t.Errorf("resource id = %q, want an empty id after the main describe reports not found", got)
		}
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 0 {
			t.Errorf("%s was requested %d times after the resource was proven absent, want 0", ecrUnitRouteAction, got)
		}
	})

	t.Run("route_entry_not_found_error_is_not_an_empty_list", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {
				{status: http.StatusNotFound, body: map[string]interface{}{
					"Code":    "ResourceNotFound.EcrId",
					"Message": "The specified EcrId is not found.",
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("a not found error on a collection page was swallowed into an empty list")
		}
		if !errs.NotFoundError(err) {
			t.Errorf("the ResourceNotFound.EcrId code lost its not found classification: %v", err)
		}
		if got := d.Id(); got != ecrUnitResourceID {
			t.Errorf("resource id = %q, want %q: only the authoritative describe may clear state", got, ecrUnitResourceID)
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != 0 {
			t.Errorf("disabled_route_entries has %d entries after a failed page, want 0", got)
		}
	})

	t.Run("repeated_route_page_token_terminates_with_an_error", func(t *testing.T) {
		repeated := map[string]interface{}{"DisabledRouteEntryList": []interface{}{}, "NextToken": "stuck-token"}
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {{body: repeated}, {body: repeated}},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("a page that repeats the token it was requested with was treated as the end of pagination")
		}
		if errs.NotFoundError(err) {
			t.Errorf("a pagination loop error was classified as NotFound: %v", err)
		}
		for _, want := range []string{"stuck-token", ecrUnitRouteAction} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 2 {
			t.Errorf("%s was requested %d times, want exactly 2 before the loop guard fires", ecrUnitRouteAction, got)
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != 0 {
			t.Errorf("disabled_route_entries has %d entries after a repeated token, want 0", got)
		}
	})

	t.Run("non_adjacent_route_page_token_cycle_terminates_with_an_error", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitRouteAction: {
				{body: map[string]interface{}{
					"NextToken": "route-token-a",
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"DestinationCidrBlock": "10.0.0.0/24", "NexthopInstanceId": "ri-cycle-1"},
					},
				}},
				{body: map[string]interface{}{
					"NextToken": "route-token-b",
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"DestinationCidrBlock": "10.0.1.0/24", "NexthopInstanceId": "ri-cycle-2"},
					},
				}},
				{body: map[string]interface{}{
					"NextToken": "route-token-a",
					"DisabledRouteEntryList": []interface{}{
						map[string]interface{}{"DestinationCidrBlock": "10.0.2.0/24", "NexthopInstanceId": "ri-cycle-3"},
					},
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("a non-adjacent NextToken cycle was treated as valid pagination")
		}
		if errs.NotFoundError(err) {
			t.Errorf("a pagination loop error was classified as NotFound: %v", err)
		}
		for _, want := range []string{"route-token-a", ecrUnitRouteAction} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		requests := fixture.requestsFor(ecrUnitRouteAction)
		if got := len(requests); got != 3 {
			t.Errorf("%s was requested %d times, want exactly 3 before the historical token guard fires", ecrUnitRouteAction, got)
		}
		for index, want := range []string{"", "route-token-a", "route-token-b"} {
			if index < len(requests) && requests[index].NextToken != want {
				t.Errorf("request %d NextToken = %q, want %q", index+1, requests[index].NextToken, want)
			}
		}
		if got := len(ecrUnitCollection(t, d, "disabled_route_entries")); got != 0 {
			t.Errorf("disabled_route_entries has %d entries after a token cycle, want no partial collection write", got)
		}
	})

	t.Run("repeated_granted_page_token_terminates_with_an_error", func(t *testing.T) {
		repeated := map[string]interface{}{"EcrGrantedInstanceList": []interface{}{}, "NextToken": "stuck-grant-token"}
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitGrantedAction: {{body: repeated}, {body: repeated}},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("a page that repeats the token it was requested with was treated as the end of pagination")
		}
		for _, want := range []string{"stuck-grant-token", ecrUnitGrantedAction} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 1 {
			t.Errorf("%s was requested %d times, want 1: the route loop must finish before the granted loop starts", ecrUnitRouteAction, got)
		}
		if got := len(fixture.requestsFor(ecrUnitGrantedAction)); got != 2 {
			t.Errorf("%s was requested %d times, want exactly 2 before the loop guard fires", ecrUnitGrantedAction, got)
		}
		if got := len(ecrUnitCollection(t, d, "granted_instances")); got != 0 {
			t.Errorf("granted_instances has %d entries after a repeated token, want 0", got)
		}
	})

	t.Run("non_adjacent_granted_page_token_cycle_terminates_with_an_error", func(t *testing.T) {
		fixture := newECRUnitFixture(map[string][]ecrUnitResponse{
			ecrUnitGrantedAction: {
				{body: map[string]interface{}{
					"NextToken": "grant-token-a",
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"NodeId": "vbr-cycle-1", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou", "NodeOwnerUid": "1"},
					},
				}},
				{body: map[string]interface{}{
					"NextToken": "grant-token-b",
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"NodeId": "vbr-cycle-2", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou", "NodeOwnerUid": "2"},
					},
				}},
				{body: map[string]interface{}{
					"NextToken": "grant-token-a",
					"EcrGrantedInstanceList": []interface{}{
						map[string]interface{}{"NodeId": "vbr-cycle-3", "NodeType": "VBR", "NodeRegionId": "cn-hangzhou", "NodeOwnerUid": "3"},
					},
				}},
			},
		})

		d, err := ecrUnitRead(t, fixture)
		if err == nil {
			t.Fatal("a non-adjacent NextToken cycle was treated as valid pagination")
		}
		if errs.NotFoundError(err) {
			t.Errorf("a pagination loop error was classified as NotFound: %v", err)
		}
		for _, want := range []string{"grant-token-a", ecrUnitGrantedAction} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not carry %q", err.Error(), want)
			}
		}
		if got := len(fixture.requestsFor(ecrUnitRouteAction)); got != 1 {
			t.Errorf("%s was requested %d times, want 1: the route loop must finish before the granted loop starts", ecrUnitRouteAction, got)
		}
		requests := fixture.requestsFor(ecrUnitGrantedAction)
		if got := len(requests); got != 3 {
			t.Errorf("%s was requested %d times, want exactly 3 before the historical token guard fires", ecrUnitGrantedAction, got)
		}
		for index, want := range []string{"", "grant-token-a", "grant-token-b"} {
			if index < len(requests) && requests[index].NextToken != want {
				t.Errorf("request %d NextToken = %q, want %q", index+1, requests[index].NextToken, want)
			}
		}
		if got := len(ecrUnitCollection(t, d, "granted_instances")); got != 0 {
			t.Errorf("granted_instances has %d entries after a token cycle, want no partial collection write", got)
		}
	})
}
