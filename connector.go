package khukuri

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"os"
	"strconv"
	"time"
)

func setupRedisConnection(timeOutSeconds int) (*redis.Client, error) {
	//TODO: Read these values from some config file of sorts.
	return redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(timeOutSeconds)*time.Second)
}

func performErrorCheck(err error) {
	if err != nil {
		//TODO: Need to use logging instead.
		fmt.Println("Error while setting a redis connection")
	}
}

func LookupAlias(alias string) (string, error) {
	c, err := setupRedisConnection(10)
	performErrorCheck(err)
	defer c.Close()
	lookupId, err := DecodeFromBase(alias)
	if err != nil {
		fmt.Println("ERROR!!!!! Can't convert string id")
		return "", err
	}
	s, err := c.Cmd("get", lookupId).Str()
	if err != nil {
		fmt.Println("ERROR!!!!! Can't convert string id")
		return "", err
	}
	return s, nil

}

func StoreUrl(baseUrl string) (string, error) {
	//Map and store this baseUrl. Return the alias string.
	//Before storing we increase the value of the global counter by 1
	c, err := setupRedisConnection(10)
	performErrorCheck(err)
	defer c.Close()
	//need to do this in a transaction
	// first get the current value
	// increment the current value
	rep := c.Cmd("multi")
	performErrorCheck(rep.Err)

	currentCounter, err := c.Cmd("get", "globalCounter").Str()

	if err != nil {
		fmt.Println("ERROR. Cannot get the current counter")
		os.Exit(1)
	}

	if currentCounter == "" {
		resp := c.Cmd("set", "globalCounter", "1")
		currentCounter = "1"
		performErrorCheck(resp.Err)
	} else {
		res := c.Cmd("incr", "globalCounter")
		performErrorCheck(res.Err)
	}

	rep = c.Cmd("exec")
	performErrorCheck(rep.Err)

	idNumber, err := strconv.ParseUint(currentCounter, 10, 64)
	performErrorCheck(err)

	res, err := c.Cmd("setnx", idNumber).Bool()
	performErrorCheck(err)

	if res == false {
		fmt.Println("The ID already exits. ERROR!!")
	}

	return EncodeToBase(idNumber)

}
