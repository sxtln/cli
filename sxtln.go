package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"sextillion.io/cli/api"
	"sextillion.io/cli/models"
)

var userFolder, _ = os.UserHomeDir()

var configFolder = fmt.Sprintf("%s/.sxtln", userFolder)
var configFilePath = fmt.Sprintf("%s/config", configFolder)

func exitIfErr(err error, httpResCode ...int) {
	var code = 0
	if len(httpResCode) > 0 {
		code = httpResCode[0]
	}
	if err != nil {
		if code == 400 {
			exitIfErr(fmt.Errorf(`invalid parameter. check your parameters`))
		} else if code == 401 {
			exitIfErr(fmt.Errorf(`unauthorized. login to sextillion with 'auth login' or 'config set-apikey'. see helm at http://https://github.com/sxtln/cli`))
		} else {
			log.Fatal(err)
		}
	}
}

func objToJson(o interface{}) models.JSON {
	raw, err := json.Marshal(o)
	exitIfErr(err)
	var ret models.JSON
	err = json.Unmarshal(raw, &ret)
	exitIfErr(err)
	return ret
}

func saveConfig(config *models.Config) {
	_, err := os.Stat(configFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(configFolder, 0755)
		exitIfErr(err)
	}

	fileContent, err := yaml.Marshal(config)
	exitIfErr(err)

	err = os.WriteFile(configFilePath, fileContent, 0644)
	exitIfErr(err)

}

func print(config *models.Config, v interface{}) {
	var out []byte
	var err error

	if config.Output == "yaml" {
		out, err = yaml.Marshal(v)
	} else {
		out, err = json.Marshal(v)
	}

	exitIfErr(err)
	fmt.Println(string(out))

}

func loadConfig() *models.Config {
	_, err := os.Stat(configFolder)
	if os.IsNotExist(err) {
		return models.NewConfig()
	}
	_, err = os.Stat(configFilePath)
	if err == nil {
		content, err := os.ReadFile(configFilePath)
		if os.IsNotExist(err) {
			return models.NewConfig()
		}
		exitIfErr(err)
		var ret models.Config
		err = json.Unmarshal(content, &ret) //try to load as json
		if err != nil {
			err = yaml.Unmarshal(content, &ret) //try to load as yaml
			exitIfErr(err)
			return &ret
		} else {
			return &ret
		}

	}
	return models.NewConfig()
}

func main() {

	var token string
	var apiKey string
	var outputFormat string

	config := loadConfig()
	//fmt.Printf("%#v\n", config)

	flag.StringVar(&outputFormat, "o", "json", "format of output, one of json,yaml default json")
	flag.StringVar(&outputFormat, "output", "json", "format of output, one of json,yaml default json")

	if outputFormat == "json" || outputFormat == "yaml" {
		config.Output = outputFormat
	}

	flag.StringVar(&token, "token", "", "the token to be used when making api call")
	flag.StringVar(&apiKey, "apikey", "", "api key to used when making api calls")

	flag.Parse()

	if token != "" {
		config.Token = token
	}
	if apiKey != "" {
		config.ApiKey = apiKey
	}

	if len(flag.Args()) < 1 {
		println("no command given. nothing todo")
		os.Exit(1)
		return
	}

	command := flag.Args()[0]
	switch command {
	case "auth":
		handleAuth(config, flag.Args()[1:])
	case "config":
		handleConfig(config, flag.Args()[1:])
	case "cluster":
		handleCluster(config, flag.Args()[1:])
	}

}

