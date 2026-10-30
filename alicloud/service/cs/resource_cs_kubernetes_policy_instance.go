// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cs

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCSKubernetesPolicyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSKubernetesPolicyInstanceCreate,
		Read:   resourceCSKubernetesPolicyInstanceRead,
		Update: resourceCSKubernetesPolicyInstanceUpdate,
		Delete: resourceCSKubernetesPolicyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Optional: true,
				//Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"warn", "deny"}, false),
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCSKubernetesPolicyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}
	csClient := CsClient{client}

	action := "DeployPolicyInstance"
	cluster_id := d.Get("cluster_id").(string)
	policy_name := d.Get("policy_name").(string)

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return errs.WrapErrorf(err, "waiting for required addons failed")
	}

	createReq := &roacs.DeployPolicyInstanceRequest{
		Action: tea.String(d.Get("action").(string)),
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convert.ConvertToInterfaceArray(v)
		stringList := make([]*string, len(namespacesMapsArray))
		for i, v := range namespacesMapsArray {
			stringList[i] = tea.String(v.(string))
		}

		createReq.Namespaces = stringList
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMap := v.(map[string]interface{})
		createReq.Parameters = convert.NormalizeMap(parametersMap)
	}

	var response *roacs.DeployPolicyInstanceResponse
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.DeployPolicyInstance(tea.String(cluster_id), tea.String(policy_name), createReq)
		if err != nil {
			if helper.NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cs_kubernetes_policy_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	var instance_name string
	if response != nil && response.Body != nil && len(response.Body.Instances) == 1 {
		instance_name = tea.StringValue(response.Body.Instances[0])
	} else {
		return errs.WrapErrorf(fmt.Errorf("no instance_name returned"), errs.DefaultErrorMsg, "alicloud_cs_kubernetes_policy_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", cluster_id, policy_name, instance_name))

	return resourceCSKubernetesPolicyInstanceRead(d, meta)
}

func resourceCSKubernetesPolicyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackServiceV2 := AckServiceV2{client}
	csClient, clientErr := client.NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}
	csClientWrapper := CsClient{csClient}

	object, err := ackServiceV2.DescribeAckPolicyInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cs_kubernetes_policy_instance DescribeAckPolicyInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	clusterId := object["cluster_id"].(string)
	// Wait for required addons to be ready
	if err := csClientWrapper.WaitForRequiredAddons(clusterId, 10*time.Minute); err != nil {
		return errs.WrapErrorf(err, "waiting for required addons failed")
	}

	d.Set("action", object["policy_action"])
	d.Set("cluster_id", object["cluster_id"])
	d.Set("instance_name", object["instance_name"])
	d.Set("policy_name", object["policy_name"])

	// Process namespaces - convert string to array
	if policyScope, ok := object["policy_scope"].(string); ok {
		if policyScope == "*" {
			// Convert "*" to empty array
			d.Set("namespaces", []interface{}{})
		} else {
			// Split comma-separated values into array
			namespaces := strings.Split(policyScope, ",")
			result := make([]interface{}, len(namespaces))
			for i, ns := range namespaces {
				result[i] = strings.TrimSpace(ns) // Trim whitespace
			}
			d.Set("namespaces", result)
		}
	}

	// Process parameters - convert YAML string to map
	if policyParameters, ok := object["policy_parameters"].(string); ok {
		paramsMap, err := convert.ConvertYamlToObject(policyParameters)
		if err == nil {
			// Apply type conversion to the parameters
			convertedParams := convert.NormalizeMap(paramsMap)
			stateParams, err := policyParametersForState(convertedParams)
			if err != nil {
				return errs.WrapError(err)
			}
			if err := d.Set("parameters", stateParams); err != nil {
				return errs.WrapError(err)
			}
		} else {
			log.Printf("[WARN] Failed to parse policy_parameters as YAML: %v", err)
		}
	}

	return nil
}

// policyParametersForState converts the heterogeneous values returned by the
// policy API to the map-of-strings representation declared by the Terraform
// schema. Composite values are JSON encoded so their state representation is
// deterministic. NormalizeMap restores the policy's supported bool, integer,
// map, and list values when the state is sent on create/update.
func policyParametersForState(parameters map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(parameters))
	for key, value := range parameters {
		stringValue, err := policyParameterValueForState(value)
		if err != nil {
			return nil, fmt.Errorf("converting policy parameter %q for state: %w", key, err)
		}
		result[key] = stringValue
	}
	return result, nil
}

