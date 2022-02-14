# AWS Metadata

The AWS Filter Enriches logs with AWS Metadata. Currently the plugin adds the EC2 instance ID and availability zone to log records. To use this plugin, you must be running in EC2 and have the [instance metadata service enabled](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html).


| Field | Description | Scheme |Default|
| ----- | ----------- | ------ | -----|
| imds_version | Specify which version of the instance metadata service to use. Valid values are 'v1' or 'v2'. | string |v2|
| az | The availability zone; for example, "us-east-1a". | bool |true|
| ec2_instance_id | The EC2 instance ID. | bool |true|
| ec2_instance_type | The EC2 instance type. | bool |false|
| private_ip | The EC2 instance private ip.| bool |false|
| ami_id | The EC2 instance image id. | bool |false|
| account_id | The account ID for current EC2 instance. | bool |false|
| hostname | The hostname for current EC2 instance.| bool |false|
| vpc_id | The VPC ID for current EC2 instance. | bool |false|


Note: If you run Fluent Bit in a container, you may have to use instance metadata v1. The plugin behaves the same regardless of which version is used.