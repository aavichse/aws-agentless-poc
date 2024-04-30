data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      "Owner" = var.Owner
      "Project"   = var.Project
    } 
  }
}