func policyParameterValueForState(value interface{}) (string, error) {
	switch typed := value.(type) {
	case nil:
		return "null", nil
	case string:
		return typed, nil
	case bool:
		return strconv.FormatBool(typed), nil
	case int:
		return strconv.Itoa(typed), nil
	case int8:
		return strconv.FormatInt(int64(typed), 10), nil
	case int16:
		return strconv.FormatInt(int64(typed), 10), nil
	case int32:
		return strconv.FormatInt(int64(typed), 10), nil
	case int64:
		return strconv.FormatInt(typed, 10), nil
	case uint:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint64:
		return strconv.FormatUint(typed, 10), nil
	case float32:
		return strconv.FormatFloat(float64(typed), 'g', -1, 32), nil
	case float64:
		return strconv.FormatFloat(typed, 'g', -1, 64), nil
	case json.Number:
		return typed.String(), nil
	}

	valueKind := reflect.ValueOf(value).Kind()
	if valueKind == reflect.Map || valueKind == reflect.Slice || valueKind == reflect.Array {
		encoded, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(encoded), nil
	}

	return "", fmt.Errorf("unsupported value type %T", value)
}

func resourceCSKubernetesPolicyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		return errs.WrapError(clientErr)
	}
	csClient := CsClient{client}

	action := "ModifyPolicyInstance"
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
	}

	cluster_id := parts[0]
	policy_name := parts[1]
	instance_name := parts[2]

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return errs.WrapErrorf(err, "waiting for required addons failed")
	}

	updateReq := &roacs.ModifyPolicyInstanceRequest{
		InstanceName: tea.String(instance_name),
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesMapsArray := convert.ConvertToInterfaceArray(v)
		stringList := make([]*string, len(namespacesMapsArray))
		for i, v := range namespacesMapsArray {
			stringList[i] = tea.String(v.(string))
		}
		updateReq.Namespaces = stringList
	}

	if v, ok := d.GetOk("action"); ok {
		updateReq.Action = tea.String(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMap := v.(map[string]interface{})
		updateReq.Parameters = convert.NormalizeMap(parametersMap)
	}

	var response *roacs.ModifyPolicyInstanceResponse
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.ModifyPolicyInstance(tea.String(cluster_id), tea.String(policy_name), updateReq)
		if err != nil {
			if helper.NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return resourceCSKubernetesPolicyInstanceRead(d, meta)
}
func resourceCSKubernetesPolicyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client, clientErr := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if clientErr != nil {
		errs.WrapError(clientErr)
	}
	csClient := CsClient{client}

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		err := errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
		return err
	}

	cluster_id := parts[0]
	policy_name := parts[1]
	action := "DeletePolicyInstance"
	deleteReq := &roacs.DeletePolicyInstanceRequest{
		InstanceName: tea.String(parts[2]),
	}

	// Wait for required addons to be ready
	if err := csClient.WaitForRequiredAddons(cluster_id, 10*time.Minute); err != nil {
		return errs.WrapErrorf(err, "waiting for required addons failed")
	}

	var response *roacs.DeletePolicyInstanceResponse
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.DeletePolicyInstance(tea.String(cluster_id), tea.String(policy_name), deleteReq)
		if err != nil {
			if helper.NeedRetry(err) || strings.Contains(err.Error(), "the object has been modified") {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, deleteReq)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"404"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}

func (s *CsClient) WaitForRequiredAddons(clusterId string, timeout time.Duration) error {
	requiredAddons := []string{"gatekeeper", "policy-template-controller"}
	for _, addonName := range requiredAddons {
		addon, err := s.GetCsKubernetesAddonInstance(clusterId, addonName)
		if err != nil {
			if !errs.NotFoundError(err) {
				return errs.WrapErrorf(err, "checking addon %s status failed", addonName)
			}
			// Addon not found, need to wait for installation
		} else if addon.Status == "running" || addon.Status == "active" {
			// Addon is already active, no need to wait
			continue
		}

		// Addon exists but not active, wait for it to become active
		stateConf := helper.BuildStateConf([]string{}, []string{"running", "active"}, timeout, 10*time.Second, s.CsKubernetesAddonStateRefreshFunc(clusterId, addonName, []string{"failed", "error"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, "waiting for addon %s to be ready failed", addonName)
		}
	}
	return nil
}
