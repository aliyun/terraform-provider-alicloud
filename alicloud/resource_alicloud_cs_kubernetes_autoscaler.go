package alicloud

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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	clusterAutoscaler          = "cluster-autoscaler"
	defaultAutoscalerNamespace = "kube-system"
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

// resourceAlicloudCSKubernetesAutoscaler defines the schema of resource
func resourceAlicloudCSKubernetesAutoscaler() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesAutoscalerCreate,
		Read:   resourceAlicloudCSKubernetesAutoscalerRead,
		Update: resourceAlicloudCSKubernetesAutoscalerUpdate,
		Delete: resourceAlicloudCSKubernetesAutoscalerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
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
					},
				},
				MaxItems: 30,
			},
			"utilization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cool_down_duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"defer_scale_in_duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceAlicloudCSKubernetesAutoscalerCreate define how to create autoscaler
func resourceAlicloudCSKubernetesAutoscalerCreate(d *schema.ResourceData, meta interface{}) error {

	var clusterId string
	if value, ok := d.GetOk("cluster_id"); ok == true {
		clusterId = value.(string)
	} else {
		return WrapError(fmt.Errorf("clusterId must be provided instead of %s", "nil"))
	}

	err := resourceAlicloudCSKubernetesAutoscalerUpdate(d, meta)

	if err != nil {
		return WrapError(err)
	}
	// set unique id of tf state
	d.SetId(fmt.Sprintf("%s-%s", clusterId, clusterAutoscaler))
	return nil
}

// resourceAlicloudCSKubernetesAutoscalerRead no need to implement
func resourceAlicloudCSKubernetesAutoscalerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceAlicloudCSKubernetesAutoscalerUpdate define how to update autoscaler configuration
func resourceAlicloudCSKubernetesAutoscalerUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	// any changes need to ready autoscaler.
	if d.HasChange("nodepools") || d.HasChange("utilization") || d.HasChange("cool_down_duration") || d.HasChange("defer_scale_in_duration") {

		regionId := client.RegionId

		var clusterId, utilization, coolDownDuration, deferScaleInDuration string
		var ok bool
		var value interface{}

		if value, ok = d.GetOk("cluster_id"); ok == true {
			clusterId = value.(string)
		} else {
			return fmt.Errorf("please provide your cluster id in region %s", regionId)
		}

		// check utilization value
		if value, ok = d.GetOk("utilization"); ok == true {
			utilization = value.(string)
		} else {
			return fmt.Errorf("please provide utilization of cluster %s", clusterId)
		}

		// check cool_down_duration value
		if value, ok = d.GetOk("cool_down_duration"); ok == true {
			coolDownDuration = value.(string)
		} else {
			return fmt.Errorf("please provide cool_down_duration of cluster %s", clusterId)
		}

		// check defer_scale_in_duration
		if value, ok = d.GetOk("defer_scale_in_duration"); ok == true {
			deferScaleInDuration = value.(string)
		} else {
			return fmt.Errorf("please provide defer_scale_in_duration of cluster %s", clusterId)
		}

		// parse nodepools
		nodePoolsParams := d.Get("nodepools").(*schema.Set)
		nodePools := nodePoolsParams.List()

		if len(nodePools) == 0 {
			return fmt.Errorf("please provide at least one node pool of %s", clusterId)
		}

		// args creation
		args := make([]string, 0)
		args = applyDefaultArgs(args)

		for _, pool := range nodePools {

			poolBytes, err := json.Marshal(pool)
			if err != nil {
				return fmt.Errorf("failed to marshal pool,because of %v", err)
			}

			pool := make(nodePool)

			err = json.Unmarshal(poolBytes, &pool)

			if err != nil {
				return fmt.Errorf("failed to unmarshal pool,because of %v", err)
			}

			// get params of node pool
			id := pool["id"]
			labels := pool["labels"]
			taints := pool["taints"]

			// get userData from cluster openapi
			userData, err := csService.GetUserData(clusterId, labels, taints)

			if err != nil {
				return fmt.Errorf("failed to get permanent token,because of %v", err)
			}

			err = UpdateScalingGroupConfiguration(client, id, userData, labels, taints)
			if err != nil {
				return fmt.Errorf("failed to update scaling group status,because of %v", err)
			}

			// get min max of scaling group
			min, max, err := GetScalingGroupSizeRange(client, id)

			if err != nil {
				return fmt.Errorf("failed to describe scaling group %s,because of %v", id, err)
			}

			nodeArgs := fmt.Sprintf("--nodes=%d:%d:%s", min, max, id)

			args = append(args, nodeArgs)
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
			return fmt.Errorf("failed to download kubeconf from cluster,because of %v", err)
		}

		err = DeployAutoscaler(args, fmt.Sprintf(defaultAutoscalerImage, regionId), kubeConfPath)

		if err != nil {
			return fmt.Errorf("failed to deploy autoscaler,because of %v", err)
		}
	}
	return nil
}

