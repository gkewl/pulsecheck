package parameters

import ()

type TokenAuthentication struct {
	Token string `json:"token"`
}


type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

//var tokenexpiration int

var DIR string


//func init(){
//    dir,_ := filepath.Abs(filepath.Dir(os.Args[0]))
//    DIR = strings.Replace(dir,"bin","",1)
//    DIR = DIR+"src/sparq/authentication/keys/private_key"
//    fmt.Println(DIR)
//    }



var settings Settings = Settings{
	      "/etc/ssl/private_key",
		"/etc/ssl/public_key.pub",
			72}

//var settings Settings = Settings{"/home/gmuthuswamy/workspace/sparq/src/sparq/authentication/keys/private_key",
//			"/home/gmuthuswamy/workspace/sparq/src/sparq/authentication/keys/public_key.pub",
//			72}

func GetPath() Settings {
return settings
  
}
