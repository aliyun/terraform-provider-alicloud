package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_vpc_peer_connection",
		&resource.Sweeper{
			Name: "alicloud_vpc_peer_connection",
			F:    testSweepVpcPeerConnection,
		})
}

func testSweepVpcPeerConnection(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListVpcPeerConnections"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := aliyunClient.NewVpcpeerClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.VpcPeerConnects", response)

		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.VpcPeerConnects", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpc Peer Connection: %s", item["Name"].(string))
				continue
			}
			action := "DeleteVpcPeerConnection"
			request := map[string]interface{}{
				"InstanceId": item["InstanceId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpc Peer Connection (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Vpc Peer Connection success: %s ", item["Name"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudVPCPeerConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_peer_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudVPCPeerConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPeerConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPCPeerConnectionBasicDependence0)
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
					"vpc_id":              "${alicloud_vpc.requesting.id}",
					"accepting_vpc_id":    "${alicloud_vpc.accepting.id}",
					"accepting_region_id": "${data.alicloud_regions.default.regions.0.id}",
					"accepting_ali_uid":   "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":              CHECKSET,
						"accepting_vpc_id":    CHECKSET,
						"accepting_region_id": CHECKSET,
						"accepting_ali_uid":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_connection_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Activated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "PeerConnection",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "PeerConnection",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudVPCPeerConnection_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_peer_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudVPCPeerConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcPeerConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPCPeerConnectionBasicDependence0)
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
					"vpc_id":               "${alicloud_vpc.requesting.id}",
					"accepting_vpc_id":     "${alicloud_vpc.accepting.id}",
					"accepting_region_id":  "${data.alicloud_regions.default.regions.0.id}",
					"accepting_ali_uid":    "${data.alicloud_account.default.id}",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"peer_connection_name": name,
					"description":          name,
					"status":               "Activated",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "PeerConnection",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"accepting_vpc_id":     CHECKSET,
						"accepting_region_id":  CHECKSET,
						"accepting_ali_uid":    CHECKSET,
						"resource_group_id":    CHECKSET,
						"peer_connection_name": name,
						"description":          name,
						"status":               "Activated",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "PeerConnection",
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

func TestAccAliCloudVPCPeerConnection_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_peer_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudVPCPeerConnectionMap0)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPCPeerConnectionBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVPCPeerConnectionDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":              "${alicloud_vpc.requesting.id}",
					"accepting_vpc_id":    "${alicloud_vpc.accepting.id}",
					"accepting_region_id": "${data.alicloud_regions.default.regions.0.id}",
					"accepting_ali_uid":   "${data.alicloud_account.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"vpc_id":              CHECKSET,
						"accepting_vpc_id":    CHECKSET,
						"accepting_region_id": CHECKSET,
						"accepting_ali_uid":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"bandwidth": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_connection_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"peer_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Activated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Rejected",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
			//		testAccCheck(map[string]string{
			//			"status": "Rejected",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "PeerConnection",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "PeerConnection",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVPCPeerConnection_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_peer_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudVPCPeerConnectionMap0)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcpeerconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudVPCPeerConnectionBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckVPCPeerConnectionDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":               "${alicloud_vpc.requesting.id}",
					"accepting_vpc_id":     "${alicloud_vpc.accepting.id}",
					"accepting_region_id":  "${data.alicloud_regions.default.regions.0.id}",
					"accepting_ali_uid":    "${data.alicloud_account.default.id}",
					"bandwidth":            "200",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"peer_connection_name": name,
					"description":          name,
					"status":               "Activated",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "PeerConnection",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCPeerConnectionExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"accepting_vpc_id":     CHECKSET,
						"accepting_region_id":  CHECKSET,
						"accepting_ali_uid":    CHECKSET,
						"bandwidth":            "200",
						"resource_group_id":    CHECKSET,
						"peer_connection_name": name,
						"description":          name,
						"status":               "Activated",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "PeerConnection",
					}),
				),
			},
		},
	})
}

