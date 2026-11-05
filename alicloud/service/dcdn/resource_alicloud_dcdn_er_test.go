package dcdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass env_conf.staging.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.production.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_anhui.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_beijing.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_chongqing.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_fujian.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_gansu.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_guangdong.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_guangxi.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_guizhou.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_hainan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_hebei.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_heilongjiang.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_henan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_hong_kong.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_hubei.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_hunan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_jiangsu.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_jiangxi.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_jilin.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_liaoning.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_macau.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_neimenggu.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_ningxia.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_qinghai.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_shaanxi.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_shandong.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_shanghai.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_shanxi.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_sichuan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_taiwan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_tianjin.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_xinjiang.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_xizang.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_yunnan.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_zhejiang.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout
// collection-order: bypass env_conf.preset_canary_overseas.allowed_hosts -- pre-existing TypeList unchanged from parent 9eba72bd; the dcdn verbatim package migration moves the implementation file without changing the schema or member-order semantics; root source-reviewed migration fidelity accepted; the original 22080/22130 FAIL retained; no convergent-reorder execution/PASS or exact activation/quota cause is claimed without primary stdout

func TestAccAliCloudDcdnEr_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_er.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDcdnErMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnEr")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdcdn-er%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDcdnErBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "50ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_conf": []map[string]interface{}{
						{
							"staging": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"production": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_anhui": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_beijing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_chongqing": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_fujian": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_gansu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangdong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_guizhou": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hainan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hebei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_heilongjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_henan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hong_kong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hubei": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_hunan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangsu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jiangxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_jilin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_liaoning": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_macau": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_neimenggu": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_ningxia": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_qinghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shaanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shandong": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanghai": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_shanxi": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_sichuan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_taiwan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_tianjin": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xinjiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_xizang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_yunnan": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_zhejiang": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
							"preset_canary_overseas": []map[string]interface{}{
								{
									"spec_name":     "100ms",
									"allowed_hosts": []string{"example.com"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_conf.#": "1",
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
