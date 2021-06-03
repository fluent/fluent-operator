## 0.6.0 / 2021-06-01

- [ENHANCEMENT] Add Kubernetes Go client.
- [ENHANCEMENT] Support syslog output.
- [CHANGE] Update the default fluent-bit version to v1.7.3.

## 0.5.0 / 2021-04-14

- [ENHANCEMENT] Support for audit logs

## 0.4.0 / 2021-04-01

- [ENHANCEMENT]  Support systemd input and Lua filter.
- [ENHANCEMENT]  Support Loki output.
- [ENHANCEMENT] Now it can set affinity and resource for fluent-bit daemonset.
- [CHANGE] Update the default fluent-bit version to v1.6.9.
- [BUGFIX] Fix some bugs.

## 0.3.0 / 2020-11-10

- [FEATURE] Support Parser plugin
[ENHANCEMENT] Support File, TCP, HTTP outputs

## 0.2.0 / 2020-08-27

- [CHANGE] Rewrite relevant CRDs. They are backwards incompatible with v0.1.0
- [CHANGE] Use kubebuilder as the building framework

## 0.1.0 / 2020-02-17

This is the first release of fluentbit operator.
