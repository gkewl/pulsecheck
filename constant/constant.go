package constant

const Newline = "\n"

//http constants
const Http_Post = "POST"
const Http_Put = "PUT"
const Http_Get = "GET"
const Http_Delete = "DELETE"

const (
	ActorType_System = "SYSTEM"
	ActorType_User   = "USER"
	Role_Guest       = "Guest"
	Role_Admin       = "Admin"
	Role_SuperUser   = "SuperUser"
)

//define commn constant here
const UserName = "Actorname"
const UserId = "UserId"
const Claims = "Claims"

const Xid = "Xid"

const DefaultAdmin = 1

const (
	ResultOk   = "Ok"
	ResultFail = "Fail"
)

const SQL_NOT_FOUND string = "sql: no rows in result set"
const SQL_DUPLICATE_ENTRY = "Error 1062: Duplicate entry"

// Operational configuration defaults
const MaxDeadlockRetries = 5

//roles
const ApplicationRole = "ApplicationRole"
const ApplicationRoleType = "actor"
const DefaultApplicationRole = "Guest"

const (
	Guest = iota
	User
	Superuser
	Admin
)

// SearchLimit is the default for search functions
const SearchLimit = 50

//MaxPageSize is the maximum page size for server-side paging
const MaxPageSize = 5000

const (
	Source_OIG = "OIG"
	Source_SAM = "SAM"
	Source_OFAC = "OFAC"
)

const (
	Reference_Test="1128a1"
)