var AliCloudVPCPeerConnectionMap0 = map[string]string{
	"bandwidth":         CHECKSET,
	"resource_group_id": CHECKSET,
	"status":            CHECKSET,
	"create_time":       CHECKSET,
}

func AliCloudVPCPeerConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_regions" "default" {
  		current = true
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "requesting" {
  		vpc_name    = var.name
  		enable_ipv6 = "true"
	}

	resource "alicloud_vpc" "accepting" {
  		vpc_name    = var.name
  		enable_ipv6 = "true"
	}
`, name)
}

func AliCloudVPCPeerConnectionBasicDependence1(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	provider "alicloud" {
  		alias  = "requesting"
  		region = "%s"
	}

	provider "alicloud" {
  		alias  = "accepting"
  		region = "cn-hangzhou"
	}

	data "alicloud_regions" "default" {
  		provider = alicloud.accepting
  		current  = true
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "requesting" {
  		provider    = alicloud.requesting
  		vpc_name    = var.name
  		enable_ipv6 = "true"
	}

	resource "alicloud_vpc" "accepting" {
  		provider    = alicloud.accepting
  		vpc_name    = var.name
  		enable_ipv6 = "true"
	}
`, name, defaultRegionToTest)
}

func testAccCheckVPCPeerConnectionDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			if err := testAccCheckVPCPeerConnectionDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckVPCPeerConnectionDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	vpcServiceV2 := VpcServiceV2{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc_peer_connection" {
			continue
		}

		_, err := vpcServiceV2.DescribeVpcPeerConnection(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func testAccCheckVPCPeerConnectionExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_vpc_peer_connection id is set")
		}

		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			vpcServiceV2 := VpcServiceV2{client}

			resp, err := vpcServiceV2.DescribeVpcPeerConnection(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}

			res = resp

			return nil
		}

		return fmt.Errorf("alicloud_vpc_peer_connection not found")
	}
}

func TestUnitAccAliCloudVpcPeerConnection(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"peer_connection_name": "CreateVpcPeerConnectionValue",
		"vpc_id":               "CreateVpcPeerConnectionValue",
		"accepting_ali_uid":    1,
		"accepting_region_id":  "CreateVpcPeerConnectionValue",
		"accepting_vpc_id":     "CreateVpcPeerConnectionValue",
		"description":          "CreateVpcPeerConnectionValue",
		"bandwidth":            100,
		"dry_run":              false,
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
		"InstanceId":        "CreateVpcPeerConnectionValue",
		"Name":              "CreateVpcPeerConnectionValue",
		"Description":       "CreateVpcPeerConnectionValue",
		"AcceptingOwnerUid": 1,
		"RegionId":          "CreateVpcPeerConnectionValue",
		"AcceptingRegionId": "CreateVpcPeerConnectionValue",
		"Bandwidth":         100,
		"Status":            "Activated",
		"BizStatus":         "Normal",
		"Vpc": map[string]interface{}{
			"VpcId": "CreateVpcPeerConnectionValue",
		},
		"AcceptingVpc": map[string]interface{}{
			"VpcId": "CreateVpcPeerConnectionValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		"InstanceId": "CreateVpcPeerConnectionValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_vpc_peer_connection", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcpeerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudVpcPeerConnectionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcPeerConnection" {
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
		err := resourceAliCloudVpcPeerConnectionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcpeerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudVpcPeerConnectionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"peer_connection_name": "UpdateVpcPeerConnectionValue",
		"description":          "UpdateVpcPeerConnectionValue",
		"bandwidth":            200,
	}
	diff, err := newInstanceDiff("alicloud_vpc_peer_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Name":        "UpdateVpcPeerConnectionValue",
		"Description": "UpdateVpcPeerConnectionValue",
		"Bandwidth":   200,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyVpcPeerConnection" {
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
		err := resourceAliCloudVpcPeerConnectionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_vpc_peer_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetVpcPeerConnectionAttribute" {
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
		err := resourceAliCloudVpcPeerConnectionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcpeerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudVpcPeerConnectionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_vpc_peer_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_vpc_peer_connection"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcPeerConnection" {
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
			if *action == "GetVpcPeerConnectionAttribute" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudVpcPeerConnectionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