// delete autoscaler deployment
func resourceAlicloudCSKubernetesAutoscalerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	regionId := client.RegionId

	clusterId := d.Get("cluster_id").(string)

	if clusterId == "" {
		return fmt.Errorf("please provide the cluster_id in region %s", regionId)
	}

	kubeConfPath, err := DownloadUserKubeConf(client, clusterId)

	if err != nil {
		return fmt.Errorf("failed to download kubeconf from cluster,because of %v", err)
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
		return fmt.Errorf("failed to parse DescribeScalingConfigurationsResponse of %s", groupId)
	}

	if err != nil {
		return WrapErrorf(err, "failed to invoke CreateDescribeScalingConfigurations,because of %v", err)
	}

	addDebug("DescribeScalingConfigurations", configurations, describeScalingConfigurationsRequest)

	if len(configurations.ScalingConfigurations.ScalingConfiguration) == 0 {
		//todo  create configuration
		return fmt.Errorf("please create the default scaling configuration of group %s", groupId)
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
			return WrapErrorf(err, "failed to invoke ModifyScalingConfiguration,because of %v", err)
		}

		addDebug("ModifyScalingConfiguration", modifyScalingConfigurationResponse, modifyScalingConfigurationRequest)
	}
	return nil
}

func GetScalingGroupSizeRange(client *connectivity.AliyunClient, groupId string) (min, max int, err error) {
	describeScalingGroupRequest := ess.CreateDescribeScalingGroupsRequest()
	describeScalingGroupRequest.RegionId = client.RegionId
	describeScalingGroupRequest.ScalingGroupId1 = groupId

	describeScalingGroupResponse, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingGroups(describeScalingGroupRequest)
	})

	if err != nil {
		return 0, 0, WrapErrorf(err, "failed to invoke DescribeScalingGroups,because of %v", err)
	}
	addDebug("DescribeScalingGroups", describeScalingGroupResponse, describeScalingGroupRequest)

	resp, ok := describeScalingGroupResponse.(*ess.DescribeScalingGroupsResponse)

	if ok != true {
		return 0, 0, fmt.Errorf("failed to parse DescribeScalingGroupsResponse of scaling group %s", groupId)
	}

	if resp.ScalingGroups.ScalingGroup == nil || len(resp.ScalingGroups.ScalingGroup) == 0 {
		return 0, 0, fmt.Errorf("the scaling group %s you specific is not found", groupId)
	}

	scalingGroup := resp.ScalingGroups.ScalingGroup[0]

	if &scalingGroup == nil {
		fmt.Errorf("the scaling group %s you specific is not found", groupId)
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
		return "", WrapErrorf(err, "failed to describe user kuberconf of cluster %s,because of %v", clusterId, err)
	}

	kubeConfResponse, ok := describeClusterUserKubeconfigResponse.(*cs.DescribeClusterUserKubeconfigResponse)

	if ok != true {
		return "", fmt.Errorf("failed to parse DescribeClusterUserKubeconfigResponse of %s", clusterId)
	}
	addDebug("DescribeClusterUserKubeconfig", describeClusterUserKubeconfigResponse, describeClusterUserKubeconfigRequest)

	// get response bytes
	kubeconfBytes := kubeConfResponse.GetHttpContentBytes()

	ukc := &userKubeConf{}
	err = json.Unmarshal(kubeconfBytes, ukc)

	if err != nil {
		return "", fmt.Errorf("failed to parse DescribeClusterUserKubeconfigResponse,because of %v", err)
	}

	content := ukc.Config

	wd, err := os.Getwd()

	if err != nil {
		return "", fmt.Errorf("failed to get current working dir,because of %v", err)
	}

	kubeConfPath := path.Join(wd, fmt.Sprintf("%s-kubeconf", clusterId))

	err = ioutil.WriteFile(kubeConfPath, []byte(content), 0755)

	if err != nil {
		return "", fmt.Errorf("failed to create kubeconf in working dir because of %v", err)
	}

	return kubeConfPath, nil
}

// delete autoscaler component
func DeleteAutoscaler(kubeconf string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)

	if err != nil {
		return fmt.Errorf("failed to build kubeconf from local path %s,because of %v", kubeconf, err)
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		return fmt.Errorf("failed to create client-go clientSet,because of %v", err)
	}

	return clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Delete(clusterAutoscaler, &metav1.DeleteOptions{})
}

// deploy cluster-autoscaler to kubernetes cluster
func DeployAutoscaler(args []string, image string, kubeconf string) error {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)
	if err != nil {
		return fmt.Errorf("failed to build kubeconf from local path %s,because of %v", kubeconf, err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create client-go clientSet,because of %v", err)
	}

	deploy, err := clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Get(clusterAutoscaler, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// create a new deploy
			deployObject := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterAutoscaler,
					Labels: map[string]string{
						"app": clusterAutoscaler,
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": clusterAutoscaler,
						},
					},

					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": clusterAutoscaler,
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
									Name:    clusterAutoscaler,
									Image:   image,
									Command: args,
								},
							},
						},
					},
				},
			}
			_, err := clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Create(deployObject)
			if err != nil {
				return fmt.Errorf("failed to create %s deployment,because of %v", clusterAutoscaler, err)
			}

		} else {
			return fmt.Errorf("failed to describe %s deployment,because of %v", clusterAutoscaler, err)
		}
	} else {
		// update deployment
		deploy.Spec.Template.Spec.Containers[0].Command = args
		_, err := clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Update(deploy)
		if err != nil {
			return fmt.Errorf("failed to update %s deployment,because of %v", clusterAutoscaler, err)
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
