resource "aws_s3_bucket" "vpc_flowlog_s3" {
    force_destroy = true
    bucket                      = "${var.Project}-flowlogs"
    object_lock_enabled         = false

    tags                        = {
        "Name" = "${var.Project}-flowlogs"
    }
}


resource "aws_s3_bucket_policy" "vpc_flowlog_s3_bucket_policy" {
    bucket = aws_s3_bucket.vpc_flowlog_s3.id
    policy =  jsonencode(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "AWSLogDeliveryWrite",
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "delivery.logs.amazonaws.com"
                    },
                    "Action": "s3:PutObject",
                    "Resource": "arn:aws:s3:::${aws_s3_bucket.vpc_flowlog_s3.bucket}/AWSLogs/aws-account-id=${data.aws_caller_identity.current.account_id}/*", 
                    "Condition": {
                        "StringEquals": {
                            "s3:x-amz-acl": "bucket-owner-full-control",
                            "aws:SourceAccount": "${data.aws_caller_identity.current.account_id}"
                        },
                        "ArnLike": {
                            "aws:SourceArn": "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
                        }
                    }
                },
                {
                    "Sid": "AWSLogDeliveryAclCheck",
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "delivery.logs.amazonaws.com"
                    },
                    "Action": "s3:GetBucketAcl",
                    "Resource": format("arn:aws:s3:::%s", aws_s3_bucket.vpc_flowlog_s3.bucket),
                    "Condition": {
                        "StringEquals": {
                            "aws:SourceAccount": "${data.aws_caller_identity.current.account_id}"
                        },
                        "ArnLike": {
                            "aws:SourceArn": "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:*"
                        }
                    }
                }
            ]
        }
)
}