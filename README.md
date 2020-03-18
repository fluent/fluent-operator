# Fluent Bit Operator

The Fluent Bit Operator for Kubernetes facilitates the deployment of Fluent Bit and provides great flexibility in building logging layer based on Fluent Bit. 

## Overview

Fluent Bit Operator defines five custom resources. Each `Input`, `Filter`, `Output` presents a Fluent Bit config 
section. `FluentBitConfig` selects them using label selector and generates the complete configuration into a ConfigMap.
 `FluentBit` defines specification for Fluent Bit Daemonset. 

- FluentBit: Defines Fluent Bit instances and associated ConfigMap
- FluentBitConfig: Select Input/Filter/Output and generates ConfigMap
- Input: Defines tail plugins
- Filter: Defines nest/modify/kubernetes plugins 
- Output: Defines forward/kafka/es/stdout plugins

## Features

- [x] Automate Fluent Bit deployment
- [x] Customize Fluent Bit config
- [x] Support common plugins
- [x] Support TLS/SSL
- [x] Read user/password from Secret
- [ ] Support parser plugins

## Contributing

The project is using Kubebuilder v2 and many files (documentation, manifests, ...) are auto-generated.