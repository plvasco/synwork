output "s3_bucket_names" {
  description = "The names of the created S3 buckets"
  value       = [for bucket in aws_s3_bucket.example : bucket.bucket]
}

output "multi_region_access_point_arn" {
  description = "The ARN of the multi-region access point"
  value       = aws_s3control_multi_region_access_point.example.arn
} 