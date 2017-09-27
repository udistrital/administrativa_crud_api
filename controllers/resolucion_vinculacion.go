package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/administrativa_crud_api/models"
)

type ResolucionVinculacionController struct {
	beego.Controller
}

func (c *ResolucionVinculacionController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)

}

func (c *ResolucionVinculacionController) GetAll() {
	listaResoluciones := models.GetAllResolucionVinculacion()
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = listaResoluciones
	c.ServeJSON()
}