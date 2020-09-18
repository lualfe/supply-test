package services

import "testing"

func TestDependsAA(t *testing.T) {
	rs := NewRuleSet()
	rs.AddDep("a", "a")
	if !rs.IsCoherent() {
		t.Error("Expected the rule set to be considered coherent, but it was not")
	}
}

func TestDependsAB_BA(t *testing.T) {
	rs := NewRuleSet()
	rs.AddDep("a", "b")
	rs.AddDep("b", "a")
	if !rs.IsCoherent() {
		t.Error("Expected the rule set to be considered coherent, but it was not")
	}
}

func TestExclusiveAB(t *testing.T) {
	rs := NewRuleSet()
	rs.AddDep("a", "b")
	rs.AddConflict("a", "b")
	if rs.IsCoherent() {
		t.Error("Expected the rule set to be considered not coherent, but it was considered coherent")
	}
}

func TestExclusiveAB_BC(t *testing.T) {
	rs := NewRuleSet()
	rs.AddDep("a", "b")
	rs.AddDep("b", "c")
	rs.AddConflict("a", "c")
	if rs.IsCoherent() {
		t.Error("Expected the rule set to be considered not coherent, but it was considered coherent")
	}
}

func TestDeepDeps(t *testing.T) {
	rs := NewRuleSet()
	rs.AddDep("a", "b")
	rs.AddDep("b", "c")
	rs.AddDep("c", "d")
	rs.AddDep("d", "e")
	rs.AddDep("a", "f")
	rs.AddConflict("e", "f")
	if rs.IsCoherent() {
		t.Error("Expected the rule set to be considered not coherent, but it was considered coherent")
	}
}