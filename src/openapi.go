package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"bytes"
	"sort"
	"time"
	"os"
	"net/url"
	"os/user"
	"crypto/sha1"
	"crypto/hmac"
	"encoding/base64"
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/gcfg"
)

const manual = `
{
	"ecs": {
		"StartInstance": ["ImageId", "RegionId"],
		"StopInstance": ["ImageId", "RegionId"]
		},
	"rds": {},
	"slb": {}
}
`

const keyFile = ".aliyuncredentials"

var accessKey string
var accessId string

type Config struct {
	Credentials struct {
		Accesskeyid     string
		Accesskeysecret string
	}
}

func request(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	data := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.Decode(&data)
	b, err := json.MarshalIndent(data, "", "    ")
	fmt.Println(string(b))
}

func ComputeSign(p *map[string]string) string {
	msg := "GET&%2F&" + url.QueryEscape(GenerateQuery(p))
	accessKey := accessKey + "&"
	h := hmac.New(sha1.New, []byte(accessKey))
	h.Write([]byte(msg))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return sign
}

func GenerateQuery(p *map[string]string) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(*p))
	for k := range *p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := (*p)[k]
		prefix := url.QueryEscape(k) + "="
		buf.WriteString(prefix)
		v = url.QueryEscape(v)
		v = strings.Replace(v, "+", "%20", -1)
		v = strings.Replace(v, "*", "%2A", -1)
		v = strings.Replace(v, "%7E", "~", -1)
		buf.WriteString(v)
		buf.WriteByte(byte('&'))
	}
	query := buf.String()
	query = query[:len(query)-1]

	return query
}

func InitAccessIdKey() {
	usr, err := user.Current()
	conf := usr.HomeDir + "/" + keyFile
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		fmt.Println("Config your accessId and accessKey through 'openapi config --id=accessId --secret=accessKey' firstly.")
		os.Exit(-1)
	}

	var cfg Config
	err = gcfg.ReadFileInto(&cfg, conf)
	if err == nil {
		accessId = cfg.Credentials.Accesskeyid
		accessKey = cfg.Credentials.Accesskeysecret
	} else {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func SaveAccessIdKey(accessId string, accessKey string) {
	usr, err := user.Current()
	conf := usr.HomeDir + "/" + keyFile
	dat := []byte(fmt.Sprintf("[Credentials]\naccesskeyid = %s\naccesskeysecret = %s\n", accessId, accessKey))
	err = ioutil.WriteFile(conf, dat, 0644)
	if err == nil {
		fmt.Println("Save accessId and accessKey successfully!")
		os.Exit(0)
	} else {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func UpdateParams(p *map[string]string) {
	const layout = "2006-01-02T15:04:05Z"

	(*p)["Format"] = "JSON"
	(*p)["SignatureVersion"] = "1.0"
	(*p)["SignatureMethod"] = "HMAC-SHA1"
	(*p)["AccessKeyId"] = accessId

	t := time.Now().UTC()
	(*p)["SignatureNonce"] = uuid.NewUUID().String()
	(*p)["TimeStamp"] = t.Format(layout)
	(*p)["Signature"] = ComputeSign(p)
}

func help(args []string) {
	var data interface{}
	decoder := json.NewDecoder(strings.NewReader(manual))
	decoder.Decode(&data)

	counter := 0
	for _, arg := range args {
		if (counter >= 2) {
			break
		}

		tmp := data.(map[string]interface{})[arg]
		if tmp != nil {
			counter++
			data = tmp
		} else {
			break;
		}
	}

	head := "Usage: openapi %s %s [property=value]\nAllow %s:"
	switch counter {
	case 0:
		head = fmt.Sprintf(head, "[production]", "[action]", "production")
	case 1:
		head = fmt.Sprintf(head, args[0], "[action]", "action")
	case 2:
		head = fmt.Sprintf(head, args[0], args[1], "property")
	}
	fmt.Println(head)

	switch e := data.(type) {
	case map[string]interface{}:
	for i, _ := range e {
		fmt.Print("\t")
		fmt.Println(i)
	}
	case []interface{}:
	for _, i := range e {
		fmt.Print("\t")
		fmt.Println(i)
	}
	}
	os.Exit(0)
}

func main() {

	if len(os.Args) <= 1 {
		help(os.Args[1:1])
		return
	} else {
		prod := os.Args[1]
		var endpoint string
		var version string
		switch prod {
		case "ecs":
			endpoint = "ecs.aliyuncs.com"
			version = "2013-01-10"
		case "slb":
			endpoint = "slb.aliyuncs.com"
			version = "2014-05-15"
		case "rds":
			endpoint = "rds.aliyuncs.com"
			version = "2013-05-28"
		case "config":
			l := len(os.Args)
			m := make(map[string]string)
			for i := 2; i < l; i++ {
				arr := strings.SplitN(os.Args[i], "=", 2)
				m[arr[0]] = arr[1]
			}
			SaveAccessIdKey(m["--id"], m["--secret"])
		case "help":
			help(os.Args[2:])
			return
		default:
			help(os.Args[1:1])
			return
		}

		if (len(os.Args) < 3) {
			help(os.Args[1:2])
			return
		}

		InitAccessIdKey()
		action := os.Args[2]
		params := map[string]string{
			"Version": version,
			"Action": action,
		}

		options := os.Args[3:]
		for _, option := range options {
			splits := strings.Split(option, "=")
			k := splits[0]
			v := splits[1]
			params[k] = v
		}

		UpdateParams(&params)
		query := GenerateQuery(&params)
		url := fmt.Sprintf("http://%s/?%s", endpoint, query)
		request(url)
	}

}
