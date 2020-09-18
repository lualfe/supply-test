package services

import (
	"log"
)

type options []option

type option struct {
	value string
	dependencies *options
	mutuallyExclusiveWith *options
}

//RuleSet represents a map of options
type RuleSet map[string]*option

//NewRuleSet instantiates a new RuleSet
func NewRuleSet() RuleSet {
	return make(map[string]*option)
}

//AddOption adds a new option to a RuleSet
func (r RuleSet) AddOption(value string) {
	r[value] = &option{
		value:                 value,
		dependencies:          &options{},
		mutuallyExclusiveWith: &options{},
	}
}

//AddDep adds a list of dependencies to a option in a RuleSet
func (r RuleSet) AddDep(value string, dependsOf ...string) {
	if _, ok := r[value]; !ok {
		r.AddOption(value)
	}
	for _, v := range dependsOf {
		if _, ok := r[v]; !ok {
			r.AddOption(v)
		}
	}

	for _, v := range dependsOf {
		if v == value {
			log.Println(value, "cannot depend on itself")
			return
		}
		*r[value].dependencies = append(*r[value].dependencies, *r[v])
	}
}

//AddConflict adds a new conflict (mutual exclusivity) to the RuleSet
func (r RuleSet) AddConflict(option, mutuallyExclusiveWith string) {
	*r[option].mutuallyExclusiveWith = append(*r[option].mutuallyExclusiveWith, *r[mutuallyExclusiveWith])
	*r[mutuallyExclusiveWith].mutuallyExclusiveWith = append(*r[mutuallyExclusiveWith].mutuallyExclusiveWith, *r[option])
}

//IsCoherent checks if a RuleSet is coherent
func (r RuleSet) IsCoherent() bool {
	allMutuallyExclusives := r.allMutuallyExclusives()
	if r.validDepsAndExclusives(allMutuallyExclusives) {
		return true
	}
	return false
}

func (r RuleSet) allMutuallyExclusives() (mutuallyExclusives []option) {
	for _, v := range r {
		mutuallyExclusives = append(mutuallyExclusives, *v.mutuallyExclusiveWith...)
	}
	return
}

func (r RuleSet) validDepsAndExclusives(mutuallyExclusives []option) bool {
	for _, v := range r {
		if !v.isCoherentWith(*v, mutuallyExclusives) {
			return false
		}
	}
	return true
}

func (origin option) isCoherentWith(option option, allMutuallyExclusives options) bool {
	for _, o := range *option.dependencies {
		if option.circularDependencyOfAny(*o.dependencies) {
			log.Println("circular dependency")
			return true
		}

		if o.hasInvalidMutuallyExclusives(allMutuallyExclusives) {
			return false
		}

		if origin.mutuallyExclusiveWith.contains(o.value) {
			return false
		}

		if !origin.isCoherentWith(o, allMutuallyExclusives) {
			return false
		}
	}
	return true
}

func (origin option) circularDependencyOfAny(dependencies options) bool {
	for _, downDep := range dependencies {
		if downDep.value == origin.value {
			return true
		}
	}
	return false
}

func (origin option) hasInvalidMutuallyExclusives(mutuallyExclusives options) bool {
	for _, v := range mutuallyExclusives {
		if origin.value == v.value {
			return true
		}
	}
	return false
}

func (o options) contains(value string) bool {
	for _, v := range o {
		if value == v.value {
			return true
		}
	}
	return false
}