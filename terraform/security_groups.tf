
resource "aws_security_group" "security_group" {
    egress      = [
        {
            cidr_blocks      = [
                "0.0.0.0/0",
            ]
            description      = ""
            from_port        = 0
            ipv6_cidr_blocks = []
            prefix_list_ids  = []
            protocol         = "-1"
            security_groups  = []
            self             = false
            to_port          = 0
        },
    ]
    ingress     = [
        {
            cidr_blocks      = [
                "0.0.0.0/0",
            ]
            description      = ""
            from_port        = 22
            ipv6_cidr_blocks = []
            prefix_list_ids  = []
            protocol         = "tcp"
            security_groups  = []
            self             = false
            to_port          = 22
        },
    ]
    name        = format("%s-sg-%s", var.Project, aws_vpc.vpc.id)
    tags        = {
      "Name" = format("%s-sg-%s", var.Project, aws_vpc.vpc.id)
    }
    vpc_id      = aws_vpc.vpc.id
}