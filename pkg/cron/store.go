package cron

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type CronStore struct {
	db *sqlx.DB
}

func NewCronStore(db *sqlx.DB) CronStore {
	return CronStore{
		db: db,
	}
}

func (c *CronStore) update(job Job) error {
	sqlStatement := `
	UPDATE cron_logs 
	SET 
		job_name=:job_name, 
		start_time=:start_time,
		end_time=:end_time 
	WHERE id=:id
	`

	if _, err := c.db.NamedExec(sqlStatement, job); err != nil {
		return err
	}

	return nil
}

func (c *CronStore) save(job Job) (int, error) {
	sqlStatement := `
	INSERT INTO cron_logs 
	(
		job_name, 
		start_time
	) 
	VALUES 
	(
		$1, 
		$2
	)
	RETURNING id
	`

	row := c.db.QueryRow(sqlStatement, job.JobName, job.StartTime)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CronStore) getLastJob(jobName string) (*Job, error) {
	var job Job
	err := c.db.Get(&job, "SELECT job_name, max(start_time) start_time FROM cron_logs WHERE job_name=$1 GROUP BY job_name", jobName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &job, nil
}

func (c *CronStore) activeJobsCount(jobName string) (int, error) {
	var count int
	err := c.db.Get(&count, "SELECT count(*) active_count FROM cron_logs WHERE job_name=$1 and end_time is null GROUP BY job_name", jobName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func (c CronStore) SetAllFinished() (int, error) {
	result, err := c.db.Exec(`
	UPDATE cron_logs 
	SET end_time=now()
	WHERE end_time is null
	`)

	if err != nil {
		return 0, err
	} else {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		return int(rowsAffected), nil
	}
}
