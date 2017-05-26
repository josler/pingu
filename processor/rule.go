package processor

import "github.com/josler/pingu/core"

type Rule interface {
	Name() string
	Match(e *core.Event) bool
}

func NewAndRule(rules ...Rule) *AndRule {
	and := AndRule{rules: map[string]Rule{}}
	for _, rule := range rules {
		and.rules[rule.Name()] = rule
	}
	return &and
}

type AndRule struct {
	rules map[string]Rule
}

func (and *AndRule) Name() string {
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
		or.rules[rule.Name()] = rule
	}
	return &or
}

type OrRule struct {
	rules map[string]Rule
}

func (or *OrRule) Name() string {
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

func (r *AllMatchRule) Name() string {
	return "all_match"
}

func (r *AllMatchRule) Match(event *core.Event) bool {
	return true
}

type NameMatchRule struct {
	MatchName string
}

func (r *NameMatchRule) Name() string {
	return "name_match"
}

func (r *NameMatchRule) Match(event *core.Event) bool {
	return event.Name == r.MatchName
}

type FieldMatchRule struct {
	FieldName string
	Value     string
}

func (r *FieldMatchRule) Name() string {
	return "field_match"
}

func (r *FieldMatchRule) Match(event *core.Event) bool {
	val, found := event.Data[r.FieldName]
	if !found {
		return false
	}
	return val == r.Value
}

type NotFieldMatchRule struct {
	FieldName string
	Value     string
}

func (r *NotFieldMatchRule) Name() string {
	return "not_field_match"
}

func (r *NotFieldMatchRule) Match(event *core.Event) bool {
	val, found := event.Data[r.FieldName]
	if !found {
		return true
	}
	return val != r.Value
}

type FieldPresentRule struct {
	FieldName string
}

func (r *FieldPresentRule) Name() string {
	return "field_present"
}

func (r *FieldPresentRule) Match(event *core.Event) bool {
	_, found := event.Data[r.FieldName]
	return found
}

type NotFieldPresentRule struct {
	FieldName string
}

func (r *NotFieldPresentRule) Name() string {
	return "not_field_present"
}

func (r *NotFieldPresentRule) Match(event *core.Event) bool {
	_, found := event.Data[r.FieldName]
	return !found
}
