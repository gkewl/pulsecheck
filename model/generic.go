package model

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

//PagingParameters provides simple paging attributes
type PagingParameters struct {
	PageNumber int64 `json:"pagenumber"`
	PageSize   int64 `json:"pagesize"`
}

//SortingParameters allows to sort by one column
//Implementor needs to validate
type SortingParameters struct {
	//asc or desc should be the only acceptable values
	SortOrder string `json:"sortorder"`
	//Name of the property/column to use for sorting
	//Eg: name.  This basic implementation supports only one column
	SortBy string `json:"sortby"`
}

//DateRangeFilter provides a from-to struct to perform filtering
type DateRangeFilter struct {
	From null.Time `json:"from"`
	To   null.Time `json:"to"`
}

//PagedSearchResponse represents the response of executing a search
type PagedSearchResponse struct {
	TotalRowCount int64       `json:"totalrowcount"`
	TotalPages    int64       `json:"totalpages"`
	PageNumber    int64       `json:"pagenumber"`
	PageSize      int64       `json:"pagesize"`
	Items         interface{} `json:"items"`
}

//NewPagingParameters returns a new default paging parameters object
func NewPagingParameters() *PagingParameters {
	return &PagingParameters{
		PageNumber: 1,
		PageSize:   100,
	}
}

//NewSortingParameters returns a new empty sorting parameters object
func NewSortingParameters() *SortingParameters {
	return &SortingParameters{
		SortBy:    "",
		SortOrder: "",
	}
}

//NewPagedResponse returns a new response with params
func NewPagedResponse(items interface{}, totalRowCount int64, pagingParams PagingParameters) PagedSearchResponse {
	return PagedSearchResponse{
		TotalRowCount: totalRowCount,
		TotalPages:    int64(math.Ceil(float64(totalRowCount) / float64(pagingParams.PageSize))),
		PageNumber:    pagingParams.PageNumber,
		PageSize:      pagingParams.PageSize,
		Items:         items,
	}
}

//Validate checks sorting parameters are valid
func (sortParam SortingParameters) Validate(validCols []string) error {
	errs := []string{}
	lowered := strings.ToLower(sortParam.SortOrder)

	if sortParam.SortBy == "" && sortParam.SortOrder == "" {
		return nil
	}

	if lowered != "asc" && lowered != "desc" {
		errs = append(errs, fmt.Sprintf("sortorder %s is invalid, only 'asc' or 'desc' is allowed", sortParam.SortOrder))
	}

	lowered = strings.ToLower(sortParam.SortBy)
	valid := false
	for _, colName := range validCols {
		if strings.EqualFold(lowered, colName) {
			valid = true
		}
	}

	if !valid {
		errs = append(errs, fmt.Sprintf("sortby %s is an invalid sortable property", sortParam.SortBy))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, "; "))
}

//Validate checks paging parameters are valid
func (pageParam PagingParameters) Validate(limit int64) error {
	errs := []string{}
	if pageParam.PageNumber <= 0 {
		errs = append(errs, "pagenumber must be greater than zero")
	}
	if pageParam.PageSize <= 0 {
		errs = append(errs, "pagesize must be greater than zero")
	}
	if pageParam.PageSize > limit {
		errs = append(errs, fmt.Sprintf("pagesize cannot be greater than %d", limit))
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.New(strings.Join(errs, "; "))
}

type SuccessResult struct {
	Success string `json:"success"`
}

type NameDescription struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
}

type NameDescriptionEdit struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	EditFlag    string      `json:"editflag"`
}

type NullableNameDescription struct {
	Id          null.Int    `db:"id" json:"id"`
	Name        null.String `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
}

// BuildNND returns a nullable name description with id field filled in
func BuildNND(id int64) NullableNameDescription {
	return NullableNameDescription{Id: null.IntFrom(id)}
}

type BasicModel struct {
	Id         int64           `db:"id" json:"id"`
	IsActive   bool            `db:"isactive" json:"isactive"`
	Createdby  NameDescription `db:"createdby" json:"createdby"`
	Created    time.Time       `db:"created" json:"created"`
	Modifiedby NameDescription `db:"modifiedby" json:"modifiedby"`
	Modified   time.Time       `db:"modified" json:"modified"`
}

type NameDescriptionPart struct {
	Id          int64  `db:"id" json:"id"`
	PartNumber  string `db:"partnumber" json:"partnumber"`
	Description string `db:"description" json:"description"`
}

type NameDescriptionArtifact struct {
	Id          null.Int    `db:"id" json:"id"`
	Name        null.String `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	ContentType null.String `db:"contenttype" json:"contenttype"`
	TargetPath  null.String `db:"targetpath" json:"targetpath"`
}

type NullableNameDescriptionPart struct {
	Id          null.Int    `db:"id" json:"id"`
	PartNumber  null.String `db:"partnumber" json:"partnumber"`
	Description null.String `db:"description" json:"description"`
}

// MakeNDP converts the nullable struct to a real one
func (nndp NullableNameDescriptionPart) MakeNDP() NameDescriptionPart {
	return NameDescriptionPart{
		Id:          nndp.Id.Int64,
		PartNumber:  nndp.PartNumber.String,
		Description: nndp.Description.String,
	}
}

// BuildNNDP returns a nullable name description part with id field filled in
func BuildNNDP(id int64) NullableNameDescriptionPart {
	return NullableNameDescriptionPart{Id: null.IntFrom(id)}
}

type NameDescriptionOptionCode struct {
	Id                     int64       `db:"id" json:"id"`
	Code                   string      `db:"code" json:"code"`
	Description            null.String `db:"description" json:"description"`
	OptionGroupID          int64       `db:"optiongroupid" json:"optiongroupid"`
	OptionGroupName        string      `db:"optiongroupname" json:"optiongroupname"`
	OptionGroupDescription null.String `db:"optiongroupdescription" json:"optiongroupdescription"`
}

//IDNameKey is an interface to objects that can be keyed by ID or name
type IDNameKey interface {
	GetID() int64
	GetName() string
	GetKey() string
}
