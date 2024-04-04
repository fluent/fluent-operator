package params

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fluent/fluent-operator/v2/pkg/utils"
)

type PluginStore struct {
	// The plugin name
	Name string
	// The key-value pairs
	Store map[string]string
	// The child plugins mounted here
	Childs []*PluginStore
	// The prefix whitespaces before this plugin mounted the parent plugin
	PrefixWhitespaces string
	// The flag whether to ignore the path field in buffer
	IgnorePath bool
	Content    string
}

func NewPluginStore(name string) *PluginStore {
	return &PluginStore{
		Name:   name,
		Store:  make(map[string]string),
		Childs: make([]*PluginStore, 0),
	}
}

func (ps *PluginStore) InsertPairs(key, value string) {
	ps.Store[key] = value
}

// The @type parameter specifies the type of the plugin
func (ps *PluginStore) InsertType(value string) {
	ps.InsertPairs("@type", value)
}

// SetIgnorePath will ignore the buffer path
func (ps *PluginStore) SetIgnorePath() {
	ps.IgnorePath = true
}

// If one label section contains a match section,
// we consider that the match section is the child of label section
func (ps *PluginStore) InsertChilds(childs ...*PluginStore) {
	if len(childs) == 0 {
		return
	}

	for _, child := range childs {
		if child != nil {
			ps.Childs = append(ps.Childs, child)
		}
	}
}

// The total hash string for this plugin store
func (ps *PluginStore) Hash() string {
	c := NewPluginStore(ps.Name)
	isNotCopyOutput := ps.Store["@type"] != "copy"

	// We must consider the tag when the output is a Copy one
	// as copy is a "flag" output: it can exist identical outputs with different tag
	for k, v := range ps.Store {
		if k == "@id" || (k == "tag" && isNotCopyOutput) {
			continue
		}
		c.Store[k] = v
	}

	// Custom plugins don't have stored properties but only contain the config
	// as plain text.
	// Set the content here to avoid generating the same hash code for all
	// custom plugins resulting in only one custom plugin being ever set for
	// one config.
	c.Content = ps.Content

	c.InsertChilds(ps.Childs...)
	return utils.HashCode(c.String())
}

// Returns tag value
func (ps *PluginStore) GetTag() string {
	return ps.Store["tag"]
}

// Returns @id value
func (ps *PluginStore) GetId() string {
	if value, ok := ps.Store["@id"]; !ok {
		return ""
	} else {
		return value
	}
}

// Returns the @label value string of this plugin store
func (ps *PluginStore) RouteLabel() string {
	if ps.Name != "route" {
		return ""
	}

	if value, ok := ps.Store["@label"]; !ok {
		return ""
	} else {
		return value
	}
}

func (ps *PluginStore) String() string {
	if ps == nil || ps.Name == "" {
		return ""
	}
	if ps.Content != "" {
		return ps.Content
	}
	var buf bytes.Buffer

	// Handles the head directive
	ps.processHead(&buf)

	// The body needs to be indented by two whitespace characters
	parentPrefixWhitespaces := ps.PrefixWhitespaces
	ps.setWhitespaces(parentPrefixWhitespaces + IntervalWhitespaces)
	ps.processBody(&buf)
	if len(ps.Childs) > 0 {
		sort.Sort(PluginStoreByNameById(ps.Childs))
		for _, child := range ps.Childs {
			child.setWhitespaces(ps.PrefixWhitespaces)
			buf.WriteString(child.String())
		}
	}

	// The tail must be indented in the same format as head directive
	ps.setWhitespaces(parentPrefixWhitespaces)
	ps.processTail(&buf)

	return buf.String()
}

func (ps *PluginStore) setWhitespaces(curentWhitespaces string) {
	ps.PrefixWhitespaces = curentWhitespaces
}

// write the head directive to buffer, i.e.: <match xx>
func (ps *PluginStore) processHead(buf *bytes.Buffer) {
	var head string
	switch PluginName(ps.Name) {
	case BufferPlugin:
		tag, ok := ps.Store[BufferTag]
		if ok {
			head = ps.headFmtSprintf(tag)
		}
	case MatchPlugin:
		head = ps.headFmtSprintf(ps.Store[MatchTag])
	case FilterPlugin:
		head = ps.headFmtSprintf(ps.Store[FilterTag])
	case TransportPlugin:
		head = ps.headFmtSprintf(ps.Store[ProtocolName])
	case LabelPlugin:
		head = ps.headFmtSprintf(ps.Store[LabelTag])
	default:
		head = fmt.Sprintf("%s<%s>\n", ps.PrefixWhitespaces, ps.Name)
	}
	buf.WriteString(head)
}

// processes the key-value pair body
func (ps *PluginStore) processBody(buf *bytes.Buffer) {
	var body string

	keys := make([]string, 0, len(ps.Store))
	for k := range ps.Store {
		// Don't add tag unless it is an input plugin
		if k == "tag" && ps.Name != "source" {
			continue
		}
		if ps.Name == string(BufferPlugin) && ps.IgnorePath {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		body += fmt.Sprintf("%s%s  %s\n", ps.PrefixWhitespaces, k, ps.Store[k])
	}

	buf.WriteString(body)
}

// write the tail directive to the buffer, i.e.: </match>
func (ps *PluginStore) processTail(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%s</%s>\n", ps.PrefixWhitespaces, ps.Name))
}

// decide to return the head directive with our without a filter - <match> or <match xx>
func (ps *PluginStore) headFmtSprintf(value string) string {
	if value != "" {
		return fmt.Sprintf("%s<%s %s>\n", ps.PrefixWhitespaces, ps.Name, value)
	}
	return fmt.Sprintf("%s<%s>\n", ps.PrefixWhitespaces, ps.Name)
}

// +kubebuilder:object:generate=false
type PluginStoreByName []*PluginStore

func (a PluginStoreByName) Len() int           { return len(a) }
func (a PluginStoreByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PluginStoreByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// +kubebuilder:object:generate=false
type PluginStoreByNameById []*PluginStore

func (a PluginStoreByNameById) Len() int      { return len(a) }
func (a PluginStoreByNameById) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PluginStoreByNameById) Less(i, j int) bool {
	if a[i].Name == a[j].Name {
		if a[i].GetTag() == "**" && a[j].GetTag() != "**" {
			return false
		}
		if a[i].GetTag() != "**" && a[j].GetTag() == "**" {
			return true
		}
		return a[i].GetId() < a[j].GetId()
	} else {
		return a[i].Name < a[j].Name
	}
}
