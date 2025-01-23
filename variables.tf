variable "aws_region" {
  description = "The AWS region to create resources in"
  default     = "us-east-1" // Change to your desired default region
}

variable "account_id" {
  description = "Your AWS account ID"
  type        = string
}

variable "destination_bucket_arns" {
  description = "List of destination bucket ARNs for replication"
  type        = list(string)
} 