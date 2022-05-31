resource "alicloud_cs_managed_kubernetes" "k8s" {
  # ... other configuration ...

  dynamic "addons" {
    for_each = var.cluster_addons
    content {
      name     = lookup(addons.value, "name", var.cluster_addons)
      config   = lookup(addons.value, "config", var.cluster_addons)
      disabled = lookup(addons.value, "disabled", var.cluster_addons)
    }
  }
}
