package cs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	aliyungocs "github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CsService struct {
	client *connectivity.AliyunClient
}

func NewCsService(client *connectivity.AliyunClient) *CsService {
	return &CsService{client: client}
}

type CsClient struct {
	client *client.Client
}

func NewCsClient(c *client.Client) *CsClient {
	return &CsClient{client: c}
}

type Component struct {
	ComponentName    string      `json:"component_name"`
	Version          string      `json:"version"`
	NextVersion      string      `json:"next_version"`
	CanUpgrade       bool        `json:"can_upgrade"`
	Required         bool        `json:"required"`
	Status           string      `json:"status"`
	ErrMessage       string      `json:"err_message"`
	Config           string      `json:"config"`
	ConfigSchema     string      `json:"config_schema"`
	Error            interface{} `json:"error"`
	SupportedActions []string    `json:"supported_actions"`
}

// Status mirrors the root package status type used by the WaitFor helpers.
// The root alicloud package cannot be imported from a service package, so the
// status type used by WaitForContainerApplication is redefined locally with the
// same underlying type, preserving the original behavior.
type Status string

const (
	COMPONENT_AUTO_SCALER      = "cluster-autoscaler"
	COMPONENT_DEFAULT_VRESION  = "v1.0.0"
	SCALING_CONFIGURATION_NAME = "kubernetes_autoscaler_autogen"
	DefaultECSTag              = "k8s.aliyun.com"
	DefaultClusterTag          = "ack.aliyun.com"
	CsPlayerAccountIdTag       = "ack.playeraccount"
	RECYCLE_MODE_LABEL         = "k8s.io/cluster-autoscaler/node-template/label/policy"
	DefaultAutoscalerTag       = "k8s.io/cluster-autoscaler"
	SCALING_GROUP_NAME         = "sg-%s-%s"
	DEFAULT_COOL_DOWN_TIME     = 300
	RELEASE_MODE               = "release"
	RECYCLE_MODE               = "recycle"

	PRIORITY_POLICY       = "PRIORITY"
	COST_OPTIMIZED_POLICY = "COST_OPTIMIZED"
	BALANCE_POLICY        = "BALANCE"

	UpgradeClusterTimeout = 30 * time.Minute
	DefaultTimeout        = 120

	// COLON_SEPARATED mirrors the root package separator constant (alicloud/common.go).
	// Redefined locally because service packages cannot import the root alicloud package;
	// value is identical, preserving the clusterId:addonName id-format behavior.
	COLON_SEPARATED = ":"

	IdMsgWithTask = errs.IdMsg + "TaskInfo: %s" // wait for async task info
)

// PayType mirrors the root package billing type (alicloud/common.go). Redefined
// locally because service packages cannot import the root alicloud package.
type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
)

// TimeType mirrors the root package period type (alicloud/common.go).
type TimeType string

const (
	Month = TimeType("Month")
)

// DiskCategory mirrors the root package disk-category type (alicloud/extension_ecs.go).
type DiskCategory string

const (
	DiskCloudEfficiency = DiskCategory("cloud_efficiency")
	DiskCloudESSD       = DiskCategory("cloud_essd")
	DiskCloudSSD        = DiskCategory("cloud_ssd")
)

// VpcService mirrors the root package VPC service (alicloud/service_alicloud_vpc.go).
// The vpc product leaf is not yet migrated, so the DescribeVSwitch method used by
// CS cluster resources is redefined locally, preserving the original behavior.
type VpcService struct {
	client *connectivity.AliyunClient
}

func NewVpcService(client *connectivity.AliyunClient) *VpcService {
	return &VpcService{client: client}
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

// KmsService mirrors the root package KMS service (alicloud/service_alicloud_kms.go).
// The kms product leaf is not yet migrated, so the Decrypt method used by CS cluster
// resources (for kms_key_id encryption context) is redefined locally.
type KmsService struct {
	client *connectivity.AliyunClient
}

func NewKmsService(client *connectivity.AliyunClient) *KmsService {
	return &KmsService{client: client}
}

func (s *KmsService) Decrypt(ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, errs.WrapError(err)
	}

	var response map[string]interface{}
	client := s.client
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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
		return plaintext, errs.WrapErrorf(err, errs.DefaultErrorMsg, context, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

// The following diff-suppress functions mirror the root package definitions
// (alicloud/diff_suppress_funcs.go); redefined locally because service packages
// cannot import the root alicloud package. Behavior is identical.
func csKubernetesMasterPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("master_instance_charge_type").(string) == "PostPaid" || !(d.Id() == "") && !d.Get("force_update").(bool)
}

func masterDiskPerformanceLevelDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("master_disk_category"); ok && v.(string) != "cloud_essd" {
		return true
	}
	return false
}

func workerDiskPerformanceLevelDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("worker_disk_category"); ok && v.(string) != "cloud_essd" {
		return true
	}
	return false
}

func csNodepoolDiskPerformanceLevelDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("system_disk_category"); ok && v.(string) != "cloud_essd" {
		return true
	}
	return false
}

func csNodepoolInstancePostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) == "PostPaid" {
		return true
	}
	return false
}

func csNodepoolScalingPolicyDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if _, ok := d.GetOk("scaling_config"); ok {
		return false
	}
	return true
}

var (
	ATTACH_SCRIPT_WITH_VERSION = `#!/bin/sh
curl http://aliacs-k8s-%s.oss-%s.aliyuncs.com/public/pkg/run/attach/%s/attach_node.sh | bash -s -- --openapi-token %s --ess true `
	NETWORK_ADDON_NAMES = []string{"terway", "kube-flannel-ds", "terway-eni", "terway-eniip"}
)

func (s *CsService) GetContainerClusterByName(name string) (cluster aliyungocs.ClusterType, err error) {
	name = helper.Trim(name)
	invoker := helper.NewInvoker()
	var clusters []aliyungocs.ClusterType
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			return csClient.DescribeClusters(name)
		})
		if e != nil {
			return e
		}
		clusters, _ = raw.([]aliyungocs.ClusterType)
		return nil
	})

	if err != nil {
		return cluster, fmt.Errorf("Describe cluster failed by name %s: %#v.", name, err)
	}

	if len(clusters) < 1 {
		return cluster, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Container Cluster", name))
	}

	for _, c := range clusters {
		if c.Name == name {
			return c, nil
		}
	}
	return cluster, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Container Cluster", name))
}

func (s *CsService) GetContainerClusterAndCertsByName(name string) (*aliyungocs.ClusterType, *aliyungocs.ClusterCerts, error) {
	cluster, err := s.GetContainerClusterByName(name)
	if err != nil {
		return nil, nil, err
	}
	var certs aliyungocs.ClusterCerts
	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			return csClient.GetClusterCerts(cluster.ClusterID)
		})
		if e != nil {
			return e
		}
		certs, _ = raw.(aliyungocs.ClusterCerts)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &cluster, &certs, nil
}

func (s *CsService) DescribeContainerApplication(clusterName, appName string) (app aliyungocs.GetProjectResponse, err error) {
	appName = helper.Trim(appName)
	cluster, certs, err := s.GetContainerClusterAndCertsByName(clusterName)
	if err != nil {
		return app, err
	}
	raw, err := s.client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *aliyungocs.ProjectClient) (interface{}, error) {
		return csProjectClient.GetProject(appName)
	})
	app, _ = raw.(aliyungocs.GetProjectResponse)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Not Found"}) {
			return app, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Container Application", appName))
		}
		return app, fmt.Errorf("Getting Application failed by name %s: %#v.", appName, err)
	}
	if app.Name != appName {
		return app, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Container Application", appName))
	}
	return
}

func (s *CsService) WaitForContainerApplication(clusterName, appName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := s.DescribeContainerApplication(clusterName, appName)
		if err != nil {
			return err
		}

		if strings.ToLower(app.CurrentState) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - helper.DefaultIntervalShort
		if timeout <= 0 {
			return errs.GetTimeErrorFromString(fmt.Sprintf("Waitting for container application %s is timeout and current status is %s.", string(status), app.CurrentState))
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CsService) DescribeCsKubernetes(id string) (cluster *aliyungocs.KubernetesClusterDetail, err error) {
	invoker := helper.NewInvoker()
	var requestInfo *aliyungocs.Client
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesClusterDetail(id)
		})
		response = raw
		return err
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, errs.WrapErrorf(err, errs.NotFoundMsg, errs.DenverdinoAliyungo)
		}
		return cluster, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeKubernetesCluster", errs.DenverdinoAliyungo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		logs.AddDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ = response.(*aliyungocs.KubernetesClusterDetail)
	if cluster.ClusterId != id {
		return cluster, errs.WrapErrorf(errs.NotFoundErr("CsKubernetes", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return
}

func (s *CsClient) DescribeClusterDetail(id string) (*client.DescribeClusterDetailResponseBody, error) {
	if id == "" {
		return nil, errs.WrapError(fmt.Errorf("cluster id is empty"))
	}

	var err error
	var response *client.DescribeClusterDetailResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = s.client.DescribeClusterDetail(tea.String(id))
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}

	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		logs.AddDebug("DescribeClusterDetail", response, id, requestMap)
	}
	return response.Body, nil
}

// DescribeClusterKubeConfigWithExpiration return cluster kube_config credential with expiration time.
// It's used for kubernetes/managed_kubernetes/serverless_kubernetes.
func (s *CsClient) DescribeClusterKubeConfigWithExpiration(clusterId string, temporaryDurationMinutes int64) (*client.DescribeClusterUserKubeconfigResponseBody, error) {
	if clusterId == "" {
		return nil, errs.WrapError(fmt.Errorf("clusterid is empty"))
	}

	request := &client.DescribeClusterUserKubeconfigRequest{
		PrivateIpAddress: tea.Bool(false),
	}
	if temporaryDurationMinutes > 0 {
		request.TemporaryDurationMinutes = tea.Int64(temporaryDurationMinutes)
	}
	var err error
	var kubeConfig *client.DescribeClusterUserKubeconfigResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		kubeConfig, err = s.client.DescribeClusterUserKubeconfig(tea.String(clusterId), request)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, errs.WrapError(err)
	}

	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = clusterId
		logs.AddDebug("DescribeClusterUserConfig", kubeConfig, request, requestMap)
	}

	return kubeConfig.Body, nil
}

