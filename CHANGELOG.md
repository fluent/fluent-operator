## 0.9.0 / 2021-08-13

- [FEATURE] Add support for `Containerd` and `CRI-O`. #112
- [FEATURE] Add support for `inotify_watcher` configuration of tail input plugin. #114
- [FEATURE] Add support for throttle filter plugin. #115
- [FEATURE] Add `runtimeClassName` support to fluentBit CRD. #116
- [FEATURE] Optional `watch-namespaces` for controller manager. #117
- [BUGFIX] Fix some bugs.

## 0.8.0 / 2021-07-23

- [FEATURE] Support setting imagePullSecrets for both operator and fluentbit #93 #94
- [FEATURE] Add switch to input.tail.memBufLimit in helm chart #87
- [ENHANCEMENT] Use hostpath instead of emptydir to store position db #72
- [ENHANCEMENT] Improve fluent-bit-watcher synchronization mechanism #74
- [ENHANCEMENT] Terminate fluent-bit process in a more elegant way in fluent-bit-watcher #90
- [ENHANCEMENT] Update README and roadmap #97 #100
- [ENHANCEMENT] Add kustomize file to manifests #99
- [BUGFIX] Fix the forward output can only use the default port problem. #89
- [BUGFIX] Fix bug it will loss log when damemonset restart. #90
- [BUGFIX] Update groupname for client-gen to logging.kubesphere.io #95

## 0.7.1 / 2021-07-08

- [ENHANCEMENT] fluent-bit-watcher: goroutine synchronization improvements. #74
- [CHANGE] add hostpath to sample configurations and manifests. #72
- [BUGFIX] Fix bug operator may crash when load plugin. #77

## 0.7.0 / 2021-06-29

- [ENHANCEMENT] Add support for plugin alias property in input and output specs. #64.
- [CHANGE] Add fluent-bit-watcher. #62.

## 0.6.2 / 2021-06-11

- [ENHANCEMENT] Update Kubernetes dependencies.
- [ENHANCEMENT] Update fluentbit resources in manifest.

## 0.6.1 / 2021-06-11

- Erroneous release with no changes from 0.6.0.

## 0.6.0 / 2021-06-01

- [ENHANCEMENT] Add Kubernetes Go client.
- [ENHANCEMENT] Support syslog output.
- [CHANGE] Update the default fluent-bit version to v1.7.3

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