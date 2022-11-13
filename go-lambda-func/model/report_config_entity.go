package model

type ReportConfigDao struct {
	PK   *string `json:"PK"`
	SK   *string `json:"SK"`
	GPK1 *string `json:"GPK1"`
	GSK1 *int64  `json:"GSK1"`
	Type *string `json:"Type"`
	*ReportConfigEntity
}

type ReportConfigEntity struct {
	Columns       *[]*ReportColumn   `json:"columns"`
	LastChangedAt *int64             `json:"lastChangedAt"`
	LastSyncAt    *int64             `json:"lastSyncAt"`
	Properties    *[]*ReportProperty `json:"properties"`
	Subtype       *string            `json:"subtype"`
	Type          *string            `json:"type"`
}

type ReportColumn struct {
	Fields *[]*ReportColumnProperty `json:"fields"`
	Id     *string                  `json:"id"`
}

type ReportColumnProperty struct {
	Align        *string `json:"align"`
	DefaultValue *string `json:"defaultValue"`
	Flex         *int    `json:"flex"`
	Key          *string `json:"key"`
	Title        *string `json:"title"`
}

type ReportProperty struct {
	BoolValue   *bool    `json:"boolValue"`
	DoubleValue *float32 `json:"doubleValue"`
	IntValue    *int64   `json:"intValue"`
	Key         *string  `json:"key"`
	StringValue *string  `json:"stringValue"`
}
