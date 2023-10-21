package cron

import (
	"log"
	"time"
)

type Job struct {
	jobName             string
	executeWithInverval time.Duration
	jobFunc             func() error
}

type Cron struct {
	Jobs []Job
}

func NewJob(jobName string, executeWithInverval time.Duration, jobFunc func() error) Job {
	return Job{
		jobName:             jobName,
		executeWithInverval: executeWithInverval,
		jobFunc:             jobFunc,
	}
}

func NewCron() *Cron {
	return &Cron{}
}

func (c *Cron) AddJob(job Job) {
	c.Jobs = append(c.Jobs, job)
}

func (c *Cron) Start() {
	for _, job := range c.Jobs {
		go c.runJob(job)
	}
}

func (*Cron) runJob(job Job) {
	log.Printf("Job '%s' started\n", job.jobName)
	if err := job.jobFunc(); err != nil {
		log.Printf("error with job '%s': %s\n", job.jobName, err.Error())
	}
	log.Printf("Job '%s' finished\n", job.jobName)
	ticker := time.NewTicker(job.executeWithInverval)

	for {
		select {
		case <-ticker.C:
			log.Printf("Job '%s' started\n", job.jobName)
			if err := job.jobFunc(); err != nil {
				log.Printf("error with job '%s': %s\n", job.jobName, err.Error())
			}
			log.Printf("Job '%s' finished\n", job.jobName)
		}
	}
}
