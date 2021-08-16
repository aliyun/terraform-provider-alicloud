package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_edas_k8s_application", &resource.Sweeper{
		Name: "alicloud_edas_k8s_application",
		F:    testSweepEdasK8sApplication,
	})
}

func testSweepEdasK8sApplication(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	applicationListRq := edas.CreateListApplicationRequest()
	applicationListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(applicationListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas k8s application in service list: %s", err)
		return nil
	}

	listApplicationResponse, _ := raw.(*edas.ListApplicationResponse)
	if listApplicationResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas k8s application in service list: %s", listApplicationResponse.Message)
		return nil
	}

	for _, v := range listApplicationResponse.ApplicationList.Application {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas application: %s", name)
			continue
		}
		log.Printf("[INFO] delete edas application: %s", name)
		// stop it before delete
		stopAppRequest := edas.CreateStopApplicationRequest()
		stopAppRequest.RegionId = region
		stopAppRequest.AppId = v.AppId

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.StopApplication(stopAppRequest)
		})
		if err != nil {
			return err
		}
		addDebug(stopAppRequest.GetActionName(), raw, stopAppRequest.RoaRequest, stopAppRequest)
		stopAppResponse, _ := raw.(*edas.StopApplicationResponse)
		changeOrderId := stopAppResponse.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, 5*time.Minute, 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return err
			}
		}

		deleteApplicationRequest := edas.CreateDeleteApplicationRequest()
		deleteApplicationRequest.RegionId = region
		deleteApplicationRequest.AppId = v.AppId

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteApplication(deleteApplicationRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteApplicationRequest.GetActionName(), raw, deleteApplicationRequest.RoaRequest, deleteApplicationRequest)
			rsp := raw.(*edas.DeleteApplicationResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = Error("Operation cannot be processed because there are running instances.")
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] DeleteApplication got an error: %s", err)
			return nil
		}
	}

	return nil
}

