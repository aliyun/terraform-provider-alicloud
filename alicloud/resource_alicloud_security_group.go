package alicloud

import (
	"fmt"
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

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(buildAliyunSecurityGroupArgs(d, meta))
	})
	if err != nil {
		actionName := buildAliyunSecurityGroupArgs(d, meta).GetActionName()
		return WrapErrorf(err, DefaultErrorMsg, "security_group:", actionName, AlibabaCloudSdkGoERROR)
	}
	resp, _ := raw.(*ecs.CreateSecurityGroupResponse)
	if resp == nil {
		return WrapError(fmt.Errorf("Creating security group got a nil response."))
	}
	d.SetId(resp.SecurityGroupId)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	var sg *ecs.DescribeSecurityGroupAttributeResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		group, e := ecsService.DescribeSecurityGroupAttribute(d.Id())
		if e != nil {
			if NotFoundError(e) || IsExceptedErrors(e, []string{InvalidSecurityGroupIdNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error DescribeSecurityGroupAttribute: %#v", e))
		}
		if group.SecurityGroupId != "" {
			sg = &group
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Create security group timeout and got an error: %#v", e))
	})

	if sg == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", sg.SecurityGroupName)
	d.Set("description", sg.Description)
	d.Set("vpc_id", sg.VpcId)
	d.Set("inner_access", sg.InnerAccessPolicy == string(GroupInnerAccept))

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
	attributeUpdate := false
	args := ecs.CreateModifySecurityGroupAttributeRequest()
	args.SecurityGroupId = d.Id()

	if err := setTags(client, TagResourceSecurityGroup, d); err != nil {
		return WrapError(err)
	} else {
		d.SetPartial("tags")
	}

	if d.HasChange("name") && !d.IsNewResource() {
		args.SecurityGroupName = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") && !d.IsNewResource() {
		args.Description = d.Get("description").(string)
		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupAttribute(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("description")
	}

	if d.HasChange("inner_access") || d.IsNewResource() {
		policy := GroupInnerAccept
		if !d.Get("inner_access").(bool) {
			policy = GroupInnerDrop
		}
		args := ecs.CreateModifySecurityGroupPolicyRequest()
		args.SecurityGroupId = d.Id()
		args.InnerAccessPolicy = string(policy)

		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("inner_access")
	}

	d.Partial(false)

	return resourceAliyunSecurityGroupRead(d, meta)
}

func resourceAliyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	req := ecs.CreateDeleteSecurityGroupRequest()
	req.SecurityGroupId = d.Id()

	return resource.Retry(6*time.Minute, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(req)
		})

		if err != nil {
			if IsExceptedError(err, SgDependencyViolation) {
				return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
			}
		}

		sg, err := ecsService.DescribeSecurityGroupAttribute(d.Id())

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		} else if sg.SecurityGroupId == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
	})

}

func buildAliyunSecurityGroupArgs(d *schema.ResourceData, meta interface{}) *ecs.CreateSecurityGroupRequest {

	args := ecs.CreateCreateSecurityGroupRequest()

	if v := d.Get("name").(string); v != "" {
		args.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		args.VpcId = v
	}
	args.ClientToken = buildClientToken("TF-CreateSecurityGroup")

	return args
}
