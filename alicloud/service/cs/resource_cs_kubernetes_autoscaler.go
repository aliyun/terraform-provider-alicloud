package cs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	ClusterAutoscaler          = "cluster-autoscaler"
	ClusterAutoscalerMeta      = "autoscaler-meta"
	DefaultAutoscalerNamespace = "kube-system"
	defaultScalingGroupTag     = "k8s.aliyun.com"
	defaultServiceAccountName  = "admin"
	defaultAutoscalerImage     = "registry-vpc.%s.aliyuncs.com/acs/autoscaler:v1.3.1-7369cf1"
	LabelPattern               = "k8s.io/cluster-autoscaler/node-template/label/"
	TaintPattern               = "k8s.io/cluster-autoscaler/node-template/taint/"
)

// nodePool defines the struct of scaling group params
type nodePool map[string]string

// userKubeConf defines the struct of response of api
type userKubeConf struct {
	Config string `json:"config"`
}

// AutoscalerMeta define the struct of autoscaler meta configmap
type AutoscalerMeta struct {
	UnneededDuration        string                          `json:"unneeded_duration"`
	CoolDownDuration        string                          `json:"cool_down_duration"`
	UtilizationThreshold    string                          `json:"utilization_threshold"`
	GpuUtilizationThreshold string                          `json:"gpu_utilization_threshold"`
	ScalingConfigurations   map[string]scalingConfiguration `json:"scaling_configurations"`
}

// scalingConfiguration define the config of scaling group
type scalingConfiguration struct {
	Id     string `json:"id"`
	Labels string `json:"labels"`
	Taints string `json:"taints"`
}

// autoscalerOptions define the options used by autoscaler creation
type autoscalerOptions struct {
	Args               []string
	RegionId           string
	AccessKeyId        string
	AccessKeySecret    string
	Image              string
	UseEcsRamRoleToken bool
}

// ResourceCSKubernetesAutoscaler defines the schema of resource
func ResourceCSKubernetesAutoscaler() *schema.Resource {
	return &schema.Resource{
		Create: resourceCSKubernetesAutoscalerCreate,
		Read:   resourceCSKubernetesAutoscalerRead,
		Update: resourceCSKubernetesAutoscalerUpdate,
		Delete: resourceCSKubernetesAutoscalerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.127.0 and will be removed in the future. Please use 'alicloud_cs_autoscaling_config' instead.",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodepools": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"taints": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"labels": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// TODO add min max support
						//"min": {
						//	Type:     schema.TypeInt,
						//	Optional: true,
						//},
						//"max": {
						//	Type:     schema.TypeInt,
						//	Optional: true,
						//},
					},
				},
				MaxItems: 30,
				MinItems: 1,
			},
			"utilization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cool_down_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"defer_scale_in_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"use_ecs_ram_role_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// resourceCSKubernetesAutoscalerCreate define how to create autoscaler
func resourceCSKubernetesAutoscalerCreate(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	// set unique id of tf state
	d.SetId(fmt.Sprintf("%s:%s", clusterId, ClusterAutoscaler))
	return resourceCSKubernetesAutoscalerUpdate(d, meta)
}

