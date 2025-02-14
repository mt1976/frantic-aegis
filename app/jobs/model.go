package jobs

type Job interface {
	Run() error
	Service() func()
	Schedule() string
	Name() string
}
