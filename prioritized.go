package prioritized

//WorkerGroup ...
type WorkerGroup struct {
	highPriorityJobQueue chan *job
	lowPriorityJobQueue  chan *job
	jobTypeRegistry      map[interface{}]func(jobData interface{})
	CloseSignal          chan struct{}
}

type job struct {
	JobType interface{}
	JobData interface{}
}

//Init ...
func Init(jobQueueCapacity int, workerConcurrency int, closeWorkerSignal chan struct{}) *WorkerGroup {
	s := &WorkerGroup{
		highPriorityJobQueue: make(chan *job, jobQueueCapacity),
		lowPriorityJobQueue:  make(chan *job, jobQueueCapacity),
		jobTypeRegistry: make(map[interface{}]func(jobData interface{})),
		CloseSignal:          closeWorkerSignal,
	}
	if closeWorkerSignal == nil {
		s.CloseSignal = make(chan struct{})
	}
	for w := 0; w < workerConcurrency; w++ {
		go s.Worker(closeWorkerSignal)
	}
	return s
}

//Worker ...
func (s *WorkerGroup) Worker(closeSignal chan struct{}) {
	for {
		select {
		case <-closeSignal:
			return
		case j := <-s.highPriorityJobQueue:
			s.handleJob(j)
		default:
			{
				select {
				case <-closeSignal:
					return
				case j := <-s.highPriorityJobQueue:
					s.handleJob(j)
				case j := <-s.lowPriorityJobQueue:
					s.handleJob(j)
				}
			}
		}
	}
}

//DefineJobType ...
func (s *WorkerGroup) DefineJobType(typeSymbol interface{}, jobHandler func(jobData interface{})) {
	s.jobTypeRegistry[typeSymbol] = jobHandler
}


//handleJob ...
func (s *WorkerGroup) handleJob(j *job) {
	s.jobTypeRegistry[j.JobType](j.JobData)
}


//PendHighPriorityJob ...
func (s *WorkerGroup) PendHighPriorityJob(typeSymbol interface{}, jobData interface{}) {
	s.highPriorityJobQueue <- &job{
		JobType: typeSymbol,
		JobData: jobData,
	}
}

//PendLowPriorityJob ...
func (s *WorkerGroup) PendLowPriorityJob(typeSymbol interface{}, jobData interface{}) {
	s.lowPriorityJobQueue <- &job{
		JobType: typeSymbol,
		JobData: jobData,
	}
}

//PendAutoFallback ...
func (s *WorkerGroup) PendAutoFallback(typeSymbol interface{}, jobData interface{}) {
	select {
		case s.highPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}: {}
		case s.lowPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}:{}
	}
}

//TryPendHighPriorityJob ...
func (s *WorkerGroup) TryPendHighPriorityJob(typeSymbol interface{}, jobData interface{}) bool {
	select {
		case s.highPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}: {
			return true
		}
		default: {
			return false
		}
	}
}

//TryPendLowPriorityJob ...
func (s *WorkerGroup) TryPendLowPriorityJob(typeSymbol interface{}, jobData interface{}) bool {
	select {
		case s.lowPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}: {
			return true
		}
		default: {
			return false
		}
	}
}

//TryPendAutoFallback ...
func (s *WorkerGroup) TryPendAutoFallback(typeSymbol interface{}, jobData interface{}) bool {
	select {
		case s.highPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}: {
			return true
		}
		case s.lowPriorityJobQueue <- &job{
			JobType: typeSymbol,
			JobData: jobData,
		}: {
			return true
		}
		default: {
			return false
		}
	}
}