// resourceCSKubernetesAutoscalerRead no need to implement
func resourceCSKubernetesAutoscalerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceCSKubernetesAutoscalerUpdate define how to update autoscaler configuration
func resourceCSKubernetesAutoscalerUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	// any changes need to ready autoscaler.
	if d.HasChange("nodepools") || d.HasChange("utilization") || d.HasChange("cool_down_duration") || d.HasChange("defer_scale_in_duration") || d.HasChange("use_ecs_ram_role_token") {

		regionId := client.RegionId

		clusterId := d.Get("cluster_id").(string)
		utilization := d.Get("utilization").(string)
		coolDownDuration := d.Get("cool_down_duration").(string)
		deferScaleInDuration := d.Get("defer_scale_in_duration").(string)
		useEcsRamroleToken := d.Get("use_ecs_ram_role_token").(bool)

		// parse nodepools
		nodePoolsParams := d.Get("nodepools").(*schema.Set)
		nodePools := nodePoolsParams.List()

		// args creation
		args := make([]string, 0)
		args = applyDefaultArgs(args)
		configmapMeta := AutoscalerMeta{
			UnneededDuration:        deferScaleInDuration,
			CoolDownDuration:        coolDownDuration,
			UtilizationThreshold:    utilization,
			GpuUtilizationThreshold: utilization,
			ScalingConfigurations:   make(map[string]scalingConfiguration),
		}

		for _, pool := range nodePools {
			poolBytes, err := json.Marshal(pool)
			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to marshal pool,because of %v", err))
			}

			pool := make(nodePool)

			err = json.Unmarshal(poolBytes, &pool)

			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to unmarshal pool,because of %v", err))
			}

			// get params of node pool
			id := pool["id"]
			labels := pool["labels"]
			taints := pool["taints"]

			// get userData from cluster openapi
			userData, err := csService.GetUserData(clusterId, labels, taints)

			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to get permanent token,because of %v", err))
			}

			err = UpdateScalingGroupConfiguration(client, id, userData, labels, taints)
			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to update scaling group status,because of %v", err))
			}

			// get min max of scaling group
			min, max, err := GetScalingGroupSizeRange(client, id)

			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to describe scaling group %s,because of %v", id, err))
			}

			nodeArgs := fmt.Sprintf("--nodes=%d:%d:%s", min, max, id)

			args = append(args, nodeArgs)
			configmapMeta.ScalingConfigurations[id] = scalingConfiguration{
				Id:     id,
				Labels: labels,
				Taints: taints,
			}
		}

		if utilization != "" {
			args = append(args, fmt.Sprintf("--scale-down-utilization-threshold=%s", utilization))
			args = append(args, fmt.Sprintf("--scale-down-gpu-utilization-threshold=%s", utilization))
		}

		if coolDownDuration != "" {
			args = append(args, fmt.Sprintf("--scale-down-delay-after-add=%s", coolDownDuration))
			args = append(args, fmt.Sprintf("--scale-down-delay-after-failure=%s", coolDownDuration))
		}

		if deferScaleInDuration != "" {
			args = append(args, fmt.Sprintf("--scale-down-unneeded-time=%s", deferScaleInDuration))
		}

		kubeConfPath, err := DownloadUserKubeConf(client, clusterId)

		if err != nil {
			return errs.WrapError(fmt.Errorf("failed to download kubeconf from cluster,because of %v", err))
		}

		options := autoscalerOptions{
			Args:               args,
			Image:              fmt.Sprintf(defaultAutoscalerImage, regionId),
			RegionId:           regionId,
			AccessKeyId:        client.AccessKey,
			AccessKeySecret:    client.SecretKey,
			UseEcsRamRoleToken: useEcsRamroleToken,
		}

		clientSet, err := GetClientSetFromKubeconf(kubeConfPath)

		if err != nil {
			return errs.WrapError(fmt.Errorf("failed to create kubernetes client,because of %v", err))
		}

		err = DeployAutoscaler(options, clientSet)

		if err != nil {
			return errs.WrapError(fmt.Errorf("failed to deploy autoscaler,because of %v", err))
		}

		err = createOrUpdateAutoscalerMeta(clientSet, configmapMeta)

		if err != nil {
			return errs.WrapError(fmt.Errorf("failed to update autoscaler meta configmap,because of %v", err))
		}
	}
	return nil
}

// delete autoscaler deployment
func resourceCSKubernetesAutoscalerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	regionId := client.RegionId

	clusterId := d.Get("cluster_id").(string)

	if clusterId == "" {
		return errs.WrapError(fmt.Errorf("please provide the cluster_id in region %s", regionId))
	}

	kubeConfPath, err := DownloadUserKubeConf(client, clusterId)

	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to download kubeconf from cluster,because of %v", err))
	}

	return DeleteAutoscaler(kubeConfPath)
}

