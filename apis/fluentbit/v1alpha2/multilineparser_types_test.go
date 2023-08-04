package v1alpha2

import (
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/multilineparser"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var multilineParserExpected = `[MULTILINE_PARSER]
    Name    multilineparser_test0
    Type    regex
    Parser    go
    Key_Content    log
[MULTILINE_PARSER]
    Name    multilineparser_test1
    Type    regex
    Rule    "start_state" "/([a-zA-Z]+ \d+ \d+\:\d+\:\d+)(.*)/" "cont"
    Rule    "cont" "/^\s+at.*/" "cont"
`

func TestMultilineParserList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
	}

	goMultilineParser := &MultilineParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "MultilineParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "multilineparser_test0",
			Labels: labels,
		},
		Spec: MultilineParserSpec{
			MultilineParser: &multilineparser.MultilineParser{
				Type:       "regex",
				Parser:     "go",
				KeyContent: "log",
			},
		},
	}

	customMultilineParser := &MultilineParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "MultilineParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "multilineparser_test1",
			Labels: labels,
		},
		Spec: MultilineParserSpec{
			MultilineParser: &multilineparser.MultilineParser{
				Type: "regex",
				Rules: []multilineparser.Rule{
					{
						Start: "start_state",
						Regex: `/([a-zA-Z]+ \d+ \d+\:\d+\:\d+)(.*)/`,
						Next:  "cont",
					},
					{
						Start: "cont",
						Regex: `/^\s+at.*/`,
						Next:  "cont",
					},
				},
			},
		},
	}

	multilineparsers := MultilineParserList{
		Items: []MultilineParser{*goMultilineParser, *customMultilineParser},
	}

	mp, err := multilineparsers.Load(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(mp).To(Equal(multilineParserExpected))
}
