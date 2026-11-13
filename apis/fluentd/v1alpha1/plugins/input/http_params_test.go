package input

import (
	"strings"
	"testing"
)

func TestHttpRespondWithEmptyImg(t *testing.T) {
	id := "http-input"
	enabled := true
	i := &Input{
		Id:   &id,
		Http: &Http{RespondsWithEmptyImg: &enabled},
	}

	ps, err := i.Params(nil)
	if err != nil {
		t.Fatalf("Params returned an error: %v", err)
	}

	if got := ps.Store["respond_with_empty_img"]; got != "true" {
		t.Errorf("respond_with_empty_img = %q, want %q in %q", got, "true", ps.String())
	}
	if strings.Contains(ps.String(), "responds_with_empty_img") {
		t.Errorf("rendered config uses a parameter in_http does not know: %q", ps.String())
	}
}
