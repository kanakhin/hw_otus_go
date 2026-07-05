package hw06pipelineexecution

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestAllStageStop(t *testing.T) {
	wg := sync.WaitGroup{}
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		wg.Wait()

		require.Len(t, result, 0)
	})
}

func newStageGen(wg *sync.WaitGroup) func(func(interface{}) interface{}) Stage {
	return func(f func(interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			if wg != nil {
				wg.Add(1)
			}
			go func() {
				if wg != nil {
					defer wg.Done()
				}
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}
}

func defaultStages(g func(func(interface{}) interface{}) Stage) []Stage {
	return []Stage{
		g(func(v interface{}) interface{} { return v }),
		g(func(v interface{}) interface{} { return v.(int) * 2 }),
		g(func(v interface{}) interface{} { return v.(int) + 100 }),
		g(func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}
}

func collectStrings(out Out) []string {
	result := make([]string, 0, 10)
	for s := range out {
		result = append(result, s.(string))
	}
	return result
}

func feedInts(in Bi, data []int) {
	go func() {
		for _, v := range data {
			in <- v
		}
		close(in)
	}()
}

func TestPipelineEmptyInput(t *testing.T) {
	g := newStageGen(nil)
	stages := defaultStages(g)

	in := make(Bi)
	close(in)

	start := time.Now()
	result := collectStrings(ExecutePipeline(in, nil, stages...))
	elapsed := time.Since(start)

	require.Empty(t, result)
	require.Less(t, elapsed, sleepPerStage+fault)
}

func TestPipelineNoStages(t *testing.T) {
	in := make(Bi)
	feedInts(in, []int{42})

	result := make([]interface{}, 0, 1)
	for v := range ExecutePipeline(in, nil) {
		result = append(result, v)
	}

	require.Equal(t, []interface{}{42}, result)
}

func TestPipelineSingleElement(t *testing.T) {
	g := newStageGen(nil)
	stages := defaultStages(g)

	in := make(Bi)
	feedInts(in, []int{1})

	start := time.Now()
	result := collectStrings(ExecutePipeline(in, nil, stages...))
	elapsed := time.Since(start)

	require.Equal(t, []string{"102"}, result)
	require.Less(t,
		int64(elapsed),
		int64(sleepPerStage)*int64(len(stages))+int64(fault))
}
