package g

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//agent上报itsm的数据结构
type CheckResult struct {
	Sid            int64     `json:"sid"`
	Domain         string    `json:"domain"`
	Target         string    `json:"target"`
	RespCode       string    `json:"resp_code"`
	RespTime       int       `json:"resp_time"`
	Status         int64     `json:"status"`
	PushTime       time.Time `json:"push_time"`
	ErrorInfo      string    `json:"error_info"`
	Port           int       `json:"port"`
	Timeout        int       `json:"timeout"`
	Uri            string    `json:"uri"`
	ExpectCode     string    `json:"expect_code"`
	Region         string    `json:"region"`
	Scheme         string    `json:"scheme"`
	Tables         string    `json:"tables"`
	Project        string    `json:"project"`
	ProbeTarget    string    `json:"probetarget"`
	QuestIP        string    `json:"questIP"`
	InstanceRegion string    `json:"instanceRegion"`
}

//下发给agent的数据结构
type DetectedItem struct {
	Sid        int64  `json:"sid"`
	Domain     string `json:"domain"`
	Target     string `json:"target"`
	Port       int    `json:"port"`
	Timeout    int    `json:"timeout"`
	Uri        string `json:"uri"`
	ExpectCode string `json:"expect_code"`
	Region     string `json:"region"`
	Scheme     string `json:"scheme"`
	Tables     string `json:"tables"`
	Project    string `json:"project"`
}

var Registry = prometheus.NewRegistry()
var RespTime = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "promhttp_metric",
		Subsystem: "urlooker",
		Name:      "respon_time",
		Help:      "show prob url respon time.",
	},
	[]string{
		"InstanceRegion",
		"QuestIP",
		"ProbeTarget",
		"sid",
		"url",
		"Domain",
		"Project",
		"Tables",
	},
)

var StatusCode = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "promhttp_metric",
		Subsystem: "urlooker",
		Name:      "status_code",
		Help:      "show prob url respon http code.",
	},
	[]string{
		"InstanceRegion",
		"QuestIP",
		"ProbeTarget",
		"sid",
		"url",
		"Domain",
		"Project",
		"Tables",
	},
)

func Init() {
	Registry.MustRegister(RespTime)
	Registry.MustRegister(StatusCode)
}