// ListAddons returns all available addons info of the cluster
func (s *CsClient) ListAddons(clusterId string) (map[string]*client.ListAddonsResponseBodyAddons, error) {
	result := make(map[string]*client.ListAddonsResponseBodyAddons)

	request := &client.ListAddonsRequest{
		ClusterId: tea.String(clusterId),
	}
	var err error
	var resp *client.ListAddonsResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.ListAddons(request)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "ListAddons", err)
	}

	if resp == nil || resp.Body == nil || resp.Body.Addons == nil {
		return nil, errs.WrapErrorf(fmt.Errorf("ListAddons response body is nil"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "ListAddons", err)
	}

	for _, addon := range resp.Body.Addons {
		ca := addon
		result[tea.StringValue(addon.Name)] = ca
	}

	return result, nil
}

// Deprecated This function returns the latest addon status information
func (s *CsClient) DescribeCsKubernetesAddonStatus(clusterId string, addonName string) (*Component, error) {
	result := &Component{}
	var err error
	var resp *client.DescribeClusterAddonsUpgradeStatusResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.DescribeClusterAddonsUpgradeStatus(&clusterId, &client.DescribeClusterAddonsUpgradeStatusRequest{
			ComponentIds: []*string{tea.String(addonName)},
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeClusterAddonsUpgradeStatus", err)
	}

	addon, ok := resp.Body[addonName]
	if !ok {
		return nil, errs.WrapErrorf(errs.NotFoundErr("alicloud_cs_kubernetes_addon", addonName), errs.ResourceNotfound)
	}

	addonInfo := addon.(map[string]interface{})["addon_info"]
	tasks := addon.(map[string]interface{})["tasks"]
	result.Version = addonInfo.(map[string]interface{})["version"].(string)
	result.CanUpgrade = addon.(map[string]interface{})["can_upgrade"].(bool)
	result.Status = tasks.(map[string]interface{})["status"].(string)
	if tErr, ok := tasks.(map[string]interface{})["error"]; ok {
		result.Error = tErr
	}
	if message, ok := tasks.(map[string]interface{})["message"]; ok {
		result.ErrMessage = message.(string)
	}

	return result, nil
}

// DescribeUserPermission is used by the acceptance-test framework to verify
// that the permissions resource still exists.
func (s *CsClient) DescribeUserPermission(uid string) ([]*client.DescribeUserPermissionResponseBody, error) {
	return describeUserPermissions(s.client, uid)
}

// DescribeCsAutoscalingConfig is used by the acceptance-test framework to
// verify the autoscaling config exists. The resource ID is the cluster ID and
// the configuration is cluster-scoped, so cluster existence is the check.
func (s *CsClient) DescribeCsAutoscalingConfig(id string) (*client.DescribeClusterDetailResponseBody, error) {
	return s.DescribeClusterDetail(id)
}

// This function returns the latest addon instance
func (s *CsClient) GetCsKubernetesAddonInstance(clusterId string, addonName string) (*Component, error) {
	component := &Component{}
	var err error
	var resp *client.GetClusterAddonInstanceResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.GetClusterAddonInstance(&clusterId, tea.String(addonName))
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AddonNotFound"}) || errs.NotFoundError(err) {
			err = errs.WrapErrorf(errs.NotFoundErr("alicloud_cs_kubernetes_addon", addonName), errs.ResourceNotfound)
			return component, err
		}
		return nil, err
	}

	if resp.Body.Name != nil {
		component.ComponentName = *resp.Body.Name
	}
	if resp.Body.Version != nil {
		component.Version = *resp.Body.Version
	}
	if resp.Body.State != nil {
		component.Status = *resp.Body.State
	}

	if resp.Body.Config != nil {
		component.Config = *resp.Body.Config
	}

	return component, nil
}

