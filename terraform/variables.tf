# Naming Convention 
variable "Owner" {
  default     = "aavichse"
}

variable "Project" {
  default     = "agentless-poc"
}

variable "IAM_User" {
  default     = "=agentless-poc-user-user"
}

variable "UserPemKey" {
  default = "aws-aavichse"
}

variable "instances" {
  type = map(object({
  }))
  default = {
    "test-vm-1" = {
    }, 
    "test-vm-2" = {
    }
  }
}