// update scaling group config
func UpdateScalingGroupConfiguration(client *connectivity.AliyunClient, groupId, userData string, labels string, taints string) (err error) {

	describeScalingConfigurationsRequest := ess.CreateDescribeScalingConfigurationsRequest()
	describeScalingConfigurationsResponse, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		describeScalingConfigurationsRequest.RegionId = client.RegionId
		describeScalingConfigurationsRequest.ScalingGroupId = groupId
		return essClient.DescribeScalingConfigurations(describeScalingConfigurationsRequest)
	})

	configurations, ok := describeScalingConfigurationsResponse.(*ess.DescribeScalingConfigurationsResponse)

	if ok != true {
		return errs.WrapError(fmt.Errorf("failed to parse DescribeScalingConfigurationsResponse of %s", groupId))
	}

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, groupId, describeScalingConfigurationsRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	logs.AddDebug("DescribeScalingConfigurations", configurations, describeScalingConfigurationsRequest)

	if len(configurations.ScalingConfigurations.ScalingConfiguration) == 0 {
		//todo  create configuration
		return errs.WrapError(fmt.Errorf("please create the default scaling configuration of group %s", groupId))
	} else {
		defaultConfiguration := configurations.ScalingConfigurations.ScalingConfiguration[0]
		// modify the default one
		modifyScalingConfigurationRequest := ess.CreateModifyScalingConfigurationRequest()
		modifyScalingConfigurationResponse, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			modifyScalingConfigurationRequest.RegionId = client.RegionId
			modifyScalingConfigurationRequest.UserData = userData
			modifyScalingConfigurationRequest.Tags = createScalingGroupTags(labels, taints)
			modifyScalingConfigurationRequest.ScalingConfigurationId = defaultConfiguration.ScalingConfigurationId
			return essClient.ModifyScalingConfiguration(modifyScalingConfigurationRequest)
		})

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, groupId, modifyScalingConfigurationRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		logs.AddDebug("ModifyScalingConfiguration", modifyScalingConfigurationResponse, modifyScalingConfigurationRequest)
	}
	return nil
}

func GetScalingGroupSizeRange(client *connectivity.AliyunClient, groupId string) (min, max int, err error) {
	describeScalingGroupRequest := ess.CreateDescribeScalingGroupsRequest()
	describeScalingGroupRequest.RegionId = client.RegionId
	describeScalingGroupRequest.ScalingGroupId = &[]string{groupId}

	describeScalingGroupResponse, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(describeScalingGroupRequest)
	})

	if err != nil {
		return 0, 0, errs.WrapErrorf(err, errs.DefaultErrorMsg, groupId, describeScalingGroupRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug("DescribeScalingGroups", describeScalingGroupResponse, describeScalingGroupRequest)

	resp, ok := describeScalingGroupResponse.(*ess.DescribeScalingGroupsResponse)

	if ok != true {
		return 0, 0, errs.WrapError(fmt.Errorf("failed to parse DescribeScalingGroupsResponse of scaling group %s", groupId))
	}

	if resp.ScalingGroups.ScalingGroup == nil || len(resp.ScalingGroups.ScalingGroup) == 0 {
		return 0, 0, errs.WrapError(fmt.Errorf("the scaling group %s you specific is not found", groupId))
	}

	scalingGroup := resp.ScalingGroups.ScalingGroup[0]

	if &scalingGroup == nil {
		return 0, 0, errs.WrapError(fmt.Errorf("the scaling group %s you specific is not found", groupId))
	}

	return scalingGroup.MinSize, scalingGroup.MaxSize, nil
}

// prepare kubeconf of kubernetes clsuter
func DownloadUserKubeConf(client *connectivity.AliyunClient, clusterId string) (string, error) {

	describeClusterUserKubeconfigRequest := cs.CreateDescribeClusterUserKubeconfigRequest()
	describeClusterUserKubeconfigResponse, err := client.WithOfficalCSClient(func(csClient *cs.Client) (interface{}, error) {
		describeClusterUserKubeconfigRequest.RegionId = client.RegionId
		describeClusterUserKubeconfigRequest.ClusterId = clusterId
		return csClient.DescribeClusterUserKubeconfig(describeClusterUserKubeconfigRequest)
	})

	if err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, clusterId, describeClusterUserKubeconfigRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	kubeConfResponse, ok := describeClusterUserKubeconfigResponse.(*cs.DescribeClusterUserKubeconfigResponse)

	if ok != true {
		return "", errs.WrapError(fmt.Errorf("failed to parse DescribeClusterUserKubeconfigResponse of %s", clusterId))
	}
	logs.AddDebug("DescribeClusterUserKubeconfig", describeClusterUserKubeconfigResponse, describeClusterUserKubeconfigRequest)

	// get response bytes
	kubeconfBytes := kubeConfResponse.GetHttpContentBytes()

	ukc := &userKubeConf{}
	err = json.Unmarshal(kubeconfBytes, ukc)

	if err != nil {
		return "", errs.WrapError(fmt.Errorf("failed to parse DescribeClusterUserKubeconfigResponse,because of %v", err))
	}

	content := ukc.Config

	wd, err := os.Getwd()

	if err != nil {
		return "", errs.WrapError(fmt.Errorf("failed to get current working dir,because of %v", err))
	}

	kubeConfPath := path.Join(wd, fmt.Sprintf("%s-kubeconf", clusterId))

	err = ioutil.WriteFile(kubeConfPath, []byte(content), 0755)

	if err != nil {
		return "", errs.WrapError(fmt.Errorf("failed to create kubeconf in working dir because of %v", err))
	}

	return kubeConfPath, nil
}

// delete autoscaler component
func DeleteAutoscaler(kubeconf string) error {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)

	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to build kubeconf from local path %s,because of %v", kubeconf, err))
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to create client-go clientSet,because of %v", err))
	}

	err = clientSet.AppsV1().Deployments(DefaultAutoscalerNamespace).Delete(ctx, ClusterAutoscaler, metav1.DeleteOptions{})

	if errors.IsNotFound(err) == true {
		return nil
	}

	return errs.WrapError(err)
}

