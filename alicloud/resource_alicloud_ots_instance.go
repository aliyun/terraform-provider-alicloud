package alicloud

import (
	ots "github.com/alibabacloud-go/tablestore-20201209/v3/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudOtsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceCreate,
		Read:   resourceAliyunOtsInstanceRead,
		Update: resourceAliyunOtsInstanceUpdate,
		Delete: resourceAliyunOtsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},

			// Expired
			"accessed_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(AnyNetwork), string(VpcOnly), string(VpcOrConsole),
				}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_type_acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_source_acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"policy_version": {
				Type:     schema.TypeInt,
				Optional: false,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  OtsHighPerformance,
				ValidateFunc: validation.StringInSlice([]string{
					string(OtsCapacity), string(OtsHighPerformance),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alias_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOTSInstanceName,
			},
			// lintignore: S006
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	instanceType := d.Get("instance_type").(string)
	request := new(ots.CreateInstanceRequest)
	request.ClusterType = StringPointer(convertInstanceType(OtsInstanceType(instanceType)))
	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return WrapError(err)
	}
	valid := false
	for _, t := range types {
		if *request.ClusterType == *t {
			valid = true
			break
		}
	}
	if !valid {
		return WrapError(Error("The instance type %s is not available in the region %s.", instanceType, client.RegionId))
	}
	instanceName := d.Get("name").(string)
	request.InstanceName = StringPointer(instanceName)
	request.ResourceGroupId = StringPointer(d.Get("resource_group_id").(string))
	request.Policy = StringPointer(d.Get("policy").(string))
	request.InstanceDescription = StringPointer(d.Get("description").(string))

	hasSetNetwork := false
	if v, ok := d.GetOk("accessed_by"); ok {
		hasSetNetwork = true
		request.Network = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(v.(string))))
	}

	hasSetACL := false
	// LIST or SET cannot set default values in schema in latest terraform version, so do it manually
	// terraform cannot handle nil and[] in list/set: https://github.com/hashicorp/terraform-plugin-sdk/issues/142
	// in terraform the zero value of list/set is [], in golang the zero value of slice is nil
	netTypeList := []*string{StringPointer(string(VpcAccess)), StringPointer(string(ClassicAccess)), StringPointer(string(InternetAccess))}
	// v not nil and [], it will be ok
	if v, ok := d.GetOk("network_type_acl"); ok {
		hasSetACL = true
		netTypeList = expandStringPointerList(v.(*schema.Set).List())
	}
	request.NetworkTypeACL = netTypeList

	netSourceList := []*string{StringPointer(string(TrustProxyAccess))}
	if v, ok := d.GetOk("network_source_acl"); ok {
		hasSetACL = true
		netSourceList = expandStringPointerList(v.(*schema.Set).List())
	}
	request.NetworkSourceACL = netSourceList

	// In order to maintain compatibility, when the Network attribute is set,
	// the ACL attribute cannot have a default value.
	if hasSetNetwork && !hasSetACL {
		request.NetworkTypeACL = nil
		request.NetworkSourceACL = nil
	}

	if tagMap, ok := d.GetOk("tags"); ok {
		var tags []*ots.CreateInstanceRequestTags

		for key, value := range tagMap.(map[string]interface{}) {
			tags = append(tags, &ots.CreateInstanceRequestTags{
				Key:   StringPointer(key),
				Value: StringPointer(value.(string)),
			})
		}
		request.Tags = tags
	}

	raw, err := client.WithOtsClient(func(client *ots.Client) (interface{}, error) {
		return client.CreateInstance(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, "CreateInstance", AlibabaCloudSdkGoERROR)
	}
	resp := raw.(*ots.CreateInstanceResponse)
	addDebug("CreateInstance", resp, request)

	d.SetId(instanceName)
	if err = otsService.WaitForOtsInstance(instanceName, toInstanceInnerStatus(Running), DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	instance, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", instance.InstanceName)
	err = d.Set("accessed_by", convertInstanceAccessedByRevert(*instance.Network))
	if err != nil {
		return err
	}
	err = d.Set("resource_group_id", instance.ResourceGroupId)
	if err != nil {
		return err
	}
	err = d.Set("network_type_acl", instance.NetworkTypeACL)
	if err != nil {
		return err
	}
	err = d.Set("network_source_acl", instance.NetworkSourceACL)
	if err != nil {
		return err
	}
	err = d.Set("policy", instance.Policy)
	if err != nil {
		return err
	}
	err = d.Set("policy_version", instance.PolicyVersion)
	if err != nil {
		return err
	}
	err = d.Set("instance_type", convertInstanceTypeRevert(*instance.InstanceSpecification))
	if err != nil {
		return err
	}
	err = d.Set("description", instance.InstanceDescription)
	if err != nil {
		return err
	}
	err = d.Set("alias_name", instance.AliasName)
	if err != nil {
		return err
	}
	err = d.Set("tags", otsRestTagsToMap(instance.Tags))
	if err != nil {
		return err
	}
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := new(ots.ChangeResourceGroupRequest)
		request.ResourceId = StringPointer(d.Id())
		request.NewResourceGroupId = StringPointer(d.Get("resource_group_id").(string))

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ChangeResourceGroup(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "ChangeResourceGroup", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.ChangeResourceGroupResponse)
		addDebug("ChangeResourceGroup", resp, request)
		d.SetPartial("resource_group_id")
	}
	hasChangeACL := false
	if !d.IsNewResource() && d.HasChange("network_type_acl") {
		request := new(ots.UpdateInstanceRequest)
		request.InstanceName = StringPointer(d.Id())

		netTypeList := expandStringPointerList(d.Get("network_type_acl").(*schema.Set).List())
		request.NetworkTypeACL = netTypeList
		// acl must set together
		netSourceList := expandStringPointerList(d.Get("network_source_acl").(*schema.Set).List())
		request.NetworkSourceACL = netSourceList

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "UpdateInstance", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.UpdateInstanceResponse)
		addDebug("UpdateInstance", resp, request)
		hasChangeACL = true
		d.SetPartial("network_type_acl")
	}

	if !d.IsNewResource() && d.HasChange("network_source_acl") {
		request := new(ots.UpdateInstanceRequest)
		request.InstanceName = StringPointer(d.Id())

		netSourceList := expandStringPointerList(d.Get("network_source_acl").(*schema.Set).List())
		request.NetworkSourceACL = netSourceList
		// acl must set together
		netTypeList := expandStringPointerList(d.Get("network_type_acl").(*schema.Set).List())
		request.NetworkTypeACL = netTypeList

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "UpdateInstance", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.UpdateInstanceResponse)
		addDebug("UpdateInstance", resp, request)
		hasChangeACL = true
		d.SetPartial("network_source_acl")
	}

	// accessed_by is a deprecated attribute, updates on accessed_by will only take effect when the ACL has not been updated.
	if !d.IsNewResource() && (d.HasChange("accessed_by") && !hasChangeACL) {
		request := new(ots.UpdateInstanceRequest)
		request.InstanceName = StringPointer(d.Id())
		request.Network = StringPointer(convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string))))

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "UpdateInstance", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.UpdateInstanceResponse)
		addDebug("UpdateInstance", resp, request)
		d.SetPartial("accessed_by")
	}

	if !d.IsNewResource() && d.HasChange("policy") {
		policy := d.Get("policy").(string)
		if policy != "" {
			request := new(ots.UpdateInstancePolicyRequest)
			request.InstanceName = StringPointer(d.Id())
			request.Policy = StringPointer(policy)
			request.PolicyVersion = Int64Pointer(int64(d.Get("policy_version").(int)))
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.UpdateInstancePolicy(request)
			})
			if err != nil {
				if NotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
				}
				return WrapErrorf(err, DefaultErrorMsg, "UpdateInstancePolicy", AlibabaCloudSdkGoERROR)
			}
			resp := raw.(*ots.UpdateInstancePolicyResponse)
			addDebug("UpdateInstancePolicy", resp, request)
		} else {
			request := new(ots.DeleteInstancePolicyRequest)
			request.InstanceName = StringPointer(d.Id())
			request.PolicyVersion = Int64Pointer(int64(d.Get("policy_version").(int)))
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.DeleteInstancePolicy(request)
			})
			if err != nil {
				if NotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
				}
				return WrapErrorf(err, DefaultErrorMsg, "DeleteInstancePolicy", AlibabaCloudSdkGoERROR)
			}
			resp := raw.(*ots.DeleteInstancePolicyResponse)
			addDebug("DeleteInstancePolicy", resp, request)
		}
		d.SetPartial("policy")
	}

	if !d.IsNewResource() && d.HasChange("description") {
		request := new(ots.UpdateInstanceRequest)
		request.InstanceName = StringPointer(d.Id())
		request.InstanceDescription = StringPointer(d.Get("description").(string))

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "UpdateInstance", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.UpdateInstanceResponse)
		addDebug("UpdateInstance", resp, request)
		d.SetPartial("description")
	}
	if d.HasChange("alias_name") {
		request := new(ots.UpdateInstanceRequest)
		request.InstanceName = StringPointer(d.Id())
		request.AliasName = StringPointer(d.Get("alias_name").(string))

		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			if NotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, "UpdateInstance", AlibabaCloudSdkGoERROR)
		}
		resp := raw.(*ots.UpdateInstanceResponse)
		addDebug("UpdateInstance", resp, request)
		d.SetPartial("alias_name")
	}

	if !d.IsNewResource() && d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			var removeKeys []*string
			for _, t := range remove {
				removeKeys = append(removeKeys, &t.Key)
			}

			request := new(ots.UntagResourcesRequest)
			request.ResourceType = StringPointer("INSTANCE")
			request.ResourceIds = expandStringPointerList([]interface{}{d.Id()})
			request.TagKeys = removeKeys

			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.UntagResources(request)
			})
			if err != nil {
				if NotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
				}
				return WrapErrorf(err, DefaultErrorMsg, "UntagResources", AlibabaCloudSdkGoERROR)
			}
			resp := raw.(*ots.UntagResourcesResponse)
			addDebug("UntagResources", resp, request)
		}

		if len(create) > 0 {
			var insertTags []*ots.TagResourcesRequestTags
			for _, t := range create {
				insertTags = append(insertTags, &ots.TagResourcesRequestTags{
					Key:   &t.Key,
					Value: &t.Value,
				})
			}

			request := new(ots.TagResourcesRequest)
			request.ResourceType = StringPointer("INSTANCE")
			request.ResourceIds = expandStringPointerList([]interface{}{d.Id()})
			request.Tags = insertTags

			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.TagResources(request)
			})
			if err != nil {
				if NotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
				}
				return WrapErrorf(err, DefaultErrorMsg, "TagResources", AlibabaCloudSdkGoERROR)
			}
			resp := raw.(*ots.TagResourcesResponse)
			addDebug("TagResources", resp, request)
		}
		d.SetPartial("tags")
	}
	if err := otsService.WaitForOtsInstance(d.Id(), toInstanceInnerStatus(Running), DefaultTimeout); err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := new(ots.DeleteInstanceRequest)
	request.InstanceName = StringPointer(d.Id())
	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.DeleteInstance(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteInstance", AlibabaCloudSdkGoERROR)
	}
	resp := raw.(*ots.DeleteInstanceResponse)
	addDebug("DeleteInstance", resp, request)

	otsService := OtsService{client}
	return WrapError(otsService.WaitForOtsInstance(d.Id(), string(Deleted), DefaultLongTimeout))
}
