package model

type AddressEntity struct {
	Address1    *string `json:"address1"`
	Address2    *string `json:"address2"`
	City        *string `json:"city"`
	Country     *string `json:"country"`
	CountryCode *string `json:"countryCode"`
	State       *string `json:"state"`
	StateCode   *string `json:"stateCode"`
	Zipcode     *string `json:"zipcode"`
}