// //////////////////////////////////////////////////////////////////
// /
// /   cluster
// /
// //////////////////////////////////////////////////////////////////////
func handleCluster(config *models.Config, args []string) {
	var clusterName string
	var nodeType string
	var nodeCount int
	var clusterId string
	var wait bool
	var forceCreate bool
	clusterFlags := flag.NewFlagSet("cluster", flag.ExitOnError)

	clusterFlags.StringVar(&clusterName, "n", "", "the name of cluster to allocate")
	clusterFlags.StringVar(&clusterName, "name", clusterName, "the name of cluster to allocate")

	clusterFlags.StringVar(&nodeType, "t", "b1", "Type of each node in the cluster. default b1")
	clusterFlags.StringVar(&nodeType, "type", nodeType, "Type of each node in the cluster. default b1")

	clusterFlags.IntVar(&nodeCount, "c", 1, "number of nodes in the cluster")
	clusterFlags.IntVar(&nodeCount, "count", nodeCount, "number of nodes in the cluster")

	clusterFlags.StringVar(&clusterId, "i", "", "id of a cluster")
	clusterFlags.StringVar(&clusterId, "id", clusterId, "id of a cluster")

	clusterFlags.BoolVar(&wait, "w", false, "waiting for operation to end")
	clusterFlags.BoolVar(&wait, "wait", wait, "waiting for operation to end")

	clusterFlags.BoolVar(&forceCreate, "f", false, "create the cluster even if there is a cluster with the same name")
	clusterFlags.BoolVar(&forceCreate, "force", forceCreate, "create the cluster even if there is a cluster with the same name")

	if len(args) < 1 {
		exitIfErr(fmt.Errorf(`missing sub command for cluster:
		ls
		create
		delete
		kubeconfig
		k8s
		kube-config
		config
		`))

	}
	subCommand := args[0]
	clusterFlags.Parse(args[1:])
	switch subCommand {
	case "ls":
		cluster_id := clusterId
		if len(clusterFlags.Args()) > 1 {
			cluster_id = clusterFlags.Args()[1:2][0]
		}
		res := listCluster(config, cluster_id)
		print(config, res)
	case "create":
		res := createCluster(config, clusterName, nodeType, nodeCount, forceCreate, wait)
		print(config, res)
	case "delete":
		cluster_id := clusterId
		if cluster_id == "" {
			if len(clusterFlags.Args()) > 1 {
				cluster_id = clusterFlags.Args()[1:2][0]
			}
		}
		if cluster_id == "" {
			exitIfErr(fmt.Errorf(`you must set cluster id. use --id or -i`))
		}
		deleteCluster(config, cluster_id)
	case "config", "kc", "kube-config", "kubeconfig":
		cluster_id := clusterId
		if cluster_id == "" {
			if len(clusterFlags.Args()) > 1 {
				cluster_id = clusterFlags.Args()[1:2][0]
			}
		}
		if cluster_id == "" {
			exitIfErr(fmt.Errorf(`you must set cluster id. use --id or -i`))
		}
		getClusterKubeConfig(config, cluster_id)
	}

}

func getClusterKubeConfig(config *models.Config, clusterId string) {
	api := api.NewApi(*config)
	code, body, err := api.Get(fmt.Sprintf(`/sc/cluster/%s/kubeconfig`, clusterId))
	exitIfErr(err, code)
	fmt.Println(string(body[`body`].([]byte))) //body of REST is byte array of text
}

func deleteCluster(config *models.Config, clusterId string) {
	api := api.NewApi(*config)

	code, _, err := api.Delete(fmt.Sprintf(`/sc/cluster/%s`, clusterId))
	exitIfErr(err, code)
}

func createCluster(config *models.Config, clusterName, nodeType string, nodeCount int, forceCreate, wait bool) models.JSON {
	api := api.NewApi(*config)
	req := make(models.JSON)
	req[`name`] = clusterName
	req[`nodeCount`] = nodeCount
	req[`nodeType`] = nodeType

	if !forceCreate && clusterName != "" {
		var clusterList = listCluster(config, ``)

		var clusters = clusterList[`clusters`].([]interface{})

		for _, cluster := range clusters {
			var clusterInfo = cluster.(models.JSON)
			if clusterInfo[`name`] == clusterName { //no need to recreate the cluster. a cluster with the name already exists
				return clusterInfo
			}
		}
	}

	code, resBody, err := api.Post(`/sc/cluster`, req)
	exitIfErr(err, code)

	if wait {
		var clusterId = resBody[`clusterId`].(string)
		endpoint := fmt.Sprintf("/sc/cluster/%s", clusterId)
		code, resBody, err = api.Get(endpoint)
		if err != nil {
			exitIfErr(fmt.Errorf(``))
		}

		var stageIndex int
		var stageIndexFloat float64
		stageIndexFloat, _ = resBody[`stageIndex`].(float64)
		stageIndex = int(stageIndexFloat)
		for code == 200 && stageIndex < 9 {
			code, resBody, err = api.Get(endpoint)
			if err != nil {
				exitIfErr(fmt.Errorf(``))
			}
			stageIndexFloat, _ = resBody[`stageIndex`].(float64)
			stageIndex = int(stageIndexFloat)
		}

		stageIndexFloat, _ = resBody[`stageIndex`].(float64)
		stageIndex = int(stageIndexFloat)
		if stageIndex != 9 {
			exitIfErr(fmt.Errorf(`cluster is not ready`))
		}
	}

	return resBody

}

