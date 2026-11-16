package expressconnectrouter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// vpcAssociationImportStep verifies that every successful configuration step is
// importable on its own, not only the final state of the lifecycle.
func vpcAssociationImportStep(resourceId string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceId,
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test ExpressConnectRouter ExpressConnectRouterVpcAssociation. >>> Resource test cases, automatically generated.
// Case 初始版本测试用例 6348
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.5.0/25", "172.16.6.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "2",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{},
				}),
				ExpectError: regexp.MustCompile(`allowed_prefixes.*at least one prefix`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          acctest.CHECKSET,
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
		},
	})
}

var AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348 = map[string]string{
	"status":         acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"association_id": acctest.CHECKSET,
}

func AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default8qAtD6" {
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_express_connect_router_express_connect_router" "defaultM9YxGW" {
  alibaba_side_asn = "65533"
}


`, name)
}

// Case 初始版本测试用例 6348  twin
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          acctest.CHECKSET,
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
		},
	})
}

// Case 初始版本测试用例 6348  raw
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_basic6348_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
					"vpc_owner_id": "1511928242963727",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
						"allowed_prefixes.#":    "4",
						"vpc_owner_id":          acctest.CHECKSET,
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.5.0/25", "172.16.6.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "2",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{},
				}),
				ExpectError: regexp.MustCompile(`allowed_prefixes.*at least one prefix`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
		},
	})
}

// Case VpcAssociation allowed_prefixes collection order 6348
// The provider ignores order-only changes while preserving import and convergence coverage.
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_allowedPrefixesOrder6348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24", "172.16.4.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Re-apply the unchanged baseline to establish an applied order without
				// a create or update. The following reorder is a separate no-op too.
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24", "172.16.4.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			{
				// The SDK requires an empty post-apply plan, which proves the
				// reordered members converged.
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.4.0/24", "172.16.3.0/24", "172.16.2.0/24", "172.16.1.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#": "4",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Reorder again in a plan-only step: a clean plan over the ascending
				// order after the descending apply keeps the convergence claim without
				// a third real prefix update.
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24", "172.16.4.0/24"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Case VpcAssociation complete scenario: description lifecycle (set, update, clear)
// and allowed_prefixes_mode switching 6348
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_complete6348_descriptionAndModeLifecycle(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationMap6348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnectrouter.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectRouterExpressConnectRouterVpcAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create with MatchMode and a description.
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24", "172.16.4.0/24"},
					"allowed_prefixes_mode": "MatchMode",
					"description":           "tf-testacc-description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"association_region_id": acctest.CHECKSET,
						"vpc_id":                acctest.CHECKSET,
						"ecr_id":                acctest.CHECKSET,
						"allowed_prefixes.#":    "4",
						"allowed_prefixes_mode": "MatchMode",
						"description":           "tf-testacc-description",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Update only the description through ModifyExpressConnectRouterAssociation.
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testacc-description-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes_mode": "MatchMode",
						"description":           "tf-testacc-description-update",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Clear the description with an explicit empty string. The API accepts a
				// 0-128 character description, so "" is a legal value that clears it.
				Config: testAccConfig(map[string]interface{}{
					"description": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Switch the mode to IncrementalMode and set a description again.
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes_mode": "IncrementalMode",
					"description":           "tf-testacc-description-incremental",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes_mode": "IncrementalMode",
						"description":           "tf-testacc-description-incremental",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				// Change only allowed_prefixes with the mode untouched: the update must
				// send the stored mode along with the prefixes, because the API contract
				// does not declare the mode used when AllowedPrefixesMode is omitted.
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allowed_prefixes.#":    "3",
						"allowed_prefixes_mode": "IncrementalMode",
					}),
				),
			},
			vpcAssociationImportStep(resourceId),
			{
				Config: testAccConfig(map[string]interface{}{
					"allowed_prefixes": []string{
						"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Case VpcAssociation invalid mode-only configuration 6348: rejected locally at plan time.
func TestAccAliCloudExpressConnectRouterExpressConnectRouterVpcAssociation_invalidModeOnly6348(t *testing.T) {
	resourceId := "alicloud_express_connect_router_vpc_association.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudExpressConnectRouterExpressConnectRouterVpcAssociationBasicDependence6348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"association_region_id": acctest.DefaultRegionToTest,
					"vpc_id":                "${alicloud_vpc.default8qAtD6.id}",
					"ecr_id":                "${alicloud_express_connect_router_express_connect_router.defaultM9YxGW.id}",
					"allowed_prefixes_mode": "MatchMode",
				}),
				ExpectError: regexp.MustCompile(`allowed_prefixes_mode.*non-empty.*allowed_prefixes`),
			},
		},
	})
}

// Test ExpressConnectRouter ExpressConnectRouterVpcAssociation. <<< Resource test cases, automatically generated.

func newVpcAssociationUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build credential: %v", err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client: %v", err)
	}
	endpoints.Store("expressconnectrouter", endpoint)
	return client
}

func writeVpcAssociationJSON(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

func writeVpcAssociationActive(t *testing.T, w http.ResponseWriter) {
	writeVpcAssociationJSON(t, w, http.StatusOK, map[string]interface{}{
		"RequestId": "unit-test",
		"AssociationList": []map[string]interface{}{{
			"AssociationId":       "assoc-1",
			"AssociationNodeType": "VPC",
			"EcrId":               "ecr-1",
			"VpcId":               "vpc-1",
			"RegionId":            "cn-hangzhou",
			"Status":              "ACTIVE",
			"Description":         "unit-description",
			"AllowedPrefixesMode": "MatchMode",
			"AllowedPrefixes":     []string{"172.16.1.0/24"},
			"VpcOwnerId":          1,
		}},
	})
}

func vpcAssociationUnitConfig() map[string]interface{} {
	return map[string]interface{}{
		"ecr_id":                "ecr-1",
		"vpc_id":                "vpc-1",
		"association_region_id": "cn-hangzhou",
	}
}

func vpcAssociationUnitState(overrides map[string]string) *terraform.InstanceState {
	state := &terraform.InstanceState{
		ID: "ecr-1:assoc-1:vpc-1",
		Attributes: map[string]string{
			"id":                    "ecr-1:assoc-1:vpc-1",
			"ecr_id":                "ecr-1",
			"vpc_id":                "vpc-1",
			"association_id":        "assoc-1",
			"association_region_id": "cn-hangzhou",
		},
	}
	for key, value := range overrides {
		state.Attributes[key] = value
	}
	return state
}

func vpcAssociationRawConfigValues(config map[string]interface{}) map[string]cty.Value {
	values := map[string]cty.Value{
		"allowed_prefixes":      cty.NullVal(cty.List(cty.String)),
		"allowed_prefixes_mode": cty.NullVal(cty.String),
		"association_id":        cty.NullVal(cty.String),
		"association_region_id": cty.NullVal(cty.String),
		"create_time":           cty.NullVal(cty.String),
		"description":           cty.NullVal(cty.String),
		"ecr_id":                cty.NullVal(cty.String),
		"status":                cty.NullVal(cty.String),
		"vpc_id":                cty.NullVal(cty.String),
		"vpc_owner_id":          cty.NullVal(cty.Number),
	}
	for _, key := range []string{"allowed_prefixes_mode", "association_region_id", "description", "ecr_id", "vpc_id"} {
		if value, ok := config[key]; ok {
			values[key] = cty.StringVal(value.(string))
		}
	}
	if raw, ok := config["allowed_prefixes"]; ok {
		switch prefixes := raw.(type) {
		case string:
			values["allowed_prefixes"] = cty.UnknownVal(cty.List(cty.String))
		case []string:
			values["allowed_prefixes"] = vpcAssociationRawPrefixList(prefixes)
		case []interface{}:
			if len(prefixes) == 0 {
				values["allowed_prefixes"] = cty.ListValEmpty(cty.String)
			} else {
				values["allowed_prefixes"] = cty.ListVal(vpcAssociationRawPrefixValues(prefixes))
			}
		}
	}
	return values
}

func vpcAssociationRawPrefixValues(prefixes []interface{}) []cty.Value {
	values := make([]cty.Value, 0, len(prefixes))
	for _, prefix := range prefixes {
		values = append(values, cty.StringVal(prefix.(string)))
	}
	return values
}

func vpcAssociationRawPrefixList(prefixes []string) cty.Value {
	if len(prefixes) == 0 {
		return cty.ListValEmpty(cty.String)
	}
	values := make([]cty.Value, 0, len(prefixes))
	for _, prefix := range prefixes {
		values = append(values, cty.StringVal(prefix))
	}
	return cty.ListVal(values)
}

func vpcAssociationRawConfig(config map[string]interface{}) cty.Value {
	return cty.ObjectVal(vpcAssociationRawConfigValues(config))
}

func vpcAssociationRawConfigWithValues(config map[string]interface{}, mode, prefixes cty.Value) cty.Value {
	values := vpcAssociationRawConfigValues(config)
	values["allowed_prefixes_mode"] = mode
	values["allowed_prefixes"] = prefixes
	return cty.ObjectVal(values)
}

func vpcAssociationApply(t *testing.T, r *schema.Resource, state *terraform.InstanceState, config map[string]interface{}, client *connectivity.AliyunClient) (*terraform.InstanceState, diag.Diagnostics) {
	t.Helper()
	rawConfig := vpcAssociationRawConfig(config)
	if state != nil {
		state.RawConfig = rawConfig
	}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if diff == nil {
		return nil, diag.FromErr(fmt.Errorf("expected a lifecycle diff"))
	}
	diff.RawConfig = rawConfig
	return r.Apply(context.Background(), state, diff, client)
}

func TestUnitVpcAssociationCreateRetryCodes(t *testing.T) {
	for _, code := range []string{
		"Conflict.Lock",
		"IncorrectStatus.Ecr",
		"IncorrectStatus.Vpc",
		"IncorrectStatus.VpcOperating",
		"IncorrectStatus.Vswitch",
		"IncorrectStatus.RouteEntry",
		"OperationDenied.ForceDeleting",
		"OperationFailed.CallingVpcYaochi",
		"InternalErrorTemp",
	} {
		code := code
		t.Run(code, func(t *testing.T) {
			var createCalls int32
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				switch r.Form.Get("Action") {
				case "CreateExpressConnectRouterAssociation":
					if atomic.AddInt32(&createCalls, 1) == 1 {
						writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
							"Code": code, "Message": "transient " + code, "RequestId": "unit-test",
						})
						return
					}
					writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{
						"AssociationId": "assoc-1", "RequestId": "unit-test",
					})
				case "DescribeExpressConnectRouterAssociation":
					writeVpcAssociationActive(t, w)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
					writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
						"Code": "InvalidParameter", "Message": "unexpected action",
					})
				}
			})

			state, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, vpcAssociationUnitConfig(), client)
			if diags.HasError() {
				t.Fatalf("create with transient %s should succeed after a retry: %v", code, diags)
			}
			if got := atomic.LoadInt32(&createCalls); got != 2 {
				t.Fatalf("CreateExpressConnectRouterAssociation calls = %d, want 2", got)
			}
			if got, want := state.ID, "ecr-1:assoc-1:vpc-1"; got != want {
				t.Fatalf("resource id = %q, want %q", got, want)
			}
		})
	}
}

func TestUnitVpcAssociationCreateTerminalCodesAreNotRetried(t *testing.T) {
	for _, code := range []string{
		"InvalidParameter",
		"QuotaExceeded.VpcCountPerEcr",
		"IllegalParamFormat.AllowedPrefixes",
		"ResourceAlreadyExist.VpcId",
		"UnsupportedRegion",
	} {
		code := code
		t.Run(code, func(t *testing.T) {
			var createCalls int32
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				if r.Form.Get("Action") != "CreateExpressConnectRouterAssociation" {
					t.Errorf("terminal error must not trigger %q", r.Form.Get("Action"))
					writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
						"Code": "InvalidParameter", "Message": "unexpected action",
					})
					return
				}
				atomic.AddInt32(&createCalls, 1)
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": code, "Message": "terminal " + code, "RequestId": "unit-test",
				})
			})

			state, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, vpcAssociationUnitConfig(), client)
			if !diags.HasError() {
				t.Fatalf("terminal error %s must be returned", code)
			}
			if got := atomic.LoadInt32(&createCalls); got != 1 {
				t.Fatalf("CreateExpressConnectRouterAssociation calls = %d, want 1", got)
			}
			if state != nil && state.ID != "" {
				t.Fatalf("failed create assigned id %q", state.ID)
			}
		})
	}
}

func TestUnitVpcAssociationCreateSendsModeOnlyWithPrefixes(t *testing.T) {
	for _, tc := range []struct {
		name     string
		setMode  bool
		mode     string
		wantMode string
	}{
		{name: "mode_with_prefixes_is_sent", setMode: true, mode: "MatchMode", wantMode: "MatchMode"},
		{name: "empty_mode_with_prefixes_is_omitted", setMode: true, mode: ""},
		{name: "unset_mode_with_prefixes_is_omitted"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var lastForm url.Values
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				switch r.Form.Get("Action") {
				case "CreateExpressConnectRouterAssociation":
					lastForm = r.Form
					writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{
						"AssociationId": "assoc-1", "RequestId": "unit-test",
					})
				case "DescribeExpressConnectRouterAssociation":
					writeVpcAssociationActive(t, w)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})

			config := vpcAssociationUnitConfig()
			config["allowed_prefixes"] = []interface{}{"172.16.1.0/24"}
			config["description"] = "unit-description"
			if tc.setMode {
				config["allowed_prefixes_mode"] = tc.mode
			}
			_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, config, client)
			if diags.HasError() {
				t.Fatalf("create with prefixes: %v", diags)
			}
			if lastForm == nil {
				t.Fatal("CreateExpressConnectRouterAssociation request was not captured")
			}
			assertFormArray(t, lastForm, "AllowedPrefixes", "172.16.1.0/24")
			if got := lastForm.Get("Description"); got != "unit-description" {
				t.Errorf("Description = %q, want unit-description", got)
			}
			values, present := lastForm["AllowedPrefixesMode"]
			if tc.wantMode == "" {
				if present {
					t.Errorf("AllowedPrefixesMode must be omitted, got %v", values)
				}
			} else if !present || len(values) != 1 || values[0] != tc.wantMode {
				t.Errorf("AllowedPrefixesMode present=%v values=%v, want %s", present, values, tc.wantMode)
			}
		})
	}

	t.Run("mode_without_prefixes_is_rejected_before_request", func(t *testing.T) {
		var requests int32
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			t.Errorf("unexpected request for an invalid mode-only configuration")
			writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
				"Code": "InvalidParameter", "Message": "unexpected request",
			})
		})

		config := vpcAssociationUnitConfig()
		config["allowed_prefixes_mode"] = "MatchMode"
		_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, config, client)
		if !diags.HasError() || !containsAll(fmt.Sprintf("%v", diags), "allowed_prefixes_mode", "non-empty", "allowed_prefixes") {
			t.Fatalf("mode-only create error = %v, want the local combination error", diags)
		}
		if got := atomic.LoadInt32(&requests); got != 0 {
			t.Fatalf("invalid mode-only create issued %d requests, want 0", got)
		}
	})
}

func TestUnitVpcAssociationPrefixOrderSuppress(t *testing.T) {
	for _, tc := range []struct {
		name       string
		oldValue   []string
		newValue   []string
		suppressed bool
	}{
		{name: "ReorderOnly", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.2.0/24", "172.16.1.0/24"}, suppressed: true},
		{name: "LongerReorder", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"}, newValue: []string{"172.16.3.0/24", "172.16.1.0/24", "172.16.2.0/24"}, suppressed: true},
		{name: "MemberAdded", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"}},
		{name: "MemberRemoved", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24"}},
		{name: "MemberReplaced", oldValue: []string{"172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.9.0/24"}},
		{name: "MultiplicityChanged", oldValue: []string{"172.16.1.0/24", "172.16.1.0/24", "172.16.2.0/24"}, newValue: []string{"172.16.1.0/24", "172.16.2.0/24", "172.16.2.0/24"}},
		{name: "OldEmpty", newValue: []string{"172.16.1.0/24"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attributes := map[string]string{
				"id":                    "ecr-1:assoc-1:vpc-1",
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_id":        "assoc-1",
				"association_region_id": "cn-hangzhou",
				"allowed_prefixes.#":    strconv.Itoa(len(tc.oldValue)),
				"allowed_prefixes_mode": "MatchMode",
			}
			for i, prefix := range tc.oldValue {
				attributes[fmt.Sprintf("allowed_prefixes.%d", i)] = prefix
			}
			state := &terraform.InstanceState{ID: "ecr-1:assoc-1:vpc-1", Attributes: attributes}
			newPrefixes := make([]interface{}, 0, len(tc.newValue))
			for _, prefix := range tc.newValue {
				newPrefixes = append(newPrefixes, prefix)
			}
			config := map[string]interface{}{
				"ecr_id":                "ecr-1",
				"vpc_id":                "vpc-1",
				"association_region_id": "cn-hangzhou",
				"allowed_prefixes":      newPrefixes,
			}
			state.RawConfig = vpcAssociationRawConfig(config)
			diff, err := expressconnectrouter.ResourceVpcAssociation().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatalf("build prefix diff: %v", err)
			}
			changed := diff != nil && !diff.Empty()
			if changed == tc.suppressed {
				t.Fatalf("changed = %t, want %t for old=%v new=%v", changed, !tc.suppressed, tc.oldValue, tc.newValue)
			}
		})
	}
}

func TestUnitVpcAssociationCreateOptionalZeroValues(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		configureDescription   bool
		description            string
		wantDescription        string
		wantDescriptionPresent bool
	}{
		{name: "description_omitted"},
		{name: "description_explicit_empty", configureDescription: true, description: "", wantDescriptionPresent: true},
		{name: "description_non_empty", configureDescription: true, description: "unit-description", wantDescription: "unit-description", wantDescriptionPresent: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var createForm url.Values
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse request: %v", err)
				}
				switch r.Form.Get("Action") {
				case "CreateExpressConnectRouterAssociation":
					createForm = r.Form
					writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{
						"AssociationId": "assoc-1", "RequestId": "unit-test",
					})
				case "DescribeExpressConnectRouterAssociation":
					writeVpcAssociationActive(t, w)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})

			config := vpcAssociationUnitConfig()
			if tc.configureDescription {
				config["description"] = tc.description
			}
			_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, config, client)
			if diags.HasError() {
				t.Fatalf("create optional values: %v", diags)
			}
			if createForm == nil {
				t.Fatal("CreateExpressConnectRouterAssociation request was not captured")
			}
			values, present := createForm["Description"]
			if present != tc.wantDescriptionPresent {
				t.Fatalf("Description presence = %t, want %t; form=%v", present, tc.wantDescriptionPresent, createForm)
			}
			if present && (len(values) != 1 || values[0] != tc.wantDescription) {
				t.Fatalf("Description values = %v, want [%q]", values, tc.wantDescription)
			}
			if formHasParameter(createForm, "AllowedPrefixes") {
				t.Fatalf("AllowedPrefixes must be omitted when unconfigured: %v", createForm)
			}
		})
	}

	for _, tc := range []struct {
		name string
		mode string
	}{
		{name: "explicit_empty_prefixes"},
		{name: "explicit_empty_prefixes_with_mode", mode: "MatchMode"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var requests int32
			client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&requests, 1)
				t.Errorf("unexpected request for an explicit empty prefix list")
			})
			config := vpcAssociationUnitConfig()
			config["allowed_prefixes"] = []interface{}{}
			if tc.mode != "" {
				config["allowed_prefixes_mode"] = tc.mode
			}
			_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), nil, config, client)
			if !diags.HasError() || !containsAll(fmt.Sprintf("%v", diags), "allowed_prefixes", "at least one prefix") {
				t.Fatalf("explicit empty prefix error = %v, want the field-specific validation error", diags)
			}
			if got := atomic.LoadInt32(&requests); got != 0 {
				t.Fatalf("explicit empty prefix create issued %d requests, want 0", got)
			}
		})
	}
}

func TestUnitVpcAssociationUpdateDescriptionRetry(t *testing.T) {
	var modifyCalls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			if atomic.AddInt32(&modifyCalls, 1) == 1 {
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": "IncorrectStatus.Ecr", "Message": "another operation is running", "RequestId": "unit-test",
				})
				return
			}
			if got := r.Form.Get("Description"); got != "unit-description" {
				t.Errorf("Description = %q, want unit-description", got)
			}
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	config := vpcAssociationUnitConfig()
	config["description"] = "unit-description"
	state := vpcAssociationUnitState(map[string]string{"description": "old-description"})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if diags.HasError() {
		t.Fatalf("description update should retry a recoverable status: %v", diags)
	}
	if got := atomic.LoadInt32(&modifyCalls); got != 2 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 2", got)
	}
}

func TestUnitVpcAssociationUpdateParameterErrorFailsFast(t *testing.T) {
	var calls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			atomic.AddInt32(&calls, 1)
			writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
				"Code": "IllegalParamFormat.Description", "Message": "description format is illegal", "RequestId": "unit-test",
			})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	config := vpcAssociationUnitConfig()
	config["description"] = "unit-description"
	state := vpcAssociationUnitState(map[string]string{"description": "old-description"})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if !diags.HasError() {
		t.Fatal("parameter errors must be returned to the caller")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 1 (parameter errors must not retry)", got)
	}
}

func TestUnitVpcAssociationUpdateSharesOneDeadline(t *testing.T) {
	var modifies, describes int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			atomic.AddInt32(&modifies, 1)
			time.Sleep(50 * time.Millisecond)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			atomic.AddInt32(&modifies, 1)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			atomic.AddInt32(&describes, 1)
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	r := expressconnectrouter.ResourceVpcAssociation()
	r.Timeouts.Update = schema.DefaultTimeout(20 * time.Millisecond)
	config := vpcAssociationUnitConfig()
	config["description"] = "new-description"
	config["allowed_prefixes"] = []interface{}{"172.16.2.0/24"}
	state := vpcAssociationUnitState(map[string]string{
		"description":        "old-description",
		"allowed_prefixes.#": "1",
		"allowed_prefixes.0": "172.16.1.0/24",
	})
	_, diags := vpcAssociationApply(t, r, state, config, client)
	if !diags.HasError() || !containsAll(fmt.Sprintf("%v", diags), "Timeout", "shared update deadline") {
		t.Fatalf("shared budget error = %v, want the shared deadline timeout", diags)
	}
	if got := atomic.LoadInt32(&modifies); got != 1 {
		t.Fatalf("modify calls = %d, want 1; a per-step budget would have issued the second write", got)
	}
	if got := atomic.LoadInt32(&describes); got != 0 {
		t.Fatalf("describe calls = %d, want 0 after the budget is exhausted", got)
	}
}

func TestUnitVpcAssociationUpdateDescriptionClearSendsEmptyString(t *testing.T) {
	var forms []url.Values
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociation":
			forms = append(forms, r.Form)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	config := vpcAssociationUnitConfig()
	config["description"] = ""
	state := vpcAssociationUnitState(map[string]string{"description": "old-description"})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if diags.HasError() {
		t.Fatalf("clear description: %v", diags)
	}
	if len(forms) != 1 {
		t.Fatalf("ModifyExpressConnectRouterAssociation calls = %d, want 1", len(forms))
	}
	values, present := forms[0]["Description"]
	if !present || len(values) != 1 || values[0] != "" {
		t.Fatalf("Description present=%v values=%v, want one explicit empty value", present, values)
	}
}

func TestUnitVpcAssociationUpdatePrefixPresenceTransitions(t *testing.T) {
	t.Run("addition_from_unset_sends_prefix_without_mode", func(t *testing.T) {
		var forms []url.Values
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse request: %v", err)
			}
			switch r.Form.Get("Action") {
			case "ModifyExpressConnectRouterAssociationAllowedPrefix":
				forms = append(forms, r.Form)
				writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
			case "DescribeExpressConnectRouterAssociation":
				writeVpcAssociationActive(t, w)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		config := vpcAssociationUnitConfig()
		config["allowed_prefixes"] = []interface{}{"172.16.2.0/24"}
		_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), vpcAssociationUnitState(nil), config, client)
		if diags.HasError() {
			t.Fatalf("add prefix configuration: %v", diags)
		}
		if len(forms) != 1 {
			t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 1", len(forms))
		}
		assertFormArray(t, forms[0], "AllowedPrefixes", "172.16.2.0/24")
		if _, present := forms[0]["AllowedPrefixesMode"]; present {
			t.Fatalf("AllowedPrefixesMode must remain omitted when unconfigured: %v", forms[0])
		}
	})

	t.Run("removal_omits_prefix_and_mode_parameters", func(t *testing.T) {
		var forms []url.Values
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse request: %v", err)
			}
			switch r.Form.Get("Action") {
			case "ModifyExpressConnectRouterAssociationAllowedPrefix":
				forms = append(forms, r.Form)
				writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
			case "DescribeExpressConnectRouterAssociation":
				writeVpcAssociationActive(t, w)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		state := vpcAssociationUnitState(map[string]string{
			"allowed_prefixes.#":    "1",
			"allowed_prefixes.0":    "172.16.1.0/24",
			"allowed_prefixes_mode": "MatchMode",
		})
		_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, vpcAssociationUnitConfig(), client)
		if diags.HasError() {
			t.Fatalf("remove prefix configuration: %v", diags)
		}
		if len(forms) != 1 {
			t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 1", len(forms))
		}
		if formHasParameter(forms[0], "AllowedPrefixes") {
			t.Fatalf("AllowedPrefixes must be omitted after configuration removal: %v", forms[0])
		}
		if _, present := forms[0]["AllowedPrefixesMode"]; present {
			t.Fatalf("AllowedPrefixesMode must be omitted after configuration removal: %v", forms[0])
		}
	})

	t.Run("explicit_empty_is_rejected_before_any_update", func(t *testing.T) {
		var requests int32
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&requests, 1)
			t.Errorf("unexpected request for an explicit empty prefix list")
		})
		config := vpcAssociationUnitConfig()
		config["allowed_prefixes"] = []interface{}{}
		config["description"] = "new-description"
		state := vpcAssociationUnitState(map[string]string{
			"allowed_prefixes.#": "1",
			"allowed_prefixes.0": "172.16.1.0/24",
			"description":        "old-description",
		})
		_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
		if !diags.HasError() || !containsAll(fmt.Sprintf("%v", diags), "allowed_prefixes", "at least one prefix") {
			t.Fatalf("explicit empty prefix update error = %v, want the field-specific validation error", diags)
		}
		if got := atomic.LoadInt32(&requests); got != 0 {
			t.Fatalf("explicit empty prefix update issued %d requests, want 0", got)
		}
	})

	t.Run("unrelated_update_does_not_modify_prefixes", func(t *testing.T) {
		var prefixCalls int32
		client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("parse request: %v", err)
			}
			switch r.Form.Get("Action") {
			case "ModifyExpressConnectRouterAssociation":
				writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
			case "ModifyExpressConnectRouterAssociationAllowedPrefix":
				atomic.AddInt32(&prefixCalls, 1)
				writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
			case "DescribeExpressConnectRouterAssociation":
				writeVpcAssociationActive(t, w)
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})
		config := vpcAssociationUnitConfig()
		config["description"] = "new-description"
		config["allowed_prefixes"] = []interface{}{"172.16.1.0/24"}
		state := vpcAssociationUnitState(map[string]string{
			"allowed_prefixes.#": "1",
			"allowed_prefixes.0": "172.16.1.0/24",
			"description":        "old-description",
		})
		_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
		if diags.HasError() {
			t.Fatalf("unrelated description update: %v", diags)
		}
		if got := atomic.LoadInt32(&prefixCalls); got != 0 {
			t.Fatalf("unrelated update issued %d prefix modifications, want 0", got)
		}
	})
}

func TestUnitVpcAssociationUpdatePrefixSendsStoredMode(t *testing.T) {
	var forms []url.Values
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			forms = append(forms, r.Form)
			if len(forms) == 1 {
				writeVpcAssociationJSON(t, w, http.StatusBadRequest, map[string]string{
					"Code": "IncorrectStatus.EcrAssociation", "Message": "association is busy", "RequestId": "unit-test",
				})
				return
			}
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	config := vpcAssociationUnitConfig()
	config["allowed_prefixes"] = []interface{}{"172.16.2.0/24"}
	config["allowed_prefixes_mode"] = "IncrementalMode"
	state := vpcAssociationUnitState(map[string]string{
		"allowed_prefixes.#":    "1",
		"allowed_prefixes.0":    "172.16.1.0/24",
		"allowed_prefixes_mode": "IncrementalMode",
	})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if diags.HasError() {
		t.Fatalf("prefix update with stored mode: %v", diags)
	}
	if len(forms) != 2 {
		t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 2", len(forms))
	}
	for i, form := range forms {
		if got := form.Get("AllowedPrefixesMode"); got != "IncrementalMode" {
			t.Errorf("request %d AllowedPrefixesMode = %q, want IncrementalMode", i+1, got)
		}
		assertFormArray(t, form, "AllowedPrefixes", "172.16.2.0/24")
	}
}

// SDKv2 treats an explicit empty Optional+Computed value as unset during Diff,
// so the Update retains and sends the stored mode with the changed prefixes.
func TestUnitVpcAssociationUpdateExplicitEmptyModeSendsStoredMode(t *testing.T) {
	var forms []url.Values
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		switch r.Form.Get("Action") {
		case "ModifyExpressConnectRouterAssociationAllowedPrefix":
			forms = append(forms, r.Form)
			writeVpcAssociationJSON(t, w, http.StatusOK, map[string]string{"RequestId": "unit-test"})
		case "DescribeExpressConnectRouterAssociation":
			writeVpcAssociationActive(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	config := vpcAssociationUnitConfig()
	config["allowed_prefixes"] = []interface{}{"172.16.2.0/24"}
	config["allowed_prefixes_mode"] = ""
	state := vpcAssociationUnitState(map[string]string{
		"allowed_prefixes.#":    "1",
		"allowed_prefixes.0":    "172.16.1.0/24",
		"allowed_prefixes_mode": "IncrementalMode",
	})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if diags.HasError() {
		t.Fatalf("explicit empty mode update: %v", diags)
	}
	if len(forms) != 1 {
		t.Fatalf("ModifyExpressConnectRouterAssociationAllowedPrefix calls = %d, want 1", len(forms))
	}
	assertFormArray(t, forms[0], "AllowedPrefixes", "172.16.2.0/24")
	values, present := forms[0]["AllowedPrefixesMode"]
	if !present || len(values) != 1 || values[0] != "IncrementalMode" {
		t.Fatalf("AllowedPrefixesMode present=%v values=%v, want the stored IncrementalMode", present, values)
	}
}

func TestUnitVpcAssociationUpdateRejectsModeOnlyChange(t *testing.T) {
	var requests int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		t.Errorf("unexpected request for a mode-only update")
	})

	config := vpcAssociationUnitConfig()
	config["allowed_prefixes_mode"] = "IncrementalMode"
	state := vpcAssociationUnitState(map[string]string{"allowed_prefixes_mode": "MatchMode"})
	_, diags := vpcAssociationApply(t, expressconnectrouter.ResourceVpcAssociation(), state, config, client)
	if !diags.HasError() || !containsAll(fmt.Sprintf("%v", diags), "allowed_prefixes_mode", "non-empty", "allowed_prefixes") {
		t.Fatalf("mode-only update error = %v, want the local combination error", diags)
	}
	if got := atomic.LoadInt32(&requests); got != 0 {
		t.Fatalf("invalid mode-only update issued %d requests, want 0", got)
	}
}

func TestUnitVpcAssociationPrefixModeConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name      string
		mode      cty.Value
		prefixes  cty.Value
		config    map[string]interface{}
		wantError []string
	}{
		{
			name:     "mode_with_prefix_is_valid",
			mode:     cty.StringVal("MatchMode"),
			prefixes: cty.ListVal([]cty.Value{cty.StringVal("172.16.1.0/24")}),
			config: map[string]interface{}{
				"allowed_prefixes_mode": "MatchMode",
				"allowed_prefixes":      []interface{}{"172.16.1.0/24"},
			},
		},
		{
			name:     "mode_with_empty_prefix_is_rejected",
			mode:     cty.StringVal("MatchMode"),
			prefixes: cty.ListValEmpty(cty.String),
			config: map[string]interface{}{
				"allowed_prefixes_mode": "MatchMode",
				"allowed_prefixes":      []interface{}{},
			},
			wantError: []string{"allowed_prefixes", "at least one prefix"},
		},
		{
			name:     "null_mode_with_empty_prefix_is_rejected",
			mode:     cty.NullVal(cty.String),
			prefixes: cty.ListValEmpty(cty.String),
			config: map[string]interface{}{
				"allowed_prefixes": []interface{}{},
			},
			wantError: []string{"allowed_prefixes", "at least one prefix"},
		},
		{
			name:     "empty_mode_with_empty_prefix_is_rejected",
			mode:     cty.StringVal(""),
			prefixes: cty.ListValEmpty(cty.String),
			config: map[string]interface{}{
				"allowed_prefixes_mode": "",
				"allowed_prefixes":      []interface{}{},
			},
			wantError: []string{"allowed_prefixes", "at least one prefix"},
		},
		{
			name:     "null_prefix_without_mode_is_valid",
			mode:     cty.NullVal(cty.String),
			prefixes: cty.NullVal(cty.List(cty.String)),
			config:   map[string]interface{}{},
		},
		{
			name:     "unknown_prefix_is_skipped",
			mode:     cty.StringVal("IncrementalMode"),
			prefixes: cty.UnknownVal(cty.List(cty.String)),
			config: map[string]interface{}{
				"allowed_prefixes_mode": "IncrementalMode",
				"allowed_prefixes":      "74D93920-ED26-11E3-AC10-0800200C9A66",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := vpcAssociationUnitConfig()
			for key, value := range tc.config {
				config[key] = value
			}
			state := vpcAssociationUnitState(nil)
			state.RawConfig = vpcAssociationRawConfigWithValues(config, tc.mode, tc.prefixes)
			_, err := expressconnectrouter.ResourceVpcAssociation().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if len(tc.wantError) > 0 && (err == nil || !containsAll(err.Error(), tc.wantError...)) {
				t.Fatalf("validation error = %v, want substrings %v", err, tc.wantError)
			}
			if len(tc.wantError) == 0 && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestUnitDescribeVpcAssociationTargetIdentity(t *testing.T) {
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		if r.Form.Get("Action") != "DescribeExpressConnectRouterAssociation" {
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
		writeVpcAssociationJSON(t, w, http.StatusOK, map[string]interface{}{
			"RequestId": "unit-test",
			"AssociationList": []map[string]interface{}{
				{
					"AssociationId":       "assoc-a",
					"AssociationNodeType": "VPC",
					"EcrId":               "ecr-1",
					"VpcId":               "vpc-a",
					"Status":              "ACTIVE",
					"Description":         "same-description",
				},
				{
					"AssociationId":       "assoc-b",
					"AssociationNodeType": "VPC",
					"EcrId":               "ecr-1",
					"VpcId":               "vpc-b",
					"Status":              "ACTIVE",
					"Description":         "same-description",
				},
			},
		})
	})
	service := expressconnectrouter.NewService(client)

	object, err := service.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation("ecr-1:assoc-b:vpc-b")
	if err != nil {
		t.Fatalf("describe association B: %v", err)
	}
	if got := object["AssociationId"]; got != "assoc-b" {
		t.Fatalf("describe B returned association %v, want assoc-b", got)
	}

	if _, err := service.DescribeExpressConnectRouterExpressConnectRouterVpcAssociation("ecr-1:assoc-c:vpc-c"); !errs.NotFoundError(err) {
		t.Fatalf("missing association C error = %v, want NotFound even though A exists", err)
	}
}

func formHasParameter(form url.Values, key string) bool {
	for formKey := range form {
		if formKey == key || strings.HasPrefix(formKey, key+".") {
			return true
		}
	}
	return false
}

func assertFormArray(t *testing.T, form url.Values, key string, want ...string) {
	t.Helper()
	got := make(url.Values)
	prefix := key + "."
	for formKey, values := range form {
		if formKey == key || strings.HasPrefix(formKey, prefix) {
			got[formKey] = append([]string(nil), values...)
		}
	}
	expected := make(url.Values, len(want))
	for i, value := range want {
		expected[fmt.Sprintf("%s.%d", key, i+1)] = []string{value}
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("%s form values = %v, want exactly %v", key, got, expected)
	}
}

func containsAll(value string, substrings ...string) bool {
	for _, substring := range substrings {
		if !strings.Contains(value, substring) {
			return false
		}
	}
	return true
}