// DescribeCsKubernetesAllAvailableAddons returns the status of all cluster addons, including both installed instances and available addons that can be installed.
func (s *CsClient) DescribeCsKubernetesAllAvailableAddons(clusterId string) (map[string]*Component, error) {
	// ListAddons
	availableAddons, err := s.ListAddons(clusterId)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Component)
	for name, addonInfo := range availableAddons {
		addon := &Component{}
		result[name] = addon

		addon.ComponentName = name
		addon.Required = tea.BoolValue(addonInfo.InstallByDefault)
		addon.SupportedActions = tea.StringSliceValue(addonInfo.SupportedActions)

		// GetClusterAddonInstance
		addonInstance, err := s.GetCsKubernetesAddonInstance(clusterId, name)
		if err != nil {
			if errs.NotFoundError(err) {
				// Update if addon not installed
				addon.NextVersion = tea.StringValue(addonInfo.Version)
				continue
			}
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesExistedAddons", err)
		}

		// DescribeAddon
		addonInfoDetail, err := s.DescribeAddon(clusterId, name, addonInstance.Version)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"AddonNotFound"}) {
				// If the specified version is not found, retry without version
				addonInfoDetail, err = s.DescribeAddon(clusterId, name, "")
				if err != nil {
					return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeAddon", err)
				}
			} else {
				return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeAddon", err)
			}
		}

		if addonInfoDetail == nil || addonInfoDetail.Body == nil {
			return nil, errs.WrapErrorf(fmt.Errorf("DescribeAddon response body is nil"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeAddon", err)
		}

		addon.NextVersion = getNextVersion(addonInfoDetail)
		// If retried without version, use the response version as next_version only if it's newer
		if addon.NextVersion == "" && addonInfoDetail.Body.Version != nil {
			targetVersion := tea.StringValue(addonInfoDetail.Body.Version)
			if cmp, err := versionCompare(addonInstance.Version, targetVersion); err == nil && cmp == 1 {
				addon.NextVersion = targetVersion
			}
		}

		// Update if addon installed
		addon.Version = addonInstance.Version
		addon.Config = addonInstance.Config
		addon.Status = addonInstance.Status
	}

	return result, nil
}

func (s *CsClient) DescribeAddon(clusterId, name, version string) (*client.DescribeAddonResponse, error) {
	var err error
	req := &client.DescribeAddonRequest{
		ClusterId: tea.String(clusterId),
		Version:   tea.String(version),
	}
	var resp *client.DescribeAddonResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.DescribeAddon(&name, req)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	return resp, err
}

func (s *CsClient) DescribeCsKubernetesAddon(id string) (*Component, error) {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	clusterId := parts[0]
	addonName := parts[1]

	addonInstance, err := s.GetCsKubernetesAddonInstance(clusterId, addonName)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AddonNotFound", "ErrorClusterNotFound"}) || errs.NotFoundError(err) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("alicloud_cs_kubernetes_addon", id), errs.ResourceNotfound)
		}
		return nil, err
	}

	addonInfo, err := s.DescribeAddon(clusterId, addonName, addonInstance.Version)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AddonNotFound"}) {
			// If the specified version is not found, retry without version
			addonInfo, err = s.DescribeAddon(clusterId, addonName, "")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if addonInfo == nil || addonInfo.Body == nil {
		return nil, errs.WrapErrorf(fmt.Errorf("DescribeAddon response body is nil"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeAddon", err)
	}

	addonInstance.Required = tea.BoolValue(addonInfo.Body.InstallByDefault)
	// compat for old behavior: if the addon cannot be upgraded, use the current version as the next version
	addonInstance.NextVersion = addonInstance.Version
	if nextVersion := getNextVersion(addonInfo); nextVersion != "" {
		addonInstance.NextVersion = nextVersion
	}
	// If retried without version, use the response version as next_version only if it's newer
	if addonInstance.NextVersion == addonInstance.Version && addonInfo.Body.Version != nil {
		targetVersion := tea.StringValue(addonInfo.Body.Version)
		if cmp, err := versionCompare(addonInstance.Version, targetVersion); err == nil && cmp == 1 {
			addonInstance.NextVersion = targetVersion
		}
	}
	addonInstance.CanUpgrade = addonInstance.Version != addonInstance.NextVersion
	addonInstance.SupportedActions = tea.StringSliceValue(addonInfo.Body.SupportedActions)

	return addonInstance, nil
}

func getNextVersion(addonInfo *client.DescribeAddonResponse) string {
	newerVersions := addonInfo.Body.NewerVersions
	for _, version := range newerVersions {
		if tea.BoolValue(version.Upgradable) {
			return *version.Version
		}
	}

	return ""
}

// versionCompare compares two semantic versions using github.com/Masterminds/semver.
// Supports formats: "1.2.3", "v1.2.3", "1.2.3-aliyun.1", "v1.31.0-apsara.6.11.6.ad796663", etc.
// Returns:
//
//	 1  if newVersion > oldVersion
//	 0  if newVersion == oldVersion
//	-1  if newVersion < oldVersion
//
// Example: versionCompare("1.20.11-aliyun.1", "1.22.3-aliyun.1") returns 1
func versionCompare(oldVersion, newVersion string) (int, error) {
	if oldVersion == "" || newVersion == "" {
		if oldVersion == "" && newVersion == "" {
			return 0, nil
		}
		if oldVersion == "" {
			return 1, nil
		}
		return -1, nil
	}

	oldSemver, err := semver.NewVersion(oldVersion)
	if err != nil {
		return -2, fmt.Errorf("failed to parse version %q: %v", oldVersion, err)
	}
	newSemver, err := semver.NewVersion(newVersion)
	if err != nil {
		return -2, fmt.Errorf("failed to parse version %q: %v", newVersion, err)
	}

	return newSemver.Compare(oldSemver), nil
}

func (s *CsClient) CsKubernetesAddonTaskRefreshFunc(clusterId string, addonName string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesAddonStatus(clusterId, addonName)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatusWithResponse, clusterId, object.ErrMessage))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CsClient) CsKubernetesAddonStateRefreshFunc(clusterId string, addonName string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.GetCsKubernetesAddonInstance(clusterId, addonName)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatusWithResponse, clusterId, object.ErrMessage))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CsClient) installAddon(d *schema.ResourceData) error {
	clusterId := d.Get("cluster_id").(string)

	body := make([]*client.InstallClusterAddonsRequestBody, 0)
	b := &client.InstallClusterAddonsRequestBody{
		Name:    tea.String(d.Get("name").(string)),
		Version: tea.String(d.Get("version").(string)),
	}

	if config, exist := d.GetOk("config"); exist {
		b.Config = tea.String(config.(string))
	}
	body = append(body, b)

	creationArgs := &client.InstallClusterAddonsRequest{
		Body: body,
	}

	var resp *client.InstallClusterAddonsResponse
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.InstallClusterAddons(&clusterId, creationArgs)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
	}

	if resp == nil || resp.Body == nil || tea.StringValue(resp.Body.TaskId) == "" {
		return errs.WrapErrorf(fmt.Errorf("failed to get taskId"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	stateConf := helper.BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, s.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.ResponseCodeMsg, d.Id(), "InstallClusterAddons", jobDetail)
	}

	return nil
}

