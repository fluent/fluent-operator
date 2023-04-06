# AWS

The AWS Filter Enriches logs with AWS Metadata. Currently the plugin adds the EC2 instance ID and availability zone to log records. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/aws-metadata**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| imdsVersion | Specify which version of the instance metadata service to use. Valid values are 'v1' or 'v2'. | string |
| az | The availability zone; for example, \"us-east-1a\". Default is true. | *bool |
| ec2InstanceID | The EC2 instance ID.Default is true. | *bool |
| ec2InstanceType | The EC2 instance type.Default is false. | *bool |
| privateIP | The EC2 instance private ip.Default is false. | *bool |
| amiID | The EC2 instance image id.Default is false. | *bool |
| accountID | The account ID for current EC2 instance.Default is false. | *bool |
| hostName | The hostname for current EC2 instance.Default is false. | *bool |
| vpcID | The VPC ID for current EC2 instance.Default is false. | *bool |
