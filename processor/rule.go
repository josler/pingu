package processor

import "github.com/josler/pingu/core"

type Rule interface {
	Type() string
	Match(e *core.Event) bool
}

func ParseRules(jsonRule *core.JsonRule) Rule {
	switch jsonRule.Type {
	case "and":
		inner := []Rule{}
		for _, rule := range jsonRule.Rules {
			inner = append(inner, ParseRules(rule))
		}
		return NewAndRule(inner...)
	case "or":
		inner := []Rule{}
		for _, rule := range jsonRule.Rules {
			inner = append(inner, ParseRules(rule))
		}
		return NewOrRule(inner...)
	case "all_match":
		return &AllMatchRule{}
	case "field_match":
		return &FieldMatchRule{Name: jsonRule.Name, Value: jsonRule.Value}
	case "not_field_match":
		return &NotFieldMatchRule{Name: jsonRule.Name, Value: jsonRule.Value}
	case "field_present":
		return &FieldPresentRule{Name: jsonRule.Name}
	case "not_field_present":
		return &NotFieldPresentRule{Name: jsonRule.Name}
	}
	return nil
}

func NewAndRule(rules ...Rule) *AndRule {
	and := AndRule{rules: map[string]Rule{}}
	for _, rule := range rules {
		and.rules[rule.Type()] = rule
	}
	return &and
}

type AndRule struct {
	rules map[string]Rule
}

func (and *AndRule) Type() string {
	return "and"
}

func (and *AndRule) Match(event *core.Event) bool {
	for _, rule := range and.rules {
		if !rule.Match(event) {
			return false
		}
	}
	return true
}

func NewOrRule(rules ...Rule) *OrRule {
	or := OrRule{rules: map[string]Rule{}}
	for _, rule := range rules {
		or.rules[rule.Type()] = rule
	}
	return &or
}

type OrRule struct {
	rules map[string]Rule
}

func (or *OrRule) Type() string {
	return "or"
}

func (or *OrRule) Match(event *core.Event) bool {
	for _, rule := range or.rules {
		if rule.Match(event) {
			return true
		}
	}
	return false
}

type AllMatchRule struct{}

func (r *AllMatchRule) Type() string {
	return "all_match"
}

func (r *AllMatchRule) Match(event *core.Event) bool {
	return true
}

type FieldMatchRule struct {
	Name  string
	Value interface{}
}

func (r *FieldMatchRule) Type() string {
	return "field_match"
}

func (r *FieldMatchRule) Match(event *core.Event) bool {
	val, found := event.Data[r.Name]
	if !found {
		return false
	}
	return val == r.Value
}

type NotFieldMatchRule struct {
	Name  string
	Value interface{}
}

func (r *NotFieldMatchRule) Type() string {
	return "not_field_match"
}

func (r *NotFieldMatchRule) Match(event *core.Event) bool {
	val, found := event.Data[r.Name]
	if !found {
		return true
	}
	return val != r.Value
}

type FieldPresentRule struct {
	Name string
}

func (r *FieldPresentRule) Type() string {
	return "field_present"
}

func (r *FieldPresentRule) Match(event *core.Event) bool {
	_, found := event.Data[r.Name]
	return found
}

type NotFieldPresentRule struct {
	Name string
}

func (r *NotFieldPresentRule) Type() string {
	return "not_field_present"
}

func (r *NotFieldPresentRule) Match(event *core.Event) bool {
	_, found := event.Data[r.Name]
	return !found
}
