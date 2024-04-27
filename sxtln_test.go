package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"testing"

	"sextillion.io/cli/models"
)

func TestCheckAuth(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	conf.ApiKey = os.Getenv("SXTLN_APIKEY")
	checkAuth(conf)
}

func TestLogin(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	conf.ApiKey = os.Getenv("SXTLN_APIKEY")
	log.Println("check login")
	resBody := doLogin(conf, "sagi.forbes@gmail.com", "Qwerty098&")
	print(conf, resBody)

	conf.ApiKey = ""
	conf.Token = resBody["token"].(string)

	checkAuth(conf)

	saveConfig(conf)
	_, err := os.Stat(configFilePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Error(err)

	}

	fmt.Print(string(configContent))

}

func TestConfig(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	var val = []string{"test value"}
	doSetConfig(conf, "token", val)
	if conf.Token != val[0] {
		t.Error("failed to set token key to ", val[0])

	}
	doSetConfig(conf, "token", nil)
	if conf.Token != "" {
		t.Error("failed to clear token key")

	}

	doSetConfig(conf, "apiKey", val)
	if conf.ApiKey != val[0] {
		t.Error("failed to set apiKey key to ", val[0])

	}

	ret := doGetConfig(conf, "apiKey")
	if ret["apiKey"] != val[0] {
		t.Error("failed to get apiKey via command")
	} else {
		print(conf, ret)
	}

	doSetConfig(conf, "apiKey", nil)
	if conf.Token != "" {
		t.Error("failed to clear apiKey key")

	}

	conf.Token = "test token"
	conf.ApiKey = "test api key"

	handleConfig(conf, []string{"view"})
}

var clusterNames = []string{"cli_testing_cluster"}

func TestClusterList(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	conf.ApiKey = os.Getenv(`SXTLN_APIKEY`)

	res := listCluster(conf, "")
	print(conf, res)

}

func TestClusterCreate(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	conf.ApiKey = os.Getenv(`SXTLN_APIKEY`)

	t.Log("creating single cluster")
	res := createCluster(conf, clusterNames[0], "b1", 1, false, true)
	print(conf, res)

}

func TestClusterDeleteAll(t *testing.T) {
	println(t.Name())
	var conf = models.NewConfig()
	conf.ApiKey = os.Getenv(`SXTLN_APIKEY`)

	res := listCluster(conf, "")

	var clusters = res[`clusters`].([]interface{})

	for _, cluster := range clusters {
		clusterInfo := cluster.(models.JSON)
		nameIdx := slices.IndexFunc(clusterNames, func(clusterName string) bool { return clusterName == clusterInfo[`name`].(string) })
		if nameIdx > -1 {
			deleteCluster(conf, clusterInfo[`clusterId`].(string))
		}
	}

}
