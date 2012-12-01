package webservice

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/landskape/application"
	"github.com/emicklei/landskape/model"
	"net/http"
)

const (
	NO_UPDATE = false
)

type ApplicationService struct {
	restful.WebService
}

func NewApplicationService() *ApplicationService {
	ws := new(ApplicationService)
	ws.Path("/applications").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)
	ws.Route(ws.GET("").To(GetAllApplications))
	ws.Route(ws.GET("/{id}").To(GetApplication))
	ws.Route(ws.PUT("/{id}").To(PutApplication))
	ws.Route(ws.POST("/{id}").To(PostApplication))
	return ws
}

func GetApplication(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	app, err := application.SharedLogic.GetApplication(id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(app)
}

func GetAllApplications(req *restful.Request, resp *restful.Response) {
	apps, err := application.SharedLogic.AllApplications()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(apps)
}

func PostApplication(req *restful.Request, resp *restful.Response) {
	app := new(model.Application)
	err := req.ReadEntity(&app)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if app.Id == "" {
		err := restful.NewError(model.MISMATCH_ID, "Id is missing")
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	_, err = application.SharedLogic.SaveApplication(app)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

func PutApplication(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	app := new(model.Application)
	err := req.ReadEntity(&app)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if app.Id != id {
		err := restful.NewError(model.MISMATCH_ID, fmt.Sprintf("Id mismatch: %v != %v", app.Id, id))
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if application.SharedLogic.ExistsApplication(id) {
		err := restful.NewError(http.StatusConflict, "Application already exists:"+id)
		resp.WriteError(http.StatusConflict, err)
		return
	}
	_, err = application.SharedLogic.SaveApplication(app)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}
