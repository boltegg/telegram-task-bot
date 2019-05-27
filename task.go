package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Task struct {
	ID          string `bson:"_id"`
	TaskID      int64
	Author      int64
	Assigned    int64
	OpenTime    time.Time
	DueTime     time.Time
	CloseTime   time.Time
	Subject     string
	Description string
	Status      string
	Project     string
}

func NewTask(author int64, assigned int64, subject string) *Task {
	return &Task{
		Author:   author,
		Assigned: assigned,
		OpenTime: time.Now(),
		DueTime:  time.Now(),
		Subject:  subject,
		Status:   "open",
	}
}

func (t *Task) Create() (i int64, err error) {
	mongo.Lock()
	defer mongo.Unlock()

	i, err = mongo.GetNextId(t.Assigned)
	if err != nil {
		return
	}

	t.TaskID = i
	t.ID = fmt.Sprintf("%d-%d", t.Assigned, i)
	t.Subject = e.EncryptString(t.Subject)
	t.Description = e.EncryptString(t.Description)

	err = mongo.ColTasks().Insert(t)
	return
}

func GetTask(assigned int64, taskId int64) (task *Task, err error) {

	err = mongo.ColTasks().FindId(fmt.Sprintf("%d-%d", assigned, taskId)).One(&task)
	if err != nil {
		return task, err
	}

	task.Subject = e.DecryptString(task.Subject)
	task.Description = e.DecryptString(task.Description)
	return task, err
}

func GetAllTask(assigned int64) (tasks []*Task, err error) {
	err = mongo.ColTasks().Find(bson.M{"assigned": assigned, "status": "open"}).All(&tasks)
	if err != nil {
		return
	}
	for _, task := range tasks {
		task.Subject = e.DecryptString(task.Subject)
		task.Description = e.DecryptString(task.Description)
	}
	return
}

func CloseTask(assigned int64, taskId int64) error {
	tid := fmt.Sprintf("%d-%d", assigned, taskId)
	err := mongo.ColTasks().UpdateId(tid, bson.M{"$set": bson.M{"status": "close", "close_time": time.Now()}})
	return err
}

func OpenTask(assigned int64, taskId int64) error {
	tid := fmt.Sprintf("%d-%d", assigned, taskId)
	err := mongo.ColTasks().UpdateId(tid, bson.M{"$set": bson.M{"status": "open"}})
	return err
}

//func NewMessage(m string) EncryptedText {
//
//	return EncryptedText("")
//}
//
//func (m *EncryptedText) Decode() string {
//
//	return ""
//}

//func init() {
//
//	enc := NewEncryptor("bimba45bvyfdtsghsdtfjsrdfgtjmfhjksrtujhtsrjftsdjfdtujrtyj")
//
//	e,err:=enc.EncryptString("pinokiodfg dfs hdsfh dffffffffffff")
//	fmt.Println(e, err)
//	d,err:=enc.DecryptString(e)
//	fmt.Println(d, err)
//
//
//	//enc := encrypt([]byte("petushok"), "sadgmjnrewiopnhgerpgher")
//	//encStr := base64.StdEncoding.EncodeToString(enc)
//	//fmt.Println("enc:", string(encStr))
//	//
//	//enc, _ = base64.StdEncoding.DecodeString(encStr)
//	//dec := decrypt(enc, "sadgmjnrewiopnhgerpgher")
//	//fmt.Println("dec:", string(dec))
//	//
//	//key2 := [32]byte{}
//	//_, err := io.ReadFull(strings.NewReader("yFV*&^$*^(*Ty807tf*)^&$R)87t780dfhdfhsdfherwtgh54g67247"), key2[:])
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//key22 := base64.StdEncoding.EncodeToString(key2[:])
//	//
//	//fmt.Println("mykey2: ", key2)
//	//
//
//
//
//	os.Exit(1)
//}
//
