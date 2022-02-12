package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/seriar-org/zed/gzc"
)

func main() {
	fmt.Println("Who's Zed?")
	repoID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Cannot convert repository id to int")
	}
	epicID, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("Cannot convert epic id to int")
	}
	timeout, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("Cannot convert timeout to int")
	}
	a := gzc.CreateAPI(&http.Client{}, "https://api.zenhub.com").WithTimeout(timeout)
	c := gzc.CreateClient(a, os.Args[1])

	e, err := c.RequestEpic(repoID, epicID)
	if err != nil {
		fmt.Printf("error %+v\n", err)
	}
	fmt.Printf("Resp: %+v\n", e)

	d, err := c.RequestDependencies(repoID)
	if err != nil {
		fmt.Printf("error %+v\n", err)
	}
	fmt.Printf("Resp: %+v\n", d)
}
