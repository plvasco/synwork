provider "aws" {
  region = "us-east-1" // Change to your desired region
}

resource "aws_s3_bucket" "example" {
  count = 3
  bucket = "my-example-bucket-${count.index}"

  // Enable default encryption
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  // Enable versioning
  versioning {
    enabled = true
  }

  // Enable replication configuration
  replication_configuration {
    role = aws_iam_role.replication_role.arn

    rules {
      id     = "replication-rule"
      status = "Enabled"

      destination {
        bucket        = var.destination_bucket_arns[count.index] // Use variable for destination bucket ARN
        storage_class = "STANDARD"
      }

      filter {
        prefix = ""
      }
    }
  }
}

resource "aws_iam_role" "replication_role" {
  name = "s3-replication-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "replication_policy" {
  role = aws_iam_role.replication_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetReplicationConfiguration",
          "s3:ListBucket"
        ]
        Resource = "arn:aws:s3:::my-example-bucket-*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObjectVersion",
          "s3:GetObjectVersionAcl"
        ]
        Resource = "arn:aws:s3:::my-example-bucket-*/*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ReplicateObject",
          "s3:ReplicateDelete"
        ]
        Resource = "arn:aws:s3:::destination-bucket-*/*"
      }
    ]
  })
}

resource "aws_s3control_multi_region_access_point" "example" {
  account_id = var.account_id // Use variable for account ID

  details {
    name = "my-multi-region-access-point"

    regions {
      bucket = aws_s3_bucket.example[0].bucket
    }

    regions {
      bucket = aws_s3_bucket.example[1].bucket
    }

    regions {
      bucket = aws_s3_bucket.example[2].bucket
    }
  }
} 