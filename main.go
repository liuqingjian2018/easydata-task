package main

import (
	//"easyerp-task/route"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
)

func init() {
	file, err := os.OpenFile("task.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func AsynTaskRun(val string) {
	url := "http://api.easyerp.test/" + val
	resp, err := http.Get(url)
	if err != nil {
		Error.Println("http get err: ", err.Error())
		//panic(err)
	}
	defer resp.Body.Close()

	Info.Println("http get url: ", url)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error.Println("http body read err: ", err.Error())
		//panic(err)
	}
	Info.Printf("http get resp: %s\n", body)
}

func ConnRedis() redis.Conn {
	conn, err := redis.Dial("tcp", "r-wz9yhnh00c5sgr28k6pd.redis.rds.aliyuncs.com:6379")
	if err != nil {
		Error.Println("redis conn error: ", err.Error())
		panic(err)
	}
	return conn
}

func main() {
	//fmt.Println("使用内部包测试：", route.Name())
	fmt.Println("cron task start.")

	c := ConnRedis()
	defer c.Close()

	key1m := "ee:task:t1s"

	//for i := 0; i < llen; i++ {
	count := 0
	for {
		count++
		_, err := c.Do("auth", "Maitui!)789")
		if err != nil {
			Info.Println("redis auth error: ", err.Error())
			time.Sleep(time.Duration(10) * time.Second)
			c = ConnRedis()
			continue
		}

		val, err := redis.String(c.Do("lpop", key1m))
		if err != nil {
			Info.Println("redis lpop error: ", err.Error(), count, " sleep 3 second")
			time.Sleep(time.Duration(3) * time.Second)
			continue
		}

		go AsynTaskRun(val)

		Info.Println("asyn task run: ", val, " sleep 1 second")

		time.Sleep(time.Duration(1) * time.Second)
	}

}

/*
values, err := redis.Values(c.Do("lrange", key1m, 0, llen))
if err != nil {
	fmt.Println("lrange err", err.Error())
	panic(err)
}
fmt.Println(key1m)
for _, v := range values {
	fmt.Printf(" %s \n", v.([]byte))
}
fmt.Println()
*/
