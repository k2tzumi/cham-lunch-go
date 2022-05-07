package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getAreaJSON() ([]byte, error) {
	url := "https://www.jma.go.jp/bosai/common/const/area.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return byteArray, nil
}

func main() {
	json, err := getAreaJSON()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// FIXME: "./src/area.json" を受け取れるようにする
	err = os.WriteFile("./src/area.json", json, 0o644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Updated src/area.json.")
}
