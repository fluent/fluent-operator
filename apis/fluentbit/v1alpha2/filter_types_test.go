package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/filter"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	. "github.com/onsi/gomega"
)

func TestFilterList_Load(t *testing.T) {
	testcases := []struct {
		name     string
		input    []Filter
		expected string
	}{
		{
			name: "a single filteritem",
			input: []Filter{
				{
					Spec: FilterSpec{
						FilterItems: []FilterItem{
							{
								Parser: &filter.Parser{
									KeyName: "log",
									Parser:  "first-parser",
								},
							},
						},
					},
				},
			},
			expected: `[Filter]
    Name    parser
    Key_Name    log
    Parser    first-parser-d41d8cd98f00b204e9800998ecf8427e
`,
		},
		{
			name: "a single filteritem, with multiple plugins",
			input: []Filter{
				{
					Spec: FilterSpec{
						FilterItems: []FilterItem{
							{
								Kubernetes: &filter.Kubernetes{
									KubeTagPrefix: "custom-tag",
								},
								Parser: &filter.Parser{
									KeyName: "log",
									Parser:  "first-parser",
								},
							},
						},
					},
				},
			},
			expected: `[Filter]
    Name    kubernetes
    Kube_Tag_Prefix    d41d8cd98f00b204e9800998ecf8427e.custom-tag
[Filter]
    Name    parser
    Key_Name    log
    Parser    first-parser-d41d8cd98f00b204e9800998ecf8427e
`,
		},
		{
			name: "multiple filteritems",
			input: []Filter{
				{
					Spec: FilterSpec{
						FilterItems: []FilterItem{
							{
								Kubernetes: &filter.Kubernetes{
									KubeTagPrefix: "custom-tag",
								},
								Parser: &filter.Parser{
									KeyName: "log",
									Parser:  "first-parser",
								},
							},
							{
								Parser: &filter.Parser{
									KeyName:     "msg",
									Parser:      "second-parser",
									ReserveData: utils.ToPtr(true),
								},
							},
							{
								Parser: &filter.Parser{
									KeyName:     "msg",
									Parser:      "third-parser",
									ReserveData: utils.ToPtr(true),
								},
							},
						},
					},
				},
			},
			expected: `[Filter]
    Name    kubernetes
    Kube_Tag_Prefix    d41d8cd98f00b204e9800998ecf8427e.custom-tag
[Filter]
    Name    parser
    Key_Name    log
    Parser    first-parser-d41d8cd98f00b204e9800998ecf8427e
[Filter]
    Name    parser
    Key_Name    msg
    Parser    second-parser-d41d8cd98f00b204e9800998ecf8427e
    Reserve_Data    true
[Filter]
    Name    parser
    Key_Name    msg
    Parser    third-parser-d41d8cd98f00b204e9800998ecf8427e
    Reserve_Data    true
`,
		},
		{
			name: "ordinal-based sorting",
			input: []Filter{
				{
					Spec: FilterSpec{
						Ordinal: 10,
						FilterItems: []FilterItem{
							{
								Parser: &filter.Parser{
									KeyName: "msg",
									Parser:  "parser-two",
								},
							},
						},
					},
				},
				{
					Spec: FilterSpec{
						Ordinal: -10,
						FilterItems: []FilterItem{
							{
								Kubernetes: &filter.Kubernetes{
									KubeTagPrefix: "custom-tag",
								},
							},
						},
					},
				},
				{
					Spec: FilterSpec{
						FilterItems: []FilterItem{
							{
								Parser: &filter.Parser{
									KeyName: "log",
									Parser:  "parser-one",
								},
							},
						},
					},
				},
			},
			expected: `[Filter]
    Name    kubernetes
    Kube_Tag_Prefix    d41d8cd98f00b204e9800998ecf8427e.custom-tag
[Filter]
    Name    parser
    Key_Name    log
    Parser    parser-one-d41d8cd98f00b204e9800998ecf8427e
[Filter]
    Name    parser
    Key_Name    msg
    Parser    parser-two-d41d8cd98f00b204e9800998ecf8427e
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			sl := plugins.NewSecretLoader(nil, "testnamespace")

			fl := FilterList{
				Items: tc.input,
			}
			renderedFilterList, err := fl.Load(sl)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(renderedFilterList).To(Equal(tc.expected))
		})
	}
}

func TestFilterList_LoadAsYaml(t *testing.T) {
	testcases := []struct {
		name     string
		input    Filter
		expected string
	}{
		{
			name: "a single filteritem",
			input: Filter{
				Spec: FilterSpec{
					FilterItems: []FilterItem{
						{
							Parser: &filter.Parser{
								KeyName: "log",
								Parser:  "first-parser",
							},
						},
					},
				},
			},
			expected: `  - name: parser
    key_name: log
    parser: first-parser-d41d8cd98f00b204e9800998ecf8427e
`,
		},
		{
			name: "a single filteritem, with multiple plugins",
			input: Filter{
				Spec: FilterSpec{
					FilterItems: []FilterItem{
						{
							Kubernetes: &filter.Kubernetes{
								KubeTagPrefix: "custom-tag",
							},
							Parser: &filter.Parser{
								KeyName: "log",
								Parser:  "first-parser",
							},
						},
					},
				},
			},
			expected: `  - name: kubernetes
    kube_tag_prefix: d41d8cd98f00b204e9800998ecf8427e.custom-tag
  - name: parser
    key_name: log
    parser: first-parser-d41d8cd98f00b204e9800998ecf8427e
`,
		},
		{
			name: "multiple filteritems",
			input: Filter{
				Spec: FilterSpec{
					FilterItems: []FilterItem{
						{
							Kubernetes: &filter.Kubernetes{
								KubeTagPrefix: "custom-tag",
							},
							Parser: &filter.Parser{
								KeyName: "log",
								Parser:  "first-parser",
							},
						},
						{
							Parser: &filter.Parser{
								KeyName:     "msg",
								Parser:      "second-parser",
								ReserveData: utils.ToPtr(true),
							},
						},
						{
							Parser: &filter.Parser{
								KeyName:     "msg",
								Parser:      "third-parser",
								ReserveData: utils.ToPtr(true),
							},
						},
					},
				},
			},
			expected: `  - name: kubernetes
    kube_tag_prefix: d41d8cd98f00b204e9800998ecf8427e.custom-tag
  - name: parser
    key_name: log
    parser: first-parser-d41d8cd98f00b204e9800998ecf8427e
  - name: parser
    key_name: msg
    parser: second-parser-d41d8cd98f00b204e9800998ecf8427e
    reserve_data: true
  - name: parser
    key_name: msg
    parser: third-parser-d41d8cd98f00b204e9800998ecf8427e
    reserve_data: true
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			sl := plugins.NewSecretLoader(nil, "testnamespace")

			fl := FilterList{
				Items: make([]Filter, 1),
			}
			fl.Items[0] = tc.input

			renderedFilterList, err := fl.LoadAsYaml(sl, 0)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(renderedFilterList).To(Equal(tc.expected))
		})
	}
}
