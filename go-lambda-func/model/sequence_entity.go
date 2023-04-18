package model

type SequenceEntityDao struct {
	PK              *string                     `json:"PK"`
	SK              *string                     `json:"SK"`
	GPK1            *string                     `json:"GPK1"`
	GSK1            *int64                      `json:"GSK1"`
	Type            *string                     `json:"Type"`
	SequenceEntitys *map[string]*SequenceEntity `json:"sequences"`
}

type SequenceEntity struct {
	CreateAt         *int64  `json:"createAt"`
	LastChangedAt    *int64  `json:"lastChangedAt"`
	LastSeqCreatedAt *int64  `json:"lastSeqCreatedAt"`
	LastSyncAt       *int64  `json:"lastSyncAt"`
	Name             *string `json:"name"`
	NextSeq          *int    `json:"nextSeq"`
	Pattern          *string `json:"pattern"`
}
