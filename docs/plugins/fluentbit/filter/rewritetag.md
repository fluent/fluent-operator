# RewriteTag

RewriteTag define a `rewrite_tag` filter, allows to re-emit a record under a new Tag. <br /> Once a record has been re-emitted, the original record can be preserved or discarded. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/rewrite-tag**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| rules | Defines the matching criteria and the format of the Tag for the matching record. The Rule format have four components: KEY REGEX NEW_TAG KEEP. | []string |
| emitterName | When the filter emits a record under the new Tag, there is an internal emitter plugin that takes care of the job. Since this emitter expose metrics as any other component of the pipeline, you can use this property to configure an optional name for it. | string |
| emitterMemBufLimit |  | string |
| emitterStorageType |  | string |
