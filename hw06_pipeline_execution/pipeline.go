package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func runWithDoneCheck(in In, done In) Out {
	out := make(Bi)

	go func() {
		closeAndClear := func() {
			close(out)
			for range in {
				continue // drain in after cancellation
			}
		}

		for {
			select {
			case <-done:
				closeAndClear()
				return
			case v, ok := <-in:
				if !ok {
					close(out)
					return
				}

				select {
				case <-done:
					closeAndClear()
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stream := in

	for _, stage := range stages {
		stream = stage(runWithDoneCheck(stream, done))
	}

	return runWithDoneCheck(stream, done)
}
