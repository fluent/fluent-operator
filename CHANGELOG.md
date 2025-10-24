# Change Log

## 3.5.0 / 2025-10-24

* Bumped chart-version by @ncauchois in https://github.com/fluent/fluent-operator/pull/1596
* Sanitize markdown by resolving linter warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1611
* Support both IPv4 and IPv6 addresses in http_listen configuration by @damyan in https://github.com/fluent/fluent-operator/pull/1616
* chore(deps): update ghcr.io/fluent/fluent-operator/fluent-operator docker tag to v3.4.0 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1617
* helm-chart: patch 3.4 release by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1623
* fix(fluent-operator): broken link about nodeselector by @leehosu in https://github.com/fluent/fluent-operator/pull/1626
* Make fluentbit livenessProbe more flexible by @CharlieR-o-o-t in https://github.com/fluent/fluent-operator/pull/1628
* docs: add link to label-router project by @csatib02 in https://github.com/fluent/fluent-operator/pull/1627
* Update fluentbit-fluentBit.yaml to address issue 1635 by @g1franc in https://github.com/fluent/fluent-operator/pull/1636
* Allow setting logfmt parser options by @cosmastech in https://github.com/fluent/fluent-operator/pull/1637
* ClusterInput/ClusterOutput helm chart changes  by @CharlieR-o-o-t in https://github.com/fluent/fluent-operator/pull/1642
* Validate region and its value in Kinesis Output plugin by @smallc2009 in https://github.com/fluent/fluent-operator/pull/1644
* Add `workers` param for s3 output by @hercynium in https://github.com/fluent/fluent-operator/pull/1647
* fix: Invalid reference by @sousa-miguel in https://github.com/fluent/fluent-operator/pull/1643
* make tls config in elastic more clearer by @smallc2009 in https://github.com/fluent/fluent-operator/pull/1645
* Bump fluent-bit-crds and fluentd-crds sub-charts to 3.4.2. by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1654
* Fix helm chart linting errors by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1656
* build(deps): bump golang to 1.24.5 by @cw-Guo in https://github.com/fluent/fluent-operator/pull/1665
* fix: fix ci check error due to shellcheck by @cw-Guo in https://github.com/fluent/fluent-operator/pull/1668
* Fix error handling by @sugaf1204 in https://github.com/fluent/fluent-operator/pull/1666
* ci: Fix "ct lint" action  by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1680
* Bump fluent-bit to 4.0.9 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1683
* Bump docker builds to Go v1.24.5 by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1684
* feat: add servicemonitor for fluent-operator to helm chart by @dennis-ge in https://github.com/fluent/fluent-operator/pull/1677
* build(deps): Bump aquasecurity/trivy-action from 0.30.0 to 0.33.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1676
* build(deps): Bump renovatebot/github-action from 41.0.22 to 43.0.10 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1675
* build(deps): Bump docker/setup-buildx-action from 3.10.0 to 3.11.1 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1650
* build(deps): Bump docker/build-push-action from 6.16.0 to 6.18.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1630
* build(deps): Bump actions/setup-go from 5.4.0 to 5.5.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1629
* Bump fluent-bit to 4.0.11 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1691
* allow s3 output plugin to get keys from secrets by @v-davegillies-upscale in https://github.com/fluent/fluent-operator/pull/1688
* Bump fluent-bit to 4.1.0 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1699
* build(deps): Bump azure/setup-helm from 4.3.0 to 4.3.1 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1695
* build(deps): Bump renovatebot/github-action from 43.0.10 to 43.0.14 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1697
* build(deps): Bump aquasecurity/trivy-action from 0.33.0 to 0.33.1 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1694
* build(deps): Bump golang.org/x/sync from 0.14.0 to 0.17.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1686
* build(deps): Bump github.com/oklog/run from 1.1.0 to 1.2.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1649
* build(deps): Bump github.com/go-logr/logr from 1.4.2 to 1.4.3 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1632
* build(deps): Bump actions/checkout from 4.2.2 to 5.0.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1698
* build(deps): Bump actions/setup-go from 5.5.0 to 6.0.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1696
* Add golangci lint and resolve linter warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1703
* Update generated files by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1705
* Add golang-ci configuration matching with latest operator SDK by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1704
* Fix Docker warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1706
* Update operator-sdk to v1.41.1 according to migrations by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1702
* Fix goconst linter warnings #1707 by @u5surf in https://github.com/fluent/fluent-operator/pull/1711
* Fix cyclomatic complexity linter warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1708
* Resolve long line length linter warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1710
* bump chart by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1715
* Resolve some duplicate code linter warnings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1718
* Fix some `lll` warnings. by @u5surf in https://github.com/fluent/fluent-operator/pull/1720
* Bump fluent-bit to 4.1.1 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1729
* Bump fluentd to v1.19.0. by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1730
* fluentd: Fix gocyclo warnings by @u5surf in https://github.com/fluent/fluent-operator/pull/1723
* build(deps): Bump github.com/go-openapi/errors from 0.22.1 to 0.22.3 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1712
* Add fields to AzureBlob output plugin by @BHYub in https://github.com/fluent/fluent-operator/pull/1727
* add ClusterFirstWithHostNet for kubeedge fluentbit by @Abirdcfly in https://github.com/fluent/fluent-operator/pull/1726
* build(deps): Bump github.com/onsi/ginkgo/v2 from 2.23.4 to 2.26.0 by @dependabot[bot] in https://github.com/fluent/fluent-operator/pull/1724
* Use unique names for Fluentbit/Fluentd controllers. by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1736
* chore(deps): update ghcr.io/fluent/fluent-operator/fluent-bit docker tag to v4.1.1 by @github-actions[bot] in https://github.com/fluent/fluent-operator/pull/1738
* Re-factor build workflow for fluent-operator. by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1737
* Fix dupl warnings by @u5surf in https://github.com/fluent/fluent-operator/pull/1735
* feat(fluentbit): add text_payload_key to stackdriver by @cw-Guo in https://github.com/fluent/fluent-operator/pull/1669

