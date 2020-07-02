package cron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"agent/g"
	"agent/utils"

	"github.com/astaxie/beego/httplib"

	"github.com/toolkits/file"
)

// 检查url
func StartCheck(w http.ResponseWriter, r *http.Request) {

	log.Println("[INFO]:", "begin--")
	beginTime := time.Now()

	items, err := GetJson()
	if err != nil {
		log.Println("[ERROR:]", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}

	// Use the common pool.
	// log.Println(items)
	var wg sync.WaitGroup
	for _, item := range items {
		// log.Println(item)
		go utils.CheckTargetStatus(item, &wg)
		wg.Add(1)
	}
	wg.Wait()

	durtime := time.Since(beginTime)
	log.Println("[INFO]:", "Spend time ：", durtime)
	log.Println("[INFO]:", "end")

}

func GetJson() ([]g.DetectedItem, error) {

	body, err := GetCacheData()

	var data []g.DetectedItem

	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("[ERROR:]", err)
	}
	log.Println(data)
	return data, err
}

func GetCacheData() ([]byte, error) {
	// 如果文件不存在
	if !file.IsExist(g.Config.CacheFile) {
		// get data from itsm
		QueryItsm()
		return file.ToBytes(g.Config.CacheFile)
	}
	// 如果文件过期
	fileTime, err := file.FileMTime(g.Config.CacheFile)
	if err != nil {
		log.Println("[ERROR:]", err)
	}

	//durtime := time.Since(fileTime)
	//beginTime := time.Now()
	// 时间转换
	log.Println("fileTime", fileTime)
	log.Println("now time", time.Now().Unix())
	//log.Println("durtime", durtime)
	durtime := time.Since(time.Unix(fileTime, 0))
	log.Println("[INFO]:", "Spend time ：", durtime.Seconds())
	log.Println("[INFO]:", "end")

	CacheTimeSec, err := strconv.ParseFloat(g.Config.CacheTimeSec, 64)
	if err != nil {
		log.Println("[ERROR:]", err)
	}

	if durtime.Seconds() > CacheTimeSec {
		log.Println("[INFO]:", "缓存过期")
		// get data from itsm
		QueryItsm()
	} else {
		log.Println("[INFO]:", "缓存有效")
	}

	return file.ToBytes(g.Config.CacheFile)
}

func QueryItsm() {
	// 过期则清除metrics
	g.StatusCode.Reset()
	g.RespTime.Reset()
	// log.Println(strings.Join([]string{g.Config.GetUrlAddr, "?region=", g.Config.Region}, ""))
	itsmUrl := strings.Join([]string{g.Config.GetUrlAddr, "?region=", g.Config.InstanceRegion}, "")
	// strings.Join([]string{"hello", "world"}, "")
	req := httplib.Get(itsmUrl)
	resp, err := req.Response()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR:]", err)
	}

	file.WriteBytes("all_url.data", body)
}
