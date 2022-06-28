package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"time"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var gNamePrefix string

func init() {
	gNamePrefix = fmt.Sprintf("tf-testacc%salicloudfcasyncinvoke", defaultRegionToTest)
	resource.AddTestSweepers("alicloud_fc_function_async_invoke_config", &resource.Sweeper{
		Name: "alicloud_fc_function_async_invoke_config",
		F:    testSweepFCFunctionAsyncInvokeConfig,
		Dependencies: []string{
			"alicloud_fc_service",
			"alicloud_fc_function",
		},
	})
}

func testSweepFCFunctionAsyncInvokeConfig(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.FcNoSupportedRegions) {
		log.Printf("[INFO] Skipping Function Compute unsupported region: %s", region)
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

	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		return fcClient.ListServices(fc.NewListServicesInput())
	})
	if err != nil {
		return fmt.Errorf("Error retrieving FC services: %s", err)
	}
	services, _ := raw.(*fc.ListServicesOutput)
	for _, v := range services.Services {
		name := *v.ServiceName
		id := *v.ServiceID
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping FC services: %s (%s)", name, id)
			continue
		}
		// Remove functions
		nextToken := ""
		for {
			args := fc.NewListFunctionsInput(name)
			if nextToken != "" {
				args.NextToken = &nextToken
			}

			raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.ListFunctions(args)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to list functions of service (%s (%s)): %s", name, id, err)
			}
			resp, _ := raw.(*fc.ListFunctionsOutput)

			if resp.Functions == nil || len(resp.Functions) < 1 {
				break
			}

			for _, function := range resp.Functions {
				// Remove triggers
				triggerDeleted := false
				triggerNext := ""
				for {
					req := fc.NewListTriggersInput(name, *function.FunctionName)
					if triggerNext != "" {
						req.NextToken = &triggerNext
					}

					raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
						return fcClient.ListTriggers(req)
					})
					if err != nil {
						log.Printf("[ERROR] Failed to list triggers of functiion (%s): %s", name, err)
					}
					resp, _ := raw.(*fc.ListTriggersOutput)

					if resp == nil || resp.Triggers == nil || len(resp.Triggers) < 1 {
						break
					}
					for _, trigger := range resp.Triggers {
						triggerDeleted = true
						if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
							return fcClient.DeleteTrigger(&fc.DeleteTriggerInput{
								ServiceName:  StringPointer(name),
								FunctionName: function.FunctionName,
								TriggerName:  trigger.TriggerName,
							})
						}); err != nil {
							log.Printf("[ERROR] Failed to delete trigger %s of function: %s.", *trigger.TriggerName, *function.FunctionName)
						}
					}

					triggerNext = ""
					if resp.NextToken != nil {
						triggerNext = *resp.NextToken
					}
					if triggerNext == "" {
						break
					}
				}
				// remove function invoke async config
				if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
					return fcClient.DeleteFunctionAsyncInvokeConfig(&fc.DeleteFunctionAsyncInvokeConfigInput{
						ServiceName:  StringPointer(name),
						FunctionName: function.FunctionName,
					})
				}); err != nil {
					log.Printf("[ERROR] Failed to delete function invoke async config %s of services: %s (%s)", *function.FunctionName, name, id)
				}
				//remove function
				if triggerDeleted {
					time.Sleep(5 * time.Second)
				}
				if _, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
					return fcClient.DeleteFunction(&fc.DeleteFunctionInput{
						ServiceName:  StringPointer(name),
						FunctionName: function.FunctionName,
					})
				}); err != nil {
					log.Printf("[ERROR] Failed to delete function %s of services: %s (%s)", *function.FunctionName, name, id)
				}
			}

			nextToken = ""
			if resp.NextToken != nil {
				nextToken = *resp.NextToken
			}
			if nextToken == "" {
				break
			}
		}

		log.Printf("[INFO] Deleting FC services: %s (%s)", name, id)
		_, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			return fcClient.DeleteService(&fc.DeleteServiceInput{
				ServiceName: StringPointer(name),
			})
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete FC services (%s (%s)): %s", name, id, err)
		}
	}

	// Delete mns queue
	_, err = client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
		queues, err := queueManager.ListQueueDetail("", 1000, gNamePrefix)
		if err != nil {
			return nil, err
		}
		for _, q := range queues.Attrs {
			err = queueManager.DeleteQueue(q.QueueName)
			if err != nil {
				break
			}
		}
		return nil, err
	})

	// Delete mns topics
	_, err = client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
		topics, err := topicManager.ListTopicDetail("", 1000, gNamePrefix)
		if err != nil {
			return nil, err
		}
		for _, t := range topics.Attrs {
			err = topicManager.DeleteTopic(t.TopicName)
			if err != nil {
				break
			}
		}
		return nil, err
	})

	return nil
}

