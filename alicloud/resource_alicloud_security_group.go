// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSecurityGroupCreate,
		Read:   resourceAliCloudEcsSecurityGroupRead,
		Update: resourceAliCloudEcsSecurityGroupUpdate,
		Delete: resourceAliCloudEcsSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inner_access_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Accept", "Drop"}, false),
			},
			"permissions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "accept",
							ValidateFunc: StringInSlice([]string{"accept", "drop"}, false),
						},
						"source_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"priority": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "1",
							ValidateFunc: StringLenBetween(1, 100),
						},
						"ipv6_source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"nic_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "intranet",
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ip_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"TCP", "UDP", "ICMP", "ICMPv6", "GRE", "ALL"}, false),
						},
						"dest_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ipv6_dest_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_group_owner_account": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_prefix_list_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_managed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEcsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("security_group_type"); ok {
		request["SecurityGroupType"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("security_group_name"); ok {
		request["SecurityGroupName"] = v
	}
	if v, ok := d.GetOkExists("service_managed"); ok {
		request["ServiceManaged"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecurityGroupId"]))

	return resourceAliCloudEcsSecurityGroupUpdate(d, meta)
}

func resourceAliCloudEcsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsSecurityGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_security_group DescribeEcsSecurityGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("security_group_name", objectRaw["SecurityGroupName"])
	d.Set("security_group_type", objectRaw["SecurityGroupType"])
	d.Set("service_managed", objectRaw["ServiceManaged"])
	d.Set("vpc_id", objectRaw["VpcId"])
	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = ecsServiceV2.DescribeDescribeSecurityGroupAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("inner_access_policy", objectRaw["InnerAccessPolicy"])
	permission1Raw, _ := jsonpath.Get("$.Permissions.Permission", objectRaw)
	permissionsMaps := make([]map[string]interface{}, 0)
	if permission1Raw != nil {
		for _, permissionChild1Raw := range permission1Raw.([]interface{}) {
			permissionsMap := make(map[string]interface{})
			permissionChild1Raw := permissionChild1Raw.(map[string]interface{})
			permissionsMap["description"] = permissionChild1Raw["Description"]
			permissionsMap["dest_cidr_ip"] = permissionChild1Raw["DestCidrIp"]
			permissionsMap["ip_protocol"] = permissionChild1Raw["IpProtocol"]
			permissionsMap["ipv6_dest_cidr_ip"] = permissionChild1Raw["Ipv6DestCidrIp"]
			permissionsMap["ipv6_source_cidr_ip"] = permissionChild1Raw["Ipv6SourceCidrIp"]
			permissionsMap["nic_type"] = permissionChild1Raw["NicType"]
			permissionsMap["policy"] = permissionChild1Raw["Policy"]
			permissionsMap["port_range"] = permissionChild1Raw["PortRange"]
			permissionsMap["priority"] = permissionChild1Raw["Priority"]
			permissionsMap["source_cidr_ip"] = permissionChild1Raw["SourceCidrIp"]
			permissionsMap["source_group_id"] = permissionChild1Raw["SourceGroupId"]
			permissionsMap["source_group_owner_account"] = permissionChild1Raw["SourceGroupOwnerAccount"]
			permissionsMap["source_port_range"] = permissionChild1Raw["SourcePortRange"]
			permissionsMap["source_prefix_list_id"] = permissionChild1Raw["SourcePrefixListId"]
			permissionsMaps = append(permissionsMaps, permissionsMap)
		}
	}
	d.Set("permissions", permissionsMaps)

	objectRaw, err = ecsServiceV2.DescribeDescribeSecurityGroupReferences(d.Id())
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudEcsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifySecurityGroupAttribute"
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("security_group_name") {
		update = true
		request["SecurityGroupName"] = d.Get("security_group_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("security_group_name")
	}
	update = false
	action = "ModifySecurityGroupPolicy"
	conn, err = client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("inner_access_policy") {
		update = true
		request["InnerAccessPolicy"] = d.Get("inner_access_policy")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("inner_access_policy")
	}

	if d.HasChange("tags") {
		ecsServiceV2 := EcsServiceV2{client}
		if err := ecsServiceV2.SetResourceTags(d, "securitygroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudEcsSecurityGroupRead(d, meta)
}

func resourceAliCloudEcsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
