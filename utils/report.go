package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	//"reflect"

	"agent/g"

	"gopkg.in/olivere/elastic.v5" //这里使用的是版本5，最新的是6，有改动
)

func Report(data *g.CheckResult) {

	// 定义recover方法，在后面程序出现异常的时候就会捕获
	defer func() {
		if r := recover(); r != nil {
			// 这里可以对异常进行一些处理和捕获
			log.Println("[Recovered]:", r)
		}
	}()

	var client *elastic.Client

	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(g.Config.Web.ReportAddr))
	if err != nil {
		log.Println("[ERROR]:", err)
		return
	}
	_, _, err = client.Ping(g.Config.Web.ReportAddr).Do(context.Background())
	if err != nil {
		log.Println("[ERROR]:", err)
	}
	//log.Println("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	_, err = client.ElasticsearchVersion(g.Config.Web.ReportAddr)
	if err != nil {
		log.Println("[ERROR]:", err)
	}

	//log.Println("Elasticsearch version %s\n", esversion)

	jsonStr, err := json.Marshal(data)

	//log.Println("[jsonStr]:", string(jsonStr))

	if err != nil {
		//log.Println("[ERROR]:", err)
		log.Println("[ERROR]:err")
	}

	e2 := string(jsonStr)

	//log.Println(reflect.TypeOf(e2))

	_, err = client.Index().
		Index(g.Config.Web.Index).
		Type(g.Config.Web.Type).
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		log.Println("[ERROR]:", err)
	}

	//log.Println("Indexed tweet %s to index s%s, type %s\n", put2.Index, put2.Type)
}