## New Contributors
* @ncauchois made their first contribution in https://github.com/fluent/fluent-operator/pull/1596
* @damyan made their first contribution in https://github.com/fluent/fluent-operator/pull/1616
* @leehosu made their first contribution in https://github.com/fluent/fluent-operator/pull/1626
* @csatib02 made their first contribution in https://github.com/fluent/fluent-operator/pull/1627
* @g1franc made their first contribution in https://github.com/fluent/fluent-operator/pull/1636
* @cosmastech made their first contribution in https://github.com/fluent/fluent-operator/pull/1637
* @hercynium made their first contribution in https://github.com/fluent/fluent-operator/pull/1647
* @sousa-miguel made their first contribution in https://github.com/fluent/fluent-operator/pull/1643
* @sugaf1204 made their first contribution in https://github.com/fluent/fluent-operator/pull/1666
* @v-davegillies-upscale made their first contribution in https://github.com/fluent/fluent-operator/pull/1688
* @u5surf made their first contribution in https://github.com/fluent/fluent-operator/pull/1711
* @BHYub made their first contribution in https://github.com/fluent/fluent-operator/pull/1727
* @Abirdcfly made their first contribution in https://github.com/fluent/fluent-operator/pull/1726

## 3.4.0 / 2025-05-08

### Features

- feat(helm/fluent-operator): add option to disable rbac creation by @gbloquel in https://github.com/fluent/fluent-operator/pull/1556
- Added support for deploying multiple fluentbit collector replicas by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1561
- feat(fluentd): add null output plugin by @cw-Guo in https://github.com/fluent/fluent-operator/pull/1578
- adding support for Syslog over TLS by @matelang in https://github.com/fluent/fluent-operator/pull/1603
- Add structured metadata support for Loki output plugin by @error9098x in https://github.com/fluent/fluent-operator/pull/1579
- expose Enable_Chunk_Trace in the crd, enabling TAP debuging by @danielpodwysocki in https://github.com/fluent/fluent-operator/pull/1588
- feat(charts): Add ability for custom `positionDB` for `FluentBit` by @kiblik in https://github.com/fluent/fluent-operator/pull/1548
- Added the ability to specify Fluentd service type by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1564
- Added the ability to set `Use_Tag_For_Meta` on fluentbit kubernetes filter by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1565
- Add support for compression to the Fluentd HTTP output plugin by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1560
- Added the ability to set `DB.locking` on fluentbit tail inputs by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1567
- Added the ability to set `Owner_References` on fluentbit kubernetes filter by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1566

### Enhancements

