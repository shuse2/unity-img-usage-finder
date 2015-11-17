package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Data struct {
	path         string
	guid         string
	dependencies []string
}

var (
	root string
	data []Data
)

func setupData(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, "png.meta") {
		guid, err := getGuid(path)
		if err != nil {
			panic(err)
		}
		data = append(data, Data{path: path, guid: guid})
	}
	return nil
}

func getGuid(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file %s", path)
		return "", err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("Failed to read file %s", path)
		return "", err
	}
	guid := m["guid"].(string)
	return guid, nil
}

func setupDependencies(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, ".prefab") {
		for i, element := range data {
			if checkContainsGuid(path, element.guid) {
				// includes
				data[i].dependencies = append(element.dependencies, path)
			}
		}
		// fmt.Printf("%v", data)
	}
	return nil
}

func checkContainsGuid(path string, guid string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file %s", path)
		return false
	}
	return strings.Contains(string(data), guid)
}

func showResult() {
	for _, element := range data {
		fmt.Printf("---------------------------------------------------\n")
		fmt.Printf("IMG: %s\n", strings.TrimSuffix(element.path, ".meta"))
		count := len(element.dependencies)
		fmt.Printf("\tCount: %d\n", count)
		if count > 0 {
			for index, dep := range element.dependencies {
				fmt.Printf("%d Used at: %s\n", index, dep)
			}
		}
		fmt.Printf("---------------------------------------------------\n")
	}
}

func main() {
	flag.Parse()
	root = flag.Arg(0)
	if root == "" {
		fmt.Printf("root path is not selected")
		os.Exit(0)
	}
	data = make([]Data, 1)
	err := filepath.Walk(root, setupData)
	if err != nil {
		fmt.Printf("filepath.Walk() returned error :%v\n", err)
	}
	err = filepath.Walk(root, setupDependencies)
	if err != nil {
		fmt.Printf("filepath.Walk() returned error :%v\n", err)
	}
	showResult()
}
