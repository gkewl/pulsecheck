package search

import (
	"fmt"
	//"github.com/jmoiron/sqlx/reflectx"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"strings"
)

type SearchInterface interface {
	Search(string, string, *common.AppContext) ([]model.NameDescription, error)
}
type DBSearch struct {
	Username string
	Userid   int64
}

func (di DBSearch) Search(entity string, term string, ctx *common.AppContext) ([]model.NameDescription, error) {

	data := []model.NameDescription{}

	//try to avoid sql injectgion
	entities := strings.Split(entity, " ")
	if len(entities) > 1 {
		return []model.NameDescription{}, eh.NewError(eh.ErrSearchMultipleEntity, "")
	}
	entities = strings.Split(entity, ";")
	if len(entities) > 1 {
		return []model.NameDescription{}, eh.NewError(eh.ErrSearchMultipleEntity, "")
	}

	term = strings.TrimSpace(term)
	arr := strings.Split(term, " ")
	var newTerm string

	if len(arr) == 0 {
		return []model.NameDescription{}, nil
	}

	params := []string{"name", "description"}

	newTerm = GetSearchCondition(params, term)

	qry := fmt.Sprintf(`select id, name, description from %s where %s and isactive=1 order by name limit 50 `, entity, newTerm)

	stmt, err := ctx.Db.Preparex(qry)
	if err != nil {
		return []model.NameDescription{}, eh.NewError(eh.ErrSearchError, "DB Error: "+err.Error())

	}
	err = stmt.Select(&data)
	if err != nil {
		return []model.NameDescription{}, eh.NewError(eh.ErrSearchError, "DB Error: "+err.Error())
	}

	return data, err

}

func GetSearchCondition(params []string, term string) string {

	term = strings.TrimSpace(term)
	arr := strings.Split(term, " ")
	var newTerm string

	if len(arr) == 0 {
		return newTerm
	}
	for index, val := range arr {
		if len(strings.TrimSpace(val)) == 0 {
			continue
		}

		if index > 0 && len(val) > 0 {
			newTerm = newTerm + " and "
		}
		var condition string
		for idx, param := range params {
			if idx > 0 {
				condition = condition + " or "
			}
			condition = condition + "(  " + param + " like '" + strings.ToLower(val) + "%' or " + param + " like '% " + strings.ToLower(val) + "%'  )"
		}
		if len(condition) > 0 {
			newTerm = newTerm + "( " + condition + " ) "
		}
	}

	return newTerm
}
