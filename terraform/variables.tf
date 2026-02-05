variable "project" {
  description = "Project name used for tags and resource naming."
  type        = string
  default     = "smctf"
}

variable "environment" {
  description = "Deployment environment (e.g. dev, staging, prod)."
  type        = string
  default     = "dev"
}

variable "aws_region" {
  description = "AWS region for resources."
  type        = string
  default     = "us-east-1"
}

variable "common_tags" {
  description = "Additional tags applied to all resources."
  type        = map(string)
  default     = {}
}

variable "dynamodb_table_name" {
  description = "DynamoDB table name used by the application."
  type        = string
  default     = "smctf-stacks"
}

variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode."
  type        = string
  default     = "PAY_PER_REQUEST"

  validation {
    condition     = contains(["PAY_PER_REQUEST", "PROVISIONED"], var.dynamodb_billing_mode)
    error_message = "dynamodb_billing_mode must be PAY_PER_REQUEST or PROVISIONED."
  }
}

variable "dynamodb_read_capacity" {
  description = "Read capacity when billing mode is PROVISIONED."
  type        = number
  default     = 5
}

variable "dynamodb_write_capacity" {
  description = "Write capacity when billing mode is PROVISIONED."
  type        = number
  default     = 5
}

variable "enable_point_in_time_recovery" {
  description = "Enable DynamoDB point-in-time recovery."
  type        = bool
  default     = true
}

variable "create_irsa_role" {
  description = "Create IAM role for service account (IRSA) to access DynamoDB from EKS workloads."
  type        = bool
  default     = true
}

variable "irsa_role_name" {
  description = "IAM role name for IRSA."
  type        = string
  default     = "smctf-container-provisioner-irsa"
}

variable "eks_oidc_provider_arn" {
  description = "OIDC provider ARN from existing EKS cluster (required when create_irsa_role=true)."
  type        = string
  default     = ""
}

variable "eks_oidc_issuer_url" {
  description = "OIDC issuer URL from existing EKS cluster (required when create_irsa_role=true)."
  type        = string
  default     = ""
}

variable "k8s_service_account_namespace" {
  description = "Kubernetes service account namespace used by container-provisioner."
  type        = string
  default     = "smctf"
}

variable "k8s_service_account_name" {
  description = "Kubernetes service account name used by container-provisioner."
  type        = string
  default     = "container-provisioner"
}
