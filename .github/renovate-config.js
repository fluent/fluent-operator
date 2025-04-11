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
  repositories: ["fluent/fluent-operator"],
  enabledManagers: ["custom.regex"],
  extends: ["config:recommended"],
  customManagers: [
    {
      customType: "regex",
      matchStringsStrategy: "any",
      fileMatch: [
        "charts/fluent-operator/Chart.yaml",
        "charts/fluent-operator/charts/fluent-bit-crds/Chart.yaml",
        "charts/fluent-operator/charts/fluentd-crds/Chart.yaml",
        "charts/fluent-operator/values.yaml",
        "cmd/fluent-watcher/fluentbit/VERSION",
        "config/.*\\.yaml",
        "docs/.*\\.yaml",
        "manifests/.*\\.yaml",
      ],
      matchStrings: [
        '# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>\\S+?)\\s+tag:\\s+"(?<currentValue>.+?)"\\s+?',
        "# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>\\S+?)\\s+tag:\\s+(?<currentValue>.+?)\\s+?",
        "# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>.+?)\\s+image:\\s*(?:.?)*:(?<currentValue>.*?)\\s+?",
        "# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>.+?)\\s+version=\\n(?<currentValue>.+?)\\n",
        '# renovate:\\s+datasource=(?<datasource>\\S+?)\\s+depName=(?<depName>.+?)\\s+appVersion:\\s+"(?<currentValue>.+?)"\\n',
      ],
      extractVersionTemplate: "^v(?<version>.*)$",
    },
  ],
};
