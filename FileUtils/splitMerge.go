package fileutils

import (
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

const chunkSize int64 = 4 << 20 /* 4 M */

// SplitIntoFiles splits file into files (fileName1, fileName2, ...), each file contains a chunck.
func SplitIntoFiles(path string) error {
	data, err := Split(path)
	if err != nil {
		return err
	}
	err = StorePartitions(data, path)
	if err != nil {
		return err
	}
	return nil
}

// Split file into array of []byte, i.e., []*[]byte
func Split(path string) ([]*[]byte, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	count := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))
	res := make([]*[]byte, count)

	fi, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	var i int64 = 1
	for ; i <= int64(count); i++ {
		b := make([]byte, chunkSize)
		fi.Seek((i-1)*(chunkSize), 0)
		if len(b) > int((fileInfo.Size() - (i-1)*chunkSize)) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		res[i-1] = &b
	}
	fi.Close()
	return res, nil
}

// SplitHash calculates the hash code of each partition of a file
func SplitHash(path string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	count := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))
	res := make([]string, count)

	fi, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	var i int64 = 1
	for ; i <= int64(count); i++ {
		b := make([]byte, chunkSize)
		fi.Seek((i-1)*(chunkSize), 0)
		if len(b) > int((fileInfo.Size() - (i-1)*chunkSize)) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		hash, err := DataHash(b)
		if err != nil {
			return nil, err
		}
		res[i-1] = hash
	}
	fi.Close()
	return res, nil
}

// StorePartitions write data chuncks into files
func StorePartitions(chunks []*[]byte, path string) error {
	for i := 0; i < len(chunks); i++ {
		f, err := os.OpenFile(path+"_part"+strconv.Itoa(int(i)), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		f.Write(*chunks[i])
		f.Close()
	}
	return nil
}

// Merge partitions into a single file, the path should be the prefix of the partition
// For example, myfile_part0, myfile_part1, myfile_part2, ... myfile_part9 would be myfile
// count should be 10
func Merge(path string, count int, mergedPath string, deletePartition bool) error {
	/* delete target file if exist */
	os.Remove(mergedPath)

	fii, err := os.OpenFile(mergedPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		filePaht := path + "_part" + strconv.Itoa(int(i))
		f, err := os.OpenFile(filePaht, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		fii.Write(b)
		f.Close()
		if deletePartition {
			os.Remove(filePaht)
		}
	}
	return nil
}