func listCluster(config *models.Config, clusterId string) models.JSON {
	api := api.NewApi(*config)
	endpoint := "/sc/cluster"
	if clusterId != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, clusterId)
	}
	code, hbody, err := api.Get(endpoint)
	exitIfErr(err, code)

	exitIfErr(err)
	return hbody
}

// //////////////////////////////////////////////////////////////////
// /
// /   config
// /
// //////////////////////////////////////////////////////////////////////
func handleConfig(config *models.Config, args []string) {
	var key string

	configFlags := flag.NewFlagSet("config", flag.ExitOnError)
	configFlags.StringVar(&key, "k", "", "key of config")
	configFlags.StringVar(&key, "key", key, "key of config")

	if len(configFlags.Args()) < 1 {
		exitIfErr(fmt.Errorf(`auth command as the following sub commands:
		get
		set
		view
		`))
	}

	subCommand := strings.ToLower(configFlags.Args()[0])
	configFlags.Parse(args[1:])

	otherArgs := make([]string, 0)
	if len(configFlags.Args()) > 1 {
		otherArgs = configFlags.Args()[1:]
	}

	switch subCommand {
	case "get":
		res := doGetConfig(config, key)
		print(config, res)
	case "set":
		res := doSetConfig(config, key, otherArgs)
		print(config, res)
	case "view":
		print(config, objToJson(config))
	}

}

func doSetConfig(config *models.Config, key string, otherArgs []string) models.JSON {
	key = strings.ToLower(key)
	var val string
	if len(otherArgs) > 0 {
		val = otherArgs[0]
	}

	wasSet := false

	switch key {
	case "token":
		config.Token = val
		wasSet = true
	case "apikey":
		config.ApiKey = val
		wasSet = true
	case "output":
		if val == "yaml" {
			config.Output = val
		} else {
			config.Output = "json"
		}
		wasSet = true
	}

	if wasSet {
		saveConfig(config)
	}

	return objToJson(config)

}

func doGetConfig(config *models.Config, key string) models.JSON {
	ret := make(models.JSON)
	key = strings.ToLower(key)
	switch key {
	case "token":
		ret["token"] = config.Token
	case "apikey":
		ret["apiKey"] = config.ApiKey
	case "output":
		ret["output"] = config.Output
	}
	return ret
}

////////////////////////////////////////////////////////////////////
///
///   Auth
///
////////////////////////////////////////////////////////////////////////

func handleAuth(config *models.Config, args []string) {
	var user string
	var pwd string
	var apiKey string
	authFlags := flag.NewFlagSet("auth", flag.ExitOnError)

	authFlags.StringVar(&user, "u", "", "user name")
	authFlags.StringVar(&user, "user", user, "user name")

	authFlags.StringVar(&pwd, "p", "", "password")
	authFlags.StringVar(&pwd, "password", pwd, "password")
	authFlags.StringVar(&apiKey, "k", "", "api key")
	authFlags.StringVar(&apiKey, "key", apiKey, "api key")

	if len(args) < 1 {
		exitIfErr(fmt.Errorf(`auth command as the following sub commands:
		check
		login
		set-apikey
		logout
		`))
	}

	subCommand := strings.ToLower(args[0])

	authFlags.Parse(args[1:])

	switch subCommand {
	case "check":
		res := checkAuth(config)
		print(config, res)
	case "set-apikey":
		res := setApiKey(config, apiKey)
		print(config, res)
	case "login":
		res := doLogin(config, user, pwd)
		print(config, res)
	case "logout":
		res := doLogout(config)
		print(config, res)
	}

}

func doLogout(config *models.Config) models.JSON {
	config.Token = ""
	saveConfig(config)
	return objToJson(config)

}

func setApiKey(config *models.Config, apiKey string) models.JSON {
	config.ApiKey = apiKey
	saveConfig(config)
	return objToJson(config)

}

func doLogin(config *models.Config, user, pwd string) models.JSON {
	api := api.NewApi(*config)
	body := models.JSON{
		"login":    user,
		"password": pwd,
	}

	code, resBody, err := api.Post("/auth/login", body)
	exitIfErr(err, code)

	config.Token = resBody[`token`].(string)

	saveConfig(config)

	return resBody
}

func checkAuth(config *models.Config) models.JSON {
	api := api.NewApi(*config)
	code, _, _ := api.Get("/auth/check")
	ret := make(models.JSON)
	if code != 200 {
		if code >= 400 && code < 500 {
			ret["ok"] = false
		} else {
			exitIfErr(fmt.Errorf("token is invalid"), code)
		}
	}

	ret["ok"] = true
	return ret
}
