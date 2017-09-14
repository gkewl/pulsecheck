package actor

import (
	"fmt"
	"github.com/gkewl/pulsecheck/apis/search"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
	"strings"
)

const (
	getQuery = `select a.id , a.name ,a.description, a.role,  a.adusername, a.email ,a.type ,a.ipconfig ,  a.isactive, a.macaddress ,
							a.created as created , a.modified as  modified, a1.id  as "manager.id" , a1.name  as "manager.name" ,
							a1.description  as "manager.description",
							a2.id  as "createdby.id",a2.name  as "createdby.name", a2.description  as "createdby.description",
							a3.id  as "modifiedby.id",a3.name  as "modifiedby.name", a3.description  as "modifiedby.description",		            	
			            	 a.lastlogintime
							from actor a
							left join actor a1 on a1.id = a.parentid
							left join actor a2 on a2.id = a.createdby
							left join actor a3 on a3.id = a.modifiedby
												
														`
)

// DLGet retreives actor with the given name
func DLGet(reqCtx common.RequestContext, name string) (actor model.Actor, err error) {
	query := getQuery + ` where a.name=? `
	actor = model.Actor{}

	err = reqCtx.Tx().Get(&actor, query, name)

	return actor, err
}

//DLGetAll retrieves all actor
func DLGetAll(reqCtx common.RequestContext, userTypes string) (actors []model.Actor, err error) {
	query := getQuery

	if len(userTypes) > 0 {
		arr := strings.Split(userTypes, ",")
		var types string
		for _, value := range arr {
			if len(strings.TrimSpace(value)) > 0 {
				if len(types) > 0 {
					types = types + ","
				}
				types = types + "'" + strings.TrimSpace(value) + "'"
			}
		}
		if len(types) > 0 {
			query = query + " where a.type in ( " + types + " )"
		}
	}

	data := []model.Actor{}

	err = reqCtx.Tx().Select(&data, query)

	return data, err
}

// DLCreate creates an actor
func DLCreate(reqCtx common.RequestContext, actor model.Actor) (model.Actor, error) {
	params := map[string]interface{}{
		"name":       actor.Name,
		"desc":       actor.Description,
		"email":      actor.EmailAddress,
		"type":       actor.Type,
		"ipaddress":  actor.IPAddress,
		"macaddress": actor.MACAddress,
		"manager":    actor.Manager.Id,
		"role":       actor.Role,
		"userid":     reqCtx.UserID(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into actor (name, description,  email, type ,ipconfig ,  macaddress , isactive, parentid, createdby,modifiedby,role)
	 		values (:name, :desc,  :email, :type, :ipaddress, :macaddress, 1, :manager, :userid, :userid, :role)`,
		params)
	if err == nil {
		actor.Id, _ = result.LastInsertId()
	}
	return actor, err
}

// DLUpdate updates some fields on a current step
func DLUpdate(reqCtx common.RequestContext, ID int64, actor model.Actor) error {
	params := map[string]interface{}{
		"id":         actor.Id,
		"name":       actor.Name,
		"desc":       actor.Description,
		"email":      actor.EmailAddress,
		"ipaddress":  actor.IPAddress,
		"macaddress": actor.MACAddress,
		"manager":    actor.Manager.Name,
		"role":       actor.Role,
		"type":       actor.Type,
		"isactive":   actor.IsActive,
		"userid":     reqCtx.UserID(),
	}

	_, err := reqCtx.Tx().NamedExec(
		`update actor a
						 left join actor a1 on a1.name= :manager
									set a.name = :name , a.description= :desc, a.email= :email , a.type= :type , a.ipconfig= :ipaddress ,
							 a.macaddress= :macaddress , a.parentid = a1.id  , a.modifiedby=:userid , a.role= :role, a.isactive =:isactive,
							 map.modifiedby=:userid
						 where a.id =:id`, params)

	return err
}

// DLUpdateLastlogin updates last login time to current time
func DLUpdateLastlogin(reqCtx common.RequestContext, ID int64) error {
	params := map[string]interface{}{
		"id":            ID,
		"lastlogintime": utilities.Now(),
		"userid":        reqCtx.UserID(),
	}

	_, err := reqCtx.Tx().NamedExec(
		`update actor 						 
			set lastlogintime=:lastlogintime, modifiedby=:userid
						 where id =:id`, params)

	return err
}

// DLMarkInactive set the isactive flag to zero for this actor
func DLMarkInactive(reqCtx common.RequestContext, actorID int64) (affect int64, err error) {
	result, err := reqCtx.Tx().Exec("update actor set isactive = 0, modifiedby = ? where id = ?",
		reqCtx.UserID(), actorID)
	affect, err = result.RowsAffected()

	return affect, err
}

//DLGetActors gets the actor by the term specified for a give type
func DLGetActors(reqCtx common.RequestContext, userType string, term string, limit int) (actors []model.ActorSearchResponse, err error) {
	data := []model.ActorSearchResponse{}

	params := []string{"a.name", "a.description"}
	searchCondition := search.GetSearchCondition(params, term)

	qry := fmt.Sprintf(` select a.id,a.name, a.adusername, a.email  as email , a.description, a.role
			             from actor  a
			             
			             where a.type=? and %s limit %d`, searchCondition, limit)

	stmt, err := reqCtx.Tx().Preparex(qry)

	if err != nil {
		return data, err

	}

	err = stmt.Select(&data, strings.ToLower(userType))

	return data, err

}
