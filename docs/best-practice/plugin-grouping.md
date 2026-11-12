

# Plugin Grouping

For fluenbit,input, filter, and output plugins are connected by label selectors. For input and output plugins, always create `ClusterInput` or `ClusterOutput` CRs for every plugin. Don't aggregate multiple inputs or outputs into one `ClusterInput` or `ClusterOutput` object, except you have a good reason to do so. Take the demo `logging stack` for example, we have one yaml file for each output.

However, for filter plugins, if you want a filter chain, the order of filters matters. You need to organize multiple filters into an array as the demo [logging stack](https://github.com/fluent/fluent-operator/blob/master/manifests/logging-stack/filter-kubernetes.yaml) suggests.For more info on various use cases of Fluent Operator Fluentd CRDs, you can refer to [Fluent-Operator-Walkthrough](https://github.com/kubesphere-sigs/fluent-operator-walkthrough#fluent-bit--fluentd-mode).

For Fluentd configs that select more than nine child plugins in one label section, set `.spec.pluginSortOrder: index` on `FluentdConfig` or `ClusterFluentdConfig` to use numeric-aware plugin ordering. The default `lexicographic` order is kept for backward compatibility.

Fluentd filters with the catch-all tag `**` are rendered after the other filters in a label section by default. To run a `**` filter between other filters, set `.spec.catchAllFilterOrder: id` on `FluentdConfig` or `ClusterFluentdConfig`, and filters are then ordered by `@id`, which follows the filter CR name and the position in its `filters` list. The @id also includes the namespace, so when ClusterFilters and namespaced Filters are mixed, the namespace is compared first: a Filter in `app` sorts before any ClusterFilter, and one in `default` after. If a CR has ten or more filters, set `pluginSortOrder: index` too so that `-10` sorts after `-9`. Matches with `**` are always rendered last, since a catch-all match would take the events meant for the matches after it.


