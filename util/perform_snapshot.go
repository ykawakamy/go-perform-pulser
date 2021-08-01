package util

import "fmt"

type PerformSnapShot struct {
	perform      uint
	no_update    uint
	skip_error   uint
	rewind_error uint
	latency      uint64

	actual_max uint32
	actual_min uint32
}

func CreatePerformSummary() PerformSnapShot {
	return PerformSnapShot{
		actual_max: const_init_max,
		actual_min: const_init_min,
	}
}

func CreatePerformSnapShot(pc *PerformCounter) PerformSnapShot {
	return PerformSnapShot{
		perform:      pc.perform,
		no_update:    pc.no_update,
		skip_error:   pc.skip_error,
		rewind_error: pc.rewind_error,
		latency:      pc.latency,

		actual_max: pc.actual_max,
		actual_min: pc.actual_min,
	}
}

func (s *PerformSnapShot) Print(elapsedNanoSec int64) {
	fmt.Printf("range: %d - %d elapsed: %d ns. process: %d op. %f ns/op. ( %f op/s ) latency: %f ms/op | errors - total: %d ( noupdate: %d rewind: %d skip: %d )\n",
		s.actual_min,
		s.actual_max,
		elapsedNanoSec,
		s.perform,
		float64(elapsedNanoSec)/float64(s.perform),
		float64(s.perform)*1e9/float64(elapsedNanoSec),
		float64(s.latency)/float64(s.perform),
		//--
		s.no_update+s.rewind_error+s.skip_error,
		s.no_update,
		s.rewind_error,
		s.skip_error,
	)
}

func (s *PerformSnapShot) add(o *PerformSnapShot) {
	s.perform += o.perform
	s.no_update += o.no_update
	s.skip_error += o.skip_error
	s.rewind_error += o.rewind_error
	s.latency += o.latency

	s.actual_max = max(s.actual_max, o.actual_max)
	s.actual_min = min(s.actual_min, o.actual_min)
}
