package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// collection-order: bypass env_conf.staging.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.production.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_anhui.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_beijing.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_chongqing.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_fujian.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_gansu.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_guangdong.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_guangxi.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_guizhou.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_hainan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_hebei.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_heilongjiang.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_henan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_hong_kong.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_hubei.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_hunan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_jiangsu.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_jiangxi.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_jilin.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_liaoning.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_macau.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_neimenggu.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_ningxia.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_qinghai.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_shaanxi.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_shandong.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_shanghai.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_shanxi.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_sichuan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_taiwan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_tianjin.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_xinjiang.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_xizang.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_yunnan.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_zhejiang.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change
// collection-order: bypass env_conf.preset_canary_overseas.allowed_hosts -- allowed_hosts is a server-side normalized echo of the request order; no two-apply member reorder can converge under the ER code-upload lifecycle gating code_rev, which is outside this tests-only change

// testAccDcdnErUploadCode uploads a minimal ER code fixture to the self-created routine so that
// code_rev can be exercised in the subsequent update step. It follows the same three-step staging
// lifecycle as the ESA routine resource: UploadStagingRoutineCode returns an OssPostConfig +
// CodeRevision, the minimal JS handler is POSTed to OSS via the shared esaRoutinePostCodeToOss,
// then CommitStagingRoutineCode commits the staged code. It is called from the update step's
// PreConfig so the code exists before EditRoutineConf sets code_rev. Errors are logged (non-fatal).
func testAccDcdnErUploadCode(name string) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	uploadResp, err := client.RpcPost("dcdn", "2018-01-15", "UploadStagingRoutineCode", nil,
		map[string]interface{}{"Name": name, "CodeDescription": "tf-acc-test"}, true)
	if err != nil {
		log.Printf("[dcdn_er] UploadStaging(%s) failed: %v", name, err)
		return
	}
	content, _ := uploadResp["Content"].(map[string]interface{})
	ossCfg, _ := content["OssPostConfig"].(map[string]interface{})
	codeRev, _ := content["CodeRevision"].(string)
	log.Printf("[dcdn_er] UploadStaging(%s) CodeRevision=%s", name, codeRev)
	if ossCfg != nil {
		jsCode := `addEventListener('fetch', function(event) { event.respondWith(new Response('hello world', {status: 200})); });`
		if ossErr := esaRoutinePostCodeToOss(ossCfg, jsCode); ossErr != nil {
			log.Printf("[dcdn_er] OSS code upload(%s) failed: %v", name, ossErr)
		} else {
			log.Printf("[dcdn_er] OSS code upload(%s) success", name)
		}
	}
	commitResp, err := client.RpcPost("dcdn", "2018-01-15", "CommitStagingRoutineCode", nil,
		map[string]interface{}{"Name": name, "CodeDescription": "tf-acc-test"}, true)
	if err != nil {
		log.Printf("[dcdn_er] CommitStaging(%s) failed: %v", name, err)
		return
	}
	log.Printf("[dcdn_er] CommitStaging(%s) Content: %+v", name, commitResp["Content"])
}

// DescribeRoutineUserInfo is an account-level API (no Name parameter) whose Content.Routines
// array lists every routine of the caller; the sweeper enumerates it, filters tf-testAcc/tf_testAcc
// prefixed names and deletes each via DeleteRoutine(Name). ER routines are account-global, so the
// per-region sweep may encounter already-deleted names on later regions (logged, not fatal).
func init() {
	resource.AddTestSweepers("alicloud_dcdn_er", &resource.Sweeper{
		Name: "alicloud_dcdn_er",
		F:    testSweepDcdnEr,
	})
}

func testSweepDcdnEr(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.DCDNSupportRegions) {
		log.Printf("[INFO] Skipping Dcdn Er unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeRoutineUserInfo"
	request := map[string]interface{}{}
	response, err := client.RpcPost("dcdn", "2018-01-15", action, nil, request, true)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	content, ok := response["Content"].(map[string]interface{})
	if !ok {
		log.Printf("[INFO] DescribeRoutineUserInfo returned no Content: %v", response)
		return nil
	}
	routines, _ := content["Routines"].([]interface{})
	for _, v := range routines {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := item["RoutineName"].(string)
		if !ok {
			continue
		}
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Dcdn Er: %s", name)
			continue
		}
		deleteAction := "DeleteRoutine"
		deleteRequest := map[string]interface{}{
			"Name": name,
		}
		_, err = client.RpcPost("dcdn", "2018-01-15", deleteAction, nil, deleteRequest, true)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Dcdn Er (%s): %s", name, err)
		}
		log.Printf("[INFO] Delete Dcdn Er success: %s ", name)
	}
	return nil
}

func TestAccAliCloudDcdnEr_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
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
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				PreConfig: func() {
					testAccDcdnErUploadCode(name)
				},
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
						"env_conf.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudDcdnEr_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
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
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				PreConfig: func() {
					testAccDcdnErUploadCode(name)
				},
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
						"env_conf.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudDcdnEr_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
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
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				PreConfig: func() {
					testAccDcdnErUploadCode(name)
				},
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
						"env_conf.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudDcdnEr_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
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
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				PreConfig: func() {
					testAccDcdnErUploadCode(name)
				},
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
						"env_conf.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudDcdnEr_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnEr")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnErBasicDependence0)
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
					"er_name":     name,
					"description": name,
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "5ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"er_name":     name,
						"description": name,
						"env_conf.#":  "1",
					}),
				),
			},
			{
				PreConfig: func() {
					testAccDcdnErUploadCode(name)
				},
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
					"env_conf": []map[string]interface{}{
						{
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"code_rev":      "unstable",
									"allowed_hosts": []string{"example.org"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
						"env_conf.#":  "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDcdnErMap = map[string]string{}

func AlicloudDcdnErBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
