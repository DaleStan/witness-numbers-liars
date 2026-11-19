package milrabliars

import (
	"encoding/binary"
	"log"
	"os"
)

type Top10ChangesAt struct {
	N       int
	Changes []Top10Movement
}

type LiarsResult struct {
	N     int
	Prime bool
	Liars []int
}

// ComputeUntil will use the given number of cores to get the running totals
// computed up to finishAt, reporting any changes to the top 10 as they come in
// to the given channel. Uses cores-1 cores for workers and 1 core to aggregate
// the data. If cores is less than 2, it will use 1 worker and 1 aggregator.
//
// The top10Changes channel is closed when this completes.
func ComputeUntil(
	runningTotals *RunningTotalLiars,
	finishAt int,
	cores int,
	top10Changes chan Top10ChangesAt,
) {
	if cores < 2 {
		cores = 2
	}

	if finishAt%2 == 0 {
		finishAt--
	}

	maxOutOfOrder := cores * 2
	outOfOrderLiars := make(map[int]LiarsResult, maxOutOfOrder)

	jobsChannel := make(chan int, maxOutOfOrder)
	liarsChannel := make(chan LiarsResult, maxOutOfOrder)
	liarsLog, err := os.OpenFile("../out/liars.dat", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Start the workers
	for i := 0; i < cores-1; i++ {
		go func() {
			for n := range jobsChannel {
				prime, liars := CalculateIsPrimeAndLiars(n)
				liarsChannel <- LiarsResult{n, prime, liars}
			}
		}()
	}

	jobsQueuedUntil := runningTotals.HighestNumberChecked
	maxQueuedJobs := maxOutOfOrder

	// Push the first set of jobs
	for i := 0; i < maxQueuedJobs && jobsQueuedUntil+2 <= finishAt; i++ {
		jobsChannel <- jobsQueuedUntil + 2
		jobsQueuedUntil += 2
	}

	for runningTotals.HighestNumberChecked < finishAt {
		// Get the next result
		result := <-liarsChannel

		// If the result is out of order, store it for later,
		// otherwise update the running totals
		if result.N != runningTotals.HighestNumberChecked+2 {
			outOfOrderLiars[result.N] = result
		} else {
			HandleLiars(liarsLog, result, runningTotals, top10Changes)
			// We may already have the next number in the out of order map
			for oooResult, ok := outOfOrderLiars[result.N+2]; ok; oooResult, ok = outOfOrderLiars[oooResult.N+2] {
				HandleLiars(liarsLog, oooResult, runningTotals, top10Changes)
				delete(outOfOrderLiars, oooResult.N)
			}
		}

		// Queue as many jobs as possible
		queuedJobs := int((jobsQueuedUntil - runningTotals.HighestNumberChecked) / 2)
		for i := 0; i < maxQueuedJobs-queuedJobs && jobsQueuedUntil+2 <= finishAt; i++ {
			jobsChannel <- jobsQueuedUntil + 2
			jobsQueuedUntil += 2
		}
	}

	// Close the jobs channel to stop the workers
	close(jobsChannel)

	// Close the changes channel so the caller knows we're done
	close(top10Changes)
}

func HandleLiars(liarsLog *os.File, result LiarsResult, runningTotals *RunningTotalLiars, top10Changes chan Top10ChangesAt) {
	if len(result.Liars) > 0 {
		binary.Write(liarsLog, binary.LittleEndian, int32(result.N))
		binary.Write(liarsLog, binary.LittleEndian, int32(len(result.Liars) >> 1))
		previousLiar := 0
		for _, liar := range result.Liars {
			if liar * 2 < result.N {
				difference := liar - previousLiar
				previousLiar = liar
				if difference >= 255 {
					binary.Write(liarsLog, binary.LittleEndian, uint8(255))
					binary.Write(liarsLog, binary.LittleEndian, uint32(difference))
				} else {
					binary.Write(liarsLog, binary.LittleEndian, uint8(difference))
				}
			}
		}
	}

	changes := runningTotals.AddNextNumber(result.N, result.Liars)
	if len(changes) > 0 {
		top10Changes <- Top10ChangesAt{
			N:       result.N,
			Changes: changes,
		}
	}
}