func (s *CsClient) upgradeAddon(d *schema.ResourceData, updateVersion, updateConfig bool) error {
	clusterId := d.Get("cluster_id").(string)

	body := make([]*client.UpgradeClusterAddonsRequestBody, 0)
	b := &client.UpgradeClusterAddonsRequestBody{
		ComponentName: tea.String(d.Get("name").(string)),
	}
	if updateVersion {
		b.NextVersion = tea.String(d.Get("version").(string))
		b.Policy = tea.String("overwrite")
	}

	if updateConfig {
		if config, exist := d.GetOk("config"); exist {
			b.Config = tea.String(config.(string))
		}
	}

	body = append(body, b)

	upgradeArgs := &client.UpgradeClusterAddonsRequest{
		Body: body,
	}

	var resp *client.UpgradeClusterAddonsResponse
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.UpgradeClusterAddons(&clusterId, upgradeArgs)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
	}

	if resp == nil || resp.Body == nil || tea.StringValue(resp.Body.TaskId) == "" {
		return errs.WrapErrorf(fmt.Errorf("failed to get taskId"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	stateConf := helper.BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.ResponseCodeMsg, d.Id(), "UpgradeClusterAddons", jobDetail)
	}

	return nil
}

func (s *CsClient) uninstallAddon(d *schema.ResourceData) error {
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	clusterId := parts[0]

	body := make([]*client.UnInstallClusterAddonsRequestAddons, 0)
	b := &client.UnInstallClusterAddonsRequestAddons{
		Name: tea.String(parts[1]),
	}
	if v, ok := d.GetOk("cleanup_cloud_resources"); ok {
		b.CleanupCloudResources = tea.Bool(v.(bool))
	}
	body = append(body, b)

	uninstallArgs := &client.UnInstallClusterAddonsRequest{
		Addons: body,
	}

	var resp *client.UnInstallClusterAddonsResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err = s.client.UnInstallClusterAddons(&clusterId, uninstallArgs)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	if resp == nil || resp.Body == nil || tea.StringValue(resp.Body.TaskId) == "" {
		return errs.WrapErrorf(fmt.Errorf("failed to get taskId"), errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	taskId := tea.StringValue(resp.Body.TaskId)
	stateConf := helper.BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, s.DescribeTaskRefreshFunc(d, taskId, []string{"fail", "failed"}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.ResponseCodeMsg, d.Id(), "UnInstallClusterAddons", jobDetail)
	}

	return nil
}

func (s *CsClient) updateAddonConfig(d *schema.ResourceData) error {
	clusterId := d.Get("cluster_id").(string)
	ComponentName := d.Get("name").(string)

	upgradeArgs := &client.ModifyClusterAddonRequest{
		Config: tea.String(d.Get("config").(string)),
	}

	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = s.client.ModifyClusterAddon(&clusterId, &ComponentName, upgradeArgs)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddonConfig", err)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.CsKubernetesAddonTaskRefreshFunc(clusterId, ComponentName, []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		status, _ := s.DescribeCsKubernetesAddonStatus(clusterId, ComponentName)
		if status != nil && status.Error != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterUpdate", status.Error)
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterUpdate", d.Id())
	}

	return nil
}

