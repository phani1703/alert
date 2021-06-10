package mongodb

import (
	// "github.com/revel/revel"
	"gopkg.in/mgo.v2"
)

type Database struct {
	s       *mgo.Session
	name    string
	session *mgo.Database
}

var modes map[string]mgo.Mode = map[string]mgo.Mode{"Monotonic": mgo.Monotonic, "Strong": mgo.Strong, "Eventual": mgo.Eventual}

func (db *Database) Connect() {

	// config_mode := revel.Config.StringDefault("mongo.mode", "")
	db.s = service.Session()
	// db.s.SetMode(modes[config_mode], true)
	session := *db.s.DB(db.name)
	db.session = &session

}

func newDBSession(name string) *Database {

	var db = Database{
		name: name,
	}
	db.Connect()
	return &db
}
