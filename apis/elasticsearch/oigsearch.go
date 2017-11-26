package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gkewl/pulsecheck/logger"

	"github.com/gkewl/pulsecheck/common"
	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gkewl/pulsecheck/model"
)

// "github.com/gkewl/pulsecheck/common"
// eh "github.com/gkewl/pulsecheck/errorhandler"

// "github.com/gkewl/pulsecheck/constant"
// "github.com/gkewl/pulsecheck/model"

// // Tweet is a structure used for serializing/deserializing data in Elasticsearch.
// type Tweet struct {
// 	User     string                `json:"user"`
// 	Message  string                `json:"message"`
// 	Retweets int                   `json:"retweets"`
// 	Image    string                `json:"image,omitempty"`
// 	Created  time.Time             `json:"created,omitempty"`
// 	Tags     []string              `json:"tags,omitempty"`
// 	Location string                `json:"location,omitempty"`
// 	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
// }

// OIGmapping -
const OIGmapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"oigsearch":{
			"properties":{
				"firstname":{
					"type":"keyword"
				},
				"middlename":{
					"type":"keyword"				
				},
				"lastname":{
					"type":"keyword"
				},
				"dateofbirth":{
					"type":"keyword"
				}			
			}
		}
	}
}`

const oigIndexName = "testoig4" //"pulsecheck"

//func (bl BLElasticSearch) PrepareQuery(oig model.OIGSearch)

// SearchOIG -
func (bl BLElasticSearch) SearchOIG(reqCtx common.RequestContext, oig model.OIGSearch) (s []model.ElasticSearchResult, err error) {

	s = []model.ElasticSearchResult{}
	q := elastic.NewBoolQuery()
	// any change in search here has to change OIGSearch ToOIG() model
	q = q.Must(elastic.NewTermQuery("firstname", strings.ToLower(oig.Firstname)),
		elastic.NewTermQuery("lastname", strings.ToLower(oig.Lastname)),
		elastic.NewTermQuery("dateofbirth", oig.DateOfBirth),
	)
	res, err := reqCtx.AppContext().Ec.Search(oigIndexName).
		Index(oigIndexName).
		Query(q).
		Pretty(true).
		//	Sort("time", true).
		Do(context.Background())

	if err != nil {
		return []model.ElasticSearchResult{}, err
	}

	// Here's how you iterate through results with full control over each step.
	if res.Hits != nil {
		// Iterate through results
		for _, hit := range res.Hits.Hits {
			// Deserialize hit.Source
			var t model.OIGSearch
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				logger.LogError(fmt.Sprintf("Error in deserializing from elastic search %+v", err), reqCtx.Xid())
			}
			s = append(s, t.ToResult())
		}
	} else {
		// No hits
		fmt.Print("No hits found\n")
	}

	return
}
