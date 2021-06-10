package models

import (
	"alert/app/models/mongodb"
	"alert/app/parsers/alert"
	"bitbucket.org/pkg/inflect"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"time"
)

type AlertTeam struct {
	Id         *bson.ObjectId `bson:"_id,omitempty"`
	Team       string         `bson:"team"`
	Developers []Developers   `bson:"developers"`
	CreatedAt  time.Time      `bson:"created_at"`
	UpdatedAt  time.Time      `bson:"updated_at"`
}

type Developers struct {
	Name        string `bson:"name"`
	PhoneNumber string `bson:"phone_number"`
}

func Collection(m interface{}) (*mgo.Collection, *mongodb.Collection) {
	typ := reflect.TypeOf(m).Elem()
	n := typ.Name()

	var c string
	c = inflect.Tableize(n)
	newCol := mongodb.NewCollectionSession(c)
	return newCol.Session, newCol
}

func (alertTeam *AlertTeam) FindOrCreateBy(params map[string]interface{}) (*AlertTeam, error) {
	col, dbs := Collection(alertTeam)
	defer dbs.Close()
	err := col.Find(&params).One(alertTeam)
	if err != nil {
		alertTeam.Save()
		return alertTeam, err
	}
	return alertTeam, nil

}

func (alertTeam *AlertTeam) GetTeamByName(name string) (AlertTeam, error) {

	o := AlertTeam{}
	col, dbs := Collection(alertTeam)
	defer dbs.Close()
	err := col.Find(bson.M{"team": name}).One(&o)
	return o, err
}

func GetNewObjectId() *bson.ObjectId {
	new_id := bson.NewObjectId()
	return &new_id
}

func (alertTeam *AlertTeam) Save() (err error) {
	alertTeam.UpdatedAt = time.Now()
	col, dbs := Collection(alertTeam)
	defer dbs.Close()
	if alertTeam.Id == nil {
		alertTeam.CreatedAt = time.Now()
		alertTeam.Id = GetNewObjectId()
		err = col.Insert(alertTeam)
	} else {
		err = col.UpdateId(alertTeam.Id, alertTeam)
	}
	return err
}

func SetTeamToDb(teamInfo alert.CreateTeamRequest) error {

	alertTeam := &AlertTeam{}
	db_order, _ := alertTeam.FindOrCreateBy(map[string]interface{}{
		"name": teamInfo.Team.Name,
	})

	db_order.Team = teamInfo.Team.Name

	for _, each := range teamInfo.Developers {
		db_order.Developers = append(db_order.Developers, Developers{Name: each.Name, PhoneNumber: each.PhoneNumber})
	}

	db_order.Save()

	return nil

}
