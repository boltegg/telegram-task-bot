package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

const (
	DB_NAME            = "fasttask"
	COLLECTION_LAST_ID = "last_id"
	COLLECTION_TASKS   = "tasks"
)

type Mongo struct {
	db *mgo.Session
	sync.Mutex
}

func NewMongo() Mongo {
	var err error
	mongo.db, err = mgo.Dial(config.MongoDbAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	mongo.db.SetMode(mgo.Monotonic, true)
	return mongo
}

func (m *Mongo) ColLastId() *mgo.Collection {
	return m.db.Copy().DB(DB_NAME).C(COLLECTION_LAST_ID)
}

func (m *Mongo) ColTasks() *mgo.Collection {
	return m.db.Copy().DB(DB_NAME).C(COLLECTION_TASKS)
}

func (m *Mongo) GetNextId(clientId int64) (int64, error) {

	var v struct {
		LastId int64 `bson:"last_id"`
	}
	err := m.ColLastId().FindId(clientId).One(&v)
	if err != nil && err.Error() != "not found" {
		return 0, fmt.Errorf("Error find id in common_vars: %v", err)
	}

	v.LastId++
	_, err = m.ColLastId().UpsertId(clientId, bson.M{"$set": bson.M{"last_id": v.LastId}})
	if err != nil {
		return 0, fmt.Errorf("Error update id in common_vars: %v", err)
	}
	return v.LastId, nil
}

//	insert record
//	c := session.DB("test").C("people")
//	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
//		&Person{"Cla", "+55 53 8402 8510"})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	get record
//	result := Person{}
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	if err != nil {
//		log.Fatal(err)
//	}