// DescribeCsKubernetesAddonMetadata returns metadata of an addon with name and version
func (s *CsClient) DescribeCsKubernetesAddonMetadata(clusterId string, name string, version string) (*Component, error) {
	resp, err := s.DescribeAddon(clusterId, name, version)

	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesExistedAddons", err)
	}

	result := &Component{
		ComponentName: *resp.Body.Name,
		Version:       *resp.Body.Version,
		ConfigSchema:  *resp.Body.ConfigSchema,
	}

	return result, nil
}

func (s *CsService) DescribeCsKubernetesNodePool(id string) (nodePool *aliyungocs.NodePoolDetail, err error) {
	invoker := helper.NewInvoker()
	var requestInfo *aliyungocs.Client
	var response interface{}

	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	clusterId := parts[0]
	nodePoolId := parts[1]

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeNodePoolDetail(clusterId, nodePoolId)
		})
		response = raw
		return err
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound", "400"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.DenverdinoAliyungo)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, nodePoolId, "DescribeNodePool", errs.DenverdinoAliyungo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = clusterId
		requestMap["NodePoolId"] = nodePoolId
		logs.AddDebug("DescribeNodepool", response, requestInfo, requestMap)
	}
	nodePool, _ = response.(*aliyungocs.NodePoolDetail)
	if nodePool.NodePoolId != nodePoolId {
		return nil, errs.WrapErrorf(errs.NotFoundErr("CsNodePool", nodePoolId), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return
}

func (s *CsService) DescribeCsManagedKubernetes(id string) (cluster *aliyungocs.KubernetesClusterDetail, err error) {
	var requestInfo *aliyungocs.Client
	invoker := helper.NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesClusterDetail(id)
		})
		response = raw
		return err
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return cluster, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeKubernetesCluster", errs.DenverdinoAliyungo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		logs.AddDebug("DescribeKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(*aliyungocs.KubernetesClusterDetail)
	if cluster.ClusterId != id {
		return cluster, errs.WrapErrorf(errs.NotFoundErr("CSManagedKubernetes", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return

}

func (s *CsService) CsKubernetesInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsKubernetesNodePoolStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesNodePool(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsServerlessKubernetesInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsServerlessKubernetes(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) DescribeCsServerlessKubernetes(id string) (*aliyungocs.ServerlessClusterResponse, error) {
	cluster := &aliyungocs.ServerlessClusterResponse{}
	var requestInfo *aliyungocs.Client
	invoker := helper.NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeServerlessKubernetesCluster(id)
		})
		response = raw
		return err
	}); err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, errs.WrapErrorf(err, errs.NotFoundMsg, errs.DenverdinoAliyungo)
		}
		return cluster, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "DescribeServerlessKubernetesCluster", errs.DenverdinoAliyungo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		logs.AddDebug("DescribeServerlessKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(*aliyungocs.ServerlessClusterResponse)
	if cluster != nil && cluster.ClusterId != id {
		return cluster, errs.WrapErrorf(errs.NotFoundErr("CSServerlessKubernetes", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return cluster, nil

}

func (s *CsService) tagsToMap(tags []aliyungocs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *CsService) ignoreTag(t aliyungocs.Tag) bool {
	filter := []string{"^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *CsService) GetPermanentToken(clusterId string) (string, error) {

	describeClusterTokensResponse, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
		return csClient.DescribeClusterTokens(clusterId)
	})
	if err != nil {
		return "", errs.WrapError(fmt.Errorf("failed to get permanent token,because of %v", err))
	}

	tokens, ok := describeClusterTokensResponse.([]*aliyungocs.ClusterTokenResponse)

	if ok != true {
		return "", errs.WrapError(fmt.Errorf("failed to parse ClusterTokenResponse of cluster %s", clusterId))
	}

	permanentTokens := make([]string, 0)

	for _, token := range tokens {
		if token.Expired == 0 && token.IsActive == 1 {
			permanentTokens = append(permanentTokens, token.Token)
			break
		}
	}

	// create a new token
	if len(permanentTokens) == 0 {
		createClusterTokenResponse, err := s.client.WithCsClient(func(csClient *aliyungocs.Client) (interface{}, error) {
			clusterTokenReqeust := &aliyungocs.ClusterTokenReqeust{}
			clusterTokenReqeust.IsPermanently = true
			return csClient.CreateClusterToken(clusterId, clusterTokenReqeust)
		})
		if err != nil {
			return "", errs.WrapError(fmt.Errorf("failed to create permanent token,because of %v", err))
		}

		token, ok := createClusterTokenResponse.(*aliyungocs.ClusterTokenResponse)
		if ok != true {
			return "", errs.WrapError(fmt.Errorf("failed to parse token of %s", clusterId))
		}
		return token.Token, nil
	}

	return permanentTokens[0], nil
}

// GetUserData of cluster
func (s *CsService) GetUserData(clusterId string, labels string, taints string) (string, error) {

	token, err := s.GetPermanentToken(clusterId)

	if err != nil {
		return "", err
	}

	if labels == "" {
		labels = fmt.Sprintf("%s=true", DefaultECSTag)
	} else {
		labels = fmt.Sprintf("%s,%s=true", labels, DefaultECSTag)
	}

	cluster, err := s.DescribeCsKubernetes(clusterId)

	if err != nil {
		return "", errs.WrapError(fmt.Errorf("failed to describe cs kuberentes cluster,because of %v", err))
	}

	extra_options := make([]string, 0)

	if len(labels) > 0 || len(taints) > 0 {

		if len(labels) != 0 {
			extra_options = append(extra_options, fmt.Sprintf("--labels %s", labels))
		}

		if len(taints) != 0 {
			extra_options = append(extra_options, fmt.Sprintf("--taints %s", taints))
		}
	}

	if network, err := GetKubernetesNetworkName(cluster); err == nil && network != "" {
		extra_options = append(extra_options, fmt.Sprintf("--network %s", network))
	}

	extra_options_in_line := strings.Join(extra_options, " ")

	version := cluster.CurrentVersion
	region := cluster.RegionId

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(ATTACH_SCRIPT_WITH_VERSION+extra_options_in_line, region, region, version, token))), nil
}