// deploy cluster-autoscaler to kubernetes cluster
func DeployAutoscaler(options autoscalerOptions, clientSet *kubernetes.Clientset) error {
	ctx := context.Background()
	deploy, err := clientSet.AppsV1().Deployments(DefaultAutoscalerNamespace).Get(ctx, ClusterAutoscaler, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			ak := options.AccessKeyId
			sk := options.AccessKeySecret
			if options.UseEcsRamRoleToken {
				ak = ""
				sk = ""
			}
			// create a new deploy
			deployObject := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: ClusterAutoscaler,
					Labels: map[string]string{
						"app": ClusterAutoscaler,
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": ClusterAutoscaler,
						},
					},

					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": ClusterAutoscaler,
							},
						},
						Spec: v1.PodSpec{
							ServiceAccountName: defaultServiceAccountName,
							Affinity: &v1.Affinity{
								NodeAffinity: &v1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
										NodeSelectorTerms: []v1.NodeSelectorTerm{
											{
												MatchExpressions: []v1.NodeSelectorRequirement{
													{
														Key:      defaultScalingGroupTag,
														Operator: v1.NodeSelectorOpNotIn,
														Values: []string{
															"true",
														},
													},
												},
											},
										},
									},
								},
							},
							Containers: []v1.Container{
								{
									Name:    ClusterAutoscaler,
									Image:   options.Image,
									Command: options.Args,
									Env: []v1.EnvVar{
										{
											Name:  "REGION_ID",
											Value: options.RegionId,
										},
										{
											Name:  "ACCESS_KEY_ID",
											Value: ak,
										},
										{
											Name:  "ACCESS_KEY_SECRET",
											Value: sk,
										},
									},
								},
							},
						},
					},
				},
			}
			_, err := clientSet.AppsV1().Deployments(DefaultAutoscalerNamespace).Create(ctx, deployObject, metav1.CreateOptions{})
			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to create %s deployment,because of %v", ClusterAutoscaler, err))
			}

		} else {
			return errs.WrapError(fmt.Errorf("failed to describe %s deployment,because of %v", ClusterAutoscaler, err))
		}
	} else {
		// update deployment
		deploy.Spec.Template.Spec.Containers[0].Command = options.Args
		_, err := clientSet.AppsV1().Deployments(DefaultAutoscalerNamespace).Update(ctx, deploy, metav1.UpdateOptions{})
		if err != nil {
			return errs.WrapError(fmt.Errorf("failed to update %s deployment,because of %v", ClusterAutoscaler, err))
		}
	}
	return nil
}