func TestAccAlicloudFCFunctionAsyncInvokeConfigUpdate(t *testing.T) {
	var v *fc.GetFunctionAsyncInvokeConfigOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("%s-%d", gNamePrefix, rand)
	var basicMap = map[string]string{
		"created_time":       CHECKSET,
		"last_modified_time": CHECKSET,
	}
	resourceId := "alicloud_fc_function_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcFunctionAsyncInvokeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":  "${alicloud_fc_service.default.name}",
					"function_name": "${alicloud_fc_function.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maximum_event_age_in_seconds": "200",
					"maximum_retry_attempts":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maximum_event_age_in_seconds": "200",
						"maximum_retry_attempts":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]string{
								{
									"destination": getMnsQueueArn(name),
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.0.on_failure.0.destination": getMnsQueueArn(name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]string{
								{
									"destination": getMnsTopicArn(name),
								},
							},
							"on_success": []map[string]string{
								{
									"destination": getMnsQueueArn(name),
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.0.on_failure.0.destination": getMnsTopicArn(name),
						"destination_config.0.on_success.0.destination": getMnsQueueArn(name),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudFCFunctionAsyncInvokeStatefulInvocationUpdate(t *testing.T) {
	var v *fc.GetFunctionAsyncInvokeConfigOutput
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("%s-%d", gNamePrefix, rand)
	var basicMap = map[string]string{
		"created_time":       CHECKSET,
		"last_modified_time": CHECKSET,
	}
	resourceId := "alicloud_fc_function_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &FcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFcFunctionAsyncInvokeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.FcNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":  "${alicloud_fc_service.default.name}",
					"function_name": "${alicloud_fc_function.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stateful_invocation": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stateful_invocation": "true",
					}),
				),
			},
		},
	})
}

func getMnsQueueArn(queueName string) string {
	region := os.Getenv("ALICLOUD_REGION")
	account := os.Getenv("ALICLOUD_ACCOUNT_ID")
	return fmt.Sprintf("acs:mns:%s:%s:/queues/%s/messages", region, account, queueName)
}

func getMnsTopicArn(topicName string) string {
	region := os.Getenv("ALICLOUD_REGION")
	account := os.Getenv("ALICLOUD_ACCOUNT_ID")
	return fmt.Sprintf("acs:mns:%s:%s:/topics/%s/messages", region, account, topicName)
}

func getFcFunctionArn(serviceName string, qualifier string, functionName string) string {
	region := os.Getenv("ALICLOUD_REGION")
	account := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if qualifier == "" {
		qualifier = "LATEST"
	}
	return fmt.Sprintf("acs:fc:%s:%s:services/%s.%s/functions/%s", region, account, serviceName, qualifier, functionName)
}

func resourceFcFunctionAsyncInvokeConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
resource "alicloud_fc_service" "default" {
  name = var.name
  role = alicloud_ram_role.default.arn
  depends_on = ["alicloud_ram_role_policy_attachment.default"]
}
resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}
resource "alicloud_oss_bucket_object" "default" {
  bucket = alicloud_oss_bucket.default.id
  key = "fc/hello.zip"
  content = <<EOF
    # -*- coding: utf-8 -*-
  def handler(event, context):
      print "hello world"
      return 'hello world'
  EOF
}
resource "alicloud_fc_function" "default" {
  service = alicloud_fc_service.default.name
  name = var.name
  oss_bucket = alicloud_oss_bucket.default.id
  oss_key = alicloud_oss_bucket_object.default.key
  memory_size = 512
  runtime = "python2.7"
  handler = "hello.handler"
  depends_on = ["alicloud_mns_queue.default", "alicloud_mns_topic.default"]
}
resource "alicloud_mns_queue" "default" {
	name = var.name
}
resource "alicloud_mns_topic" "default" {
	name = var.name
}
resource "alicloud_ram_role" "default" {
	name = var.name
	document = <<DEFINITION
	{
		"Statement": [
		  {
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {
			  "Service": [
				"fc.aliyuncs.com"
			  ]
			}
		  }
		],
		"Version": "1"
	}
	DEFINITION
	description = "this is a test"
	force = true
}
resource "alicloud_ram_policy" "default" {
	name = var.name
	document = <<DEFINITION
	{
		"Version": "1",
		"Statement": [
		  {
			"Action": "mns:*",
			"Resource": "*",
			"Effect": "Allow"
		  }
		]
	  }
	DEFINITION
}
resource "alicloud_ram_role_policy_attachment" "default" {
	role_name = alicloud_ram_role.default.name
	policy_name = alicloud_ram_policy.default.name
	policy_type = "Custom"
}
`, name)
}
