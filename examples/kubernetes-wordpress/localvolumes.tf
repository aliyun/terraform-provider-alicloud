provider "kubernetes" {}

resource "kubernetes_persistent_volume" "mysql" {
  metadata {
    name = "local-pv-mysql"
    labels {
      type = "local"
    }
  }
  spec {
    capacity {
      storage = "20Gi"
    }
    access_modes = ["ReadWriteOnce"]
    persistent_volume_source {
      host_path {
        path = "/tmp/data/pv-mysql"
      }
    }
  }
}


resource "kubernetes_persistent_volume" "wordpress" {
  metadata {
    name = "local-pv-wordpress"
    labels {
      type = "local"
    }
  }
  spec {
    capacity {
      storage = "20Gi"
    }
    access_modes = ["ReadWriteOnce"]
    persistent_volume_source {
      host_path {
        path = "/tmp/data/pv-wordpress"
      }
    }
  }
}