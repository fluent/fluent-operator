## 1.0.1 / 2022-05-10

### ENHANCEMENT
- Add utils related unit tests  (#251) 
- Change the original namespace to fluent (#267)
- Support logstashPrefix to ClusterOutput (#268)
- Add its shortname for each CRD  (#269)

### BUGFIX
- Fix the bug of feature request issue will closed by mistake (#257) 
- Fix crio init container mount path (#260) 
- Fix version error (#261) 
- Fix Helm Chart for Fluentd resources settings (#271) 
- Fix reload error (#277) 

## 1.0.0 / 2022-03-25

### Breaking Changes
- All Fluent Bit CRDs has been changed from namespaced scope to cluster scope
- CRDs and controller for Fluentd have been added

### Features
- Add priority class to Fluent Bit type (#146) 
- Add support for Fluent Bit RetryLimit in outputs (#148) 
- Add Fluent Bit Datadog output (#149)  
- Add support for Fluent Bit rewrite tag (#155)  
- Add Fluent Bit multiline logs support (#172)
- Add Fluent Bit aws filter plugin (#173)
- Add Fluent Bit multiline filter plugin (#176) 
- Add Fluent Bit Firehose plugin support (#178)
- Add Fluent Bit volume crd (#186)
- Renaming fluentbit-operator to fluent-operator (#189 #190)
- Add more fluentd examples (#194)
- Add Fluentd to helm charts (#204 #208 )
- Encrypt sensitive information for Fluentd output plugin (#219) 
- Enable multi-workers in one Fluentd pod (#194)
- Integrate e2e/function tests for generating Fluentd configuration (#203 #206 )
- Refine docs (#199 #228)
- Refactor multi images/binaries build, add github CI (#152 #154 #213 #214)
- Add CI templates (#248)
- Add Time_Key_Nanos field (#250)

### ENHANCEMENT
- Set the crictl path to a variable (#181) 
- Improved Fluent Bit kafka plugin (#182)

### BUGFIX
- Fix the incorrect key of the Fluent Bit es parser plugin (#164) 
- Fix the incorrect keys of the Fluent Bit es output plugin (#160) 
- Fix initcontainer script (#202) 
- Refine Fluentd CRs status (#225)
- Fix ci and make the repository importable and downloadable (#229)
- Fix codegen && add support for verifying codegen (#234 #238)
- Fix helm && Optimize helm (#236 #245 #246)

## 0.13.0 / 2022-3-14

- [FEATURE] Add priority class to Fluent Bit type #146
- [FEATURE] Add support for Fluent Bit RetryLimit in outputs #148
- [FEATURE] Add Fluent Bit Datadog output. #149
- [FEATURE] Add main workflow actions. #152
- [FEATURE] Add support for rewrite tag. #155
- [FEATURE] Add aws filter plugin. #173
- [FEATURE] Add multiline filter plugin. #176
- [FEATURE] Add Firehose plugin support. #178
- [FEATURE] Add volume crd. #186
- [ENHANCEMENT] Upgrade layout from Kubebuilder v2 to Kubebuilder v3.1. #147
- [ENHANCEMENT] Set the crictl path to a variable #181
- [ENHANCEMENT] Improved Fluent Bit kafka plugin #182
- [BUGFIX] Fix the incorrect keys of the Fluent Bit es output plugin #160
- [BUGFIX] Fix the incorrect key of the Fluent Bit es parser plugin #164

## 0.12.0 / 2021-09-07

- [FEATURE] Add support for collecting contained and cri-o service log. #142
- [ENHANCEMENT] Optionally enable namespace scoped RBAC. #137
- [BUGFIX] Update CRD YAML. #136

## 0.11.0 / 2021-09-01

- [FEATURE] Add support for Read_From_Head to the tail input plugin #129
- [ENHANCEMENT] Canonical Config #131
- [BUGFIX] Adjust FluentBit permissions #133

## 0.10.0 / 2021-08-20

- [FEATURE] Add support for polling in fluent-bit-watcher. #126
- [FEATURE] Add support for Amazon ElasticSearch Service and Elastic's Elasticsearch Service. #125

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
- [BUGFIX] Update groupname for client-gen to fluentbit.fluent.io #95

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