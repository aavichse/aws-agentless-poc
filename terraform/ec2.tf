resource "aws_instance" "instance" {
    for_each                             = var.instances
    ami                                  = "ami-07d9b9ddc6cd8dd30"
    associate_public_ip_address          = true
    disable_api_stop                     = false
    disable_api_termination              = false
    ebs_optimized                        = false
    get_password_data                    = false
    hibernation                          = false
    instance_type                        = "t2.micro"
    ipv6_address_count                   = 0
    key_name                             = var.UserPemKey
    monitoring                           = false
    source_dest_check                    = true
    subnet_id                            = aws_subnet.subnet.id
    tags                                 = {
        Name    = "${var.Project}-vm-${each.key}"
    }
    tenancy                              = "default"
    vpc_security_group_ids               = [
        aws_security_group.security_group.id,
    ]

    capacity_reservation_specification {
        capacity_reservation_preference = "open"
    }

    credit_specification {
        cpu_credits = "standard"
    }

    enclave_options {
        enabled = false
    }

    maintenance_options {
        auto_recovery = "default"
    }

    metadata_options {
        http_endpoint               = "enabled"
        http_protocol_ipv6          = "disabled"
        http_put_response_hop_limit = 2
        http_tokens                 = "required"
        instance_metadata_tags      = "disabled"
    }

    private_dns_name_options {
        enable_resource_name_dns_a_record    = false
        enable_resource_name_dns_aaaa_record = false
        hostname_type                        = "ip-name"
    }

    root_block_device {
        delete_on_termination = true
        encrypted             = false
        tags                  = {}
        tags_all              = {}
        throughput            = 0
        volume_size           = 8
        volume_type           = "gp2"
    }
}


