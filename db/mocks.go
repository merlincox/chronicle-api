package db

import (
	"io/ioutil"
	"fmt"
	"os"
)

func getLocalSkeletonData() []byte {

	skeletonBytes, err := ioutil.ReadFile("./mock/data.json")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return skeletonBytes
}

