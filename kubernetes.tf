resource "kubernetes_namespace" "bitbucket_runner_ns" {
  metadata {
    name = "bitbucket-runner"
  }
}

resource "kubernetes_service_account" "runner_sa" {
  metadata {
    name      = "bitbucket-runner-sa"
    namespace = kubernetes_namespace.bitbucket_runner_ns.metadata[0].name
  }
}

resource "kubernetes_deployment" "bitbucket_runner" {
  metadata {
    name      = "bitbucket-runner-deployment"
    namespace = kubernetes_namespace.bitbucket_runner_ns.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "bitbucket-runner"
      }
    }

    template {
      metadata {
        labels = {
          app = "bitbucket-runner"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.runner_sa.metadata[0].name

        container {
          name  = "bitbucket-runner"
          image = "docker-public.packages.atlassian.com/sox/atlassian/bitbucket-pipelines-runner:latest"

          env {
            name  = "BITBUCKET_PIPELINES_RUNNER_UUID"
            value = var.runner_uuid
          }

          env {
            name  = "BITBUCKET_PIPELINES_RUNNER_TOKEN"
            value = var.runner_token
          }

          ports {
            container_port = 8080
          }
        }
      }
    }
  }
}
