package alicloud

import (
	"fmt"

	cr20181201 "github.com/alibabacloud-go/cr-20181201/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCrEEBuildRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrEEBuildRuleCreate,
		Read:   resourceAlicloudCrEEBuildRuleRead,
		Delete: resourceAlicloudCrEEBuildRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"artifact_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"build_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCrEEBuildRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceID := d.Get("instance_id").(string)
	scopeType := d.Get("scope_type").(string)
	scopeID := d.Get("scope_id").(string)
	artifactType := d.Get("artifact_type").(string)
	response := &cr20181201.CreateArtifactBuildRuleResponse{}
	request := cr20181201.CreateArtifactBuildRuleRequest{
		ArtifactType: tea.String(artifactType),
		InstanceId:   tea.String(instanceID),
		ScopeId:      tea.String(scopeID),
		ScopeType:    tea.String(scopeType),
	}
	raw, err := crService.client.WithCr20181201Client(func(c *cr20181201.Client) (interface{}, error) {
		return c.CreateArtifactBuildRule(&request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_build_rule", "CreateArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}
	addDebug("CreateArtifactBuildRule", raw, request)
	response, _ = raw.(*cr20181201.CreateArtifactBuildRuleResponse)
	if response.Body == nil || !tea.BoolValue(response.Body.IsSuccess) || tea.StringValue(response.Body.BuildRuleId) == "" {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_cr_ee_build_rule", "CreateArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(instanceID, ":", scopeID, ":", tea.StringValue(response.Body.BuildRuleId), ":", scopeType, ":", artifactType))

	return resourceAlicloudCrEEBuildRuleRead(d, meta)
}

func resourceAlicloudCrEEBuildRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	id := d.Id()
	parts, err := ParseResourceId(id, 5)
	if err != nil {
		return WrapError(err)
	}
	instanceID := parts[0]
	scopeID := parts[1]
	buildRuleID := parts[2]
	scopeType := parts[3]
	artifactType := parts[4]
	response := &cr20181201.GetArtifactBuildRuleResponse{}
	request := cr20181201.GetArtifactBuildRuleRequest{
		ArtifactType: tea.String(artifactType),
		BuildRuleId:  tea.String(buildRuleID),
		InstanceId:   tea.String(instanceID),
		ScopeId:      tea.String(scopeID),
		ScopeType:    tea.String(scopeType),
	}
	raw, err := crService.client.WithCr20181201Client(func(c *cr20181201.Client) (interface{}, error) {
		return c.GetArtifactBuildRule(&request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, "GetArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}
	addDebug("GetArtifactBuildRule", raw, request)

	response, _ = raw.(*cr20181201.GetArtifactBuildRuleResponse)
	if response.Body == nil || !tea.BoolValue(response.Body.IsSuccess) || tea.StringValue(response.Body.ScopeId) != scopeID ||
		tea.StringValue(response.Body.BuildRuleId) != buildRuleID || tea.StringValue(response.Body.ScopeType) != scopeType ||
		tea.StringValue(response.Body.ArtifactType) != artifactType {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, id, "GetArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}
	d.Set("instance_id", instanceID)
	d.Set("scope_id", response.Body.ScopeId)
	d.Set("build_rule_id", response.Body.BuildRuleId)
	d.Set("artifact_type", response.Body.ArtifactType)
	d.Set("scope_type", response.Body.ScopeType)

	return nil
}

func resourceAlicloudCrEEBuildRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	parts, err := ParseResourceId(d.Id(), 5)
	if err != nil {
		return WrapError(err)
	}
	instanceID := parts[0]
	scopeID := parts[1]
	scopeType := parts[3]
	artifactType := parts[4]
	request := &DeleteArtifactBuildRuleRequest{
		InstanceID:   instanceID,
		ScopeType:    scopeType,
		ArtifactType: artifactType,
		ScopeID:      scopeID,
	}
	response, err := crService.DeleteArtifactBuildRule(request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_build_rule", "DeleteArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}
	addDebug("DeleteArtifactBuildRule", response, request)
	if !response.Data.IsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), "DeleteArtifactBuildRule", AlibabaCloudSdkGoERROR)
	}

	return nil
}
