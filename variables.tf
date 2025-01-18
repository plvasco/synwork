variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

# variable "runner_uuid" {
#   description = "UUID for the Bitbucket runner"
#   type        = string
# }

# variable "runner_token" {
#   description = "Token for the Bitbucket runner"
#   type        = string
#   sensitive   = true
# }

## updated k8s.
variable "oauth_client_id" {
  description = "OAuth Client ID for the Bitbucket runner"
  type        = string
}

variable "oauth_client_secret" {
  description = "OAuth Client Secret for the Bitbucket runner"
  type        = string
}

variable "account_uuid" {
  description = "Account UUID for the Bitbucket runner"
  type        = string
}

variable "repository_uuid" {
  description = "Repository UUID for the Bitbucket runner"
  type        = string
}

variable "runner_uuid" {
  description = "Runner UUID for the Bitbucket runner"
  type        = string
  }