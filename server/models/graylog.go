package models

import (
	"time"
)

type SystemInfo struct {
	ClusterId string `json:"cluster_id"`
	NodeId    string `json:"node_id"`
	Version   string `json:"version"`
}

type Stream struct {
	ID string `json:"id,omitempty"`
}

type IndexRangeSummary struct {
	TookMs       int       `json:"took_ms"`
	Begin        time.Time `json:"begin"`
	End          time.Time `json:"end"`
	CalculatedAt time.Time `json:"calculated_at"`
	IndexName    string    `json:"index_name"`
}

type SearchDecorationStats struct {
	AddedFields   []string `json:"added_fields"`
	RemovedFields []string `json:"removed_fields"`
	ChangedFields []string `json:"changed_fields"`
}

type DecorationStats struct {
	AddedFields   map[string]any `json:"added_fields"`
	RemovedFields map[string]any `json:"removed_fields"`
	ChangedFields map[string]any `json:"changed_fields"`
}

type ResultMessage struct {
	Index           string            `json:"index"`
	HighlightRanges map[string][]any  `json:"highlight_ranges"`
	Message         map[string]string `json:"message"`
	DecorationStats DecorationStats   `json:"decoration_stats"`
}

type Message struct {
	Query           string                `json:"query"`
	BuiltQuery      string                `json:"built_query"`
	Time            int                   `json:"time"`
	TotalResults    int                   `json:"total_results"`
	From            time.Time             `json:"from"`
	To              time.Time             `json:"to"`
	Fields          []string              `json:"fields"`
	UsedIndices     []IndexRangeSummary   `json:"used_indices"`
	DecorationStats SearchDecorationStats `json:"decoration_stats"`
	Messages        []ResultMessage       `json:"messages"`
}

// == Custom models ==

type QueryRequestInput struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Filter string `json:"filter"`
	Fields string `json:"fields"`
	Sort   string `json:"sort"`
	Range  int    `json:"range"`
}
