package model

import (
	"reflect"
	"sort"
	"strings"
)

type TaxGroupDao struct {
	PK        *string                     `json:"PK"`
	SK        *string                     `json:"SK"`
	GPK1      *string                     `json:"GPK1"`
	GSK1      *int64                      `json:"GSK1"`
	Type      *string                     `json:"Type"`
	TaxGroups *map[string]*TaxGroupEntity `json:"taxGroups"`
}

// PK - STORE#<store_id>
// SK - TAX_GROUP#<tax_group_id>

type TaxGroupEntity struct {
	Description   *string    `json:"description"`
	GroupId       *string    `json:"groupId"`
	Name          *string    `json:"name"`
	TaxRules      []*TaxRule `json:"taxRules"`
	LastUpdatedAt *int64     `json:"lastChangedAt"`
	LastUpdatedBy *string    `json:"lastUpdatedBy"`
}

type TaxRule struct {
	Amount                  *float64 `json:"amount"`
	AuthorityId             *string  `json:"authorityId"`
	AuthorityName           *string  `json:"authorityName"`
	AuthorityType           *string  `json:"authorityType"`
	EffectiveDateTimeStamp  *float64 `json:"effectiveDateTimeStamp"`
	ExpirationDateTimeStamp *float64 `json:"expirationDateTimeStamp"`
	GroupId                 *string  `json:"groupId"`
	LocationId              *string  `json:"locationId"`
	MaximumTaxableAmount    *float64 `json:"maximumTaxableAmount"`
	MinimumTaxableAmount    *float64 `json:"minimumTaxableAmount"`
	Percent                 *float64 `json:"percent"`
	RuleId                  *string  `json:"ruleId"`
	RuleName                *string  `json:"ruleName"`
	RuleSequence            *int     `json:"ruleSequence"`
	LastUpdatedAt           *int64   `json:"lastChangedAt"`
	LastUpdatedBy           *string  `json:"lastUpdatedBy"`
}

type TaxRuleSorter []*TaxRule

func (a TaxRuleSorter) Len() int {
	return len(a)
}

func (a TaxRuleSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TaxRuleSorter) Less(i, j int) bool {
	return strings.Compare(*a[i].RuleId, *a[j].RuleId) < 0
}

type TaxGroupEntitySorter []*TaxGroupEntity

func (a TaxGroupEntitySorter) Len() int {
	return len(a)
}

func (a TaxGroupEntitySorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TaxGroupEntitySorter) Less(i, j int) bool {
	return strings.Compare(*a[i].GroupId, *a[j].GroupId) < 0
}

func (t *TaxGroupEntity) Merge(o TaxGroupEntity) TaxGroupEntity {
	if reflect.DeepEqual(t, &o) {
		return o
	}

	server := t.TaxRules
	client := o.TaxRules
	sort.Sort(TaxRuleSorter(server))
	sort.Sort(TaxRuleSorter(client))

	// Merge the tax rules
	var merged []*TaxRule
	var i, j int
	for i < len(server) && j < len(client) {
		if strings.Compare(*server[i].RuleId, *client[j].RuleId) < 0 {
			merged = append(merged, server[i])
			i++
		} else if strings.Compare(*server[i].RuleId, *client[j].RuleId) > 0 {
			merged = append(merged, client[j])
			j++
		} else {
			mergedRule := server[i].Merge(*client[j])
			merged = append(merged, &mergedRule)
			i++
			j++
		}
	}

	for i < len(server) {
		merged = append(merged, server[i])
		i++
	}

	for j < len(client) {
		merged = append(merged, client[j])
		j++
	}
	o.TaxRules = merged
	return o
}

func (t *TaxRule) Merge(o TaxRule) TaxRule {
	if reflect.DeepEqual(t, &o) {
		return o
	}
	if t.LastUpdatedAt != nil && *t.LastUpdatedAt > *o.LastUpdatedAt {
		return *t
	}
	return o
}
