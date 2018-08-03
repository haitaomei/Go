package DistributedDS

type DArray interface {
	Get(int64) int
	Set(int int64)
}

type DIntArray struct {
	dataLocations []IDIntArrayPartition
}

// func CreateDIntArray(size int64, chunk int64) DIntArray {

// }

func (sp *DIntArray) Get(index int64) int {
	for _, dataStore := range sp.dataLocations {
		s, e := dataStore.getRange()
		if index >= s && index <= e {
			return dataStore.get(index)
		}
	}
	return 0
}

func (sp *DIntArray) Set(data int, index int64) {
	for _, dataStore := range sp.dataLocations {
		s, e := dataStore.getRange()
		if index >= s && index <= e {
			dataStore.set(data, index)
			break
		}
	}
}
