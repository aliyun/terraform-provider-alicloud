package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"golang.org/x/net/context"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func TestApplyDefaultArgs(t *testing.T) {
	args := make([]string, 0)
	args = applyDefaultArgs(args)
	if len(args) != 0 {
		t.Log("pass TestApplyDefaultArgs")
		return
	}
	t.Error("TestApplyDefaultArgs failed to apply default args")
}

func TestCreateScalingGroupTags(t *testing.T) {
	validLabels := "a=b,c=d"
	validTaints := "e=f:NoSchedule"
	tags := createScalingGroupTags(validLabels, validTaints)

	validLabelsArr := strings.Split(validLabels, ",")

	validTaintsArr := strings.Split(validTaints, ",")

	for _, label := range validLabelsArr {
		labelKeyValue := strings.Split(label, "=")
		if ok := strings.Contains(tags, fmt.Sprintf("%s%s", LabelPattern, labelKeyValue[0])); ok != true {
			t.Error("failed to pass TestCreateScalingGroupTags,because convert labels failure")
		}
	}

	for _, taint := range validTaintsArr {
		taintKeyValue := strings.Split(taint, "=")
		if ok := strings.Contains(tags, fmt.Sprintf("%s%s", TaintPattern, taintKeyValue[0])); ok != true {
			t.Error("failed to pass TestCreateScalingGroupTags,because convert taints failure")
		}
	}
	t.Log("pass TestCreateScalingGroupTags")
}

// lintignore: AT001
func TestAccAliCloudCSKubernetesAutoscaler_basic(t *testing.T) {
	resourceId := "alicloud_cs_kubernetes_autoscaler.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSKubernetesAutoscaler-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKubernetesAutoscalerConfigDependence)

	nodepools := func(value string) []map[string]string {
		return []map[string]string{
			{
				"id":     "${alicloud_cs_kubernetes_node_pool.autoscaler.scaling_group_id}",
				"labels": fmt.Sprintf("autoscaler-acc=%s", value),
				"taints": fmt.Sprintf("autoscaler-acc=%s:NoSchedule", value),
			},
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":              "${alicloud_cs_managed_kubernetes.default.id}",
					"utilization":             "0.5",
					"cool_down_duration":      "10m",
					"defer_scale_in_duration": "10m",
					"use_ecs_ram_role_token":  true,
					"nodepools":               nodepools("initial"),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "cluster_id"),
					resource.TestCheckResourceAttr(resourceId, "utilization", "0.5"),
					resource.TestCheckResourceAttr(resourceId, "cool_down_duration", "10m"),
					resource.TestCheckResourceAttr(resourceId, "defer_scale_in_duration", "10m"),
					resource.TestCheckResourceAttr(resourceId, "use_ecs_ram_role_token", "true"),
					resource.TestCheckResourceAttr(resourceId, "nodepools.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "nodepools.*", map[string]string{
						"labels": "autoscaler-acc=initial",
						"taints": "autoscaler-acc=initial:NoSchedule",
					}),
					testAccCheckCSKubernetesAutoscalerRuntime(
						resourceId,
						"0.5",
						"10m",
						"10m",
						"autoscaler-acc=initial",
						"autoscaler-acc=initial:NoSchedule",
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"utilization":             "0.6",
					"cool_down_duration":      "5m",
					"defer_scale_in_duration": "6m",
					"nodepools":               nodepools("updated"),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "utilization", "0.6"),
					resource.TestCheckResourceAttr(resourceId, "cool_down_duration", "5m"),
					resource.TestCheckResourceAttr(resourceId, "defer_scale_in_duration", "6m"),
					resource.TestCheckResourceAttr(resourceId, "nodepools.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "nodepools.*", map[string]string{
						"labels": "autoscaler-acc=updated",
						"taints": "autoscaler-acc=updated:NoSchedule",
					}),
					testAccCheckCSKubernetesAutoscalerRuntime(
						resourceId,
						"0.6",
						"5m",
						"6m",
						"autoscaler-acc=updated",
						"autoscaler-acc=updated:NoSchedule",
					),
				),
			},
			{
				// Remove the autoscaler while its cluster remains available so that the
				// deployment deletion can be verified before the final fixture cleanup.
				Config: resourceCSKubernetesAutoscalerConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cs_kubernetes_node_pool.autoscaler", "scaling_group_id"),
					testAccCheckCSKubernetesAutoscalerDeleted("alicloud_cs_managed_kubernetes.default"),
				),
			},
		},
	})
}

