package hw06pipelineexecution

// Channel type aliases.
type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

// Stage function that starts inner goroutine to do some work with data and returns channel.
type Stage func(in In) (out Out)

// ExecutePipeline executes a pipeline built from the provided stages.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	if in == nil { // устраняет deadlock в кейсе, когда вместо входного канала передан nil
		rChan := make(Bi)
		defer close(rChan)
		out = rChan
	}

	for _, stage := range stages {
		out = stage(cancellableChan(out, done))
	}

	return out
}

func cancellableChan(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				<-in
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}
