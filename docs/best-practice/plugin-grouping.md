

# Plugin Grouping

For fluenbit,input, filter, and output plugins are connected by label selectors. For input and output plugins, always create `ClusterInput` or `ClusterOutput` CRs for every plugin. Don't aggregate multiple inputs or outputs into one `ClusterInput` or `ClusterOutput` object, except you have a good reason to do so. Take the demo `logging stack` for example, we have one yaml file for each output.

However, for filter plugins, if you want a filter chain, the order of filters matters. You need to organize multiple filters into an array as the demo [logging stack](https://github.com/fluent/fluent-operator/blob/master/manifests/logging-stack/filter-kubernetes.yaml) suggests.For more info on various use cases of Fluent Operator Fluentd CRDs, you can refer to [Fluent-Operator-Walkthrough](https://github.com/kubesphere-sigs/fluent-operator-walkthrough#fluent-bit--fluentd-mode).



