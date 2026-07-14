variable "aws_region" {
  description = "AWS region"
  default     = "eu-central-1"
}

variable "instance_type" {
  description = "EC2 instance type"
  default     = "t3.micro"
}

variable "project_name" {
  description = "Project name for resource tags"
  default     = "markdownkb"
}

variable "volume_size" {
  description = "Root volume size in GB"
  default     = 8
}
