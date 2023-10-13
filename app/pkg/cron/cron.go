package cron

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/justty/golang-units/app/pkg/utils"
)

type Job struct {
	ID                  int    `db:"id"`
	JobName             string `db:"job_name"`
	executeWithInverval *time.Duration
	executeAtTime       []time.Duration
	jobFunc             func()
	concurrency         int
	StartTime           time.Time  `db:"start_time"`
	EndTime             *time.Time `db:"end_time"`
}

func NewJob() Job {
	return Job{concurrency: 1}
}

func (j Job) Name(name string) Job {
	j.JobName = name
	return j
}

func (j Job) Every(interval time.Duration) Job {
	j.executeWithInverval = &interval
	return j
}

func (j Job) At(at time.Duration) Job {
	j.executeAtTime = append(j.executeAtTime, at)
	return j
}

func (j Job) Task(f func()) Job {
	j.jobFunc = f
	return j
}

type Cron struct {
	store *CronStore
	Jobs  []Job
}

func NewCron(store *CronStore) *Cron {
	return &Cron{store: store}
}

func (c *Cron) AddJob(job Job) {
	c.Jobs = append(c.Jobs, job)
}

func (c *Cron) Start() {
	fineshedJobs, err := c.store.SetAllFinished()
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("error finished cron: %s\n", err)
	}
	fmt.Printf("Finished %d active jobs\n", fineshedJobs)

	for {
		for _, job := range c.Jobs {
			if c.creationAllowed(job) {
				go c.runJob(job)
			}
		}
		time.Sleep(time.Minute)
	}
}

func (c *Cron) creationAllowed(job Job) bool {
	lastJob, err := c.store.getLastJob(job.JobName)
	if err != nil {
		fmt.Printf("error getting cron jobs: %s\n", err)
		return false
	}

	activeCount, err := c.store.activeJobsCount(job.JobName)
	if err != nil {
		fmt.Printf("error counting cron jobs: %s\n", err)
		return false
	}

	if len(job.executeAtTime) > 0 {
		for _, at := range job.executeAtTime {
			atTime := utils.DayToDuration(at)
			if time.Now().Hour() == atTime.Hour() && (lastJob == nil || lastJob.StartTime.Format("2006-01-02") != time.Now().Format("2006-01-02")) && activeCount < job.concurrency {
				return true
			}
		}
		return false
	}

	if activeCount >= job.concurrency {
		return false
	}

	if lastJob == nil || (time.Now().UnixNano()-lastJob.StartTime.UnixNano()) > int64(*job.executeWithInverval) {
		return true
	}

	return false
}

func (c *Cron) runJob(job Job) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("error with job %s %d\n", job.JobName, job.ID)
		}
		job.EndTime = utils.TimePointer(time.Now())
		err := c.store.update(job)
		if err != nil {
			fmt.Printf("error with update job %s %d\n", job.JobName, job.ID)
		}
	}()

	job.StartTime = time.Now()
	id, err := c.store.save(job)
	if err != nil {
		fmt.Printf("error saving job %s\n", job.JobName)
	}
	job.ID = id

	job.jobFunc()
}