func (s *CsClient) DescribeTaskRefreshFunc(d *schema.ResourceData, taskId string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		var taskInfo *client.DescribeTaskInfoResponse
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			taskInfo, err = s.client.DescribeTaskInfo(tea.String(taskId))
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if errs.NotFoundError(err) {
				return taskInfo, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		currentState := tea.StringValue(taskInfo.Body.State)
		for _, failState := range failStates {
			if currentState == failState {
				if taskInfo.Body.Error != nil {
					return taskInfo.Body.Error, currentState, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentState))
				}

				return taskInfo, currentState, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentState))
			}
		}
		return taskInfo, currentState, nil
	}
}

func GetKubernetesNetworkName(cluster *aliyungocs.KubernetesClusterDetail) (network string, err error) {

	metadata := make(map[string]interface{})
	if err := json.Unmarshal([]byte(cluster.MetaData), &metadata); err != nil {
		return "", fmt.Errorf("unmarshal metaData failed. error: %s", err)
	}

	for _, name := range NETWORK_ADDON_NAMES {
		if _, ok := metadata[fmt.Sprintf("%s%s", name, "Version")]; ok {
			return name, nil
		}
	}
	return "", fmt.Errorf("no network addon found")
}

// setCertificateAuthorityAttribute must be false for resources whose schema
// does not declare certificate_authority. Terraform Plugin SDK v2 rejects
// writes to undeclared state addresses.
func setCerts(d *schema.ResourceData, meta interface{}, skipSetCertificateAuthority bool, setCertificateAuthorityAttribute bool) error {
	client := meta.(*connectivity.AliyunClient)
	roaClient, err := client.NewRoaCsClient()
	if err != nil {
		return errs.WrapError(err)
	}
	csClient := CsClient{roaClient}
	kubeConfig, err := csClient.DescribeClusterKubeConfigWithExpiration(d.Id(), 0)
	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to get kubeconfig for cluster %s: %s", d.Id(), err))
	}
	m, err := FlattenAlicloudCSCertificate(kubeConfig)
	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to parse kubeconfig for cluster %s: %s", d.Id(), err))
	}
	if len(m) >= 3 {
		if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
			if err := helper.WriteToFile(ce.(string), m["client_cert"]); err != nil {
				return errs.WrapError(err)
			}
		}
		if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
			if err := helper.WriteToFile(key.(string), m["client_key"]); err != nil {
				return errs.WrapError(err)
			}
		}
		if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
			if err := helper.WriteToFile(ca.(string), m["cluster_cert"]); err != nil {
				return errs.WrapError(err)
			}
		}
	}
	// kube_config
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		if err := helper.WriteToFile(file.(string), tea.StringValue(kubeConfig.Config)); err != nil {
			return errs.WrapError(err)
		}
	}

	if setCertificateAuthorityAttribute {
		if skipSetCertificateAuthority {
			if err := d.Set("certificate_authority", []interface{}{
				map[string]string{
					"cluster_cert": "",
					"client_cert":  "",
					"client_key":   "",
				},
			}); err != nil {
				return errs.WrapError(fmt.Errorf("error setting certificate_authority: %s", err))
			}
		} else {
			if err := d.Set("certificate_authority", []interface{}{m}); err != nil {
				return errs.WrapError(fmt.Errorf("error setting certificate_authority: %s", err))
			}
		}
	}

	return nil
}

type AckServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewAckServiceV2(client *connectivity.AliyunClient) *AckServiceV2 {
	return &AckServiceV2{client: client}
}

// DescribeAckNodepool <<< Encapsulated get interface for Ack Nodepool.

