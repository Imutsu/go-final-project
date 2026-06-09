package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DateFormat       = "20060102"
	SearchDateFormat = "02.01.2006"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
		`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {

	var rows *sql.Rows
	var err error

	if search == "" {

		rows, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date
			LIMIT ?
		`, limit)

	} else if date, errDate := time.Parse(SearchDateFormat, search); errDate == nil {

		rows, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date
			LIMIT ?
		`, date.Format(DateFormat), limit)

	} else {

		search = "%" + search + "%"

		rows, err = DB.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ?
			OR comment LIKE ?
			ORDER BY date
			LIMIT ?
		`, search, search, limit)

	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task

	err := DB.QueryRow(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(`
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`
		DELETE FROM scheduler
		WHERE id = ?
	`, id)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	_, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	return err
}
