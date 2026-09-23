output "slb_ip" {
  value = kubernetes_service.wordpress.load_balancer_ingress[0].ip
}

