package jobs

import (
	"github.com/mt1976/frantic-core/timing"
)

type template struct {
}

func (job template) Run() error {
	jobNotifications()
	NextRun(job)
	return nil
}

func (job template) Service() func() {
	return func() {
		job.Run()
	}
}

func (job template) Schedule() string {
	return "10 7 * * *"
}

func (job template) Name() string {
	return "Template Job"
}

func jobNotifications() {
	// Do something every day at midnight

	j := timing.Start(domain, "Send", "Service")

	j.Stop(0)
}
