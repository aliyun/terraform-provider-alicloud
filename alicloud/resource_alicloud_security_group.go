package alicloud

import (
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

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"inner_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateCreateSecurityGroupRequest()

	if v := d.Get("name").(string); v != "" {
		request.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		request.VpcId = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "security_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ecs.CreateSecurityGroupResponse)
	d.SetId(response.SecurityGroupId)
	addDebug(request.GetActionName(), raw)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeSecurityGroupAttribute(d.Id())
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

	if d.HasChange("inner_access") || d.IsNewResource() {
		if !(d.Get("inner_access").(bool) && d.IsNewResource()) {
			policy := GroupInnerAccept
			if !d.Get("inner_access").(bool) {
				policy = GroupInnerDrop
			}
			request := ecs.CreateModifySecurityGroupPolicyRequest()
			request.SecurityGroupId = d.Id()
			request.InnerAccessPolicy = string(policy)
			request.ClientToken = buildClientToken(request.GetActionName())

			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifySecurityGroupPolicy(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
			d.SetPartial("inner_access")
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunSecurityGroupRead(d, meta)
	}

	update := false
	request := ecs.CreateModifySecurityGroupAttributeRequest()
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
		addDebug(request.GetActionName(), raw)
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
	request.SecurityGroupId = d.Id()

	return resource.Retry(6*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(request)
		})

		if err != nil {
			if IsExceptedError(err, SgDependencyViolation) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)

		_, err = ecsService.DescribeSecurityGroupAttribute(d.Id())

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})

}
