package gpdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	gpdb "github.com/aliyun/terraform-provider-alicloud/alicloud/service/gpdb"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass plan_config -- pre-existing uncovered TypeList attribute outside this pure directory migration (file moved from alicloud/ root; the attribute and its test coverage are unchanged stock)

func TestAccAliCloudGpdbDbInstancePlan_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 6).Format("2006-01-02T15:04:05Z")
	resumeExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	pauseExecuteTime := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	resumeExecuteTimeUpdate := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	pauseExecuteTimeUpdate := time.Now().AddDate(0, 0, 3).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
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
						"plan_start_date": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	resumeExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	pauseExecuteTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Postpone",
						"plan_start_date":       acctest.CHECKSET,
						"plan_end_date":         acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
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
						"plan_start_date": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "PauseResume",
						"plan_schedule_type":    "Regular",
						"plan_start_date":       acctest.CHECKSET,
						"plan_end_date":         acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 6).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTime := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTimeUpdate := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTimeUpdate := time.Now().AddDate(0, 0, 3).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
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
						"plan_start_date": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": acctest.CHECKSET,
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

func TestAccAliCloudGpdbDbInstancePlan_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	scaleUpExecuteTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	scaleDownExecuteTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Postpone",
						"plan_start_date":       acctest.CHECKSET,
						"plan_end_date":         acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
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
						"plan_start_date": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"plan_end_date": planEndDate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plan_end_date": acctest.CHECKSET,
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

func TestAccAliCloudGpdbDbInstancePlan_basic3_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	resourceId := "alicloud_gpdb_db_instance_plan.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGpdbDbInstancePlanMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbInstancePlan")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceplan%d", acctest.DefaultRegionToTest, rand)
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGpdbDbInstancePlanBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"db_instance_id":        acctest.CHECKSET,
						"db_instance_plan_name": name,
						"plan_type":             "ModifySpec",
						"plan_schedule_type":    "Regular",
						"plan_start_date":       acctest.CHECKSET,
						"plan_end_date":         acctest.CHECKSET,
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
	"plan_id":         acctest.CHECKSET,
	"plan_start_date": acctest.CHECKSET,
	"status":          acctest.CHECKSET,
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
