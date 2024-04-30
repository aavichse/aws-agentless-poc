resource "aws_iam_user" "iam_user" {
    name      = var.IAM_User
    path      = "/"
    tags      = {
        Name = format("%s-iam_user.%s", var.Project, var.IAM_User)
    }
}
resource "aws_iam_policy" "agentless-poc-user_iam_policy_s3_flowlogs" {
    description = "Allow read buckets contains flow logs"
    name        = format("%s_agentless-poc-user_iam_policy_s3_flowlogs", var.Project)
    path        = "/"
    policy      = jsonencode(
        {
            Statement = [
                {
                    Action   = [
                        "s3:GetObject",
                        "s3:ListBucket",
                    ]
                    Effect   = "Allow"
                    Resource = [
                        "arn:aws:s3:::{var.Project}-flowlogs",
                        "arn:aws:s3:::{var.Project}-flowlogs/*",
                    ]
                },
            ]
            Version   = "2012-10-17"
        }
    )
    tags        = {
        Name = format("%s_policy_s3_flowlogs", var.Project)
    }
    
}


resource "aws_iam_user_policy_attachment" "agentless-poc-user_iam_policy_s3_flowlogs_attach" {
  user       = aws_iam_user.iam_user.name
  policy_arn = aws_iam_policy.agentless-poc-user_iam_policy_s3_flowlogs.arn
}
