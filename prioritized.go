package prioritized

type PrioritizedWorkers struct {
	highPriorityJobQueue chan *Job
	midPriorityJobQueue  chan *Job
	lowPriorityJobQueue  chan *Job
	jobTypeRegistry      map[interface{}]func(jobData interface{})
	CloseSignal          chan struct{}
}

type Job struct {
	JobType interface{}
	JobData interface{}
}

func Init(jobQueueCapacity int, jobWorkerConcurrency int, closeWorkerSignal chan struct{}) *PrioritizedWorkers {
	s := &PrioritizedWorkers{
		highPriorityJobQueue: make(chan *Job, jobQueueCapacity),
		midPriorityJobQueue:  make(chan *Job, jobQueueCapacity),
		lowPriorityJobQueue:  make(chan *Job, jobQueueCapacity),
		CloseSignal:          closeWorkerSignal,
	}
	if closeWorkerSignal == nil {
		s.CloseSignal = make(chan struct{})
	}
	for w := 0; w < jobWorkerConcurrency; w++ {
		go s.Worker(closeWorkerSignal)
	}
	return s
}

func (s *PrioritizedWorkers) Worker(closeSignal chan struct{}) {
	for {
		select {
		case <-closeSignal:
			return
		case j := <-s.highPriorityJobQueue:
			s.HandleJob(j)
		default:
			{
				select {
				case <-closeSignal:
					return
				case j := <-s.midPriorityJobQueue:
					s.HandleJob(j)
				default:
					{
						select {
						case <-closeSignal:
							return
						case j := <-s.lowPriorityJobQueue:
							s.HandleJob(j)
						default:
							{
								select {
								case <-closeSignal:
									return
								case j := <-s.highPriorityJobQueue:
									s.HandleJob(j)
								case j := <-s.midPriorityJobQueue:
									select {
									case <-closeSignal:
										return
									case j := <-s.highPriorityJobQueue:
										s.HandleJob(j)
									default:
									}
									s.HandleJob(j)
								case j := <-s.lowPriorityJobQueue:
									select {
									case <-closeSignal:
										return
									case j := <-s.highPriorityJobQueue:
										s.HandleJob(j)
									case j := <-s.midPriorityJobQueue:
										s.HandleJob(j)
									default:
									}
									s.HandleJob(j)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (s *PrioritizedWorkers) DefineJobType(typeSymbol interface{}, jobHandler func(jobData interface{})) {
	s.jobTypeRegistry[typeSymbol] = jobHandler
}

func (s *PrioritizedWorkers) HandleJob(j *Job) {
	s.jobTypeRegistry[j.JobType](j.JobData)
}

func (s *PrioritizedWorkers) PendHighPriorityJob(typeSymbol interface{}, jobData interface{}) {
	s.highPriorityJobQueue <- &Job{
		JobType: typeSymbol,
		JobData: jobData,
	}
}

func (s *PrioritizedWorkers) PendMidPriorityJob(typeSymbol interface{}, jobData interface{}) {
	s.midPriorityJobQueue <- &Job{
		JobType: typeSymbol,
		JobData: jobData,
	}
}

func (s *PrioritizedWorkers) PendLowPriorityJob(typeSymbol interface{}, jobData interface{}) {
	s.lowPriorityJobQueue <- &Job{
		JobType: typeSymbol,
		JobData: jobData,
	}
}
