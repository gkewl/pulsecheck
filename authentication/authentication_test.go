package authentication_test

import (
	"net/http/httptest"

	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/model"

	"github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/routehandler"
)

var (
	httpCaller protocol.HTTPTestCaller
	mqttCaller protocol.MQTTTestCaller
	server     *httptest.Server
	router     *routehandler.HTTPRouter
	baseUrl    string
	mqttClient *protocol.MQTTClient
)

var apis = common.APIRoutes{
	authentication.GetRoutes(),
}

func init() {
	config.LoadConfigurations()

	ctx := common.AppContext{}
	ctx.Db, _ = dbhandler.CreateConnection()
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenInfo, _ = authBackend.GenerateToken("rgunari@gmail.com", model.UserCompany{UserID: 1, CompanyID: 1}, "USER")
	router = routehandler.NewRouter(&ctx, apis, "/api/v1")
	server = httptest.NewServer(router) //Creating new server with the user handlers
	baseUrl = server.URL + "/api/v1"
	mqttClient = protocol.NewMQTTClient(&ctx, apis)
}

// func TestMachineToken123(t *testing.T) {

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/machinetoken/UnitTestUserName", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")

// 	request, err := http.NewRequest("POST", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if res.StatusCode != 200 {
// 		t.Errorf("Expected: 200  Actual : %d", res.StatusCode)
// 	}

// 	var output model.TokenInfo
// 	err = utilities.ScanResponseObject(res.Body, &output)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(output.Token) == 0 {
// 		t.Fatal("Token is empty")
// 	}

// 	if output.Exp != 0 {
// 		t.Error("Expiry has value and exepcted to be 0")
// 	}
// }

// func TestMachineTokenWithInvalidUser(t *testing.T) {

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/machinetoken/UnitTestUserNameNotFOund", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")

// 	request, err := http.NewRequest("POST", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if res.StatusCode != 404 {
// 		t.Errorf("Expected: 404  Actual : %d", res.StatusCode)
// 	}

// }

// func TestValidationOfToken(t *testing.T) {

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/validatetoken", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")

// 	request, err := http.NewRequest("GET", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 		fmt.Println(err)
// 	}

// 	if res.StatusCode != 200 {
// 		t.Errorf("Expected: 200  Actual : %d", res.StatusCode)
// 	}

// }

// func TestValidationOfTokenFailure(t *testing.T) {

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/validatetoken", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")

// 	request, err := http.NewRequest("GET", apiURL, reader)
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 		fmt.Println(err)
// 	}

// 	if res.StatusCode != 401 {
// 		t.Errorf("Expected: 401  Actual : %d", res.StatusCode)
// 	}

// }
// func TestTestAuthExForAdmin(t *testing.T) {
// 	authBackend := authentication.InitJWTAuthenticationBackend()
// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Admin")

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/validatetoken", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")
// 	//token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY3RvcmlkIjoyNDQwNiwiaWF0IjoxNDY0NzM2NzE3LCJzY29wZXMiOiJvcGVyYXRvciIsInN1YiI6IlVuaXRUZXN0TWFjaGluZV9kYWlPM09JPSJ9.E3skgYuAkIC-shoOaBnv4uy-YlsTUtdjCqEfsvPnYkkOAIbjAE6RWJw6FEj-akpEKx3-9JsmrxTgIj547_QVvTxJWxTM9xcGpSYGXdzMR8_slBUVUITQmhOSnGQq5wn7Rv973GVO3R8QaZlH3XuV_eLZFwKSsyy8roq8dpm-x39lMlBARa-QVDbZfhpCTBmqpRM0FRjT9jVCH3wlASIU0Am5YcWoUIESpvYVTG39KKd2HnQsdmfmNwUPbzM3F3EnQMVjsER6oNtyCbRehnmlA06UFBJygJ6T2rGr5uz5gbyIg9kpUKRrPYKEX7GcUIt_JGRoy6tMpo20fpcPhc-gbg"
// 	request, err := http.NewRequest("GET", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 		fmt.Println(err)
// 	}

