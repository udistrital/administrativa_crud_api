package controllers

import (
	"github.com/udistrital/administrativa_crud_api/models"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"fmt"
	"github.com/astaxie/beego"
)

// VinculacionDocenteController oprations for VinculacionDocente
type VinculacionDocenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *VinculacionDocenteController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("InsertarVinculaciones", c.InsertarVinculaciones)
}

// GetTotalContratosXResolucion ...
// @Title GetTotalContratosXResolucion
// @Description Retorna el valor total de la contratación para la resolución
// @Param id_resolucion query string false "nomina a listar"
// @Success 201 {object} int
// @Failure 403 body is empty
// @router /get_total_contratos_x_resolucion/:id_resolucion [get]
func (c *VinculacionDocenteController) GetTotalContratosXResolucion() {
	idStr := c.Ctx.Input.Param(":id_resolucion")
	v, err := models.GetTotalContratosXResolucion(idStr)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Post ...
// @Title Post
// @Description create VinculacionDocente
// @Success 201 {int}
// @Failure 403 body is empty
// @router /InsertarVinculaciones [post]
func (c *VinculacionDocenteController) InsertarVinculaciones() {
	var v []models.VinculacionDocente
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if id,err := models.AddConjuntoVinculaciones(v); err == nil {

			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = id
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Post ...
// @Title Post
// @Description create VinculacionDocente
// @Param	body		body 	models.VinculacionDocente	true		"body for VinculacionDocente content"
// @Success 201 {int} models.VinculacionDocente
// @Failure 403 body is empty
// @router / [post]
func (c *VinculacionDocenteController) Post() {
	var v models.VinculacionDocente
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddVinculacionDocente(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {

		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get VinculacionDocente by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.VinculacionDocente
// @Failure 403 :id is empty
// @router /:id [get]
func (c *VinculacionDocenteController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetVinculacionDocenteById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get VinculacionDocente
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.VinculacionDocente
// @Failure 403
// @router / [get]
func (c *VinculacionDocenteController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllVinculacionDocente(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the VinculacionDocente
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.VinculacionDocente	true		"body for VinculacionDocente content"
// @Success 200 {object} models.VinculacionDocente
// @Failure 403 :id is not int
// @router /:id [put]
func (c *VinculacionDocenteController) Put() {
	fmt.Println("edicion")
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.VinculacionDocente{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateVinculacionDocenteById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		fmt.Println("rro")
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the VinculacionDocente
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *VinculacionDocenteController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteVinculacionDocente(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetVinculacionesAgrupadas ...
// @Title GetVinculacionesAgrupadas
// @Description get vinculaciones agrupadas por docente
// @Param	id_resolucion		path 	string	true
// @Success 200 {object} models.VinculacionDocente
// @Failure 403 :id is empty
// @router /get_vinculaciones_agrupadas/:id_resolucion [get]
func (c *VinculacionDocenteController) GetVinculacionesAgrupadas() {
	idStr := c.Ctx.Input.Param(":id_resolucion")
	v, err := models.GetVinculacionesAgrupadas(idStr)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetValoresTotalesPorDisponibilidad ...
// @Title GetValoresTotalesPorDisponibilidad
// @Description se obtiene el valor total que ha sido apartado para cierta disponibilidad
// @Param	anio	    path 	string	true
// @Param	periodo		path 	string	true
// @Param	id_disponibilidad	path 	string	true
// @Success 200 {int}
// @Failure 403
// @router /get_valores_totales_x_disponibilidad/:anio/:periodo/:id_disponibilidad [get]
func (c *VinculacionDocenteController) GetValoresTotalesPorDisponibilidad() {
	anio := c.Ctx.Input.Param(":anio")
	periodo := c.Ctx.Input.Param(":periodo")
	id_disponibilidad := c.Ctx.Input.Param(":id_disponibilidad")
	fmt.Println("asdf",anio, periodo, id_disponibilidad)
	v, err := models.GetValoresTotalesPorDisponibilidad(anio, periodo, id_disponibilidad)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}
