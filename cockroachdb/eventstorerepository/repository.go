package eventstorerepository

import (
		"context"
		"database/sql"
		"fmt"

		"eventstore"
)

type repository struct {
		db *sql.DB
}

// New creates a new event store repository
func New(db *sql.DB) (eventstore.Repository, error) {
		return &repository{
			db: db,
		}, nil
}

// CreateEvent creates a new event to the event store
func (repo *repository) CreateEvent(ctx context.Context, event *eventstore.Event) error {
		// insert example row in events table
		var err error
		var sql string
		if event.EventData != "" {
			sql = `INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, stream) VALUES ($1, $2, $3, $4, $5, $6)`
			_, err = repo.db.ExecContext(ctx, sql, event.ID, event.EventType, event.AggregateID, event.AggregateType, event.EventData, event.Stream)
		}
		if err != nil {
			return fmt.Errorf("error creating event: %w", err)
		}
		return nil
}

// GetEvents gets all events for the given aggregate and event
func (repo repository) GetEvents(ctx context.Context, filter *eventstore.GetEventsRequest) ([]*eventstore.Event, error) {
		var rows *sql.Rows
		var err error
		var sql string
		if filter.EventId != "" && filter.AggregateId == "" {
			sql = `SELECT id, eventtype, aggregateid, aggregatetype, eventdata, stream FROM events`
			rows, err = repo.db.QueryContext(ctx, sql)
		}

		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var events []*eventstore.Event
		// Loop through rows, using Scan to assign column data to struct fields
		for rows.Next() {
			var e eventstore.Event
			err := rows.Scan(&e.ID, &e.EventType, &e.AggregateID, &e.AggregateType, &e.EventData, &e.Stream)
			if err != nil {
				return nil, err
			}
			events = append(events, &e)
		}
		if err = rows.Err(); err != nil {
				return events, err
		}
		return events, nil
}

