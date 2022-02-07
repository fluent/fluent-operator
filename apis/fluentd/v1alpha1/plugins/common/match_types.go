package common

type Match struct {
	// buffer section
	Buffer *Buffer `json:"buffer,omitempty"`
	// format section
	Format *Format `json:"format,omitempty"`
	// inject section
	Inject *Inject `json:"inject,omitempty"`
}
