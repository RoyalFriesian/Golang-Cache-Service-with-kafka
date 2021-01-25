package main

import (
	"encoding/json"
	"context"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"gopkg.in/mgo.v2"

	"github.com/segmentio/kafka-go"
)

const (
	hosts      = "localhost:9999"
	database   = "cachedb"
	username   = "myTester"
	password   = "mytesterpass"
	collection = "data"

	hosts_bk      = "localhost:9999"
	database_bk   = "backupdb"
	username_bk   = "myTester2"
	password_bk   = "mytester2pass"
	collection_bk = "data"

	kafkaBroker	  = "localhost:9092"
)

type Data struct {
	Key       string
	Data      string
	Timestamp time.Time
}

type Response struct {
	StatusCode int
	Msg        string
}

func main() {
	go StartKafka()
	//pagination()
	//http.HandleFunc("/reload", reload)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/data/page/", pagination)
	//http.HandleFunc("/", HelloServer)
	fmt.Println("Server is running!")

	http.ListenAndServe(":8080", nil)

}

func pagination(w http.ResponseWriter, r *http.Request) {
	var page int = 1
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := r.ParseForm(); err != nil {
			page , _ = strconv.Atoi(r.FormValue("page"))
		}
		info := &mgo.DialInfo{
			Addrs:    []string{hosts},
			Timeout:  60 * time.Second,
			Database: database,
			Username: username,
			Password: password,
		}
		session, err1 := mgo.DialWithInfo(info)
		if err1 != nil {
			panic(err1)
		}
		session.SetSafe(&mgo.Safe{})
		c := bootstrap(session)
		var limit int = 20
		//var page int = 3
		var data []Data
		q := c.Find(nil).Sort("timestamp")
		q = q.Skip((page - 1) * limit).Limit(limit)
		err := q.All(&data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Norm Find Data: %+v\n", data)
		json.NewEncoder(w).Encode(data)
	} else {
		fmt.Printf("[{\"status\" : \"Wrong Request!\"}])
	}
}

func StartKafka() {
	conf := kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    "myTopic",
		GroupID:  "g1",
		MaxBytes: 10,
	}
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Some error occured", err)
			continue
		}
		if string(m.Value) == "reload" {
			reload()
			
		}
		fmt.Println("Message is ", string(m.Value))
	}
}

func reload() {
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: "cachedb",
		Username: username,
		Password: password,
	}

	session, err1 := mgo.DialWithInfo(info)
	if err1 != nil {
		panic(err1)
	}
	session.SetSafe(&mgo.Safe{})
	c := bootstrap(session)

	info_bk := &mgo.DialInfo{
		Addrs:    []string{hosts_bk},
		Timeout:  60 * time.Second,
		Database: database_bk,
		Username: username_bk,
		Password: password_bk,
	}

	session_bk, err_bk := mgo.DialWithInfo(info_bk)

	if err_bk != nil {
		panic(err1)
	}
	var results []Data
	err := c.Find(nil).All(&results)
	if err != nil {
		panic(err)
	}
	c2 := bootstrapCleanDB(session_bk)
	for i := range results {
		doc := results[i]
		err = c2.Insert(doc)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Database Re-loaded")
	//w.Write([]byte("[{'status':200, 'msg':'SUCCESS'}]"))
}

func Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("No content found"))
			return
		}
		w.WriteHeader(http.StatusOK)
		key := r.FormValue("key")
		data := r.FormValue("data")
		info := &mgo.DialInfo{
			Addrs:    []string{hosts},
			Timeout:  60 * time.Second,
			Database: database,
			Username: username,
			Password: password,
		}
		session, err1 := mgo.DialWithInfo(info)
		if err1 != nil {
			panic(err1)
		}
		session.SetSafe(&mgo.Safe{})
		c := bootstrap(session)
		err := c.Insert(
			&Data{Key: key, Data: data, Timestamp: time.Now()},
		)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("\n\n")
		w.Write([]byte("[{\"status\":200, \"msg\":\"success\"}]"))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func bootstrap(s *mgo.Session) *mgo.Collection {
	c := s.DB(database).C("data")
	index := mgo.Index{
		Key:        []string{"key"},
		Unique:     true,
		Background: true,
	}

	c.EnsureIndex(index)
	return c
}
			   
//At the time of reload cache from backup database we first need to clean the DB. 
func bootstrapCleanDB(s *mgo.Session) *mgo.Collection {
	s.DB(database).DropDatabase()
	c := s.DB(database).C("backup")
	index := mgo.Index{
		Key:        []string{"key"},
		Unique:     true,
		Background: true,
	}

	c.EnsureIndex(index)

	return c
}

