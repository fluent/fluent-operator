## 2.5.0 / 2023-09-13
### Features
- Feat: add support to run Fluentd as a DaemonSet (#839)
- Feat: Add Prometheus exporter output plugin for Fluent Bit (#840)
- Feat: Add Fluent Bit forward input plugin (#843)
- Feat: support fluentd grok parser plugin (#861) 
- Feat: add gelf output plugin to fluentbit (#882) 
- Feat: add fluentbit opentelemetry plugin (#890) 
- Feat: Add serviceAccount Annotations to FluentBit resource (#898) 
- Feat: Add Helm priorityClassName mapping to Fluentd (#902) 
- Feat: add fluentbit http plugin (#904) 
- Feat: add fluentbit mqtt plugin (#911) 
- Feat: add fluentbit collectd plugin (#914) 

### ENHANCEMENT
- Remove Duplicate Cluster parsers in Fluent-bit config. (#853) 
- Add option to configure hostNetwork daemonset propertie (#863)
- Helm chart generation of ClusterOutput for loki (#865) (#906)
- Add SuppressTypeName option to helm, as newer ES needs it for bulk operation (#869) 
- Adjusting the parameters of fluentbit (#880) @wenchajun
- Add an extras section to the chart allowing new and less used features of the CRD to be used from the chart (#889)
- Add ImagePullSecret for fluentd (#891) 
- Add compress in fluentbit output es (#899)
- Expose ports on collector statefulset (#917) 
- Bump fluentbit to 2.1.9 (#921) 
- build(deps): Bump github.com/onsi/gomega from 1.27.8 to 1.27.10 (#844) 
- build(deps): Bump k8s.io/apimachinery from 0.27.3 to 0.27.4 (#847) 
- build(deps): Bump helm/kind-action from 1.7.0 to 1.8.0 (#859) 
- build(deps): Bump golang from 1.20.6-alpine3.17 to 1.20.7-alpine3.17 in /cmd/fluent-manager (#860) 
- build(deps): Bump github.com/go-kit/kit from 0.12.0 to 0.13.0 (#896) 
- build(deps): Bump golang from 1.20.7-alpine3.17 to 1.21.1-alpine3.17 in /cmd/fluent-manager (#913) 

### BUGFIX 
- Fix: Fluentd's s3 output plugin is compatible with minio (#858) 
- Fix: add hostNetwork option (#866) 
- Fix azure blob name & boolean parameters. (#887) 
- Fix: fluentd watchedNamespaces Helm mapping (#901) 

## 2.4.0 / 2023-07-19
### Features
- Feat: add scheduler support for fluentbit collector (#776)
- Users can upgrade fluentbit-operator to fluent-operator using the script (#779)
- Feat: adds the fluentd output plugin for Datadog (#803)
- Feat: add DefaultFilterSelector and DefaultOutputSelector to fluentd (#804)
- Feat: add S3 output plugin for Fluent Bit (#819)
- Support file system as storage layer in service section of fluenbit (#825)

### ENHANCEMENT
- Introduce stripUnderscores in fluent-operator helm values for fluentbit Input Systemd (#782)
- Add options in FluentOperator helm chart to add more systemdFilter in Fluentbit Systemd Input and a condition on systemdFilter to enable/disable (#785)
- Add podSecurityContext for fluentbit in fluent operator helm chart (#788)
- Rename field podSecurityContext to securityContext in Fluent Operator Deployment (#790)
- Add fluent operator security context at container level (#792)
- Add security context for fluenbit container (#796)
- Specify init container resources for fluent-operator deployment (#817)
- Update fluentd base Dockerfile (#820)
- MountPropagation option for internal mounts (#834)
- Fluent-bit upgrade to v2.1.7 (#836)
- build(deps): Bump golang from 1.20.4 to 1.20.6 in /docs/best-practice/forwarding-logs-via-http (#831)
- build(deps): Bump golang from 1.20.5-alpine3.17 to 1.20.6-alpine3.17 in /cmd/fluent-manager (#830)
- build(deps): Bump k8s.io/apimachinery from 0.27.2 to 0.27.3 (#828)
- build(deps): Bump golang from 1.20.4-alpine3.17 to 1.20.5-alpine3.17 in /cmd/fluent-manager (#783)
- build(deps): Bump github.com/onsi/gomega from 1.27.7 to 1.27.8 (#794)
- build(deps): Bump github.com/go-openapi/errors from 0.20.3 to 0.20.4 (#795)

### BUGFIX 
- Fix: resource deletion and adoption for 3 controllers (#777)
- Fix: Correct fluentd prase TimeFormat config key (#780)
- Fixes #798 storageClassName field not taken into account (#799)
- Fix: plugins document index (#822)
- Update the _helpers.tpl file (#823)
- Fix: incorrect field names in fluentd buffer plugin (#824)

## 2.3.0 / 2023-06-05
### Features
- Feat: Adding influxdb plugin (#690)
- Feat: Add EnvVars support to FluentD (#697)
- Feat: Add Pod Annotations support to FluentD (#698)
- Feat: Fluent operator & fluentbit: Added tolerations, nodeSelector + more (#704)
- Feat: Add fluentbit.affinity configuration (#726)
- Feat: Allow setting fluentd PodSecurityContext (#744)
- Feat: Fluentd in_tail plugin (#753)
- Feat: Add missing fluentd buffer options (#757)
- Feat: Add AWS Kinesis Data Streams output plugin for Fluent Bit (#768)
- Feat: Add global log_level support for fluentd (#770)
- Feat: Add scheduler support for fluentbit & fluentd (#771)

### ENHANCEMENT
- EnvVars support in fluentbit helm template (#706) 
- Add uri field for each telemetry type in opentelemetry plugin, remove old uri field (#708)
- Adjust fluentd watcher dependabot (#716)
- remove the deprecated -i flag in go build (#720) 
- Adjust fluentd arm64 image build timeout (#721)
- Adjust edge metrics collection config (#736) 
- Add some fluentbit helm opts (#743)
- Align CRDs and Operator with the fluentbit loki output (#756)
- Fluent-bit upgrade to v2.1.4 (#767)
- build(deps): Bump k8s.io/apimachinery from 0.26.3 to 0.27.1 (#701)
- build(deps): Bump helm/chart-testing-action from 2.1.0 to 2.4.0 (#710)
- build(deps): Bump k8s.io/klog/v2 from 2.90.1 to 2.100.1 (#712)
- build(deps): Bump golang from 1.20.3-alpine3.17 to 1.20.4-alpine3.17 in /cmd/fluent-manager (#713)
- build(deps): Bump golang from 1.20.3-alpine3.16 to 1.20.4-alpine3.16 in /cmd/fluent-watcher/fluentbit (#714)
- build(deps): Bump golang from 1.20.2 to 1.20.4 in /docs/best-practice/forwarding-logs-via-http (#715) 
- build(deps): Bump golang from 1.19.2-alpine3.16 to 1.20.4-alpine3.16 in /cmd/fluent-watcher/fluentd (#717)
- build(deps): Bump arm64v8/ruby from 3.1-slim-bullseye to 3.2-slim-bullseye in /cmd/fluent-watcher/fluentd (#718)
- build(deps): Bump alpine from 3.16 to 3.17 in /cmd/fluent-watcher/fluentd (#719)
- build(deps): Bump github.com/onsi/gomega from 1.27.6 to 1.27.7 (#748)
- build(deps): Bump k8s.io/apimachinery from 0.27.1 to 0.27.2 (#751)
- build(deps): Bump helm/kind-action from 1.5.0 to 1.7.0 (#765)

### BUGFIX 
- Fix: Fix missing log level  (#691) 
- Fix: Fix rewrite_tag match rule and trim start of string pattern (#692) 
- Fix(docs): Update cluster outputs docs link (#724)
- Fix: dereference pointers in parser filter plugin for fluentd (#745)
- Fix: fluentbit namespace-logging: only generate rewrite tag config once (#746)
- Fix: minor typo fix for firehose (#764)
- Fix: fix typo for estimate_current_event in fluentd (#769) 

## 2.2.0 / 2023-04-07
### Features
- Feat: Adding Fluentd cloudwatch plugin (#586)
- Feat: Adding an argument for disabling unused controllers (#621) 
- Feat: Namespace level CRDs and logging with FluentBit Daemonset (#630)
- Feat: Add service configurations for the components (#657)
- Add support for collecting edge metrics in Helm chart (#668)

### ENHANCEMENT
- Update controller-gen to v0.11.3 & update Makefile CRD_OPTIONS (#624)
- Fluentd add volumes & volumeClaimTemplates (#633)
- Make v2 importable (#631)
- Create or patch rbac objects (#635)
- Remove the deprecated -i flag for go build (#638)
- Add charts test (#639)
- Disable stdout output by default in helm (#641)
- Explain how to disable component controllers in the README (#642)
- Change gelfShortMessgeKey to gelfShortMessageKey (#643)
- Enable annotations for service account in fluentbit/fluentd (#647)
- Update README to mention new namespace FluentBit CRDs (#649)
- Add fluent-operator.drawio (#650)
- Updated image of operator architecture (#651)
- Automatic pushing of fluentbit-debug version images (#656)
- Make RBAC comptaible with multiple instances (#658)
- Upgrade Fluent Bit to v2.0.11 (#684)
- Update github runner to ubuntu 22.04 (#677)
- Build(deps): Bump golang from 1.20.1-alpine3.17 to 1.20.2-alpine3.17 in /cmd/fluent-manager (#606)
- Build(deps): Bump golang from 1.19.5-alpine3.16 to 1.20.2-alpine3.16 in /cmd/fluent-watcher/fluentbit (#607)
- Build(deps): Bump k8s.io/klog/v2 from 2.90.0 to 2.90.1 (#615) 
- Build(deps): Bump k8s.io/client-go from 0.26.2 to 0.26.3 (#626)
- Build(deps): Bump k8s.io/api from 0.26.2 to 0.26.3 (#628)
- Build(deps): Bump github.com/onsi/gomega from 1.27.2 to 1.27.5 (#637)
- Build(deps): Bump alpine from 3.17.2 to 3.17.3 in /cmd/fluent-watcher/fluentd/base (#648)
- Build(deps): Bump actions/checkout from 2 to 3 (#660)
- Build(deps): Bump golang from 1.20.1 to 1.20.2 in /docs/best-practice/forwarding-logs-via-http (#661)
- Build(deps): Bump actions/setup-go from 3 to 4 (#662)
- Build(deps): Bump actions/setup-python from 2 to 4 (#663)
- Build(deps): Bump azure/setup-helm from 1 to 3 (#664)
- Build(deps): Bump helm/kind-action from 1.2.0 to 1.5.0 (#665)
- Build(deps): Bump github.com/go-logr/logr from 1.2.3 to 1.2.4 (#671)
- Build(deps): Bump github.com/onsi/gomega from 1.27.5 to 1.27.6 (#672)
- Build(deps): Bump golang from 1.20.2-alpine3.16 to 1.20.3-alpine3.16 in /cmd/fluent-watcher/fluentbit (#679)
- Build(deps): Bump sigs.k8s.io/controller-runtime from 0.14.5 to 0.14.6 (#673)
- Build(deps): Bump golang from 1.20.2-alpine3.17 to 1.20.3-alpine3.17 in /cmd/fluent-manager (#678)

### BUGFIX 
- Fix: Properly exclude fluentbit output when stdout output is enabled (#618)
- Fix: Fix helm chart lint errors (#634)
- Fix: Fix segfault with DisableBuferVollume, rename to disableBufferVolume (#644)
- Fix: Fix the permissions in the cluster roles and bindings in helm (#667)
- Fix: Namespace level secret loader for namespaced FluentBit configs (#674)
- Fix: Add missing record modifier options (#675)
- Fix: Rename plugin docs directories to fix doc generation, add missing docs, minor doc improvements (#681)
- Fix: Fixing unit testing bugs (#682 #683)
- Fix: Fix null pointer error when creating namespace level CR (#686)

## 2.1.0 / 2023-03-13
### Features
- Feat: Adding Azure Blob output plugin (#549)
- Feat: Generic custom plugin type for Fluentd CRDs (#555) 
- Feat: Adding azureLogAnalytics output plugin for fluentbit (#563)
- Feat: Add ability to customize metrics port (#587)
- Feat: Enable fluentbit healthcheck (#598)
- Feat: Adding GCP Stackdriver Fluentbit Output Plugin (#605)
- Feat: Adding Cloudwatch for Fluentbit Output Plugin (#609) 

### ENHANCEMENT
- Support multi-architecture compilation, add platform amd64 compilation (#566)
- Update kubebuilder and kubectl (#574) 
- Config: run "make manifests" to generate metricsPort (#593)
- Make default ClusterInputs optional and configurable (#595)
- Bump kustomize from 4.5.7 to 5.0.0 (#572)
- Bump k8s.io/client-go from 0.25.4 to 0.26.1 (#573)
- build(deps): Bump k8s.io/klog/v2 from 2.80.1 to 2.90.0 (#551)
- build(deps): Bump github.com/joho/godotenv from 1.4.0 to 1.5.1 (#552)
- build(deps): Bump github.com/go-kit/log from 0.2.0 to 0.2.1 (#553) 
- build(deps): Bump alpine from 3.17.1 to 3.17.2 in /cmd/fluent-watcher/fluentd/base (#569)
- build(deps): Bump golang from 1.19.5-alpine3.17 to 1.20.1-alpine3.17 in /cmd/fluent-manager (#571)
- build(deps): Bump golang from 1.19.5 to 1.20.1 in /docs/best-practice/forwarding-logs-via-http (#596)
- build(deps): Bump sigs.k8s.io/controller-runtime from 0.14.4 to 0.14.5 (#599)
- build(deps): Bump github.com/onsi/gomega from 1.26.0 to 1.27.2 (#600)
- build(deps): Bump k8s.io/client-go from 0.26.1 to 0.26.2 (#602)

### BUGFIX 
- Fix: Fix Code format (including comment) (#565) 
- Fix: Update CRDs description / Documentation, conform to code (#591)
- Fix: Set the `path` field in fluentd to optional (#592)
- Fix: Add /finalizers to fluent-operator-clusterRole.yaml to fix openshift (#608)

## 2.0.1 / 2023-02-08
### ENHANCEMENT
- Upgrade Fluentd to v1.15.3 (#556) 
- Upgrade Fluentbit to v2.0.9 (#557)

### BUGFIX 
- Fix: Fix the bug of adding `label` (#548)

## 2.0.0 / 2023-02-03
### Features
- Feat: Support adding annotations to the fluent-operator deployment (#467)
- Feat: Support adding labels to the fluent-operator and the fluent-bit pods (#468)
- Feat: Add external plugin flag in the Fluent-Bit watcher (#469)
- Feat: Support adding annotations to the fluent-bit DaemonSet (#474)
- Feat: Add the `Collector` CRD and controller to support deploying Fluent Bit as a StatefulSet (#484)
- Feat: Add process termination timeout to fluent-bit-watcher (#512) 
- Feat: Add `dnsPolicy` and other Kubernetes filter options to the FluentBit CRD (#528) 

### ENHANCEMENT
- Add the `DockerModeParser` parameter to the fluentbit tail plugin (#486)
- Increase operator memory limit to 60Mi (#496)
- Refines the fluent-operator chart (#526)
- Update definition of flushThreadCount (#527)
- Upgrade Fluent Bit to v2.0.8 (#531)
- Refines e2e test script (#535)
- Dependabot: Update schedule and fix typo (#493)
- Build(deps): Bump k8s.io/client-go from 0.25.2 to 0.25.4 (#475) 
- Build(deps): Bump sigs.k8s.io/controller-runtime from 0.13.0 to 0.13.1 (#476)
- Build(deps): Bump github.com/fsnotify/fsnotify from 1.5.4 to 1.6.0 (#477)
- Build(deps): Bump golang from 1.19.2 to 1.19.3 in /docs/best-practice/forwarding-logs-via-http (#478)
- Build(deps): Bump alpine from 3.16.2 to 3.17.0 in /cmd/fluent-watcher/fluentd/base (#479) 
- Build(deps): Bump golang from 1.19.2-alpine3.16 to 1.19.3-alpine3.16 in /cmd/fluent-manager (#480) 
- Build(deps): Bump github.com/onsi/gomega from 1.21.1 to 1.24.1 (#481)
- Build(deps): Bump golang from 1.19.3 to 1.19.4 in /docs/best-practice/forwarding-logs-via-http (#497)
- Build(deps): Bump alpine from 3.17.0 to 3.17.1 in /cmd/fluent-watcher/fluentd/base (#507)
- Build(deps): Bump golang from 1.19.3-alpine3.16 to 1.19.5-alpine3.16 in /cmd/fluent-manager (#508)
- Build(deps): Bump golang from 1.19.2-alpine3.16 to 1.19.5-alpine3.16 in /cmd/fluent-watcher/fluentbit (#509) 
- Build(deps): Bump k8s.io/api from 0.25.4 to 0.26.1 (#519)
- Build(deps): Bump k8s.io/apimachinery from 0.25.4 to 0.26.1 (#520)
- Build(deps): Bump github.com/onsi/gomega from 1.24.1 to 1.26.0 (#530)
- Build(deps): Bump roots/issue-closer-action from 1.1 to 1.2 (#538) 
- Build(deps): Bump golang from 1.19.4 to 1.19.5 in /docs/best-practice/forwarding-logs-via-http (#539) 

### BUGFIX 
- Fix: Add Collector CRD to kustomization & Helm ClusterRole template (#515)
- Fix: Adjust fluentd-loki-output-plugin params (#523)
- Fix: Fix adding labels to the fluent-bit pods (#537) 

## 1.7.0 / 2022-11-23
### Features
- Feat: adding retry_limit to http-outputs (#445)
- Add environment variable support to the FluentBit CRD (#449)
- Make more fluent-bit configurations configurable via the FluentBit resource (#452)
- Feat: control/configure default ClusterFilters (helm chart) (#453)
- Add fluent-bit service and option to extend the RBAC configurations (#462)

### ENHANCEMENT
- Splunk make eventfield plural (#447)

### BUGFIX 
- Fix: intendation corrected in fluentbit-fluentBit.yaml (#454)
- Fix: fluentbit template render error - fixes #457 (#458)

## 1.6.1 / 2022-10-31
### BUGFIX 
- Fix: add missing config attributes for splunk output (#437)
- Fix(go): Update go version from 1.19.1 to 1.19.2 to resolve vulnerabilities. (#438)
- Revert "build: Enhance binary" (#439)

## 1.6.0 / 2022-10-25
### Features
- Add Fluent Bit Splunk output plugin (#417)

### ENHANCEMENT
- Bump github.com/go-kit/kit from 0.9.0 to 0.12.0 (#412)
- Bump github.com/joho/godotenv from 1.3.0 to 1.4.0 (#413)
- Bump github.com/go-openapi/errors from 0.19.2 to 0.20.3 (#414)
- Bump actions/checkout from 2 to 3 (#415)
- Build: Enhance binary (#416)
- Chore(deps): bump github.com/onsi/gomega from 1.20.1 to 1.21.1 (#419)
- Added support for Time_Offset parameter in regex parser (#423)
- Changing type of SplunkToken from string to secret (#427)
- Upgrade docker image version (#432)

### BUGFIX 
- Fixing a typo 'Spklunk' to 'Splunk' (#420)
- Helm: Fixing error in fluentbit-FluentBit (#422) 
- Fix clusterParser to ClusterParser (#426)
- Fix: Handling optional bool parameters for Splunk ClusterOutput (#428)

## 1.5.1 / 2022-09-30

### ENHANCEMENT
- Add Dependabot (#386)
- Bump azure/setup-helm from 1 to 3 (#387)
- Bump alpine from 3.13 to 3.16.2 in /cmd/fluent-watcher/fluentd/base (#388)
- Bump golang from 1.17.10-alpine3.16 to 1.19.1-alpine3.16 in /cmd/fluent-manager (#389) 
- Bump docker/setup-buildx-action from 1 to 2 (#390)
- Bump docker/login-action from 1 to 2 (#391) 
- Bump golang from 1.14 to 1.19.1 in /docs/best-practice/forwarding-logs-via-http (#392) 
- Bump actions/setup-go from 2 to 3 (#393)
- Bump actions/cache from 2 to 3 (#394)
- Bump sigs.k8s.io/yaml from 1.2.0 to 1.3.0 (#396)
- update go mod (#402) 
- Upgrade fluentbit to v1.9.9 (#403)
- Upgrade go version (#405)
- Upgrade golang image version (#406)

## 1.5.0 / 2022-09-24

### Features
- Add SecurityContext to FluentBit CRD (#344)
- Add OpenTelemetry output plugin (#345)
- Add Node Exporter Metrics input plugin (#345)
- Add Fluentd Loki output plugin (#346)
- Add Prometheus scrape metrics input plugin (#362)
- Add Prometheus remote write output plugin (#362)
- Add Fluent Bit metrics input plugin (#366)
- Add alias for the filter plugin (#370)
- Support custom plugins (#377, #380)
- Support receiving Non-K8s format log by Fluentd (#382)
- Add HostNetwork support for the Fluent Bit DaemonSet (#369)

### ENHANCEMENT
- Add node label to the Prometheus remote write metrics(#372)
- Simplify the steps of the issue report (#334)
- Add Fluentd Loki output plugin docs (#349)
- Add guide for node exporter metrics plugin (#353)
- Docs: update the index of the Fluent Bit plugins (#354)
- Add release drafter (#379)
- Add docs for the Prometheus scrape metrics input plugin and the Prometheus remote write output plugin (#381)
- Upgrade fluentbit to v1.9.8 (#384)

### BUGFIX 
- Fix the bug of feature request issue will be closed by mistake. (#341) 
- Correct invalid links (#347)

## 1.1.0 / 2022-06-15

### Features
- Add OpenSearch plugin for Fluent Bit (#298)
- Support custom annotations (#313)
- Add OpenSearch plugin for Fluentd (#324)
- Add helm & docs for OpenSearch plugin (#329)

### ENHANCEMENT
- Move some docs to fluent operator walkthrough (#290) 
- Docs refactoring (#291 #292 #293 #303 #314)
- Update go version (#316)
- Use a single systemd input plugin for various components (#323)

## 1.0.2 / 2022-05-17

### ENHANCEMENT
- Change reload signal from SIGUSR2 to SIGHUP  (#288) 

## 1.0.1 / 2022-05-12

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