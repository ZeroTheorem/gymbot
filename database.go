package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func writeData(lvl, exp int64) {
	f, err := os.Create("data.csv")
	if err != nil {
		fmt.Printf("error opened file: %v", err)
		return
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	err = writer.Write([]string{fmt.Sprint(lvl), fmt.Sprint(exp)})
	if err != nil {
		fmt.Printf("error write to file: %v", err)
		return
	}

}

func getData() [2]int64 {
	f, err := os.OpenFile("data.csv", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	row, err := reader.Read()
	if err != nil {
		panic(err)
	}
	var data [2]int64
	for i, v := range row {
		numberConv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("err convert data: %v", err)
			return [2]int64{}
		}
		data[i] = numberConv
	}
	return data
}
