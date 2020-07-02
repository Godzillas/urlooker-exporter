package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"agent/g"

	"github.com/astaxie/beego/httplib"
)

const (
	NO_ERROR          = 0
	REQ_TIMEOUT       = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH   = 3
	DNS_ERROR         = 4
	ERROR_TIME_OUT    = "Request time exceeds threshold"
)

// CheckTargetStatus
func CheckTargetStatus(item g.DetectedItem, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	// 逻辑判断 如果是ip直接探测 如果是域名则解析所有对应的IP再探测
	if IsIP(item.Domain) {
		targetIP := item.Domain
		probeTarget := item.Target
		_ = checkTargetStatus(item, targetIP, probeTarget)

		//log.Println("[INFO]:", item.Domain, "is IP")
	} else {
		//log.Println("[INFO]:", item.Domain, "is Domain")
		// 解析IP
		targetIPs, err := LookupIP(item.Domain, 1000)
		if err != nil {
			log.Println("[ERROR]:", item.Domain, err)
			return
		}
		for _, targetIP := range targetIPs {
			//log.Println("[INFO]:", item.Domain, string(targetIP))
			// 组装
			probeTarget := item.Scheme
			probeTarget += "://"
			probeTarget += targetIP
			probeTarget += item.Uri
			//log.Println("ips-------", probeTarget)
			_ = checkTargetStatus(item, targetIP, probeTarget)

		}
	}
	// log.Println(checkResult)
}

func checkTargetStatus(item g.DetectedItem, targetIP string, probeTarget string) (itemCheckResult *g.CheckResult) {
	// log.Println("[INFO]:", item.Domain, string(targetIP), probeTarget)
	itemCheckResult = &g.CheckResult{
		Sid:            item.Sid,
		Domain:         item.Domain,
		Target:         item.Target,
		RespTime:       item.Timeout,
		Port:           item.Port,
		Uri:            item.Uri,
		ExpectCode:     item.ExpectCode,
		Region:         item.Region,
		Scheme:         item.Scheme,
		Tables:         item.Tables,
		Project:        item.Project,
		RespCode:       "0",
		ProbeTarget:    probeTarget,
		QuestIP:        targetIP,
		InstanceRegion: g.Config.InstanceRegion,
	}
	reqStartTime := time.Now()

	req := httplib.Get(probeTarget)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(3*time.Second, 10*time.Second)
	req.Header("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.SetHost(item.Domain)

	resp, err := req.Response()
	//itemCheckResult.PushTime = time.Now().Unix()
	itemCheckResult.PushTime = time.Now()

	if err != nil {
		//log.Println("[ERROR]:", item.Sid, item.Domain, err)
		itemCheckResult.Status = REQ_TIMEOUT
		// 添加错误信息
		itemCheckResult.ErrorInfo = err.Error()
		go Report(itemCheckResult)
		return
	}
	defer resp.Body.Close()
	log.Println(item.Domain)
	respCode := strconv.Itoa(resp.StatusCode)
	itemCheckResult.RespCode = respCode

	respTime := int(time.Now().Sub(reqStartTime).Nanoseconds() / 1000000)
	itemCheckResult.RespTime = respTime

	metricsLabels := make(map[string]string)

	metricsLabels["InstanceRegion"] = itemCheckResult.InstanceRegion
	metricsLabels["QuestIP"] = itemCheckResult.QuestIP
	metricsLabels["ProbeTarget"] = itemCheckResult.ProbeTarget
	metricsLabels["sid"] = strconv.Itoa(int(itemCheckResult.Sid))
	metricsLabels["url"] = itemCheckResult.Target
	metricsLabels["Domain"] = itemCheckResult.Domain
	metricsLabels["Project"] = itemCheckResult.Project
	metricsLabels["Tables"] = itemCheckResult.Tables
	g.RespTime.With(metricsLabels).Set(float64(respTime))
	g.StatusCode.With(metricsLabels).Set(float64(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	if respTime > item.Timeout {
		itemCheckResult.Status = REQ_TIMEOUT
		//log.Println("[ERROR]:", item.Sid, item.Domain, REQ_TIMEOUT)
		// 添加错误信息
		itemCheckResult.ErrorInfo = "响应超时，设定时间: "
		itemCheckResult.ErrorInfo += strconv.Itoa(item.Timeout)
		itemCheckResult.ErrorInfo += "\r\n 响应时间: "
		itemCheckResult.ErrorInfo += strconv.Itoa(respTime)
		itemCheckResult.ErrorInfo += "\r\n"
		itemCheckResult.ErrorInfo += string(body)
		log.Println(itemCheckResult.ErrorInfo)
		go Report(itemCheckResult)
		return
	}

	if strings.Index(respCode, item.ExpectCode) == 0 || (len(item.ExpectCode) == 0 && respCode == "200") {
		itemCheckResult.Status = NO_ERROR
		return

	} else {
		itemCheckResult.Status = INVALID_RESP_CODE
		//log.Println("[ERROR]:", item.Sid, item.Domain, "respCode:", respCode)
		// 状态码不等于200 尝试输出 body

		// 添加错误信息
		itemCheckResult.ErrorInfo = string(body)
		//log.Println("[ERROR]:", itemCheckResult.ErrorInfo)
		go Report(itemCheckResult)

	}
	return
}

func IsIP(ip string) bool {
	if ip != "" {
		isOk, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
		if isOk {
			return isOk
		}
	}
	return false
}