func testAccCheckCSKubernetesAutoscalerRuntime(resourceId, utilization, coolDownDuration, deferScaleInDuration, labels, taints string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("not found: %s", resourceId)
		}

		clusterId := rs.Primary.Attributes["cluster_id"]
		if clusterId == "" {
			return fmt.Errorf("cluster_id of %s is empty", resourceId)
		}
		if expectedId := fmt.Sprintf("%s:%s", clusterId, clusterAutoscaler); rs.Primary.ID != expectedId {
			return fmt.Errorf("unexpected autoscaler id: got %s, want %s", rs.Primary.ID, expectedId)
		}

		scalingGroupId := ""
		for key, value := range rs.Primary.Attributes {
			if strings.HasPrefix(key, "nodepools.") && strings.HasSuffix(key, ".id") && value != "" {
				scalingGroupId = value
				break
			}
		}
		if scalingGroupId == "" {
			return fmt.Errorf("scaling group id of %s is empty", resourceId)
		}

		clientSet, cleanup, err := testAccCSKubernetesAutoscalerClientSet(clusterId)
		if err != nil {
			return err
		}
		defer cleanup()

		deployment, err := clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Get(context.Background(), clusterAutoscaler, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get %s deployment: %w", clusterAutoscaler, err)
		}
		if len(deployment.Spec.Template.Spec.Containers) == 0 {
			return fmt.Errorf("deployment %s has no containers", clusterAutoscaler)
		}

		commands := deployment.Spec.Template.Spec.Containers[0].Command
		expectedCommands := []string{
			fmt.Sprintf("--scale-down-utilization-threshold=%s", utilization),
			fmt.Sprintf("--scale-down-gpu-utilization-threshold=%s", utilization),
			fmt.Sprintf("--scale-down-delay-after-add=%s", coolDownDuration),
			fmt.Sprintf("--scale-down-delay-after-failure=%s", coolDownDuration),
			fmt.Sprintf("--scale-down-unneeded-time=%s", deferScaleInDuration),
		}
		for _, expected := range expectedCommands {
			if !testAccCSKubernetesAutoscalerContainsCommand(commands, expected) {
				return fmt.Errorf("deployment %s command does not contain %q: %v", clusterAutoscaler, expected, commands)
			}
		}

		nodesCommandFound := false
		for _, command := range commands {
			if strings.HasPrefix(command, "--nodes=") && strings.HasSuffix(command, ":"+scalingGroupId) {
				nodesCommandFound = true
				break
			}
		}
		if !nodesCommandFound {
			return fmt.Errorf("deployment %s command does not contain a node range for scaling group %s: %v", clusterAutoscaler, scalingGroupId, commands)
		}

		configMap, err := clientSet.CoreV1().ConfigMaps(defaultAutoscalerNamespace).Get(context.Background(), clusterAutoscalerMeta, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get %s configmap: %w", clusterAutoscalerMeta, err)
		}
		config, ok := configMap.Data[clusterAutoscalerMeta]
		if !ok || config == "" {
			return fmt.Errorf("configmap %s does not contain key %s", clusterAutoscalerMeta, clusterAutoscalerMeta)
		}

		var meta autoscalerMeta
		if err := json.Unmarshal([]byte(config), &meta); err != nil {
			return fmt.Errorf("failed to unmarshal %s configmap: %w", clusterAutoscalerMeta, err)
		}
		if meta.UtilizationThreshold != utilization {
			return fmt.Errorf("unexpected utilization threshold: got %s, want %s", meta.UtilizationThreshold, utilization)
		}
		if meta.GpuUtilizationThreshold != utilization {
			return fmt.Errorf("unexpected gpu utilization threshold: got %s, want %s", meta.GpuUtilizationThreshold, utilization)
		}
		if meta.CoolDownDuration != coolDownDuration {
			return fmt.Errorf("unexpected cool down duration: got %s, want %s", meta.CoolDownDuration, coolDownDuration)
		}
		if meta.UnneededDuration != deferScaleInDuration {
			return fmt.Errorf("unexpected unneeded duration: got %s, want %s", meta.UnneededDuration, deferScaleInDuration)
		}

		scalingConfig, ok := meta.ScalingConfigurations[scalingGroupId]
		if !ok {
			return fmt.Errorf("configmap %s does not contain scaling group %s", clusterAutoscalerMeta, scalingGroupId)
		}
		if scalingConfig.Id != scalingGroupId {
			return fmt.Errorf("unexpected scaling configuration id: got %s, want %s", scalingConfig.Id, scalingGroupId)
		}
		if scalingConfig.Labels != labels {
			return fmt.Errorf("unexpected scaling configuration labels: got %s, want %s", scalingConfig.Labels, labels)
		}
		if scalingConfig.Taints != taints {
			return fmt.Errorf("unexpected scaling configuration taints: got %s, want %s", scalingConfig.Taints, taints)
		}

		return nil
	}
}