- Update fluent-operator-clusterRole.yaml by @duj4 in https://github.com/fluent/fluent-operator/pull/1502
- Pin GitHub actions on commit hash according to best practices by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1513
- Use go.mod version in workflows by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1514
- Move dash so that labelKeys and removeKeys on separate line by @heytrav in https://github.com/fluent/fluent-operator/pull/1509
- makefile: Remove chmod+x by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1521
- gitignore: remove gitignore file and move content to .gitignore by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1519
- editorconfig: Add .editorconfig to ensure files are formatted consistently by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1518
- go-vet: Fix the Go vet findings by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1517
- actions: Remove cache action for Go by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1528
- hack: Fix shellcheck findings in bash scripts by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1527
- chore: update watcher start log by @cw-Guo in https://github.com/fluent/fluent-operator/pull/1529
- Re-factors CI workflow for building & publishing fluent-bit image by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1531
- Update formatting based on prettier plugin by @truongnht in https://github.com/fluent/fluent-operator/pull/1536
- build-fb-image: Update release documentation by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1515
- Bump fluent-bit to 4.0.0 by @github-actions in https://github.com/fluent/fluent-operator/pull/1544
- Add renovate workflow to bump fluent-bit version by @truongnht in https://github.com/fluent/fluent-operator/pull/1535
- renovate wf: Runs renovate job on ubuntu-latest by @truongnht in https://github.com/fluent/fluent-operator/pull/1549
- helm-chart: Improve templates by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1555
- helm-chart: Streamline image values and usage by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1550
- Support setting fluent-bit bufferChunkSize for tail input by @truongnht in https://github.com/fluent/fluent-operator/pull/1569
- makefile: Update Makefile by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1571
- ci: Ensure all generated code is committed by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1570
- helm-chart: Prevent few more occasions of template injection by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1572
- charts/fluent-operator: Add `namespace` to the ServiceAccount by @TeddyAndrieux in https://github.com/fluent/fluent-operator/pull/1590
- ci: Ensure helm tests run on changes to the chart by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1597
- dependabot: Group the k8s.io/- dependency updates in single PR by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1594
- dependabot: Refactor docker ecosystem to new syntax by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1595
- Fix(doc): open_search.md referencing elasticsearch name by @Anghille in https://github.com/fluent/fluent-operator/pull/1408
- fix: Update outdated crds by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1605

### Dependencies

- actions: Pin setup-helm to v4.3.0 + Bump Helm to v3.17.2 by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1520
- Bump fluent-bit to 3.2.9 by @github-actions in https://github.com/fluent/fluent-operator/pull/1511
- images: Align Go version to be 1.24.1 based on go.mod defined version by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1516
- build(deps): Bump actions/setup-go from 5.3.0 to 5.4.0 by @dependabot in https://github.com/fluent/fluent-operator/pull/1525
- build(deps): Bump actions/cache from 4.2.2 to 4.2.3 by @dependabot in https://github.com/fluent/fluent-operator/pull/1524
- build(deps): Bump aquasecurity/trivy-action from 0.29.0 to 0.30.0 by @dependabot in https://github.com/fluent/fluent-operator/pull/1523
- build(deps): Bump github.com/onsi/gomega from 1.34.2 to 1.36.3 by @dependabot in https://github.com/fluent/fluent-operator/pull/1533
- build(deps): Bump helm/chart-testing-action from 2.6.1 to 2.7.0 by @dependabot in https://github.com/fluent/fluent-operator/pull/1487
- fluent-bit: Bump fluent-bit from 3.2.9 to 3.2.10 by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1537
- dependencies: Update go dependencies by @marcofranssen in https://github.com/fluent/fluent-operator/pull/1552
- Bump fluent-bit to 4.0.1 by @github-actions in https://github.com/fluent/fluent-operator/pull/1584
- build(deps): Bump docker/build-push-action from 6.15.0 to 6.16.0 by @dependabot in https://github.com/fluent/fluent-operator/pull/1599
- build(deps): Bump renovatebot/github-action from 41.0.13 to 41.0.22 by @dependabot in https://github.com/fluent/fluent-operator/pull/1598
- chore(deps): update ghcr.io/fluent/fluent-operator/fluent-bit docker tag to v4.0.1 by @github-actions in https://github.com/fluent/fluent-operator/pull/1585
- build(deps): Bump golang.org/x/sync from 0.13.0 to 0.14.0 by @dependabot in https://github.com/fluent/fluent-operator/pull/1604

### Bugfixes

- Fix Kubernetes Events DBSync config by @ZephireNZ in https://github.com/fluent/fluent-operator/pull/1546
- Add missing multiline ClusterFilter values by @discostur in https://github.com/fluent/fluent-operator/pull/1581
- (ci/fluentbit) fix: Adds permission to publish packages  by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1538
- (ci/fluentbit) fix: provide packages:write permission by @joshuabaird in https://github.com/fluent/fluent-operator/pull/1539
- Fixed nil pointer dereference (panic) when port numbers are unset by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1563
- Fix ClusterFluentBitConfig rendering in helm chart when using yaml configFileFormat by @truongnht in https://github.com/fluent/fluent-operator/pull/1573
- fix(crd): disallow null values for logfmt parser to prevent fluent-bit crash by @sandy2008 in https://github.com/fluent/fluent-operator/pull/1543
- Fix fluentbit service selector not using pod labels when defined by @solidDoWant in https://github.com/fluent/fluent-operator/pull/1575

## New Contributors

