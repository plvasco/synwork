resource "kubernetes_namespace_v1" "runner" {
  metadata {
    name = "bitbucket-runner-ns"
  }
}

resource "kubernetes_secret_v1" "runner_oauth_credentials" {
  metadata {
    name      = "runner-oauth-credentials"
    namespace = kubernetes_namespace_v1.runner.metadata[0].name
    labels = {
      accountUuid = var.account_uuid
      repositoryUuid = var.repository_uuid
      runnerUuid = var.runner_uuid
    }
  }

  data = {
    oauthClientId     = base64encode(var.oauth_client_id)
    oauthClientSecret = base64encode(var.oauth_client_secret)
  }

  type = "Opaque"
}


resource "kubernetes_deployment_v1" "runner" {
  metadata {
    name      = "bitbucket-runner"
    namespace = kubernetes_namespace_v1.runner.metadata[0].name
  }

  spec {
    selector {
      match_labels = {
        app = "runner"
      }
    }

    replicas = 1

    template {
      metadata {
        labels = {
          app           = "runner"
          accountUuid = var.account_uuid
          repositoryUuid = var.repository_uuid
          runnerUuid = var.runner_uuid
        }
      }

      spec {
        container {
          name  = "runner"
          image = "docker-public.packages.atlassian.com/sox/atlassian/bitbucket-pipelines-runner"

          env {
            name  = "ACCOUNT_UUID"
            value = var.account_uuid
          }

          env {
            name  = "REPOSITORY_UUID"
            value = var.repository_uuid
          }

          env {
            name = "OAUTH_CLIENT_ID"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.runner_oauth_credentials.metadata[0].name
                key  = "oauthClientId"
              }
            }
          }

          env {
            name = "OAUTH_CLIENT_SECRET"
            value_from {
              secret_key_ref {
                name = kubernetes_secret_v1.runner_oauth_credentials.metadata[0].name
                key  = "oauthClientSecret"
              }
            }
          }

          env {
            name  = "WORKING_DIRECTORY"
            value = "/tmp"
          }

          volume_mount {
            name       = "tmp"
            mount_path = "/tmp"
          }

          volume_mount {
            name       = "docker-containers"
            mount_path = "/var/lib/docker/containers"
            read_only  = true
          }

          volume_mount {
            name       = "var-run"
            mount_path = "/var/run"
          }
        }

        container {
          name  = "docker-in-docker"
          image = "docker:20.10.5-dind"

          security_context {
            privileged = true
          }

          volume_mount {
            name       = "tmp"
            mount_path = "/tmp"
          }

          volume_mount {
            name       = "docker-containers"
            mount_path = "/var/lib/docker/containers"
          }

          volume_mount {
            name       = "var-run"
            mount_path = "/var/run"
          }
        }

        volume {
          name = "tmp"
          empty_dir {}
        }

        volume {
          name = "docker-containers"
          empty_dir {}
        }

        volume {
          name = "var-run"
          empty_dir {}
        }
      }
    }
  }
}
