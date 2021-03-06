package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/udistrital/utils_oas/time_bogota"
)

type Organizacion struct {
	Id                int               `orm:"column(id);pk;auto"`
	Nombre            string            `orm:"column(nombre)"`
	Ente              int               `orm:"column(ente)"`
	TipoOrganizacion  *TipoOrganizacion `orm:"column(tipo_organizacion);rel(fk)"`
	FechaCreacion     string            `orm:"column(fecha_creacion);null"`
	FechaModificacion string            `orm:"column(fecha_modificacion);null"`
}

func (t *Organizacion) TableName() string {
	return "organizacion"
}

func init() {
	orm.RegisterModel(new(Organizacion))
}

// AddOrganizacion insert a new Organizacion into database and returns
// last inserted Id on success.
func AddOrganizacion(m *Organizacion) (id int64, err error) {
	o := orm.NewOrm()
	o.Begin()
	var en = &Ente{Id: 0, TipoEnte: &TipoEnte{Id: 2}} //id del tipo ente para organizacion
	iden, err := o.Insert(en)
	if err == nil {
		m.FechaCreacion = time_bogota.TiempoBogotaFormato()
		m.FechaModificacion = time_bogota.TiempoBogotaFormato()
		m.Ente = int(iden)
		m.Id = int(iden)
		id, err = o.Insert(m)
		if err == nil {
			o.Commit()
			return
		}
	}
	o.Rollback()
	return
}

// GetOrganizacionById retrieves Organizacion by Id. Returns error if
// Id doesn't exist
func GetOrganizacionById(id int) (v *Organizacion, err error) {
	o := orm.NewOrm()
	v = &Organizacion{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllOrganizacion retrieves all Organizacion matches certain condition. Returns empty list if
// no records exist
func GetAllOrganizacion(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Organizacion)).RelatedSel(3)
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Organizacion
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateOrganizacion updates Organizacion by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrganizacionById(m *Organizacion) (err error) {
	o := orm.NewOrm()
	v := Organizacion{Id: m.Id}
	m.FechaModificacion = time_bogota.TiempoBogotaFormato()
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "Nombre", "Ente", "TipoOrganizacion", "FechaModificacion"); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrganizacion deletes Organizacion by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrganizacion(id int) (err error) {
	o := orm.NewOrm()
	v := Organizacion{Id: id}
	// ascertain id exists in the database
	o.Begin()
	if err = o.Read(&v); err == nil {
		if _, err = o.Delete(&Ente{Id: v.Ente}); err == nil {
			if _, err = o.Delete(&Organizacion{Id: id}); err == nil {
				o.Commit()
				return
			}
		}
	}
	o.Rollback()
	return
}
