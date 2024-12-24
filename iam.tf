resource "aws_iam_role" "bitbucket_runner_role" {
  name = "bitbucket-runner-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "eks.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "bitbucket_runner_policy" {
  name        = "bitbucket-runner-policy"
  description = "Policy for Bitbucket runner to access AWS services"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:*",
          "s3:*",
          "ec2:*"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "attach_runner_policy" {
  role       = aws_iam_role.bitbucket_runner_role.name
  policy_arn = aws_iam_policy.bitbucket_runner_policy.arn
}
