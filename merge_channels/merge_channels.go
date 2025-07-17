package merge_channels

import (
	"slices"
	"sync"
)

func MergeChannels(channels ...chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	for _, currentChan := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range currentChan {
				out <- i
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func MergeChannelPattern() []int {
	ch1, ch2, ch3 := make(chan int), make(chan int), make(chan int)
	go func() {
		for _, v := range []int{1, 2, 3} {
			ch1 <- v
		}
		close(ch1)
	}()
	go func() {
		for _, v := range []int{4, 5, 6} {
			ch2 <- v
		}
		close(ch2)
	}()
	go func() {
		for _, v := range []int{7, 8, 9} {
			ch3 <- v
		}
		close(ch3)
	}()
	var res []int
	for mergedChannelValues := range MergeChannels(ch1, ch2, ch3) {
		res = append(res, mergedChannelValues)
	}
	slices.Sort(res)
	return res
}