func testAccCheckCSKubernetesAutoscalerDeleted(clusterResourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster, ok := s.RootModule().Resources[clusterResourceId]
		if !ok || cluster.Primary == nil || cluster.Primary.ID == "" {
			return fmt.Errorf("not found: %s", clusterResourceId)
		}

		clientSet, cleanup, err := testAccCSKubernetesAutoscalerClientSet(cluster.Primary.ID)
		if err != nil {
			return err
		}
		defer cleanup()

		return resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, err := clientSet.AppsV1().Deployments(defaultAutoscalerNamespace).Get(context.Background(), clusterAutoscaler, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			}
			if err != nil {
				return resource.RetryableError(fmt.Errorf("failed to get %s deployment while checking deletion: %w", clusterAutoscaler, err))
			}
			return resource.RetryableError(fmt.Errorf("deployment %s still exists", clusterAutoscaler))
		})
	}
}

func testAccCSKubernetesAutoscalerClientSet(clusterId string) (*kubernetes.Clientset, func(), error) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	kubeConfigPath, err := DownloadUserKubeConf(client, clusterId)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to download kubeconfig for cluster %s: %w", clusterId, err)
	}

	cleanup := func() {
		_ = os.Remove(kubeConfigPath)
	}
	clientSet, err := getClientSetFromKubeconf(kubeConfigPath)
	if err != nil {
		cleanup()
		return nil, func() {}, fmt.Errorf("failed to create kubernetes client for cluster %s: %w", clusterId, err)
	}
	return clientSet, cleanup, nil
}

func testAccCSKubernetesAutoscalerContainsCommand(commands []string, expected string) bool {
	for _, command := range commands {
		if command == expected {
			return true
		}
	}
	return false
}

func resourceCSKubernetesAutoscalerConfigDependence(name string) string {
	return resourceCSAuoscalingConfigDependence(name) + `

resource "alicloud_cs_kubernetes_node_pool" "autoscaler" {
  cluster_id            = alicloud_cs_managed_kubernetes.default.id
  node_pool_name        = var.name
  vswitch_ids           = [local.vswitch_id]
  instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
  security_group_ids    = [alicloud_cs_managed_kubernetes.default.security_group_id]
  password              = "Terraform1234"
  image_type            = "AliyunLinux3ContainerOptimized"
  system_disk_category  = "cloud_efficiency"
  system_disk_size      = 40
  install_cloud_monitor = false

  scaling_config {
    enable   = true
    min_size = 0
    max_size = 1
    type     = "cpu"
  }
}
`
}
