package alicloud

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssAlbNlbServerGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAlbNlbServerGroupsCreate,
		Read:   resourceAliyunEssAlbNlbServerGroupsRead,
		Update: resourceAliyunEssAlbNlbServerGroupsUpdate,
		Delete: resourceAliyunEssAlbNlbServerGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					if v, ok := m["server_group_id"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := m["port"]; ok {
						buf.WriteString(fmt.Sprintf("%d-", v.(int)))
					}
					return hashcode.String(buf.String())
				},
			},
			"force_attach": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunEssAlbNlbServerGroupsCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("scaling_group_id").(string))
	return resourceAliyunEssAlbNlbServerGroupsUpdate(d, meta)
}

func resourceAliyunEssAlbNlbServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("scaling_group_id", object.ScalingGroupId)
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("server_groups", essService.flattenServerGroupList(object.ServerGroups.ServerGroup))
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliyunEssAlbNlbServerGroupsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	serverGroupsMapFromScalingGroup := serverGroupMapFromScalingGroup(object.ServerGroups.ServerGroup)
	serverGroupsMapFromConfig := serverGroupMapFromConfig(d.Get("server_groups").(*schema.Set))
	attachMap, detachMap := attachOrDetachServerGroupMap(serverGroupsMapFromConfig, serverGroupsMapFromScalingGroup)

	err = detachServerGroups(d, client, detachMap)
	if err != nil {
		return WrapError(err)
	}
	d.SetPartial("server_groups")

	err = attachServerGroups(d, client, attachMap)
	if err != nil {
		return WrapError(err)
	}
	d.SetPartial("server_groups")
	d.Partial(false)
	return resourceAliyunEssAlbNlbServerGroupsRead(d, meta)
}

func resourceAliyunEssAlbNlbServerGroupsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	serverGroupsFromConfig := serverGroupMapFromConfig(d.Get("server_groups").(*schema.Set))
	_, detachMap := attachOrDetachServerGroupMap(make(map[string]string, 0), serverGroupsFromConfig)

	err := detachServerGroups(d, client, detachMap)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

func serverGroupMapFromScalingGroup(serverGroup []ess.ServerGroup) map[string]string {
	serverGroupMap := make(map[string]string)
	if serverGroup != nil && len(serverGroup) > 0 {
		for _, v := range serverGroup {
			key := fmt.Sprintf("%s_%s_%d_%d", v.ServerGroupId, v.Type, v.Port, v.Weight)
			serverGroupMap[key] = key
		}
	}
	return serverGroupMap
}

func serverGroupMapFromConfig(serverGroups *schema.Set) map[string]string {
	serverGroupMap := make(map[string]string)
	serverGroupList := serverGroups.List()
	if len(serverGroupList) > 0 {
		for _, v := range serverGroupList {
			serverGroup := v.(map[string]interface{})
			serverGroupId := serverGroup["server_group_id"].(string)
			serverType := serverGroup["type"].(string)
			port := serverGroup["port"].(int)
			weight := serverGroup["weight"].(int)
			key := fmt.Sprintf("%s_%s_%d_%d", serverGroupId, serverType, port, weight)
			serverGroupMap[key] = key
		}
	}
	return serverGroupMap
}

func attachOrDetachServerGroupMap(newMap map[string]string, oldMap map[string]string) (map[string]string, map[string]string) {
	attachMap := make(map[string]string)
	detachMap := make(map[string]string)
	for k, v := range newMap {
		if _, ok := oldMap[k]; !ok {
			attachMap[k] = v
		}
	}
	for k, v := range oldMap {
		if _, ok := newMap[k]; !ok {
			detachMap[k] = v
		}
	}
	return attachMap, detachMap
}

func buildEssServerGroupListMap(serverGroupMap map[string]string) map[string][]string {
	serverGroupRequestMap := make(map[string][]string, 0)
	for _, v := range serverGroupMap {
		attrs := strings.Split(v, "_")
		serverGroupId := attrs[0]
		if _, ok := serverGroupRequestMap[serverGroupId]; !ok {
			serverGroupAttributes := make([]string, 0)
			serverGroupAttributes = append(serverGroupAttributes, v)
			serverGroupRequestMap[serverGroupId] = serverGroupAttributes
		} else {
			serverGroupAttributes := serverGroupRequestMap[serverGroupId]
			serverGroupAttributes = append(serverGroupAttributes, v)
			serverGroupRequestMap[serverGroupId] = serverGroupAttributes
		}
	}
	return serverGroupRequestMap
}

func attachServerGroups(d *schema.ResourceData, client *connectivity.AliyunClient, attachMap map[string]string) error {
	if len(attachMap) > 0 {
		serverGroupListMap := buildEssServerGroupListMap(attachMap)
		request := map[string]interface{}{}
		var response map[string]interface{}
		attachScalingGroupServerGroups := make([]map[string]interface{}, 0)
		for k, v := range serverGroupListMap {
			serverGroup := map[string]interface{}{
				"ServerGroupId": k,
			}
			for _, e := range v {
				attrs := strings.Split(e, "_")
				serverGroup["Type"] = attrs[1]
				serverGroup["Port"] = attrs[2]
				serverGroup["Weight"] = attrs[3]
			}
			attachScalingGroupServerGroups = append(attachScalingGroupServerGroups, serverGroup)
		}
		request["RegionId"] = client.RegionId
		request["ScalingGroupId"] = d.Id()
		request["ForceAttach"] = requests.NewBoolean(d.Get("force_attach").(bool))
		request["ServerGroup"] = attachScalingGroupServerGroups

		action := "AttachServerGroups"
		conn, err := client.NewEssClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"VServerGroupProcessing", "BackendServer.configuring"}) {
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
	return nil
}

func detachServerGroups(d *schema.ResourceData, client *connectivity.AliyunClient, detachMap map[string]string) error {
	if len(detachMap) > 0 {
		request := map[string]interface{}{}
		var response map[string]interface{}
		serverGroupListMap := buildEssServerGroupListMap(detachMap)
		detachScalingGroupServerGroups := make([]map[string]interface{}, 0)
		for k, v := range serverGroupListMap {
			serverGroup := map[string]interface{}{
				"ServerGroupId": k,
			}
			for _, e := range v {
				attrs := strings.Split(e, "_")
				serverGroup["Type"] = attrs[1]
				serverGroup["Port"] = attrs[2]
				serverGroup["Weight"] = attrs[3]
			}
			detachScalingGroupServerGroups = append(detachScalingGroupServerGroups, serverGroup)
		}

		request["RegionId"] = client.RegionId
		request["ScalingGroupId"] = d.Id()
		request["ForceAttach"] = requests.NewBoolean(d.Get("force_attach").(bool))
		request["ServerGroup"] = detachScalingGroupServerGroups

		action := "DetachServerGroups"
		conn, err := client.NewEssClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken(action)
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-28"), StringPointer("AK"), nil, request, &runtime)
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
	return nil
}
