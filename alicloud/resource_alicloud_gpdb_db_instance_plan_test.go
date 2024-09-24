package alicloud

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestAccAliCloudGpdbDbInstancePlan_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 6).Format("2006-01-02T15:04:05Z")
	resumeExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	pauseExecuteTime := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	resumeExecuteTimeUpdate := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	pauseExecuteTimeUpdate := time.Now().AddDate(0, 0, 3).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Postpone",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"execute_time": resumeExecuteTime,
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"execute_time": pauseExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Postpone",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"execute_time": resumeExecuteTimeUpdate,
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"execute_time": pauseExecuteTimeUpdate,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	resumeExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	pauseExecuteTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Postpone",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"execute_time": resumeExecuteTime,
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"execute_time": pauseExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Postpone",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Regular",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 2 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Regular",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 1 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "PauseResume",
					"plan_schedule_type":    "Regular",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"resume": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"pause": []interface{}{
								map[string]interface{}{
									"plan_cron_time": "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Regular",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBServerlessSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 6).Format("2006-01-02T15:04:05Z")
	scaleInExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleOutExecuteTime := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleInExecuteTimeUpdate := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	scaleOutExecuteTimeUpdate := time.Now().AddDate(0, 0, 3).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence1)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "Resize",
					"plan_schedule_type":    "Postpone",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "2",
									"execute_time":     scaleInExecuteTime,
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"execute_time":     scaleOutExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "Resize",
						"plan_schedule_type":    "Postpone",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"execute_time":     scaleInExecuteTimeUpdate,
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "8",
									"execute_time":     scaleOutExecuteTimeUpdate,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBServerlessSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleInExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleOutExecuteTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence1)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "Resize",
					"plan_schedule_type":    "Postpone",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "2",
									"execute_time":     scaleInExecuteTime,
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"execute_time":     scaleOutExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "Resize",
						"plan_schedule_type":    "Postpone",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic3(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBServerlessSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence1)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "Resize",
					"plan_schedule_type":    "Regular",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "2",
									"plan_cron_time":   "0 0 0 1/1 * ? ",
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"plan_cron_time":   "0 0 2 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "Resize",
						"plan_schedule_type":    "Regular",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"plan_cron_time":   "0 0 1 1/1 * ? ",
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "8",
									"plan_cron_time":   "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic3_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBServerlessSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence1)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "Resize",
					"plan_schedule_type":    "Regular",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_in": []interface{}{
								map[string]interface{}{
									"segment_node_num": "2",
									"plan_cron_time":   "0 0 0 1/1 * ? ",
								},
							},
							"scale_out": []interface{}{
								map[string]interface{}{
									"segment_node_num": "6",
									"plan_cron_time":   "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "Resize",
						"plan_schedule_type":    "Regular",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic4(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 6).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTime := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTimeUpdate := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTimeUpdate := time.Now().AddDate(0, 0, 3).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "ModifySpec",
					"plan_schedule_type":    "Postpone",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec": "4C32G",
									"execute_time":  scaleUpExecuteTime,
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec": "2C16G",
									"execute_time":  scaleDownExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Postpone",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec": "8C64G",
									"execute_time":  scaleUpExecuteTimeUpdate,
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec": "4C32G",
									"execute_time":  scaleDownExecuteTimeUpdate,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic4_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "ModifySpec",
					"plan_schedule_type":    "Postpone",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec": "4C32G",
									"execute_time":  scaleUpExecuteTime,
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec": "2C16G",
									"execute_time":  scaleDownExecuteTime,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Postpone",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic5(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "ModifySpec",
					"plan_schedule_type":    "Regular",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec":  "4C32G",
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec":  "2C16G",
									"plan_cron_time": "0 0 2 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Regular",
						"plan_config.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_plan_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_plan_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_start_date": planStartDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_start_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_desc": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "cancel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "cancel",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec":  "8C64G",
									"plan_cron_time": "0 0 1 1/1 * ? ",
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec":  "4C32G",
									"plan_cron_time": "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_config.#": "1",
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

func TestAccAliCloudGpdbDbInstancePlan_basic5_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := resourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDbInstancePlan")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", defaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
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
					"db_instance_id":        "${alicloud_gpdb_instance.default.id}",
					"db_instance_plan_name": name,
					"plan_type":             "ModifySpec",
					"plan_schedule_type":    "Regular",
					"plan_start_date":       planStartDate,
					"plan_end_date":         planEndDate,
					"plan_desc":             name,
					"status":                "active",
					"plan_config": []interface{}{
						map[string]interface{}{
							"scale_up": []interface{}{
								map[string]interface{}{
									"instance_spec":  "4C32G",
									"plan_cron_time": "0 0 0 1/1 * ? ",
								},
							},
							"scale_down": []interface{}{
								map[string]interface{}{
									"instance_spec":  "2C16G",
									"plan_cron_time": "0 0 6 1/1 * ? ",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":        CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Regular",
						"plan_start_date":       CHECKSET,
						"plan_end_date":         CHECKSET,
						"plan_desc":             name,
						"status":                "active",
						"plan_config.#":         "1",
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

var AliCloudGpdbDbInstancePlanMap0 = map[string]string{
	"plan_id":         CHECKSET,
	"plan_start_date": CHECKSET,
	"status":          CHECKSET,
}

func AliCloudGpdbDbInstancePlanBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_gpdb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_gpdb_zones.default.ids.0
	}

	resource "alicloud_gpdb_instance" "default" {
  		db_instance_category  = "HighAvailability"
  		db_instance_class     = "gpdb.group.segsdx1"
  		db_instance_mode      = "StorageElastic"
  		description           = var.name
  		engine                = "gpdb"
  		engine_version        = "6.0"
  		zone_id               = data.alicloud_gpdb_zones.default.ids.0
  		instance_network_type = "VPC"
  		instance_spec         = "2C16G"
  		payment_type          = "PayAsYouGo"
  		seg_storage_type      = "cloud_essd"
  		seg_node_num          = 4
  		storage_size          = 50
  		vpc_id                = data.alicloud_vpcs.default.ids.0
  		vswitch_id            = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

func AliCloudGpdbDbInstancePlanBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "ap-southeast-1c"
	}

	resource "alicloud_gpdb_instance" "default" {
  		db_instance_mode      = "Serverless"
  		description           = var.name
		engine                = "gpdb"
  		engine_version        = "6.0"
  		zone_id               = "ap-southeast-1c"
  		instance_network_type = "VPC"
  		instance_spec         = "4C16G"
  		payment_type          = "PayAsYouGo"
  		seg_node_num          = 2
  		vpc_id                = data.alicloud_vpcs.default.ids.0
  		vswitch_id            = data.alicloud_vswitches.default.ids.0
	}
`, name)
}

func TestUnitAccAliCloudGpdbDbInstancePlan(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_instance_plan_name": "CreateGpdbDbInstancePlanValue",
		"plan_desc":             "CreateGpdbDbInstancePlanValue",
		"plan_type":             "CreateGpdbDbInstancePlanValue",
		"plan_schedule_type":    "CreateGpdbDbInstancePlanValue",
		"plan_start_date":       "CreateGpdbDbInstancePlanValue",
		"plan_end_date":         "CreateGpdbDbInstancePlanValue",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "CreateGpdbDbInstancePlanValue",
						"execute_time":   "CreateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "CreateGpdbDbInstancePlanValue",
						"execute_time":     "CreateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
		"db_instance_id": "CreateGpdbDbInstancePlanValue",
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
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"DBInstanceId":     "CreateGpdbDbInstancePlanValue",
					"PlanId":           "CreateGpdbDbInstancePlanValue",
					"PlanScheduleType": "CreateGpdbDbInstancePlanValue",
					"PlanType":         "CreateGpdbDbInstancePlanValue",
					"PlanDesc":         "CreateGpdbDbInstancePlanValue",
					"PlanName":         "CreateGpdbDbInstancePlanValue",
					"PlanStartDate":    "CreateGpdbDbInstancePlanValue",
					"PlanEndDate":      "CreateGpdbDbInstancePlanValue",
					"PlanConfig":       "{\"resume\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"CreateGpdbDbInstancePlanValue\",\"executeTime\":\"CreateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":       "active",
				},
			},
		},
		"Status": "success",
	}
	CreateMockResponse := map[string]interface{}{
		"PlanId": "CreateGpdbDbInstancePlanValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_gpdb_db_instance_plan", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGpdbDbInstancePlanCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstancePlan" {
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
		err := resourceAliCloudGpdbDbInstancePlanCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGpdbDbInstancePlanUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"db_instance_plan_name": "UpdateGpdbDbInstancePlanValue",
		"plan_desc":             "UpdateGpdbDbInstancePlanValue",
		"plan_start_date":       "UpdateGpdbDbInstancePlanValue",
		"plan_end_date":         "UpdateGpdbDbInstancePlanValue",
		"status":                "cancel",
		"plan_config": []interface{}{
			map[string]interface{}{
				"resume": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"pause": []interface{}{
					map[string]interface{}{
						"plan_cron_time": "UpdateGpdbDbInstancePlanValue",
						"execute_time":   "UpdateGpdbDbInstancePlanValue",
					},
				},
				"scale_in": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
				"scale_out": []interface{}{
					map[string]interface{}{
						"plan_cron_time":   "UpdateGpdbDbInstancePlanValue",
						"execute_time":     "UpdateGpdbDbInstancePlanValue",
						"segment_node_num": "2",
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"PlanList": []interface{}{
				map[string]interface{}{
					"PlanDesc":      "UpdateGpdbDbInstancePlanValue",
					"PlanName":      "UpdateGpdbDbInstancePlanValue",
					"PlanStartDate": "UpdateGpdbDbInstancePlanValue",
					"PlanEndDate":   "UpdateGpdbDbInstancePlanValue",
					"PlanConfig":    "{\"resume\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"pause\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\"},\"scaleOut\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"},\"scaleIn\":{\"planCronTime\":\"UpdateGpdbDbInstancePlanValue\",\"executeTime\":\"UpdateGpdbDbInstancePlanValue\",\"segmentNodeNum\":\"2\"}}",
					"PlanStatus":    "cancel",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDBInstancePlan" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "SetDBInstancePlanStatus" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstancePlanUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstancePlans" {
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
		err := resourceAliCloudGpdbDbInstancePlanRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGpdbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGpdbDbInstancePlanDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_gpdb_db_instance_plan", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_gpdb_db_instance_plan"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstancePlan" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGpdbDbInstancePlanDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
