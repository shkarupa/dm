package main

func MergeSort(items int, compare func(i, j int) int, indices chan int) {
	go MergeSortRec(compare, 0, items - 1, indices)
}

func MergeSortRec(compare func(i, j int) int, low, high int, indices chan int) {
	merge := func(left, right, out chan int) {
		i, open1 := <-left
		j, open2 := <-right
		for open1 && open2 {
			if compare(i, j) <= 0 {
				out <- i
				i, open1 = <-left
			} else {
				out <- j
				j, open2 = <-right
			}
		}
		for open1 {
			out <- i
			i, open1 = <-left
		}
		for open2 {
			out <- j
			j, open2 = <-right
		}
		close(out)
	}
	if low < high {
		med := (low + high) / 2
		left, right := make(chan int), make(chan int)
		go MergeSortRec(compare, low, med, left)
		go MergeSortRec(compare, med + 1, high, right)
		go merge(left, right, indices)
	} else {
		if high >= 0 {
			indices <- low
		}
		close(indices)
	}
}

func main() {
}
