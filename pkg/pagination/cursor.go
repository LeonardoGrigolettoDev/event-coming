package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

// Cursor represents a pagination cursor
type Cursor struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"ca"`
}

// PageInfo holds pagination metadata
type PageInfo struct {
	HasNextPage     bool    `json:"has_next_page"`
	HasPreviousPage bool    `json:"has_previous_page"`
	StartCursor     *string `json:"start_cursor,omitempty"`
	EndCursor       *string `json:"end_cursor,omitempty"`
	TotalCount      *int64  `json:"total_count,omitempty"`
}

// CursorRequest holds pagination request parameters
type CursorRequest struct {
	First  *int    `json:"first,omitempty" form:"first"`   // Forward pagination: number of items
	After  *string `json:"after,omitempty" form:"after"`   // Forward pagination: cursor
	Last   *int    `json:"last,omitempty" form:"last"`     // Backward pagination: number of items
	Before *string `json:"before,omitempty" form:"before"` // Backward pagination: cursor
}

// Validate validates the cursor request
func (r *CursorRequest) Validate() error {
	// Can't use both forward and backward pagination
	if (r.First != nil || r.After != nil) && (r.Last != nil || r.Before != nil) {
		return errors.New("cannot use forward and backward pagination together")
	}

	// Validate limits
	if r.First != nil && (*r.First < 1 || *r.First > 100) {
		return errors.New("first must be between 1 and 100")
	}
	if r.Last != nil && (*r.Last < 1 || *r.Last > 100) {
		return errors.New("last must be between 1 and 100")
	}

	return nil
}

// GetLimit returns the requested limit with a default
func (r *CursorRequest) GetLimit() int {
	if r.First != nil {
		return *r.First
	}
	if r.Last != nil {
		return *r.Last
	}
	return 20 // default
}

// IsForward returns true if this is forward pagination
func (r *CursorRequest) IsForward() bool {
	return r.Last == nil && r.Before == nil
}

// EncodeCursor encodes a cursor to a base64 string
func EncodeCursor(id string, createdAt time.Time) string {
	cursor := Cursor{
		ID:        id,
		CreatedAt: createdAt,
	}
	data, _ := json.Marshal(cursor)
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeCursor decodes a base64 cursor string
func DecodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return nil, nil
	}

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("invalid cursor format")
	}

	var cursor Cursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, errors.New("invalid cursor data")
	}

	return &cursor, nil
}

// NewPageInfo creates a PageInfo from a list of items
func NewPageInfo[T any](
	items []T,
	getCursor func(item T) (string, time.Time),
	hasMore bool,
	totalCount *int64,
	isForward bool,
) PageInfo {
	info := PageInfo{
		TotalCount: totalCount,
	}

	if len(items) == 0 {
		return info
	}

	if isForward {
		info.HasNextPage = hasMore
		info.HasPreviousPage = false // Would need to check database
	} else {
		info.HasPreviousPage = hasMore
		info.HasNextPage = false // Would need to check database
	}

	// Set cursors
	firstID, firstTime := getCursor(items[0])
	lastID, lastTime := getCursor(items[len(items)-1])

	startCursor := EncodeCursor(firstID, firstTime)
	endCursor := EncodeCursor(lastID, lastTime)

	info.StartCursor = &startCursor
	info.EndCursor = &endCursor

	return info
}

// Connection represents a paginated connection response
type Connection[T any] struct {
	Edges    []Edge[T] `json:"edges"`
	PageInfo PageInfo  `json:"page_info"`
}

// Edge represents an edge in the connection
type Edge[T any] struct {
	Node   T      `json:"node"`
	Cursor string `json:"cursor"`
}

// NewConnection creates a connection from items
func NewConnection[T any](
	items []T,
	getCursor func(item T) (string, time.Time),
	hasMore bool,
	totalCount *int64,
	isForward bool,
) Connection[T] {
	edges := make([]Edge[T], len(items))
	for i, item := range items {
		id, createdAt := getCursor(item)
		edges[i] = Edge[T]{
			Node:   item,
			Cursor: EncodeCursor(id, createdAt),
		}
	}

	return Connection[T]{
		Edges:    edges,
		PageInfo: NewPageInfo(items, getCursor, hasMore, totalCount, isForward),
	}
}
