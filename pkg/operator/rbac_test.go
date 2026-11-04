package operator

import (
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"
)

// Each instance writes its own rbacRules into its ClusterRole, so two instances of a component must not share
// one: the instance reconciled last would replace the other's rules.
func TestMakeRBACObjectsUsesAClusterRolePerInstance(t *testing.T) {
	events := []rbacv1.PolicyRule{
		{APIGroups: []string{""}, Resources: []string{"events"}, Verbs: []string{"list", "watch"}},
	}

	eventsRole, _, eventsBinding := MakeRBACObjects("events", "logging", "fluent-bit", events, nil)
	logsRole, _, logsBinding := MakeRBACObjects("logs", "logging", "fluent-bit", nil, nil)

	if eventsRole.Name == logsRole.Name {
		t.Fatalf("both instances use ClusterRole %q", eventsRole.Name)
	}
	if eventsBinding.RoleRef.Name != eventsRole.Name || logsBinding.RoleRef.Name != logsRole.Name {
		t.Fatalf("bindings refer to %q and %q, want %q and %q",
			eventsBinding.RoleRef.Name, logsBinding.RoleRef.Name, eventsRole.Name, logsRole.Name)
	}
	if len(eventsRole.Rules) != 2 || len(logsRole.Rules) != 1 {
		t.Fatalf("rules = %d and %d, want the default rule plus each instance's own rules",
			len(eventsRole.Rules), len(logsRole.Rules))
	}

	fluentdRole, _, _ := MakeRBACObjects("events", "logging", "fluentd", nil, nil)
	if fluentdRole.Name == eventsRole.Name {
		t.Fatalf("a Fluentd and a FluentBit named the same share ClusterRole %q", fluentdRole.Name)
	}
}
