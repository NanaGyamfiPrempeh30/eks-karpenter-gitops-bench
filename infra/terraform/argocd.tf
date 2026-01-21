resource "helm_release" "argocd" {
  namespace        = "argocd"
  create_namespace = true
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = "5.53.0"
  wait             = true
  timeout          = 600

  values = [
    <<-EOT
    configs:
      params:
        server.insecure: true
    server:
      service:
        type: LoadBalancer
    applicationSet:
      enabled: true
    EOT
  ]

  depends_on = [
    module.eks
  ]
}
