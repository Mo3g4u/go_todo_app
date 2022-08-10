package store

import (
	"context"

	"github.com/Mo3g4u/go_todo_app/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer, id entity.UserID,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
				id, title,
                status, created, modified
              FROM task
			  WHERE user_id = ?;`
	if err := db.SelectContext(ctx, &tasks, sql, id); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	sql := `INSERT INTO task
			(user_id, title, status, created, modified)
	VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql, t.UserID, t.Title, t.Status,
		r.Clocker.Now(), r.Clocker.Now(),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}