// 	if res.StatusCode != 200 {
// 		t.Errorf("Expected: 200  Actual : %d   ", res.StatusCode)

// 	}
// }

// func TestTestAuthExForSuperuser(t *testing.T) {
// 	authBackend := authentication.InitJWTAuthenticationBackend()
// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Superuser")

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/validatetoken", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")
// 	//token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY3RvcmlkIjoyNDQwNiwiaWF0IjoxNDY0NzM2NzE3LCJzY29wZXMiOiJvcGVyYXRvciIsInN1YiI6IlVuaXRUZXN0TWFjaGluZV9kYWlPM09JPSJ9.E3skgYuAkIC-shoOaBnv4uy-YlsTUtdjCqEfsvPnYkkOAIbjAE6RWJw6FEj-akpEKx3-9JsmrxTgIj547_QVvTxJWxTM9xcGpSYGXdzMR8_slBUVUITQmhOSnGQq5wn7Rv973GVO3R8QaZlH3XuV_eLZFwKSsyy8roq8dpm-x39lMlBARa-QVDbZfhpCTBmqpRM0FRjT9jVCH3wlASIU0Am5YcWoUIESpvYVTG39KKd2HnQsdmfmNwUPbzM3F3EnQMVjsER6oNtyCbRehnmlA06UFBJygJ6T2rGr5uz5gbyIg9kpUKRrPYKEX7GcUIt_JGRoy6tMpo20fpcPhc-gbg"
// 	request, err := http.NewRequest("GET", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 		fmt.Println(err)
// 	}

// 	if res.StatusCode != 200 {
// 		t.Errorf("Expected: 200  Actual : %d   ", res.StatusCode)

// 	}

// }

// func TestTestAuthExForGuest_Failed(t *testing.T) {
// 	authBackend := authentication.InitJWTAuthenticationBackend()
// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Guest")

// 	apiURL := fmt.Sprintf("%s/api/v1/auth/validatetokenex", server.URL) //Grab the address for the API endpoint
// 	reader = strings.NewReader("")
// 	//token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY3RvcmlkIjoyNDQwNiwiaWF0IjoxNDY0NzM2NzE3LCJzY29wZXMiOiJvcGVyYXRvciIsInN1YiI6IlVuaXRUZXN0TWFjaGluZV9kYWlPM09JPSJ9.E3skgYuAkIC-shoOaBnv4uy-YlsTUtdjCqEfsvPnYkkOAIbjAE6RWJw6FEj-akpEKx3-9JsmrxTgIj547_QVvTxJWxTM9xcGpSYGXdzMR8_slBUVUITQmhOSnGQq5wn7Rv973GVO3R8QaZlH3XuV_eLZFwKSsyy8roq8dpm-x39lMlBARa-QVDbZfhpCTBmqpRM0FRjT9jVCH3wlASIU0Am5YcWoUIESpvYVTG39KKd2HnQsdmfmNwUPbzM3F3EnQMVjsER6oNtyCbRehnmlA06UFBJygJ6T2rGr5uz5gbyIg9kpUKRrPYKEX7GcUIt_JGRoy6tMpo20fpcPhc-gbg"
// 	request, err := http.NewRequest("GET", apiURL, reader)
// 	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenInfo.Token))
// 	res, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		t.Error(err)
// 		fmt.Println(err)
// 	}

// 	if res.StatusCode != 401 {
// 		t.Errorf("Expected: 401  Actual : %d   ", res.StatusCode)

// 	}

// }

// func testGenerateTokens(t *testing.T) {
// 	authBackend := authentication.InitJWTAuthenticationBackend()

// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Admin")
// 	fmt.Println("Admin Token: ", tokenInfo)

// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Superuser")
// 	fmt.Println("Superuser Token: ", tokenInfo)

// 	tokenInfo, _ = authBackend.GenerateToken("test", 13066, "USER", "Guest")
// 	fmt.Println("Guest Token: ", tokenInfo)

// }
