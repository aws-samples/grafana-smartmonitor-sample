package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Job struct {
	ID      uuid.UUID
	Summary string
}

var (
	scheduler   gocron.Scheduler
	jobsMap     = make(map[string]Job)
	jobsMapLock sync.RWMutex
)

func StartScheduler() error {

	var err error
	scheduler, err = gocron.NewScheduler()

	if err != nil {
		return err
	}

	go func() {
		logger.Println("scheduler start")
		scheduler.Start()
		select {
		case <-time.After(time.Minute):
		}

	}()

	time.Sleep(2 * time.Second)
	LoadEnableJob()
	return nil
}

func LoadEnableJob() {
	logger.Println("LoadEnableJob starting......")
	jobs, err := GetEnableMonitorJobs()
	if err != nil {
		panic(err)
	}
	for _, job := range jobs {
		projectID := job.Project
		cronExpr := job.Cron
		_, err = AddJob(projectID, cronExpr)
		if err != nil {
			// Handle error
			logger.Printf("Error adding job for project %s: %v", projectID, err)
			continue
		}

		logger.Printf("Added job for project %s with cron expression %s", projectID, cronExpr)
	}
}

func AddJob(projectID, cron string) (string, error) {
	// Check if a job already exists for the projectID
	if ExistsJob(projectID) {
		// Remove the existing job
		logger.Printf("Project: %s exits , remove first ", projectID)
		err := RemoveJobByProject(projectID)
		if err != nil {
			return "", err
		}
	}

	define, task, err := CreateJob(projectID, cron)
	if err != nil {
		return "", err
	}

	j, err := scheduler.NewJob(*define, task)
	if err != nil {
		return "", err
	}

	jobID, err := uuid.Parse(j.ID().String())
	if err != nil {
		return "", err
	}

	logger.Printf("jobID: %s\n", jobID)

	jobsMapLock.Lock()
	defer jobsMapLock.Unlock()
	jobsMap[projectID] = Job{ID: jobID, Summary: ""}

	return jobID.String(), nil
}

func GetJobByProject(projectID string) (*Job, bool) {
	jobsMapLock.RLock()
	defer jobsMapLock.RUnlock()
	job, ok := jobsMap[projectID]
	return &job, ok
}

func ExistsJob(projectID string) bool {
	jobsMapLock.RLock()
	defer jobsMapLock.RUnlock()
	_, ok := jobsMap[projectID]
	return ok
}

func RemoveJobByProject(projectID string) error {
	jobsMapLock.Lock()
	defer jobsMapLock.Unlock()
	job, ok := jobsMap[projectID]
	if !ok {
		return fmt.Errorf("no job found for project %s", projectID)
	}

	delete(jobsMap, projectID)
	logger.Printf("RemoveJob: %s", projectID)
	return scheduler.RemoveJob(job.ID)
}

func CreateJob(projectID, cron string) (*gocron.JobDefinition, gocron.Task, error) {
	// create new job
	task := gocron.NewTask(
		func() {

			logger.Printf("Running job for project %s", projectID)

			//get metrics items
			items, err := GetMetricsByProject(projectID)
			if err != nil {
				panic(err)
			}

			if len(items) == 0 {
				return
			}

			logger.Printf("Project: %s load metrics : %d", projectID, len(items))

			ids := make([]int64, 0, len(items))
			for _, item := range items {
				ids = append(ids, item.ID)
			}

			newMonitorItems, err := BatchRunItems(ids)
			if err != nil {
				logger.Println(err)
				return
			}

			logger.Printf("Project: %s check metrics completed ", projectID)

			output, err := RunProjectSummary(newMonitorItems)
			if err != nil {
				panic(err)
			}
			logger.Printf("Project: %s run summary completed ", projectID)
			job, _ := GetJobByProject(projectID)
			if job != nil {
				jobsMap[projectID] = Job{
					ID:      job.ID,
					Summary: output,
				}

			}

		},
	)

	// 创建一个新的作业
	job := gocron.CronJob(
		cron,
		true,
	)

	// job := gocron.DurationJob(
	// 	time.Second * 60,
	// )

	// 返回 gocron.JobDefine 和 gocron.Task
	return &job, task, nil
}
