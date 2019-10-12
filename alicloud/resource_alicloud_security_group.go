package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupCreate,
		Read:   resourceAliyunSecurityGroupRead,
		Update: resourceAliyunSecurityGroupUpdate,
		Delete: resourceAliyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupDescription,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"inner_access": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'inner_access' has been deprecated from provider version 1.55.3. Use 'inner_access_policy' replaces it.",
			},
			"inner_access_policy": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"inner_access"},
				ValidateFunc:  validateAllowedStringValue([]string{"Accept", "Drop"}),
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateCreateSecurityGroupRequest()
	request.RegionId = client.RegionId

	if v := d.Get("name").(string); v != "" {
		request.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_security_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateSecurityGroupResponse)
	d.SetId(response.SecurityGroupId)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeSecurityGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("name", object.SecurityGroupName)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("inner_access", object.InnerAccessPolicy == string(GroupInnerAccept))
	d.Set("inner_access_policy", object.InnerAccessPolicy)

	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.SecurityGroupId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)

	if len(response.SecurityGroups.SecurityGroup) == 1 &&
		response.SecurityGroups.SecurityGroup[0].SecurityGroupId == d.Id() {
		d.Set("resource_group_id", response.SecurityGroups.SecurityGroup[0].ResourceGroupId)
	}

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceSecurityGroup)
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if err := setTags(client, TagResourceSecurityGroup, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	if d.HasChange("inner_access_policy") || d.HasChange("inner_access") || d.IsNewResource() {
		policy := GroupInnerAccept
		if v, ok := d.GetOk("inner_access_policy"); ok && v.(string) != "" {
			policy = GroupInnerAccessPolicy(v.(string))
		} else if v, ok := d.GetOkExists("inner_access"); ok && !v.(bool) {
			policy = GroupInnerDrop
		}

		request := ecs.CreateModifySecurityGroupPolicyRequest()
		request.RegionId = client.RegionId
		request.SecurityGroupId = d.Id()
		request.InnerAccessPolicy = string(policy)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("inner_access")
		d.SetPartial("inner_access_policy")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunSecurityGroupRead(d, meta)
	}

	update := false
	request := ecs.CreateModifySecurityGroupAttributeRequest()
	request.RegionId = client.RegionId
	request.SecurityGroupId = d.Id()
	if d.HasChange("name") {
		request.SecurityGroupName = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
		d.SetPartial("description")
	}

	d.Partial(false)

	return resourceAliyunSecurityGroupRead(d, meta)
}

func resourceAliyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.RegionId = client.RegionId
	request.SecurityGroupId = d.Id()

	err := resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(request)
		})

		if err != nil {
			if IsExceptedError(err, SgDependencyViolation) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR)
	}
	return WrapError(ecsService.WaitForSecurityGroup(d.Id(), Deleted, DefaultTimeoutMedium))

}
