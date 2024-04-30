
resource "aws_route_table" "route_table" {
  vpc_id = aws_vpc.vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }

  tags = {
    Name  = format("%s-rt", var.Project)
  }
}

resource "aws_internet_gateway" "gateway" {
    tags     = {
        "Name"  = format("%s-gateway", var.Project)
    }
    vpc_id                          = aws_vpc.vpc.id
}


resource "aws_vpc" "vpc" {
    assign_generated_ipv6_cidr_block = false
    cidr_block                       = "10.0.0.0/16"
    enable_dns_hostnames             = true
    enable_dns_support               = true
    instance_tenancy                 = "default"
    tags                             = {
        "Name"  = format("%s-vpc", var.Project)
        "Tracker" = "cloud-app"
    }
}

resource "aws_subnet" "subnet" {
    assign_ipv6_address_on_creation = false
    availability_zone               = "us-east-1a"
    cidr_block                      = "10.0.1.0/24"
    map_public_ip_on_launch         = false
    tags                            = {
        "Name"  = format("%s-subnet", var.Project)
        "Tracker" = "cloud-app"
    }
    vpc_id                          = aws_vpc.vpc.id
}


resource "aws_route_table_association" "route_table_app_association" {
  subnet_id      = aws_subnet.subnet.id
  route_table_id = aws_route_table.route_table.id
}


resource "aws_flow_log" "vpc_flowlog" {
    log_destination          = "arn:aws:s3:::${aws_s3_bucket.vpc_flowlog_s3.bucket}"
    log_destination_type     = "s3"
    log_format               = "$${account-id} $${action} $${az-id} $${bytes} $${dstaddr} $${dstport} $${end} $${flow-direction} $${instance-id} $${interface-id} $${log-status} $${packets} $${pkt-dst-aws-service} $${pkt-dstaddr} $${pkt-src-aws-service} $${pkt-srcaddr} $${protocol} $${region} $${srcaddr} $${srcport} $${start} $${sublocation-id} $${sublocation-type} $${subnet-id} $${tcp-flags} $${traffic-path} $${type} $${version} $${vpc-id}"
    max_aggregation_interval = 60
    tags                     = {
        "Name"    = "${var.Project}-vpc-flowlogs"
    }
    traffic_type             = "ALL"
    vpc_id                   = aws_vpc.vpc.id

    destination_options {
        file_format                = "plain-text"
        hive_compatible_partitions = true
        per_hour_partition         = false
    }
}