- @duj4 made their first contribution in https://github.com/fluent/fluent-operator/pull/1502
- @marcofranssen made their first contribution in https://github.com/fluent/fluent-operator/pull/1513
- @heytrav made their first contribution in https://github.com/fluent/fluent-operator/pull/1509
- @sandy2008 made their first contribution in https://github.com/fluent/fluent-operator/pull/1543
- @ZephireNZ made their first contribution in https://github.com/fluent/fluent-operator/pull/1546
- @kiblik made their first contribution in https://github.com/fluent/fluent-operator/pull/1548
- @gbloquel made their first contribution in https://github.com/fluent/fluent-operator/pull/1556
- @solidDoWant made their first contribution in https://github.com/fluent/fluent-operator/pull/1563
- @discostur made their first contribution in https://github.com/fluent/fluent-operator/pull/1581
- @danielpodwysocki made their first contribution in https://github.com/fluent/fluent-operator/pull/1588
- @TeddyAndrieux made their first contribution in https://github.com/fluent/fluent-operator/pull/1590
- @matelang made their first contribution in https://github.com/fluent/fluent-operator/pull/1603
- @error9098x made their first contribution in https://github.com/fluent/fluent-operator/pull/1579
- @Anghille made their first contribution in https://github.com/fluent/fluent-operator/pull/1408

## 3.3.0 / 2025-02-27

### Features

