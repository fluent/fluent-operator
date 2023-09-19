# Nginx

NGINX Exporter Metrics input plugin scrapes metrics from the NGINX stub status handler. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/nginx**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Name of the target host or IP address to check, default: localhost | string |
| port | Port of the target nginx service to connect to, default: 80 | *int32 |
| statusURL | The URL of the Stub Status Handler,default: /status | string |
| nginxPlus | Turn on NGINX plus mode,default: true | *bool |
