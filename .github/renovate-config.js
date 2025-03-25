const helmRegex = {
  customType: "regex",
  datasourceTemplate: "helm",
  matchStringsStrategy: "combination",
};

const customRegex = {
  customType: "regex",
  matchStringsStrategy: "combination",
};

module.exports = {
  username: "renovate[bot]",
  gitAuthor: "Renovate Bot <bot@renovateapp.com>",
  onboarding: false,
  platform: "github",
  dryRun: null,
  repositories: ["truongnht/fluent-operator"],
  enabledManagers: ["custom.regex"],
  extends: ["config:recommended"],
  customManagers: [
    {
      customType: "regex",
      matchStringsStrategy: "any",
      fileMatch: [
        "charts/fluent-operator/values.yaml",
        "config/samples/fluentbit_v1alpha2_fluentbit.yaml",
        "docs/best-practice/forwarding-logs-via-http/deploy/fluentbit-fluentBit.yaml",
        "manifests/kubeedge/fluentbit-fluentbit-edge.yaml",
        "manifests/logging-stack/fluentbit-fluentBit.yaml",
        "manifests/quick-start/fluentbit.yaml",
        "manifests/regex-parser/fluentbit-fluentBit.yaml",
        "cmd/fluent-watcher/fluentbit/VERSION",
      ],
      matchStrings: [
        '# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>\\S+?)\\s+tag:\\s+"(?<currentValue>.+?)"\\s+?',
        "# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>.+?)\\s+image:\\s*(?:.?)*:(?<currentValue>.*?)\\s+?",
        "# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>.+?)\\s+version=\\n(?<currentValue>.+?)\\n",
      ],
      extractVersionTemplate: "^v(?<version>.*)$",
    },
  ],
};