func TestAccAlicloudEdasK8sApplication_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_k8s_application.default"
	ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sappb%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
	region := os.Getenv("ALICLOUD_REGION")
	image := fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region)
	updateImg := fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/provider:1.0", region)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":      "${var.name}",
					"cluster_id":            "${alicloud_edas_k8s_cluster.default.id}",
					"package_type":          "Image",
					"image_url":             image,
					"replicas":              "2",
					"internet_slb_protocol": "TCP",
					"internet_slb_port":     "8080",
					"internet_target_port":  "18082",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_url":             image,
						"replicas":              "2",
						"internet_slb_protocol": "TCP",
						"internet_slb_port":     "8080",
						"internet_target_port":  "18082",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"internet_slb_protocol", "internet_slb_port", "internet_target_port", "package_version"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_url": updateImg,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_url": updateImg,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":      "/bin/sh",
					"command_args": []string{"-c", "sleep 1000"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command":        "/bin/sh",
						"command_args.#": "2",
						"command_args.0": "-c",
						"command_args.1": "sleep 1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"envs": map[string]string{"a": "b"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"envs.%": "1",
						"envs.a": "b",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"limit_m_cpu":    "500",
					"limit_mem":      "1000",
					"requests_m_cpu": "100",
					"requests_mem":   "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"limit_m_cpu":    "500",
						"limit_mem":      "1000",
						"requests_m_cpu": "100",
						"requests_mem":   "100",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"replicas":       "3",
					"image_url":      updateImg,
					"command":        "/bin/sh",
					"command_args":   []string{"-c", "sleep 1001"},
					"envs":           map[string]string{"a": "c"},
					"limit_m_cpu":    "501",
					"limit_mem":      "1001",
					"requests_m_cpu": "101",
					"requests_mem":   "101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_url":      updateImg,
						"replicas":       "3",
						"command":        "/bin/sh",
						"command_args.#": "2",
						"command_args.0": "-c",
						"command_args.1": "sleep 1001",
						"envs.%":         "1",
						"envs.a":         "c",
						"limit_m_cpu":    "501",
						"limit_mem":      "1001",
						"requests_m_cpu": "101",
						"requests_mem":   "101",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudEdasK8sApplicationJar_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_k8s_application.default"
	ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sappb%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
	packageUrl := "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/SPRING_CLOUD_PROVIDER.jar"
	updateUrl := "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/DUBBO_PROVIDER.jar"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":     "${var.name}",
					"cluster_id":           "${alicloud_edas_k8s_cluster.default.id}",
					"package_type":         "FatJar",
					"package_url":          packageUrl,
					"jdk":                  "Open JDK 8",
					"replicas":             "1",
					"readiness":            `{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"host\":\"\", \"port\":18081}}`,
					"liveness":             `{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"host\":\"\", \"port\":18081}}`,
					"empty_dir":            `[{\"name\":\"emptydirvolume-0\",\"mountPath\":\"/home/emptydir/test\"}]`,
					"local_volume":         `[{\"mountPath\":\"/home/hostpath/tmp\",\"name\":\"localvolume-0\",\"nodePath\":\"/tmp\",\"opsAuth\":0,\"type\":\"Directory\"}]`,
					"deploy_across_zones":  "true",
					"deploy_across_nodes":  "true",
					"java_start_up_config": `{\"InitialHeapSize\":{\"original\":512,\"startup\":\"-Xms512m\"},\"MaxHeapSize\":{\"original\":512,\"startup\":\"-Xmx512m\"},\"CustomParams\":{\"original\":\"-Dtestkey=testval\",\"startup\":\"-Dtestkey=testval\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type":         "FatJar",
						"package_url":          packageUrl,
						"replicas":             "1",
						"jdk":                  "Open JDK 8",
						"readiness":            CHECKSET,
						"liveness":             CHECKSET,
						"empty_dir":            CHECKSET,
						"local_volume":         CHECKSET,
						"deploy_across_zones":  "true",
						"deploy_across_nodes":  "true",
						"java_start_up_config": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"empty_dir":            `[{\"name\":\"emptydirvolume-0\",\"mountPath\":\"/home/emptydir/test2\"}]`,
					"local_volume":         `[{\"mountPath\":\"/home/hostpath/tmp2\",\"name\":\"localvolume-0\",\"nodePath\":\"/tmp\",\"opsAuth\":0,\"type\":\"Directory\"}]`,
					"deploy_across_zones":  "false",
					"deploy_across_nodes":  "false",
					"custom_tolerations":   `[{\"key\":\"toleration-test2\",\"operator\":\"Equal\",\"value\":\"test2\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300},{\"key\":\"toleration-test\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"}]`,
					"custom_affinity":      `{\"nodeAffinity\":{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchExpressions\":[{\"key\":\"beta.kubernetes.io/arch\",\"operator\":\"NotIn\",\"values\":[\"arm64\",\"arm\"]}]}]},\"preferredDuringSchedulingIgnoredDuringExecution\":[{\"weight\":1,\"preference\":{\"matchExpressions\":[{\"key\":\"beta.kubernetes.io/arch\",\"operator\":\"In\",\"values\":[\"amd64\"]}]}}]},\"podAffinity\":{\"requiredDuringSchedulingIgnoredDuringExecution\":[{\"labelSelector\":{\"matchExpressions\":[{\"key\":\"edas.oam.acname\",\"operator\":\"NotIn\",\"values\":[\"test\"]}]},\"topologyKey\":\"kubernetes.io/hostname\"}]},\"podAntiAffinity\":{\"preferredDuringSchedulingIgnoredDuringExecution\":[{\"weight\":1,\"podAffinityTerm\":{\"labelSelector\":{\"matchExpressions\":[{\"key\":\"edas.appid\",\"operator\":\"DoesNotExist\"}]},\"namespaces\":[\"default\"],\"topologyKey\":\"failure-domain.beta.kubernetes.io/region\"}}]}}`,
					"java_start_up_config": `{\"InitialHeapSize\":{\"original\":512,\"startup\":\"-Xms512m\"},\"MaxHeapSize\":{\"original\":512,\"startup\":\"-Xmx512m\"},\"CustomParams\":{\"original\":\"-Dtestkey2=testval2\",\"startup\":\"-Dtestkey2=testval2\"}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"empty_dir":            CHECKSET,
						"local_volume":         CHECKSET,
						"deploy_across_zones":  "false",
						"deploy_across_nodes":  "false",
						"custom_tolerations":   CHECKSET,
						"custom_affinity":      CHECKSET,
						"java_start_up_config": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"readiness":           "{}",
					"liveness":            "{}",
					"empty_dir":           "[]",
					"local_volume":        "[]",
					"deploy_across_zones": "true",
					"deploy_across_nodes": "true",
					"custom_tolerations":  "[]",
					"custom_affinity":     "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readiness":           "{}",
						"liveness":            "{}",
						"empty_dir":           "[]",
						"local_volume":        "[]",
						"deploy_across_zones": "true",
						"deploy_across_nodes": "true",
						"custom_tolerations":  "[]",
						"custom_affinity":     "{}",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"package_url", "package_version"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"jdk":         "Dragonwell JDK 8",
					"replicas":    "2",
					"package_url": updateUrl,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"jdk":         "Dragonwell JDK 8",
						"replicas":    "2",
						"package_url": updateUrl,
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

func TestAccAlicloudEdasK8sApplication_multi(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alicloud_edas_k8s_application.default.1"
	ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(100, 999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sappm%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
	region := os.Getenv("ALICLOUD_REGION")
	image := fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"application_name": "${var.name}-${count.index}",
					"cluster_id":       "${alicloud_edas_k8s_cluster.default.id}",
					"replicas":         "1",
					"package_type":     "Image",
					"image_url":        image,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasK8sApplicationBasicMap = map[string]string{
	"application_name": CHECKSET,
	"cluster_id":       CHECKSET,
}

func testAccCheckEdasK8sApplicationDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasK8sApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alicloud_zones" default {
		  available_resource_creation = "VSwitch"
		}
		
		data "alicloud_instance_types" "default" {
		  availability_zone = data.alicloud_zones.default.zones.0.id
		  cpu_core_count = 2
		  memory_size = 4
		  kubernetes_node_role = "Worker"
		}
		
		resource "alicloud_vpc" "default" {
		  name = var.name
		  cidr_block = "10.1.0.0/21"
		}
		
		resource "alicloud_vswitch" "default" {
		  name = var.name
		  vpc_id = alicloud_vpc.default.id
		  cidr_block = "10.1.1.0/24"
		  availability_zone = data.alicloud_zones.default.zones.0.id
		}
		
		resource "alicloud_cs_managed_kubernetes" "default" {
		  worker_instance_types = [data.alicloud_instance_types.default.instance_types.0.id]
		  name = var.name
		  worker_vswitch_ids = [alicloud_vswitch.default.id]
		  worker_number = 				"2"
		  password =                    "Test12345"
		  pod_cidr =                   	"172.20.0.0/16"
		  service_cidr =               	"172.21.0.0/20"
		  worker_disk_size =            "50"
		  worker_disk_category =        "cloud_ssd"
		  worker_data_disk_size =       "20"
		  worker_data_disk_category =   "cloud_ssd"
		  worker_instance_charge_type = "PostPaid"
		  slb_internet_enabled =        "true"
          cluster_spec =                "ack.pro.small"
		}
		
		resource "alicloud_edas_k8s_cluster" "default" {
		  cs_cluster_id = alicloud_cs_managed_kubernetes.default.id
		}
		`, name)
}
