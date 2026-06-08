package api

import (
	"net/http"

	"scheduler/pkg/db"
)

const TasksLimit = 50

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(TasksLimit, search)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, TaskResp{
		Tasks: tasks,
	})
}
