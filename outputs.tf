output "bitbucket_runner_namespace" {
  value = kubernetes_namespace.bitbucket_runner_ns.metadata[0].name
}

output "iam_role_arn" {
  value = aws_iam_role.bitbucket_runner_role.arn
}
