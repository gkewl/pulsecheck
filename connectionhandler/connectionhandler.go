package connectionhandler

import (
	"crypto/tls"
	//"encoding/json"
	"os"
	"fmt"
	"net/http"
	//"reflect"
	//"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

func CreateElasticConnection() *elastic.Client {
	// Starting with elastic.v5, you must pass a context to execute each service
	//ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.

	c := &http.Client{
		Transport: &http.Transport{
			// TLSClientConfig: &tls.Config{},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	client, err := elastic.NewClient(
		elastic.SetHttpClient(c),
		elastic.SetSniff(false), //this is needed attribute to connect to elastic cloud
		elastic.SetURL(os.Getenv("ELASTIC_URL")),
		elastic.SetBasicAuth(os.Getenv("ELASTIC_UNAME"), os.Getenv("ELASTIC_PWD"))
		//elastic.SetURL("https://aade6f5cd32cedd31ee3a3c61384275f.us-central1.gcp.cloud.es.io:9243"),
		//elastic.SetBasicAuth("elastic", "zcAIW0nyX6AOeBQNGEwFPXaA"))

	//client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		fmt.Println("failed to connect to Elastic")
		panic(err)
	}
	fmt.Println("connected to elastic")
	return client
}
