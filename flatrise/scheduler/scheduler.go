package scheduler

type Scheduler interface {
	Run()

	SearchAll() error
	Search(engine string) error
}
