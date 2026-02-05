output "dynamodb_table_name" {
  description = "Provisioned DynamoDB table name."
  value       = aws_dynamodb_table.stacks.name
}

output "dynamodb_table_arn" {
  description = "Provisioned DynamoDB table ARN."
  value       = aws_dynamodb_table.stacks.arn
}

output "app_dynamodb_policy_arn" {
  description = "IAM policy ARN for application DynamoDB access."
  value       = aws_iam_policy.app_dynamodb.arn
}

output "irsa_role_arn" {
  description = "IRSA IAM role ARN for container-provisioner workload (null when create_irsa_role=false)."
  value       = var.create_irsa_role ? aws_iam_role.irsa[0].arn : null
}
