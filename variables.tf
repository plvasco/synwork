variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "runner_uuid" {
  description = "UUID for the Bitbucket runner"
  type        = string
}

variable "runner_token" {
  description = "Token for the Bitbucket runner"
  type        = string
  sensitive   = true
}
