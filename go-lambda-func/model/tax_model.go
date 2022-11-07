package model

type TaxGroupDao struct {
	PK *string `json:"PK"`
	SK *string `json:"SK"`
	*TaxGroup
}

// PK - STORE#<store_id>
// SK - TAX_GROUP#<tax_group_id>

type TaxGroup struct {
	StoreId     *string    `json:"storeId"`
	Description *string    `json:"description"`
	GroupId     *string    `json:"groupId"`
	Name        *string    `json:"name"`
	TaxRules    []*TaxRule `json:"taxRules"`
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
}
