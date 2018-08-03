package DistributedDS

type IDIntArrayPartition interface {
	set(int, int64)
	get(int64) int
	getRange() (int64, int64)
}

type DIntArrayPartition struct {
	localArray []int
	start      int64
	end        int64
}

func CreateDIntArrayPartition(start int64, end int64) DIntArrayPartition {
	return DIntArrayPartition{
		localArray: make([]int, end-start),
		start:      start,
		end:        end,
	}
}

func (sp *DIntArrayPartition) set(data int, index int64) {
	sp.localArray[index-sp.start] = data
}

func (sp *DIntArrayPartition) get(index int64) int {
	return sp.localArray[index-sp.start]
}

func (sp *DIntArrayPartition) getRange() (int64, int64) {
	return sp.start, sp.end
}
