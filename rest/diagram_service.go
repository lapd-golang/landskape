package rest

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/landskape/application"
	"github.com/emicklei/landskape/model"
)

var DotConfig = map[string]string{}

type DiagramResource struct {
	service application.Logic
}

func NewDiagramService(s application.Logic) *restful.WebService {
	ws := new(restful.WebService)
	d := DiagramResource{service: s}
	tags := []string{"diagrams"}

	ws.Path("/v1/diagrams").
		Produces("text/plain")
	ws.Route(ws.GET("/").To(d.computeDiagram).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Doc(`Compute a graphical diagram with all (filtered) connections for all systems and the given scope`).
		Param(ws.QueryParameter("from", "comma separated list of system ids")).
		Param(ws.QueryParameter("to", "comma separated list of system ids")).
		Param(ws.QueryParameter("type", "comma separated list of known connection types")).
		Param(ws.QueryParameter("center", "comma separated list of system ids")).
		Param(ws.QueryParameter("cluster", "show clusters based on the values of the give system attribute")).
		Param(ws.QueryParameter("format", "svg (default), png")))
	return ws
}

func (d DiagramResource) computeDiagram(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()
	filter := model.ConnectionsFilter{
		Froms:   asFilterParameter(req.QueryParameter("from")),
		Tos:     asFilterParameter(req.QueryParameter("to")),
		Types:   asFilterParameter(req.QueryParameter("type")),
		Centers: asFilterParameter(req.QueryParameter("center"))}
	dotOnly := req.QueryParameter("format") == "dot"
	connections, err := d.service.AllConnections(ctx, filter)
	if err != nil {
		log.Printf("AllConnections failed:%#v", err)
		resp.WriteError(500, err)
		return
	}
	format := req.QueryParameter("format")
	if "" == format {
		format = "svg"
	}
	id, err := model.GenerateUUID()
	if err != nil {
		log.Printf("GenerateUUID failed:%v", err)
		resp.WriteError(500, err)
		return
	}
	input := fmt.Sprintf("%v/%v.dot", DotConfig["tmp"], id)
	output := fmt.Sprintf("%v/%v.%v", DotConfig["tmp"], id, format)

	dotBuilder := application.NewDotBuilder()
	dotBuilder.ClusterBy(req.QueryParameter("cluster"))
	dotBuilder.BuildFromAll(connections)

	if dotOnly {
		resp.AddHeader("Content-Type", "text/plain")
		dotBuilder.WriteDot(resp)
		return
	}
	dotBuilder.WriteDotFile(input)

	cmd := exec.Command(DotConfig["binpath"],
		fmt.Sprintf("-T%v", format),
		fmt.Sprintf("-o%v", output),
		input)
	err = cmd.Start()
	if err != nil {
		log.Printf("Dot command start failed:%v", err)
		resp.WriteError(500, err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Dot did not complete:%v", err)
		resp.WriteError(500, err)
		return
	}
	// resp.AddHeader("Content-Type", "image/svg+xml")
	http.ServeFile(resp, req.Request, output)
}
