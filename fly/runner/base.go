package runner

import "go.mongodb.org/mongo-driver/bson"

type TaskBase struct {
	Collection string
}

type BM = bson.M
type BD = bson.D
