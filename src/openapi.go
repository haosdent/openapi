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
        "AddDisk": ["RegionId", "InstanceId", "Size", "SnapshotId", "ClientToken"],
        "AllocatePublicIpAddress": ["RegionId", "InstanceId"],
        "AuthorizeSecurityGroup": ["RegionId", "SecurityGroupId", "IpProtocol", "PortRange", "SourceGroupId", "SourceCidrIp", "Policy", "NicType"],
        "CreateImage": ["RegionId", "SnapshotId", "ImageVersion", "Description", "OSName", "ClientToken"],
        "CreateInstance": ["RegionId", "ImageId", "InstanceType", "SecurityGroupId", "InstanceName", "InternetChargeType", "InternetMaxBandwidthIn", "InternetMaxBandwidthOut", "HostName", "Password", "SystemDisk.Category", "DataDisk.n.Size", "DataDisk.n.Category", "DataDisk.n.SnapshotId", "ClientToken"],
        "CreateSecurityGroup": ["RegionId", "Description", "ClientToken"],
        "CreateSnapshot": ["RegionId", "InstanceId", "DiskId", "SnapshotName", "ClientToken"],
        "DescribeImages": ["RegionId", "PageNumber", "PageSize", "ImageId", "ImageOwnerAlias"],
        "DescribeInstanceAttribute": ["RegionId", "InstanceId"],
        "DescribeInstanceDisks": ["RegionId", "InstanceId"],
        "DescribeInstanceStatus": ["RegionId", "PageNumber", "PageSize"],
        "DescribeInstanceTypes": ["RegionId"],
        "DescribeRegions": [],
        "DescribeSecurityGroupAttribute": ["RegionId", "SecurityGroupId", "NicType"],
        "DescribeSecurityGroups": ["RegionId", "PageNumber", "PageSize"],
        "DescribeSnapshotAttribute": ["RegionId", "SnapshotId"],
        "DescribeSnapshots": ["RegionId", "InstanceId", "DiskId"],
        "DeleteDisk": ["RegionId", "InstanceId", "DiskId"],
        "DeleteImage": ["RegionId", "ImageId"],
        "DeleteInstance": ["RegionId", "InstanceId"],
        "DeleteSecurityGroup": ["RegionId", "SecurityGroupId"],
        "DeleteSnapshot": ["RegionId", "DiskId", "InstanceId", "SnapshotId"],
        "JoinSecurityGroup": ["RegionId", "InstanceId", "SecurityGroupId"],
        "LeaveSecurityGroup": ["RegionId", "InstanceId", "SecurityGroupId"],
        "ModifyInstanceAttribute": ["RegionId", "InstanceId", "InstanceName", "Password", "HostName", "SecurityGroupId"],
        "RebootInstance": ["RegionId", "InstanceId", "ForceStop"],
        "ResetDisk": ["RegionId", "InstanceId", "DiskId", "SnapshotId"],
        "RevokeSecurityGroup": ["RegionId", "SecurityGroupId", "IpProtocol", "PortRange", "SourceGroupId", "SourceCidrIp", "Policy", "NicType"],
        "StartInstance": ["RegionId", "InstanceId", "ForceStop"],
        "StopInstance": ["RegionId", "InstanceId", "ForceStop"]
    },
    "rds": {
        "DescribeDBInstances": ["RegionId", "DBInstanceId", "DBInstanceStatus", "Engine", "DBInstanceNetType"],
        "SwitchDBInstanceNetType": ["RegionId", "DBInstanceId", "ConnectionStringPrefix"],
        "RestartDBInstance": ["RegionId", "DBInstanceId"],
        "CreateDatabase": ["RegionId", "DBInstanceId", "DBName", "CharacterSetName", "AccountName", "AccountPrivilege", "DBDescription"],
        "DeleteDatabase": ["RegionId", "DBInstanceId", "DBName"],
        "DescribeDatabases": ["RegionId", "DBInstanceId", "DBName", "DBStatus"],
        "CreateImportDataUpload": ["RegionId", "DBInstanceId", "DBName"],
        "ImportData": ["RegionId", "DBInstanceId", "FileName"],
        "DescribeDataFiles": ["RegionId", "DBInstanceId", "StartTime", "EndTime"],
        "CreateAccount": ["RegionId", "DBInstanceId", "AccountName", "AccountPassword", "DBName", "AccountPrivilege", "AccountDescription"],
        "ModifyAccountAttribute": ["RegionId", "DBInstanceId", "AccountName", "AccountPassword", "OldAccountPassword", "AccountPrivilege"],
        "GrantAccountPrivilege": ["RegionId", "DBInstanceId", "AccountName", "DBName"],
        "RevokeAccountPrivilege": ["RegionId", "DBInstanceId", "AccountName", "DBName"],
        "DescribeAccounts": ["RegionId", "DBInstanceId", "AccountName", "DBName"],
        "ModifySecurityIps": ["RegionId", "DBInstanceId", "SecurityIps"],
        "DescribeSecurityIps": ["RegionId", "DBInstanceId"],
        "CreateBackup": ["RegionId", "DBInstanceId"],
        "DescribeBackups": ["RegionId", "DBInstanceId", "BackupStatus", "BackupMode"],
        "RestoreDBInstance": ["RegionId", "DBInstanceId", "BackupId"],
        "DescribeResourceUsage": ["RegionId", "DBInstanceId", "StartTime", "EndTime"],
        "DescribeDBInstancePerformance": ["RegionId", "DBInstanceId", "Key", "StartTime", "EndTime"]
    },
    "slb": {
        "CreateLoadBalancer": ["RegionId", "LoadBalancerName", "AddressType", "InternetChargeType", "Bandwidth", "ClientToken"],
        "DeleteLoadBalancer": ["RegionId", "LoadBalancerId"],
        "ModifyLoadBalancerInternetSpec": ["RegionId", "LoadBalancerId", "InternetChargeType", "Bandwidth"],
        "SetLoadBalancerStatus": ["RegionId", "LoadBalancerId", "LoadBalancerStatus"],
        "SetLoadBalancerName": ["RegionId", "LoadBalancerId", "LoadBalancerName"],
        "DescribeLoadBalancers": ["RegionId", "LoadBalancerId", "AddressType", "InternetChargeType", "ServerId"],
        "DescribeLoadBalancerAttribute": ["RegionId", "LoadBalancerId"],
        "DescribeRegions": [],
        "CreateLoadBalancerHTTPListener": ["RegionId", "LoadBalancerId", "ListenerPort", "BackendServerPort", "Bandwidth", "XForwardedFor", "Scheduler", "StickySession", "StickySessionType", "CookieTimeout", "Cookie", "HealthCheck", "HealthCheckDomain", "HealthCheckURI", "HealthCheckConnectPort", "HealthyThreshold", "UnhealthyThreshold", "UnhealthyThreshold", "HealthCheckInterval"],
        "CreateLoadBalancerTCPListener": ["RegionId", "LoadBalancerId", "ListenerPort", "BackendServerPort", "Bandwidth", "Scheduler", "PersistenceTimeout", "HealthCheckConnectPort", "HealthyThreshold", "UnhealthyThreshold", "HealthCheckConnectTimeout", "HealthCheckInterval"],
        "DeleteLoadBalancerListener": ["RegionId", "LoadBalancerId", "ListenerPort"],
        "StopLoadBalancerListener": ["RegionId", "LoadBalancerId", "ListenerPort"],
        "StartLoadBalancerListener": ["RegionId", "LoadBalancerId", "ListenerPort"],
        "SetLoadBalancerHTTPListenerAttribute": ["RegionId", "LoadBalancerId", "ListenerPort", "BackendServerPort", "Bandwidth", "XForwardedFor", "Scheduler", "StickySession", "StickySessionType", "CookieTimeout", "Cookie", "HealthCheck", "HealthCheckDomain", "HealthCheckURI", "HealthCheckConnectPort", "HealthyThreshold", "UnhealthyThreshold", "UnhealthyThreshold", "HealthCheckInterval"],
        "SetLoadBalancerTCPListenerAttribute": ["RegionId", "LoadBalancerId", "ListenerPort", "BackendServerPort", "Bandwidth", "Scheduler", "PersistenceTimeout", "HealthCheckConnectPort", "HealthyThreshold", "UnhealthyThreshold", "HealthCheckConnectTimeout", "HealthCheckInterval"],
        "DescribeLoadBalancerHTTPListenerAttribute": ["RegionId", "LoadBalancerId", "ListenerPort"],
        "DescribeLoadBalancerTCPListenerAttribute": ["RegionId", "LoadBalancerId", "ListenerPort"],
        "AddBackendServers": ["RegionId", "LoadBalancerId", "BackendServers"],
        "RemoveBackendServers": ["RegionId", "LoadBalancerId", "BackendServers"],
        "DescribeHealthStatus": ["RegionId", "LoadBalancerId", "ListenerPort"]
    },
    "ess": {
        "CreateScalingGroup": ["RegionId", "MaxSize", "MinSize", "ScalingGroupName", "DefaultCooldown", "RemovalPolicy.N", "LoadBalancerId", "DBInstanceId.N"],
        "ModifyScalingGroup": ["ScalingGroupId", "ScalingGroupName", "ActiveScalingConfigurationId", "MinSize", "MaxSize", "DefaultCooldown", "RemovalPolicy.N"],
        "DescribeScalingGroups": ["RegionId", "ScalingGroupId.N", "ScalingGroupName.N", "PageNumber", "PageSize"],
        "EnableScalingGroup": ["ScalingGroupId", "ActiveScalingConfigurationId", "InstanceId.N"],
        "DisableScalingGroup": ["ScalingGroupId"],
        "DeleteScalingGroup": ["ScalingGroupId", "ForceDelete"],
        "DescribeScalingInstances": ["RegionId", "ScalingGroupId", "ScalingConfigurationId", "InstanceId.N", "HealthStatus", "LifecycleState", "CreationType", "PageNumber", "PageSize"],
        "CreateScalingConfiguration": ["ScalingGroupId", "ImageId", "InstanceType", "SecurityGroupId", "ScalingConfigurationName", "InternetChargeType", "InternetMaxBandwidthIn", "InternetMaxBandwidthOut", "SystemDisk.Category"],
        "DescribeScalingConfigurations": ["RegionId", "ScalingGroupId", "ScalingConfigurationId.N", "ScalingConfigurationName.N", "PageNumber", "PageSize"],
        "DeleteScalingConfiguration": ["ScalingConfigurationId"],
        "CreateScalingRule": ["ScalingGroupId", "AdjustmentType", "AdjustmentValue", "ScalingRuleName", "Cooldown"],
        "ModifyScalingRule": ["ScalingRuleId", "AdjustmentType", "AdjustmentValue", "ScalingRuleName", "Cooldown"],
        "DescribeScalingRules": ["RegionId", "ScalingGroupId", "ScalingRuleId.N", "ScalingRuleName.N", "ScalingRuleAri.N", "PageNumber", "PageSize"],
        "DeleteScalingRule": ["ScalingRuleId"],
        "ExecuteScalingRule": ["ScalingRuleAri", "ClientToken"],
        "AttachInstances": ["ScalingGroupId", "InstanceId.N"],
        "RemoveInstances": ["ScalingGroupId", "InstanceId.N"],
        "CreateScheduledTask": ["RegionId", "ScheduledAction", "LaunchTime", "ScheduledTaskName", "Description", "LaunchExpirationTime", "RecurrenceType", "RecurrenceValue", "RecurrenceEndTime", "TaskEnabled"],
        "ModifyScheduledTask": ["ScheduledTaskId", "ScheduledTaskName", "Description", "ScheduledAction", "LaunchTime", "LaunchExpirationTime", "RecurrenceType", "RecurrenceValue", "RecurrenceEndTime", "TaskEnabled"],
        "DescribeScheduledTasks": ["RegionId", "ScheduledTaskId.N", "ScheduledTaskName.N", "ScheduledAction.N", "PageNumber", "PageSize"],
        "DeleteScheduledTask": ["ScheduledTaskId"]
    }
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
			version = "2014-05-26"
		case "slb":
			endpoint = "slb.aliyuncs.com"
			version = "2014-05-15"
		case "rds":
			endpoint = "rds.aliyuncs.com"
			version = "2013-05-28"
		case "ess":
			endpoint = "ess.aliyuncs.com"
			version = "2014-08-28"
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
