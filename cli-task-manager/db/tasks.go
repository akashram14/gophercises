package db

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key           int
	Value         string
	Completed     bool
	TimeCompleted time.Time
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})

}

func CreateTask(taskvalue string) error {
	task := Task{
		Value:     taskvalue,
		Completed: false,
	}
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id := int(id64)
		key := itob(id)

		task.Key = id

		if buf, err := json.Marshal(task); err != nil {
			return err
		} else if err := b.Put(key, buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			if !task.Completed {
				tasks = append(tasks, task)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func AllCompletedTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		var task Task
		for k, v := c.First(); k != nil; k, v = c.Next() {

			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			if task.Completed &&
				!task.TimeCompleted.IsZero() &&
				time.Since(task.TimeCompleted).Hours() <= 12 {
				tasks = append(tasks, task)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func CompleteTask(key int) error {
	var task Task
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		v := b.Get(itob(key))

		if err := json.Unmarshal(v, &task); err != nil {
			return err
		}

		task.Completed = true
		task.TimeCompleted = time.Now()

		if buf, err := json.Marshal(task); err != nil {
			return err
		} else if err := b.Put(itob(task.Key), buf); err != nil {
			return err
		}
		return nil
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
