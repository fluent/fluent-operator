package common

import (
	"strings"
	"testing"
)

func TestServerParams(t *testing.T) {
	host := "fluentd-aggregator-0"
	name := "aggregator-0"
	port := "24224"
	standby := "true"
	weight := "60"
	s := &Server{
		Host:       &host,
		ServerName: &name,
		Port:       &port,
		Standby:    &standby,
		Weight:     &weight,
	}

	ps, err := s.Params(nil)
	if err != nil {
		t.Fatalf("Params returned an error: %v", err)
	}

	want := map[string]string{
		"host":    host,
		"name":    name,
		"port":    port,
		"standby": standby,
		"weight":  weight,
	}
	for key, value := range want {
		if got := ps.Store[key]; got != value {
			t.Errorf("%s = %q, want %q", key, got, value)
		}
	}
	if _, ok := ps.Store["server"]; ok {
		t.Errorf("unexpected parameter server in %q", ps.String())
	}
	if !strings.Contains(ps.String(), "<server>") {
		t.Errorf("rendered config has no <server> section: %q", ps.String())
	}
}
