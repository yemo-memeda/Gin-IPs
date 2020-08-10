package dao

import (
	"Gin-IPs/src/configure"
	"Gin-IPs/src/utils/database/mongodb"
	"Gin-IPs/src/utils/log"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Logger  *logrus.Logger
	MgoPool *mongo.Client
	MgoDb   string
}

var ModelClient = new(Models)

func Init() error {
	if logger, err := mylog.New(
		configure.GinConfigValue.Log.Path, configure.GinConfigValue.Log.Name,
		configure.GinConfigValue.Log.Level, nil, configure.GinConfigValue.Log.Count); err != nil {
		return err
	} else {
		ModelClient.Logger = logger
	}
	var err error
	ModelClient.MgoPool, err = mongodb.CreatePool(configure.GinConfigValue.Mgo.Uri, configure.GinConfigValue.Mgo.PoolSize)
	if err != nil {
		ModelClient.Logger.Errorf("Collection Client Pool With Uri %s Create Failed: %s", configure.GinConfigValue.Mgo.Uri, err)
		return err
	}
	ModelClient.MgoDb = configure.GinConfigValue.Mgo.Database
	ModelClient.Logger.Infof("Collection Client Pool Created successful With Uri %s", configure.GinConfigValue.Mgo.Uri)
	ModelClient.Logger.Infof("Models Created Success")
	return nil
}

func (m *Models) LogMongo() {
	for log := range mongodb.MongoLogChannel {
		logField := map[string]interface{}{
			"Database":   log.Database,
			"Collection": log.Collection,
			"Action":     log.Action,
			"Result":     log.Result,
		}
		switch log.Documents.(type) {
		case map[string]interface{}:
			docBytes, _ := json.Marshal(log.Documents)
			logField["Documents"] = string(docBytes)
		default:
			logField["Documents"] = log.Documents
		}
		if log.Ok {
			m.Logger.WithFields(logField).Info("")
		} else {
			m.Logger.WithFields(logField).Error(log.ErrMsg)
		}
	}
}

func Start() {
	go ModelClient.LogMongo()
	// MockTest()  // 插入初始化数据
}