func (s *AckServiceV2) DescribeAckNodepool(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	ClusterId := parts[0]
	NodepoolId := parts[1]
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)

	action := fmt.Sprintf("/clusters/%s/nodepools/%s", ClusterId, NodepoolId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"ErrorNodePoolNotFound", "ErrorClusterNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Nodepool", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AckServiceV2) AckNodepoolStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AckNodepoolStateRefreshFuncWithApi(id, field, failStates, s.DescribeAckNodepool)
}

func (s *AckServiceV2) AckNodepoolStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *AckServiceV2) DescribeAsyncAckNodepoolStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTaskInfo(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAckNodepool >>> Encapsulated.

// DescribeAsyncDescribeTaskInfo <<< Encapsulated for Ack.
func (s *AckServiceV2) DescribeAsyncDescribeTaskInfo(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)
	task_id, err := jsonpath.Get("$.task_id", res)

	action := fmt.Sprintf("/tasks/%s", task_id)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, header, nil)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeTaskInfo >>> Encapsulated.

// DescribeAckPolicyInstance <<< Encapsulated get interface for Ack PolicyInstance.
func (s *AckServiceV2) DescribeAckPolicyInstance(id string) (object map[string]interface{}, err error) {
	aliyunClient := s.client
	// Using the official SDK method DescribePolicyInstances
	csClient, clientErr := aliyunClient.NewRoaCsClient()
	if clientErr != nil {
		return object, errs.WrapError(clientErr)
	}

	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 3, len(parts)))
		return object, err
	}

	cluster_id := parts[0]
	describeReq := &client.DescribePolicyInstancesRequest{
		PolicyName:   tea.String(parts[1]),
		InstanceName: tea.String(parts[2]),
	}
	action := "DescribePolicyInstances"
	var response map[string]interface{}
	var rawResponse interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)

	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		// Call the official DescribePolicyInstances method
		var resp *client.DescribePolicyInstancesResponse
		resp, err = csClient.DescribePolicyInstances(tea.String(cluster_id), describeReq)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		// Convert response to map[string]interface{} for compatibility
		respMap := make(map[string]interface{})
		if resp.Body != nil {
			byt, _ := json.Marshal(resp.Body)
			d := json.NewDecoder(strings.NewReader(string(byt)))
			d.UseNumber()
			d.Decode(&respMap)
		}
		response = respMap
		rawResponse = resp.Body

		return nil
	})
	logs.AddDebug(action, response, describeReq)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidResource.NotFound", "404"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PolicyInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	// Handle the case where response itself is an array (not wrapped in body)
	// This happens when the API returns a direct array instead of an object with a body field
	if rawResponse != nil {
		// Try to marshal and unmarshal to handle different response formats
		b, _ := json.Marshal(rawResponse)
		var directArray []interface{}
		if json.Unmarshal(b, &directArray) == nil && len(directArray) > 0 {
			if item, ok := directArray[0].(map[string]interface{}); ok {
				return item, nil
			}
		}
	}

	// Handle the case where response is an array wrapped in body field
	if bodyArray, ok := response["body"].([]interface{}); ok && len(bodyArray) > 0 {
		// Return first element if it's a map
		if firstElement, ok := bodyArray[0].(map[string]interface{}); ok {
			return firstElement, nil
		}
		return object, errs.WrapErrorf(fmt.Errorf("unexpected response format: first element is not a map"), errs.FailedGetAttributeMsg, id, "$[*]", response)
	}

	// Handle the case where response is a map with items in body field
	if bodyMap, ok := response["body"].(map[string]interface{}); ok {
		if itemsArray, ok := bodyMap["items"].([]interface{}); ok && len(itemsArray) > 0 {
			if firstElement, ok := itemsArray[0].(map[string]interface{}); ok {
				return firstElement, nil
			}
		}
	}

	// Fallback to original logic with jsonpath for backward compatibility
	v, err := jsonpath.Get("$[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$[*]", response)
	}

	// Check if the result is a slice and has at least one element
	resultSlice, ok := v.([]interface{})
	if !ok || len(resultSlice) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PolicyInstance", id), errs.NotFoundMsg, response)
	}

	// Check if the first element is a map
	firstElement, ok := resultSlice[0].(map[string]interface{})
	if !ok {
		return object, errs.WrapErrorf(fmt.Errorf("unexpected response format: first element is not a map"), errs.FailedGetAttributeMsg, id, "$[*]", response)
	}

	return firstElement, nil
}

// DescribeAckPolicyInstance >>> Encapsulated.

// DescribeAckCluster <<< Encapsulated get interface for Ack Cluster.

func (s *AckServiceV2) DescribeAckCluster(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	ClusterId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)

	action := fmt.Sprintf("/clusters/%s", ClusterId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("CS", "2015-12-15", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *AckServiceV2) AckClusterStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.AckClusterStateRefreshFuncWithApi(id, field, failStates, s.DescribeAckCluster)
}

func (s *AckServiceV2) AckClusterStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *AckServiceV2) DescribeAsyncAckClusterStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeTaskInfo(d, res)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeAckCluster >>> Encapsulated.
