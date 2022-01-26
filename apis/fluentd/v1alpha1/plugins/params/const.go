package params

type PluginName string

const (
	InputPlugin     PluginName = "input"
	ForwardPlugin   PluginName = "forward"
	HttpPlugin      PluginName = "http"
	TransportPlugin PluginName = "transport"
	InjectPlugin    PluginName = "inject"
	FormatPlugin    PluginName = "format"
	TimePlugin      PluginName = "time"
	SecurityPlugin  PluginName = "security"
	AuthPlugin      PluginName = "auth"
	UserPlugin      PluginName = "user"
	ClientPlugin    PluginName = "client"
	ServerPlugin    PluginName = "server"

	FilterPlugin            PluginName = "filter"
	GrepPlugin              PluginName = "grep"
	RecordTransformerPlugin PluginName = "record_transformer"
	ParserPlugin            PluginName = "parser"
	StdoutPlugin            PluginName = "stdout"

	ReLabelPlugin       PluginName = "relabel"
	LabelPlugin         PluginName = "label"
	LabelRouterPlugin   PluginName = "label_router"
	S3Plugin            PluginName = "s3"
	KafkaPlugin         PluginName = "kafka2"
	ElasticsearchPlugin PluginName = "elasticsearch"
	MatchPlugin         PluginName = "match"
	BufferPlugin        PluginName = "buffer"

	BufferTag    string = "tag"
	LabelTag     string = "tag"
	MatchTag     string = "tag"
	FilterTag    string = "tag"
	ProtocolName string = "protocol"
	// Default interval whitespaces between parent and child
	IntervalWhitespaces string = "  "
	DefaultFmtExpr      string = "  %s"
)

var (
	DeftaultTag = "**"
)
