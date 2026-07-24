package cron

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
	"github.com/robfig/cron/v3"
)

// CronManager manages scheduled cron jobs.
type CronManager struct {
	repo      *repository.CronRepository
	scheduler *cron.Cron
	entryMap  map[uint]cron.EntryID
}

// NewCronManager creates a new CronManager.
func NewCronManager(repo *repository.CronRepository) *CronManager {
	return &CronManager{
		repo:     repo,
		scheduler: cron.New(cron.WithSeconds()),
		entryMap: make(map[uint]cron.EntryID),
	}
}

// Create creates a new cron job.
func (m *CronManager) Create(job *model.CronJob) error {
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	job.Status = "stopped"

	if err := m.repo.Create(job); err != nil {
		return fmt.Errorf("failed to create cron job: %w", err)
	}
	return nil
}

// Update updates an existing cron job.
func (m *CronManager) Update(job *model.CronJob) error {
	job.UpdatedAt = time.Now()

	// Stop existing job if running
	if entryID, ok := m.entryMap[job.ID]; ok {
		m.scheduler.Remove(entryID)
		delete(m.entryMap, job.ID)
	}

	if err := m.repo.Update(job); err != nil {
		return fmt.Errorf("failed to update cron job: %w", err)
	}

	if job.Status == "running" {
		return m.startJob(job)
	}
	return nil
}

// Delete deletes a cron job.
func (m *CronManager) Delete(id uint) error {
	if entryID, ok := m.entryMap[id]; ok {
		m.scheduler.Remove(entryID)
		delete(m.entryMap, id)
	}

	if err := m.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete cron job: %w", err)
	}
	return nil
}

// List returns all cron jobs.
func (m *CronManager) List() ([]model.CronJob, error) {
	return m.repo.List()
}

// Start starts a cron job by adding it to the scheduler.
func (m *CronManager) Start(id uint) error {
	job, err := m.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cron job: %w", err)
	}

	job.Status = "running"
	job.UpdatedAt = time.Now()
	if err := m.repo.Update(job); err != nil {
		return fmt.Errorf("failed to update cron job status: %w", err)
	}

	return m.startJob(job)
}

// Stop stops a cron job by removing it from the scheduler.
func (m *CronManager) Stop(id uint) error {
	job, err := m.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cron job: %w", err)
	}

	if entryID, ok := m.entryMap[id]; ok {
		m.scheduler.Remove(entryID)
		delete(m.entryMap, id)
	}

	job.Status = "stopped"
	job.UpdatedAt = time.Now()
	if err := m.repo.Update(job); err != nil {
		return fmt.Errorf("failed to update cron job status: %w", err)
	}
	return nil
}

// RunNow immediately executes a cron job.
func (m *CronManager) RunNow(id uint) error {
	job, err := m.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cron job: %w", err)
	}

	now := time.Now()
	job.LastRun = &now
	m.repo.Update(job)

	go m.executeJob(job)
	return nil
}

// StartScheduler starts the cron scheduler and loads all running jobs.
func (m *CronManager) StartScheduler() error {
	jobs, err := m.repo.List()
	if err != nil {
		return fmt.Errorf("failed to list cron jobs: %w", err)
	}

	for _, job := range jobs {
		if job.Status == "running" {
			if err := m.startJob(&job); err != nil {
				continue // log error but continue loading other jobs
			}
		}
	}

	m.scheduler.Start()
	return nil
}

// StopScheduler stops the cron scheduler.
func (m *CronManager) StopScheduler() {
	m.scheduler.Stop()
}

// startJob adds a job to the scheduler.
func (m *CronManager) startJob(job *model.CronJob) error {
	entryID, err := m.scheduler.AddFunc(job.Spec, func() {
		m.executeJob(job)
	})
	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}
	m.entryMap[job.ID] = entryID
	return nil
}

// executeJob executes a cron job command.
func (m *CronManager) executeJob(job *model.CronJob) {
	now := time.Now()
	job.LastRun = &now

	switch job.Type {
	case "shell":
		cmd := exec.Command("sh", "-c", job.Command)
		if err := cmd.Run(); err != nil {
			// Log error in production
			_ = err
		}
	case "script":
		cmd := exec.Command(job.Command)
		if err := cmd.Run(); err != nil {
			_ = err
		}
	case "backup":
		cmd := exec.Command("sh", "-c", job.Command)
		if err := cmd.Run(); err != nil {
			_ = err
		}
	case "http":
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", job.Command)
		if err := cmd.Run(); err != nil {
			_ = err
		}
	}

	// Calculate next run
	entryID, ok := m.entryMap[job.ID]
	if ok {
		if entry := m.scheduler.Entry(entryID); !entry.Next.IsZero() {
			job.NextRun = &entry.Next
		}
	}

	m.repo.Update(job)
}