// set replicas
func int32Ptr(i int32) *int32 { return &i }

// apply default params of autoscaler command
func applyDefaultArgs(args []string) []string {
	args = append(args, "./cluster-autoscaler")
	args = append(args, "--v=5")
	args = append(args, "--stderrthreshold=info")
	args = append(args, "--cloud-provider=alicloud")
	args = append(args, "--scan-interval=30s")
	args = append(args, "--ok-total-unready-count=1000")
	args = append(args, "--max-empty-bulk-delete=50")
	args = append(args, "--expander=least-waste")
	args = append(args, "--leader-elect=false")
	args = append(args, "--skip-nodes-with-local-storage=false")
	return args
}

// convert labels and taints to scaling group tags
func createScalingGroupTags(labels string, taints string) string {

	tags := make(map[string]string)
	tags[defaultScalingGroupTag] = "true"

	labelArr := strings.Split(labels, ",")

	for _, label := range labelArr {
		labelKeyValue := strings.Split(label, "=")
		if len(labelKeyValue) == 2 {
			tags[fmt.Sprintf("%s%s", LabelPattern, labelKeyValue[0])] = labelKeyValue[1]
		}
	}

	taintsArr := strings.Split(taints, ",")
	for _, taint := range taintsArr {
		taintKeyValue := strings.Split(taint, "=")
		if len(taintKeyValue) == 2 {
			tags[fmt.Sprintf("%s%s", TaintPattern, taintKeyValue[0])] = taintKeyValue[1]
		}
	}
	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return ""
	}
	return string(tagsBytes)
}

// GetClientSetFromKubeconf return the clientSet from kubeconf
func GetClientSetFromKubeconf(kubeconf string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)
	if err != nil {
		return nil, errs.WrapError(fmt.Errorf("failed to build kubeconf from local path %s,because of %v", kubeconf, err))
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, errs.WrapError(fmt.Errorf("failed to create client-go clientSet,because of %v", err))
	}
	return clientSet, nil
}

// update autoscaler meta configmap
func createOrUpdateAutoscalerMeta(clientSet *kubernetes.Clientset, meta AutoscalerMeta) error {
	ctx := context.Background()
	meta_bytes, err := json.Marshal(meta)
	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to marshal autoscaler meta,because of %v", err))
	}
	meta_map := make(map[string]string)
	meta_map[ClusterAutoscalerMeta] = string(meta_bytes)
	cm := &v1.ConfigMap{
		Data: meta_map,
	}
	cm.Name = ClusterAutoscalerMeta

	configmap, err := clientSet.CoreV1().ConfigMaps(DefaultAutoscalerNamespace).Get(ctx, ClusterAutoscalerMeta, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = clientSet.CoreV1().ConfigMaps(DefaultAutoscalerNamespace).Create(ctx, cm, metav1.CreateOptions{})
			if err != nil {
				return errs.WrapError(fmt.Errorf("failed to create configmap of autoscaler meta,because of %v", err))
			}
			return nil
		} else {
			// return errror
			return errs.WrapError(fmt.Errorf("failed to describe configmap autoscaler meta,because of %v", err))
		}
	}
	configmap.Data = meta_map
	// update configmapa
	_, err = clientSet.CoreV1().ConfigMaps(DefaultAutoscalerNamespace).Update(ctx, configmap, metav1.UpdateOptions{})

	if err != nil {
		return errs.WrapError(fmt.Errorf("failed to update configmap autoscaler meta,because of %v", err))
	}

	return nil
}
