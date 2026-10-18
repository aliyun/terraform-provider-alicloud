package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cr_ee_repo", &resource.Sweeper{
		Name: "alicloud_cr_ee_repo",
		F:    testSweepCrEERepo,
	})
}

// testSweepCrEERepo cleans up Container Registry Enterprise Edition
// repositories left on CR EE instances by acceptance tests. Repositories are
// swept before namespaces (alicloud_cr_ee_namespace DependsOn
// alicloud_cr_ee_repo) so a namespace with test repos under it can be
// removed. Repositories on namespaces or instances marked NODELETING are
// preserved; in non-SWEEP_ALL mode only names with the standard test prefixes
// are removed.
func testSweepCrEERepo(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(fmt.Errorf("error getting Alicloud client: %s", err))
	}
	client := rawClient.(*connectivity.AliyunClient)
	crService := &CrService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftestAcc",
	}

	// Repositories live under namespaces on CR EE instances; enumerate the
	// instances and their namespaces first.
	pageNo, pageSize := 1, PageSizeLarge
	var instanceIds []string
	for {
		resp, err := crService.ListCrEEInstances(pageNo, pageSize)
		if err != nil {
			log.Printf("[ERROR] Failed to ListCrEEInstances: %s", err)
			return nil
		}
		for _, inst := range resp.Instances {
			instanceIds = append(instanceIds, inst.InstanceId)
		}
		if len(resp.Instances) < pageSize {
			break
		}
		pageNo++
	}

	for _, instanceId := range instanceIds {
		nsPageNo, nsPageSize := 1, PageSizeLarge
		for {
			nsResp, err := crService.ListCrEENamespaces(instanceId, nsPageNo, nsPageSize)
			if err != nil {
				log.Printf("[ERROR] Failed to ListCrEENamespaces for instance %s: %s", instanceId, err)
				break
			}
			for _, ns := range nsResp.Namespaces {
				repoPageNo, repoPageSize := 1, PageSizeLarge
				for {
					repoResp, err := crService.ListCrEERepos(instanceId, ns.NamespaceName, repoPageNo, repoPageSize)
					if err != nil {
						log.Printf("[ERROR] Failed to ListCrEERepos for instance %s namespace %s: %s", instanceId, ns.NamespaceName, err)
						break
					}
					for _, repo := range repoResp.Repositories {
						name := repo.RepoName
						if protectedSweepVpcName(name) || protectedSweepVpcName(repo.RepoNamespaceName) {
							log.Printf("[INFO] Skipping CR EE repo %s:%s:%s: protected (NODELETING)", instanceId, repo.RepoNamespaceName, name)
							continue
						}
						if !sweepAll() {
							skip := true
							for _, p := range prefixes {
								pl := strings.ToLower(p)
								if strings.HasPrefix(strings.ToLower(name), pl) ||
									strings.HasPrefix(strings.ToLower(repo.RepoNamespaceName), pl) {
									skip = false
									break
								}
							}
							if skip {
								continue
							}
						}
						id := fmt.Sprintf("%s:%s:%s", instanceId, repo.RepoNamespaceName, name)
						log.Printf("[INFO] Deleting CR EE repo: %s (repoId=%s)", id, repo.RepoId)
						if _, derr := crService.DeleteCrEERepo(id, repo.RepoId); derr != nil {
							if NotFoundError(derr) {
								continue
							}
							log.Printf("[ERROR] Failed to delete CR EE repo %s: %s", id, derr)
							continue
						}
					}
					if len(repoResp.Repositories) < repoPageSize {
						break
					}
					repoPageNo++
				}
			}
			if len(nsResp.Namespaces) < nsPageSize {
				break
			}
			nsPageNo++
		}
	}
	return nil
}

func TestAccAliCloudCREERepo_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
					"namespace":   "${alicloud_cr_ee_namespace.default.name}",
					"name":        name,
					"repo_type":   "PUBLIC",
					"summary":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"namespace":   CHECKSET,
						"name":        name,
						"repo_type":   "PUBLIC",
						"summary":     name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detail": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_immutability": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_immutability": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_immutability": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_immutability": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREERepo_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      "${alicloud_cr_ee_instance.default.id}",
					"namespace":        "${alicloud_cr_ee_namespace.default.name}",
					"name":             name,
					"repo_type":        "PUBLIC",
					"summary":          name,
					"detail":           name,
					"tag_immutability": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":      CHECKSET,
						"namespace":        CHECKSET,
						"name":             name,
						"repo_type":        "PUBLIC",
						"summary":          name,
						"detail":           name,
						"tag_immutability": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCREERepo_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_ee_repo.default.5"
	ra := resourceAttrInit(resourceId, AliCloudCREERepoMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEERepo")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCREERepoBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "6",
					"instance_id": "${alicloud_cr_ee_instance.default.id}",
					"namespace":   "${alicloud_cr_ee_namespace.default.name}",
					"name":        name + "-${count.index}",
					"repo_type":   "PUBLIC",
					"summary":     name,
					"detail":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"namespace":   CHECKSET,
						"name":        name + fmt.Sprint(-5),
						"repo_type":   "PUBLIC",
						"summary":     name,
						"detail":      name,
					}),
				),
			},
		},
	})
}

var AliCloudCREERepoMap0 = map[string]string{
	"repo_id": CHECKSET,
}

func AliCloudCREERepoBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cr_ee_instance" "default" {
		payment_type   = "Subscription"
		period         = 1
		renewal_status = "ManualRenewal"
		instance_type  = "Basic"
		instance_name  = var.name
	}

	resource "alicloud_cr_ee_namespace" "default" {
  		instance_id        = alicloud_cr_ee_instance.default.id
  		name               = var.name
  		auto_create        = false
  		default_visibility = "PRIVATE"
	}
`, name)
}
