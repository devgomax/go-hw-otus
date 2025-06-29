package sqlstorage

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/devgomax/go-hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const eventsTable = "events"

// Repository модель БД типа sql.
type Repository struct {
	pool *pgxpool.Pool
}

// New конструктор БД типа sql.
func New() *Repository {
	return &Repository{}
}

// Connect открывает соединение с БД.
func (r *Repository) Connect(ctx context.Context, dsn string) error {
	if r.pool != nil {
		return errors.New("[sqlstorage::NewConnection]: can't call Connect on established connection")
	}

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::NewConnection]: can't establish connection to DB")
	}

	r.pool = pool

	return nil
}

// Close закрывает соединение с БД.
func (r *Repository) Close() {
	r.pool.Close()
}

// CreateEvent сохраняет событие в БД.
func (r *Repository) CreateEvent(ctx context.Context, event *storage.Event) error {
	m, err := storage.Serialize(event)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::CreateEvent]: can't serialize event")
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(eventsTable).
		SetMap(m)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::CreateEvent]: can't build sql query")
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::CreateEvent]: can't execute sql query")
	}

	return nil
}

// UpdateEvent обновляет событие в БД.
func (r *Repository) UpdateEvent(ctx context.Context, event *storage.Event) error {
	m, err := storage.Serialize(event)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::UpdateEvent]: can't serialize event")
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(eventsTable).
		SetMap(m).
		Where(sq.Eq{"id": event.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::UpdateEvent]: can't build sql query")
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::UpdateEvent]: can't execute sql query")
	}

	return nil
}

// DeleteEvent удаляет событие из БД.
func (r *Repository) DeleteEvent(ctx context.Context, eventID string) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(eventsTable).
		Where(sq.Eq{"id": eventID})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::DeleteEvent]: can't build sql query")
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::DeleteEvent]: can't execute sql query")
	}

	return nil
}

// ReadDailyEvents читает события за указанную дату.
func (r *Repository) ReadDailyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "title", "starts_at", "ends_at", "description", "user_id", "notify_interval").
		From(eventsTable).
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Gt{"ends_at": start},
			sq.Lt{"starts_at": end},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadDailyEvents]: can't build sql query")
	}

	var events []*storage.Event

	if err = pgxscan.Select(ctx, r.pool, &events, query, args...); err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadDailyEvents]: can't execute sql query")
	}

	return events, nil
}

// ReadWeeklyEvents читает события за неделю, начиная с указанной даты.
func (r *Repository) ReadWeeklyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(7 * 24 * time.Hour)

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "title", "starts_at", "ends_at", "description", "user_id", "notify_interval").
		From(eventsTable).
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Gt{"ends_at": start},
			sq.Lt{"starts_at": end},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadWeeklyEvents]: can't build sql query")
	}

	var events []*storage.Event

	if err = pgxscan.Select(ctx, r.pool, &events, query, args...); err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadWeeklyEvents]: can't execute sql query")
	}

	return events, nil
}

// ReadMonthlyEvents читает события за месяц, начиная с указанной даты.
func (r *Repository) ReadMonthlyEvents(ctx context.Context, userID string, date time.Time) ([]*storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month()+1, date.Day()+1, 0, 0, 0, 0, date.Location())

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "title", "starts_at", "ends_at", "description", "user_id", "notify_interval").
		From(eventsTable).
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Gt{"ends_at": start},
			sq.Lt{"starts_at": end},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadMonthlyEvents]: can't build sql query")
	}

	var events []*storage.Event

	if err = pgxscan.Select(ctx, r.pool, &events, query, args...); err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadMonthlyEvents]: can't execute sql query")
	}

	return events, nil
}

// ReadEventsToNotify читает события, у которых (starts_at - now()) <= notify_interval.
func (r *Repository) ReadEventsToNotify(ctx context.Context) ([]*storage.Event, error) {
	now := time.Now().UTC()

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "title", "starts_at", "user_id").
		From(eventsTable).
		Where(sq.And{
			sq.Eq{"processed": false},
			sq.LtOrEq{"starts_at - notify_interval": now},
			sq.GtOrEq{"ends_at": now},
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadEventsToNotify]: can't build sql query")
	}

	var events []*storage.Event

	if err = pgxscan.Select(ctx, r.pool, &events, query, args...); err != nil {
		return nil, errors.Wrap(err, "[sqlstorage::ReadEventsToNotify]: can't execute sql query")
	}

	return events, nil
}

// SetEventsProcessedStatus помечает события, как обработанные (уведомления отправлены).
func (r *Repository) SetEventsProcessedStatus(ctx context.Context, ids ...string) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(eventsTable).
		Set("processed", true).
		Where(sq.Eq{"id": ids})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::SetEventsProcessedStatus]: can't build sql query")
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "[sqlstorage::SetEventsProcessedStatus]: can't execute sql query")
	}

	return nil
}
