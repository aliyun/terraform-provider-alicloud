package cloudsso

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	convert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudCloudSsoAccessConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSsoAccessConfigurationCreate,
		Read:   resourceAliCloudCloudSsoAccessConfigurationRead,
		Update: resourceAliCloudCloudSsoAccessConfigurationUpdate,
		Delete: resourceAliCloudCloudSsoAccessConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringMatch(regexp.MustCompile(`^[a-zA-z0-9-]{1,32}$`), "The name of the resource. The name can be up to `32` characters long and can contain letters, digits, and hyphens (-)"),
			},
			"session_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(900, 43200),
			},
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringLenBetween(0, 1024),
			},
			"force_remove_permission_policies": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"permission_policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					policy := v.(map[string]interface{})
					if v, ok := policy["permission_policy_type"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := policy["permission_policy_name"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := policy["permission_policy_document"]; ok {
						document := make(map[string]interface{})
						err := json.Unmarshal([]byte(v.(string)), &document)
						if err == nil {
							documentString, _ := json.Marshal(document)
							buf.WriteString(fmt.Sprintf("%s-", documentString))
						}
					}
					return int(crc32.ChecksumIEEE([]byte(buf.String())))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_policy_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"System", "Inline"}, false),
						},
						"permission_policy_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"permission_policy_document": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
					},
				},
			},
			"access_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudSsoAccessConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessConfiguration"
	request := make(map[string]interface{})
	var err error

	request["DirectoryId"] = d.Get("directory_id")
	request["AccessConfigurationName"] = d.Get("access_configuration_name")

	if v, ok := d.GetOk("session_duration"); ok {
		request["SessionDuration"] = v
	}

	if v, ok := d.GetOk("relay_state"); ok {
		request["RelayState"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_sso_access_configuration", action, errs.AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.AccessConfiguration", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_cloud_sso_access_configuration")
	} else {
		accessConfigurationId := resp.(map[string]interface{})["AccessConfigurationId"]
		d.SetId(fmt.Sprintf("%v:%v", request["DirectoryId"], accessConfigurationId))
	}

	return resourceAliCloudCloudSsoAccessConfigurationUpdate(d, meta)
}

func resourceAliCloudCloudSsoAccessConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}

	object, err := cloudssoService.DescribeCloudSsoAccessConfiguration(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_access_configuration cloudssoService.DescribeCloudSsoAccessConfiguration Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("directory_id", parts[0])
	d.Set("access_configuration_name", object["AccessConfigurationName"])
	d.Set("relay_state", object["RelayState"])
	d.Set("description", object["Description"])
	d.Set("access_configuration_id", object["AccessConfigurationId"])

	if sessionDuration, ok := object["SessionDuration"]; ok && fmt.Sprint(sessionDuration) != "0" {
		session_durationVal, err := convert.FormatInt(sessionDuration)
		if err != nil {
			return err
		}
		d.Set("session_duration", session_durationVal)
	}

	listPermissionPoliciesInAccessConfigurationObject, err := cloudssoService.ListPermissionPoliciesInAccessConfiguration(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}

	if permissionPolicies, ok := listPermissionPoliciesInAccessConfigurationObject["PermissionPolicies"]; ok && permissionPolicies != nil {
		permissionPoliciesMaps := make([]map[string]interface{}, 0)
		for _, permissionPoliciesList := range permissionPolicies.([]interface{}) {
			permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})
			permissionPoliciesMap := make(map[string]interface{})

			if permissionPolicyType, ok := permissionPoliciesArg["PermissionPolicyType"]; ok {
				permissionPoliciesMap["permission_policy_type"] = permissionPolicyType
			}

			if permissionPolicyName, ok := permissionPoliciesArg["PermissionPolicyName"]; ok {
				permissionPoliciesMap["permission_policy_name"] = permissionPolicyName
			}

			if permissionPolicyDocument, ok := permissionPoliciesArg["PermissionPolicyDocument"]; ok {
				permissionPoliciesMap["permission_policy_document"] = permissionPolicyDocument
			}

			permissionPoliciesMaps = append(permissionPoliciesMaps, permissionPoliciesMap)
		}

		d.Set("permission_policies", permissionPoliciesMaps)
	}

	return nil
}

func resourceAliCloudCloudSsoAccessConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	var response map[string]interface{}
	d.Partial(true)

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	update := false
	updateAccessConfigurationReq := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	if !d.IsNewResource() && d.HasChange("session_duration") {
		update = true

		if v, ok := d.GetOk("session_duration"); ok {
			updateAccessConfigurationReq["NewSessionDuration"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("relay_state") {
		update = true
	}
	if v, ok := d.GetOk("relay_state"); ok {
		updateAccessConfigurationReq["NewRelayState"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		updateAccessConfigurationReq["NewDescription"] = v
	}

	if update {
		action := "UpdateAccessConfiguration"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, updateAccessConfigurationReq, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationConflict.Task"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateAccessConfigurationReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}

	if d.HasChange("permission_policies") {
		oldPermissionPolicies, newPermissionPolicies := d.GetChange("permission_policies")
		removed := oldPermissionPolicies.(*schema.Set).Difference(newPermissionPolicies.(*schema.Set)).List()
		added := newPermissionPolicies.(*schema.Set).Difference(oldPermissionPolicies.(*schema.Set)).List()

		if len(removed) > 0 {
			action := "RemovePermissionPolicyFromAccessConfiguration"

			removePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range removed {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
					response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, removePermissionPolicyFromAccessConfigurationReq, false)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationConflict.Task"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, removePermissionPolicyFromAccessConfigurationReq)

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(added) > 0 {
			action := "AddPermissionPolicyToAccessConfiguration"

			addPermissionPolicyToAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range added {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				addPermissionPolicyToAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				addPermissionPolicyToAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				if addPermissionPolicyToAccessConfigurationReq["PermissionPolicyType"] == "Inline" {

					if inlinePolicyDocument, ok := permissionPoliciesArg["permission_policy_document"]; ok {
						addPermissionPolicyToAccessConfigurationReq["InlinePolicyDocument"] = inlinePolicyDocument
					}
				}

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
					response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, addPermissionPolicyToAccessConfigurationReq, false)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationConflict.Task"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, addPermissionPolicyToAccessConfigurationReq)

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		// Provisioning access configuration when permission policies has changed.
		objects, err := cloudssoService.DescribeCloudSsoAccessConfigurationProvisionings(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]))
		if err != nil {
			return errs.WrapError(err)
		}

		for _, object := range objects {
			err = cloudssoService.CloudssoServicAccessConfigurationProvisioning(fmt.Sprint(parts[0]), fmt.Sprint(parts[1]), fmt.Sprint(object["TargetType"]), fmt.Sprint(object["TargetId"]))
			if err != nil {
				return errs.WrapError(err)
			}
		}

	}

	d.Partial(false)

	return resourceAliCloudCloudSsoAccessConfigurationRead(d, meta)
}

func resourceAliCloudCloudSsoAccessConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	if v, ok := d.GetOk("permission_policies"); ok {
		removed := v.(*schema.Set).List()

		if len(removed) > 0 {
			action := "RemovePermissionPolicyFromAccessConfiguration"

			removePermissionPolicyFromAccessConfigurationReq := map[string]interface{}{
				"DirectoryId":           parts[0],
				"AccessConfigurationId": parts[1],
			}

			for _, permissionPoliciesList := range removed {
				permissionPoliciesArg := permissionPoliciesList.(map[string]interface{})

				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyType"] = permissionPoliciesArg["permission_policy_type"]
				removePermissionPolicyFromAccessConfigurationReq["PermissionPolicyName"] = permissionPoliciesArg["permission_policy_name"]

				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
					response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, removePermissionPolicyFromAccessConfigurationReq, false)
					if err != nil {
						if errs.IsExpectedErrors(err, []string{"OperationConflict.Task"}) || helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, removePermissionPolicyFromAccessConfigurationReq)

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	action := "DeleteAccessConfiguration"

	request := map[string]interface{}{
		"DirectoryId":           parts[0],
		"AccessConfigurationId": parts[1],
	}

	if v, ok := d.GetOkExists("force_remove_permission_policies"); ok {
		request["ForceRemovePermissionPolicies"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeletionConflict.AccessConfiguration.Provisioning", "DeletionConflict.AccessConfiguration.AccessAssignment", "OperationConflict.Task", "DeletionConflict.AccessConfiguration.Task"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.AccessConfiguration"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
