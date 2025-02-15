// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"time"

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
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"inner_access"},
				ValidateFunc:  StringInSlice([]string{"Accept", "Drop"}, false),
				// The InnerAccessPolicy attribute of enterprise level security group can't be modified.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("security_group_type").(string) == "enterprise"
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"normal", "enterprise"}, false),
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//There is a bug in the SecurityGroupName attribute in CreateSecurityGroup
				//ValidateFunc: StringMatch(regexp.MustCompile("^[a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9:_-]{2,128}$"), "The name must be 2 to 128 characters in length. It must start with a letter and cannot start with `http://` or `https://`. It can contain letters, digits, colons (:), underscores (_), and hyphens (-)."),
				Deprecated: "Field `name` has been deprecated from provider version 1.239.0. New field `security_group_name` instead.",
			},
			"inner_access": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `inner_access` has been deprecated from provider version 1.55.3. New field `inner_access_policy` instead.",
			},
		},
	}
}

func resourceAliCloudEcsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("security_group_type"); ok {
		request["SecurityGroupType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	if v, ok := d.GetOk("security_group_name"); ok {
		request["SecurityGroupName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["SecurityGroupName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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

	if objectRaw["CreationTime"] != nil {
		d.Set("create_time", objectRaw["CreationTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["ResourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	}
	if objectRaw["SecurityGroupName"] != nil {
		d.Set("security_group_name", objectRaw["SecurityGroupName"])
		d.Set("name", objectRaw["SecurityGroupName"])
	}
	if objectRaw["SecurityGroupType"] != nil {
		d.Set("security_group_type", objectRaw["SecurityGroupType"])
	}
	if objectRaw["VpcId"] != nil {
		d.Set("vpc_id", objectRaw["VpcId"])
	}

	tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = ecsServiceV2.DescribeSecurityGroupDescribeSecurityGroupAttribute(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	if objectRaw["InnerAccessPolicy"] != nil {
		d.Set("inner_access_policy", objectRaw["InnerAccessPolicy"])
		d.Set("inner_access", fmt.Sprint(objectRaw["InnerAccessPolicy"]) == string(GroupInnerAccept))
	}

	return nil
}

func resourceAliCloudEcsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	action := "ModifySecurityGroupAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("security_group_name") {
		update = true

		if v, ok := d.GetOk("security_group_name"); ok {
			request["SecurityGroupName"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true

		if v, ok := d.GetOk("name"); ok {
			request["SecurityGroupName"] = v
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ModifySecurityGroupPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("inner_access_policy") {
		update = true

		if v, ok := d.GetOk("inner_access_policy"); ok {
			request["InnerAccessPolicy"] = v
		}
	}

	objectRaw, err := ecsServiceV2.DescribeSecurityGroupDescribeSecurityGroupAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}

	innerAccess, ok := d.GetOkExists("inner_access")
	if ok && (innerAccess != (fmt.Sprint(objectRaw["InnerAccessPolicy"]) == string(GroupInnerAccept))) {
		update = true

		switch innerAccess {
		case true:
			request["InnerAccessPolicy"] = "Accept"
		case false:
			request["InnerAccessPolicy"] = "Drop"
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "JoinResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ResourceType"] = "securitygroup"
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsServiceV2.SetResourceTags(d, "securitygroup"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudEcsSecurityGroupRead(d, meta)
}

func resourceAliCloudEcsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SecurityGroupId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"DependencyViolation"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
