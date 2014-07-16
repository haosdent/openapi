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
	"desc": "Usage: openapi [production] [action] [property=value]\nAllow production:",
	"child": {
		"ecs": {
			"desc": "Usage: openapi ecs [action] [property=value]\nAllow action:",
			"child": {
				"DescribeInstanceAttribute": {
					"desc": "Usage: openapi ecs DescribeInstanceAttribute [property=value]\nAllow property:",
					"child": {
						"InstanceId": {
							"desc": ""
						},
						"RegionId": {
							"desc": ""
						}
					}
				}
			}
		},
		"slb": {
			"desc": "",
			"child": {
				"DescribeInstanceAttribute": {
					"desc": "",
					"child": {
						"InstanceId": {
							"desc": "Usage: openapi [production] [action] [property=value]\nAllow production:"
						},
						"RegionId": {
							"desc": "Usage: openapi [production] [action] [property=value]\nAllow production:"
						}
					}
				}
			}
		},
		"rds": {
			"desc": "",
			"child": {
				"DescribeInstanceAttribute": {
					"desc": "",
					"child": {
						"InstanceId": {
							"desc": "Usage: openapi [production] [action] [property=value]\nAllow production:"
						},
						"RegionId": {
							"desc": "Usage: openapi [production] [action] [property=value]\nAllow production:"
						}
					}
				}
			}
		}
	}
}
`

const keyFile = ".aliyuncredentials"
var accessKey string
var accessId string

type Config struct {
	Credentials struct {
		Accesskeyid string
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
	//fmt.Println(msg)
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
	var cfg Config
	err = gcfg.ReadFileInto(&cfg, usr.HomeDir + "/" + keyFile)
	if err == nil {
		//fmt.Println(.Get("Credentials", "accesskeyid"))
		//fmt.Println(cfg.Credentials.Accesskeyid)
		//fmt.Println(cfg.Credentials.Accesskeysecret)
		accessId = cfg.Credentials.Accesskeyid
		accessKey = cfg.Credentials.Accesskeysecret
		fmt.Println(accessId)
		fmt.Println(accessKey)
		os.Exit(-1)
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

	t := time.Now()
	(*p)["SignatureNonce"] = uuid.NewUUID().String()
	(*p)["TimeStamp"] = t.Format(layout)
	(*p)["TimeStamp"] = "2014-07-09T15:25:06Z"
	(*p)["Signature"] = ComputeSign(p)
}

func help(args []string) {
	//fmt.Println(args)
	data := map[string]interface{}{}
	decoder := json.NewDecoder(strings.NewReader(manual))
	decoder.Decode(&data)

	for _, arg := range args {
		tmp := data["child"].(map[string]interface{})[arg]
		if tmp != nil {
			data = tmp.(map[string]interface{})
		} else {
			break
		}
	}

	fmt.Println(data["desc"])
	for i, _ := range data["child"].(map[string]interface{}) {
		fmt.Print("\t")
		fmt.Println(i)
	}
}

func main() {
	fmt.Println(os.Args)
	InitAccessIdKey()

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