- Add skip empty lines in tail input (#1352) @smallc2009
- Update the module path to github.com/fluent/fluent-operator/v3 (#1355) @jiuxia211
- Update fluentd image references (#1357) @reegnz
- Update fluent-bit to 3.1.8 (#1356) @reegnz
- Remove influxdb host validation in ClusterOutput (#1363) @smallc2009
- Add release-tool workflow to generate release PR (#1362) @cw-Guo
- Align CRDs with fluentbit out_azure (#1371) @felfa01
- Support rollout restart for daemonset and statefulset (#1375) @cw-Guo
- Add namespaceClusterFbCfg to ClusterFluentBitConfig custom resource (#1382) @btalakola
- Add affinity support for fluent-operator deployment (#1401) @smallc2009
- Add support for logs_body_key parameter on Opentelemetry output (#1411) @yilmazo
- Add HostAliases support to Fluent Bit and Fluentd specifications (#1413) @MioOgbeni
- Add rdkafka gem installation to Dockerfiles for fluent-watcher (#1415) @MioOgbeni
- Add filter ordinals (#1386) @reegnz
- Support fluentbit tail offsetKey parameters (#1437) @cw-Guo
- Add additional params for input & output APIs and retry_limit for default loki output (#1442) @chrono2002
- Add VERSION file for fluentbit image (#1447) @cw-Guo
- Add pipeline to bump fluent bit version (#1448) @cw-Guo
- Add envFrom support for fluentd daemonset and statefulsets (#1458) @thapabishwa
- Add livenessProbe to FluentBit template (#1460) @CharlieR-o-o-t
- Add support for reload_after, sniffer_class_name es output parameters (#1462) @penekk

### ENHANCEMENT

- Fix indentation bug (#1360) @harshvora10101
- Fix documentation links (#1361) @reegnz
- Update helm chart version and clarify helm-related docs (#1378) @joshuabaird
- Add README to Helm chart (#1381) @joshuabaird
- Change fluentbit flushSeconds type to float64 (#1406) @jjsiv
- Improve pipelines and add documents (#1450) @cw-Guo
- Support setting configFileFormat in helmchart (#1466) @truongnht
- Bump fluent-bit to 3.2.5 (#1464) @github-actions
- Re-factor fluentd CI workflows (#1472) @joshuabaird
- Update fluentd to v1.17.1 (#1478) @joshuabaird

### Dependencies

- Bump golang from 1.22.6-alpine3.19 to 1.23.2-alpine3.19 in /cmd/fluent-manager (#1369) @dependabot
- Bump golang from 1.22.6-alpine3.19 to 1.23.2-alpine3.19 in /cmd/fluent-watcher/fluentbit (#1368) @dependabot
- Bump golang from 1.22.6-alpine3.19 to 1.23.2-alpine3.19 in /cmd/fluent-watcher/fluentd (#1367) @dependabot
- Bump golang from 1.22.4 to 1.23.1 in /docs/best-practice/forwarding-logs-via-http (#1366) @dependabot
- Bump github.com/onsi/gomega from 1.34.1 to 1.34.2 (#1331) @dependabot
- Bump golang from 1.23.1 to 1.23.4 in /docs/best-practice/forwarding-logs-via-http (#1446) @dependabot
- Bump helm/kind-action from 1.10.0 to 1.12.0 (#1445) @dependabot
- Bump aquasecurity/trivy-action from 0.24.0 to 0.29.0 (#1426) @dependabot
- Bump golang from 1.23.2-alpine3.19 to 1.23.4-alpine3.19 in /cmd/fluent-watcher/fluentbit (#1429) @dependabot
- Bump golang from 1.23.2-alpine3.19 to 1.23.4-alpine3.19 in /cmd/fluent-manager (#1430) @dependabot
- Bump golang from 1.23.2-alpine3.19 to 1.23.4-alpine3.19 in /cmd/fluent-watcher/fluentd (#1431) @dependabot
- Bump alpine from 3.20 to 3.21 in /cmd/fluent-watcher/fluentd (#1434) @dependabot
- Bump arm64v8/ruby from 3.3-slim-bullseye to 3.4-slim-bullseye in /cmd/fluent-watcher/fluentd (#1443) @dependabot
- Bump actions/checkout from 3 to 4 (#1394) @dependabot
- Bump github.com/fsnotify/fsnotify from 1.7.0 to 1.8.0 (#1396) @dependabot
- Bump golang.org/x/net from 0.28.0 to 0.33.0 (#1467) @dependabot

### BUGFIX

- Fix missing cloudAuth/cloudId pair inserts (#1463) @penekk

## New Contributors

- @harshvora10101 made their first contribution in #1360
- @btalakola made their first contribution in #1382
- @yilmazo made their first contribution in #1411
- @thapabishwa made their first contribution in #1458
- @CharlieR-o-o-t made their first contribution in #1460
- @penekk made their first contribution in #1462
- @truongnht made their first contribution in #1466

## 3.2.0 / 2024-09-21

### Features

- Expose args and command attributes for FluentBit CRD (#1350) @reegnz
- Add option to disable operator resources in Helm chart (#1348) @jiuxia211
- Support lua filter in namespaced CRD (#1342) @cw-Guo
- Add cloudAuthSecret & awsAuthSecret (#1338) @cw-Guo
- Add exec wasi input plugin (#1326) @jiuxia211
- Add wasm filter piugin (#1325) @jiuxia211
- Expose lua filter type_array_key parameter (#1323) @reegnz
- Support storage.total_limit_size in syslog plugin (#1318) @jk-mob
- Expose fluentbit init-container values in helm chart (#1320) @RajatPorwal5
- Add logs_body_key_attributes option for OpenTelemetry output plugin (#1322) @LKummer
- Add log to metrics plugin (#1305) @Athishpranav200

### ENHANCEMENT

- Update fluentbit to 3.1.7 (#1329) @jiuxia211

### BUGFIX

- Fix assignment to entry in nil map when --watch-namespaces flag is provided (#1334) @alexandrevilain
- Fix annotations too long issue (#1309) @cw-Guo

## 3.1.0 / 2024-08-14

### Features

- Render ConfigMap only if key is not empty string (#1299) @dex4er
- Set explicit fluent-bit name label selector (#1293) @rmvangun
- Allow fluent-operator to watch Kubernetes events (#1277) @thomasgouveia
- add fluent bit config-reload via HTTP (#1286) @jiuxia211
- feat(fluentbit): add fluentbit input_udp plugin (#1267) @cw-Guo
- add tag and tag_from_uri for opentelemetry input plugin (#1255) @smallc2009
- add compression to opensearch output plugin (#1258) @smallc2009
- Support for patch release tags. (#1246) @joshuabaird
- Add missing fluent-bit config parameters (#1244) @reegnz

### ENHANCEMENT

- build(deps): Bump github.com/onsi/gomega from 1.33.1 to 1.34.1 (#1275) @dependabot
- build(deps): Bump github.com/go-logr/logr from 1.4.1 to 1.4.2 (#1271) @dependabot
- build(deps): Bump k8s.io/klog/v2 from 2.120.1 to 2.130.1 (#1272) @dependabot
- build(deps): Bump aquasecurity/trivy-action from 0.23.0 to 0.24.0 (#1279) @dependabot
- build(deps): Bump golang.org/x/sync from 0.7.0 to 0.8.0 (#1283) @dependabot
- build(deps): Bump golang from 1.22.5-alpine3.19 to 1.22.6-alpine3.19 in /cmd/fluent-manager (#1288) @dependabot
- build(deps): Bump sigs.k8s.io/controller-runtime from 0.18.4 to 0.18.5 (#1295) @dependabot
- build(deps): Bump golang from 1.22.5-alpine3.19 to 1.22.6-alpine3.19 in /cmd/fluent-watcher/fluentbit (#1289) @dependabot
- build(deps): Bump golang from 1.22.5-alpine3.19 to 1.22.6-alpine3.19 in /cmd/fluent-watcher/fluentd (#1290) @dependabot
- update fluentbit to v3.1.4 (#1282) @jiuxia211
- Mention multiline parser CRDs in README (#1280) @reegnz
- build(deps): Bump github.com/go-openapi/errors from 0.20.4 to 0.22.0 (#1274) @dependabot
- build(deps): Bump golang from 1.22.0-alpine3.19 to 1.22.5-alpine3.19 in /cmd/fluent-watcher/fluentd (#1260) @dependabot
- build(deps): Bump golang from 1.22.0-alpine3.19 to 1.22.5-alpine3.19 in /cmd/fluent-watcher/fluentbit (#1261) @dependabot
- build(deps): Bump golang from 1.22.0-alpine3.19 to 1.22.5-alpine3.19 in /cmd/fluent-manager (#1262) @dependabot
- build(deps): Bump fluent/fluent-bit from 3.1.3-debug to 3.1.4-debug in /cmd/fluent-watcher/fluentbit (#1266) @dependabot
- build(deps): Bump fluent/fluent-bit from 3.1.2-debug to 3.1.3-debug in /cmd/fluent-watcher/fluentbit (#1245) @dependabot
- build(deps): Bump golang from 1.22.3 to 1.22.4 in /docs/best-practice/forwarding-logs-via-http (#1227) @dependabot
- build(deps): Bump k8s.io/client-go from 0.26.3 to 0.30.3 (#1254) @dependabot
- build(deps): bump k8s.io/client-go, k8s.io/apimachinery, k8s.io/api, … (#1251) @jiuxia211
- Bump fluentbit to 3.1.2. (#1240) @joshuabaird
- build(deps): Bump fluent/fluent-bit from 3.1.0-debug to 3.1.2-debug in /cmd/fluent-watcher/fluentbit (#1238) @dependabot

### BUGFIX

- bug: Allows to render net properties for outputs (#1298) @dex4er
- BUG: re-add accidently removed flag.Parse (#1294) @developer-guy
- Fix service monitor label selector scope (#1284) @rmvangun
- fix(rbac): Add missing rbac rules for namespaced fluentbit (#1265) @alexandrevilain
- Fixes "build fluent operator" CI workflow (#1263) @joshuabaird
- Fixes fluentd/fluent-bit image build CI workflows (#1259) @joshuabaird
- Fix release name on fluentbit output loki (#1248) @yildizozan

## 3.0.0 / 2024-07-09

### Features

- Feat: add daemonset terminationGracePeriodSeconds (#1204) @smallc2009
- Feat: add kubernetes events input plugin (#1209) @smallc2009
- Feat: support yaml config file (#1208) @cw-Guo
- Feat(helm): respect helm release namespace setting (#1214) @reegnz
- Feat: Adding Fluentbit's unified networking interface (#1217) @localleon
- Feat: add elasticsearch options (#1220) @bakervos
- Feat: add rbacRules to values.yaml with events watching as fixed permissions (#1223) @SvenThies
- Feat: add 'sslVerify' to opensearch output (#1226) @zmw85
- Feat: add bearer token auth for loki (#1224) @raynay-r

### ENHANCEMENT

- Adding instructions on how to set run operator for developement (#1216) @localleon
- Templatize ServiceMonitor (#1218) @smallc2009
- Update fluentbit to v3.1.0 (#1233) @wenchajun
- build(deps): Bump docker/build-push-action from 5 to 6 (#1228) @dependabot
- build(deps): Bump aquasecurity/trivy-action from 0.21.0 to 0.23.0 (#1229) @dependabot

### BUGFIX

- Fix: rewrite tag nil pointer reference (#1232) @cw-Guo
- Fix: fix the parsing of disableComponentControllers in helm (#1222) @mritunjaysharma394
- Fix(fluentd): add securityContext and podSecurityContext in values.yaml (#1230) @SvenThies
- Fix: fix parsersfile default parsers.conf path (#1225) @cw-Guo

## 2.9.0 / 2024-06-13

### Features

- Feat: Support elastic_data_stream (#1190) @fschlager
- Feat: Add storage total limit size to es plugin (#1196) @smallc2009
- Feat: Add tag parameter to forward output plugin (#1167) @fschlager
- Feat: Support datadog plugin api key to allow for secret injection (#1070) @nitintecg
- Feat: Add cloudId and cloudAuth parameters to elastic (#1169) @fschlager

### ENHANCEMENT

- Updates setup manifests to be compatible with v2.8.0 (#1161) @joshuabaird
- Update index.md (#1180) @lansaloni
- Upgrade fluentd to 1.17.0. (#1198) @joshuabaird
- Update fluentbit to v3.0.7 (#1199) @joshuabaird
- Update references to fluentd:1.17.0 image (#1200) @joshuabaird
- build(deps): Bump golang.org/x/net from 0.17.0 to 0.23.0 (#1140) @dependabot
- build(deps): Bump helm/kind-action from 1.9.0 to 1.10.0 (#1156) @dependabot
- build(deps): Bump alpine from 3.19 to 3.20 in /cmd/fluent-watcher/fluentd (#1179) @dependabot
- build(deps): Bump golang from 1.22.0 to 1.22.3 in /docs/best-practice/forwarding-logs-via-http (#1191) @dependabot
- build(deps): Bump aquasecurity/trivy-action from 0.13.1 to 0.21.0 (#1192) @dependabot

### BUGFIX

- Fix: missing inputs and clusterInputs CRDs in setup.yaml (#1144) (#1145) @antrema
- Fix: bugfix namespaced filters (#1143) @MarkusFreitag
- Fix: fix release cycles for fluentd and fluentbit images manually. (#1183) @sarathchandra24
- Fix: fix quotes for disable-component-controller argument string in fluent-operator deployment template. (#1160) @nickytd
- Fix: fix fluentd path issues. (#1195) @sarathchandra24
- Fix: fix fluent-bit image name. (#1201) @joshuabaird

## 2.8.0 / 2024-04-22

### Features

- Feat: feat: add multiline parser support for fluentbit (#1100) @ksdpmx
- Feat: feat: enforce Fluentd tests (#1110) @antrema
- Feat: feat: implement SecretLoader as interface and enforce Fluentd tests u… (#1109) @antrema
- Feat: namespaced tag re_emitter parameters support (#1085) @chrono2002
- Feat: LUA filter inline code support (#1081) @chrono2002
- Feat: feat(disableLogVolumes): expose disableLogVolumes in helm chart (#1082) @L1ghtman2k
- Feat: feat(chart): Add operator.extraArgs to add extra args to the fluent-operator container (#1074) @alexandrevilain
- Feat: feat(tls): allow overwriting tls for s3  (#1078) @L1ghtman2k
- Feat: feat: Emitter storage_type and mem_buf_limit config (#1069) @chrono2002
- Feat: Rework fluent-bit-watcher and make use of the hot-reload mechanism (#1051) @markusthoemmes
- Feat: Add clusterinput tail resource to support setting bufferMaxSize (#1052) @opencmit2
- Feat: feat(1062): Configure logLevel in ClusterFluentBitConfig (#1063) @dennis-ge
- Feat: feat: support s3 server side encryption (#1039) @cw-Guo

### ENHANCEMENT

- Helm multiline passer template and usage (#1138) @onecer
- make lua scripts `code` and `script` optional (#1129) @onecer
- MultilineParser achieve an effect similar to embedding by using anonymous structs (#1133) @onecer
- helm template fluentbit output es support logstashPrefixKey (#1119) @onecer
- Add Profile field for Fluent bit S3 output (#1127) @jeff303
- auto gen plugins documentation directory (#1121) @onecer
- build(deps): Bump google.golang.org/protobuf from 1.28.1 to 1.33.0 (#1091) @dependabot
- build(deps): Bump actions/cache from 3 to 4 (#1055) @dependabot
- build(deps): Bump actions/setup-go from 4 to 5 (#1034) @dependabot
- build(deps): Bump sigs.k8s.io/yaml from 1.3.0 to 1.4.0 (#1028) @dependabot
- build(deps): Bump golang from 1.21.4 to 1.22.0 in /docs/best-practice/forwarding-logs-via-http (#1077) @dependabot
- build(deps): Bump helm/kind-action from 1.8.0 to 1.9.0 (#1076) @dependabot
- build(deps): Bump azure/setup-helm from 3 to 4 (#1075) @dependabot
- build(deps): Bump actions/setup-python from 4 to 5 (#1035) @dependabot
- build(deps): Bump arm64v8/ruby from 3.2-slim-bullseye to 3.3-slim-bullseye in /cmd/fluent-watcher/fluentd (#1032) @dependabot
- HostNetwork DNS Policy (#1101) @opp-svega
- adds servicemonitor (#1089) @chrono2002
- push image to multiple registry (#1079) @sarathchandra24
- fluentbit/output/elasticsearch: Add writeOperation option (#1080) @icy
- Bump fluent-bit to 2.2.2 (#1053) @markusthoemmes
- Bump Golang version to 1.21 and replace slice utils with stdlib (#1049) @markusthoemmes
- Add stackdriver output to the Helm Chart (#1040) @UgurcanAkkok

### BUGFIX

- Fix: Passing variables to es output config is fixed (#1099) @aido93
- fix fluent-operator clusterrole in manifests directory (#1098) @Cajga
- fix: delete remaining debug traces (#1107) @antrema
- fix: ordered fluentd-config #1104 (#1106) @antrema
- Fix 1090: Path issues (#1093) @sarathchandra24
- fix: no Path default value for memory buffer #1103 (#1105) @antrema
- fix build args fluentd image (#1095) @sarathchandra24
- Fix: attach latest tag to images - DockerHub synchronize with GHCR (#1086) @sarathchandra24
- fix: incorrect fields in syslog input plugin parameters (#1084) @lukasboettcher
- fix(fluentd): Use custom plugin content for hash generation (#1059) @MisterMX
- fix prometheus-remote-write-edge templates error (#1036) @JiaweiGithub
- Fixed the .spec.loki.tls map rendering in FluentBit loki ClusterOutput (#1031) @isemichastnov

## 2.7.0 / 2023-12-19

### Features

- Feat: Add copy output plugin for fluentd #1017 (#1018) @antrema

### ENHANCEMENT

- Update fluentd-filter-kafka.yaml (#1016) @blackshy
- build(deps): Bump alpine in /cmd/fluent-watcher/fluentd (#1014) @dependabot
- build(deps): Bump golang in /cmd/fluent-manager (#1009) @dependabot
- build(deps): Bump helm/chart-testing-action from 2.6.0 to 2.6.1 (#1006) @dependabot
- build(deps): Bump golang in /docs/best-practice/forwarding-logs-via-http (#1004) @dependabot
- build(deps): Bump github.com/onsi/gomega from 1.28.0 to 1.30.0 (#1002) @dependabot

### BUGFIX

- fix: Add SSL/TLS settings feature for fluentd output Elasticsearch #418 (#1011) @antrema
- fix: Add RBAC permissions for input and clusterinput (#1019) @MisterMX
- fix: missing CRD entries and documentation #1020 (#1022) @antrema

## 2.6.0 / 2023-11-22

### Features

- Feat: Add fluentbit nginx  plugin (#924)
- Feat: Add fluentbit statsd plugin (#925)
- Feat: Add fluentbit syslog plugin (#931)
- Feat: Add fluentbit tcp plugin (#936)
- Feat: Add in_sample plugin to fluentd to facilitate tests. (#937)
- Feat: Adds the fluent-plugin-prometheus plugin to fluentd. (#966)
- Feat: Adds fluentd monitor_agent input plugin (#967)
- Feat(fluentd): Input plugin CRs (#972)
- Feat: Add readiness & liveness probe for fluentd (#980)

### ENHANCEMENT

- Add fluentbit daemonset hostPath toggle (#926)
- Allow setting dnsPolicy for fluentbit (#951)
- Default cri parser should contain Time_Keep On, otherwise no time tag exists at output (#958)
- Chore: Replace deprecated command with environment file (#970)
- Upgrade chart-testing-action to v2.6.0 (#976)
- Allow passing env vars using the chart (#977)
- Sort custom resources by metadata.name (#988)
- Bump fluentbit to 2.2.0 (#994)
- build(deps): Bump docker/login-action from 2 to 3 (#939)
- build(deps): Bump actions/checkout from 3 to 4 (#940)
- build(deps): Bump docker/setup-buildx-action from 2 to 3 (#941)
- build(deps): Bump github.com/onsi/gomega from 1.27.10 to 1.28.0 (#944)
- build(deps): Bump golang from 1.21.1-alpine3.17 to 1.21.2-alpine3.17 in /cmd/fluent-manager (#950)
- build(deps): Bump golang from 1.21.2-alpine3.17 to 1.21.3-alpine3.17 in /cmd/fluent-manager (#953)
- build(deps): Bump golang.org/x/net from 0.14.0 to 0.17.0 (#954)
- build(deps): Bump github.com/fsnotify/fsnotify from 1.6.0 to 1.7.0 (#981)
- build(deps): Bump golang from 1.21.1 to 1.21.3 in /docs/best-practice/forwarding-logs-via-http (#974)
- build(deps): Bump golang from 1.21.3-alpine3.17 to 1.21.4-alpine3.17 in /cmd/fluent-manager (#983)
- build(deps): Bump github.com/go-logr/logr from 1.2.4 to 1.3.0 (#989)

### BUGFIX

- Fix(doc): Fluentbit splunk output docs (#935) @Macbet
- Fix: Add parserSelector to clusterFluentBitConfig chart templates (#956)
- Fix: fd record transformer parameters (#960)
- Fix: fluentd in_http plugin keepalive_timeout option (#968)
- Fix: fluentd parser keep_time_key (#987)

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
- [ENHANCEMENT] Support File, TCP, HTTP outputs

## 0.2.0 / 2020-08-27

- [CHANGE] Rewrite relevant CRDs. They are backwards incompatible with v0.1.0
- [CHANGE] Use kubebuilder as the building framework

## 0.1.0 / 2020-02-17

This is the first release of fluentbit operator.
