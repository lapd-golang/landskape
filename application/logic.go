package application

import (
	"errors"
	"github.com/emicklei/landskape/dao"
	"github.com/emicklei/landskape/model"
	"log"
	"time"
)

var SharedLogic Logic

type Logic struct {
	SystemDao     dao.SystemDao
	ConnectionDao dao.ConnectionDao
}

func (self Logic) AllSystems() (model.Systems, error) {
	apps, err := self.SystemDao.FindAll()
	if err != nil {
		return model.Systems{}, err
	}
	return model.Systems{apps}, nil
}

func (self Logic) AllConnections(filter model.ConnectionsFilter) (model.Connections, error) {
	cons, err := self.ConnectionDao.FindAllMatching(filter)
	if err != nil {
		return model.Connections{}, err
	}
	return model.Connections{cons}, nil
}

func (self Logic) DeleteConnection(con model.Connection) error {
	return self.ConnectionDao.Remove(con)
}

func (self Logic) SaveConnection(con model.Connection) error {
	log.Printf("logic.save:%#v", con)
	// Check from and to for existence
	if con.From == "" || !self.ExistsSystem(con.From) {
		return errors.New("Invalid from (empty or non-exist):" + con.From)
	}
	if con.To == "" || !self.ExistsSystem(con.To) {
		return errors.New("Invalid to (empty or non-exist):" + con.To)
	}
	if con.Type == "" {
		return errors.New("Invalid type (empty)")
	}
	return self.ConnectionDao.Save(con)
}

func (self Logic) GetSystem(id string) (model.System, error) {
	return self.SystemDao.FindById(id)
}

func (self Logic) DeleteSystem(id string) error {
	// TODO remove all its connections
	return self.SystemDao.RemoveById(id)
}

func (self Logic) ExistsSystem(id string) bool {
	return false
	//	result, _ := self.SystemDao.Exists(id)
	//	return result.Id == id
}

func (self Logic) SaveSystem(app *model.System) (*model.System, error) {
	app.Modified = time.Now()
	return app, self.SystemDao.Save(app)
}
