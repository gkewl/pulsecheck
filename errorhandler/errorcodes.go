package errorhandler

import (
	"errors"
	"net/http"
)

const (
	json_conv_err_text = "Error during json conversion after processing: "
)

// ErrDBNoRows used to match against mysql "error" when no rows returned
var ErrDBNoRows = errors.New("sql: no rows in result set")

// ErrDBDuplicate used to match against mysql duplicate key error
var ErrDBDuplicate = "Error 1062: Duplicate entry"

// ErrDeadlockText used to detect the mysql deadlock error
var ErrDeadlockText = "Error 1213: Deadlock found"

// TestingNamedError is for use in testing
var TestingNamedError = newNamedError(0, "testing", http.StatusInternalServerError)

var (
	ErrInternalAppError        = newNamedError(101, "Internal application error", http.StatusInternalServerError)
	ErrInternalDBError         = newNamedError(102, "Internal database error", http.StatusInternalServerError)
	ErrNotYetImplemented       = newNamedError(103, "Not yet implemented", http.StatusInternalServerError)
	ErrDBContention            = newNamedError(104, "Database contention blocked request", http.StatusInternalServerError)
	ErrUnknownAction           = newNamedError(105, "Requested action not known", http.StatusNotFound)
	ErrUnsupportedLegacyAction = newNamedError(106, "Requested action has legacy interface", http.StatusConflict)
	ErrInvalidDate             = newNamedError(107, "Date was invalid", http.StatusPreconditionFailed)
	ErrInvalidDataFound        = newNamedError(108, "Invalid data was found in the DB", http.StatusPreconditionFailed)
	ErrMisconfiguration        = newNamedError(109, "Misconfiguration", http.StatusConflict)
)

/***************************************
	Common error across api starts from: 801
****************************************/
var (
	ErrOwnerNotAvailable         = newNamedError(801, "Owner not available in actor", http.StatusPreconditionFailed)
	ErrJsonEncodeFail            = newNamedError(802, json_conv_err_text, http.StatusInternalServerError)
	ErrJsonDecodeFail            = newNamedError(803, "JSON decoding failure", http.StatusBadRequest)
	ErrActorMissingOrNotFound    = newNamedError(804, "Actor not provided or not found", http.StatusPreconditionFailed)
	ErrLocationMissingOrNotFound = newNamedError(805, "Location not provided or not found", http.StatusPreconditionFailed)
	ErrMissingMeta               = newNamedError(806, "_meta_ key not found or malformed", http.StatusBadRequest)
	ErrSenderNotAllowed          = newNamedError(807, "Sender may not call this API", http.StatusForbidden)
	ErrXmlDecodeFail             = newNamedError(808, "XML decoding failure", http.StatusBadRequest)
	ErrPartMissingOrNotFound     = newNamedError(809, "PartNumber not provided or not found", http.StatusPreconditionFailed)
)

/***************************************
	DB  error starts from: 901
****************************************/
var (
	ErrDBConnectingDB         = newNamedError(900, "Error Connecting DB", http.StatusInternalServerError)
	ErrDBCreatingTransactions = newNamedError(901, "Error creating transaction", http.StatusNotFound)
	ErrDBCommitTransactions   = newNamedError(902, "Error committing transaction", http.StatusNotFound)
)

/***************************************
	Actor id starts from: 10001
****************************************/
var (
	ErrActorInvalidEntity  = newNamedError(10001, "Actor Invalid entity", http.StatusPreconditionFailed)
	ErrActorInsert         = newNamedError(10003, "Inserting entity to database failed", http.StatusConflict)
	ErrActorUpdate         = newNamedError(10004, "Update actor failed", http.StatusPreconditionFailed)
	ErrActorDelete         = newNamedError(10005, "Delete actor failed", http.StatusConflict)
	ErrActorDataNotFound   = newNamedError(10006, "Data not found", http.StatusNotFound)
	ErrActorSearch         = newNamedError(10007, "Search actor failed", http.StatusConflict)
	ErrActorParentNotFound = newNamedError(10008, "Parent name not found", http.StatusNotFound)
	ErrActorUserType       = newNamedError(10009, "Creating USER type actor not allowed", http.StatusPreconditionFailed)
	ErrActorInvalidRole    = newNamedError(10010, "Role update not allowed", http.StatusPreconditionFailed)
)

//authentication
var (
	ErrAuthGenerateToken         = newNamedError(16004, "Generate token failed", http.StatusUnauthorized)
	ErrUnAuthorizedUserForAPI    = newNamedError(16006, "User is not authorized to use current API", http.StatusUnauthorized)
	ErrAuthUserNameNotAvailable  = newNamedError(16007, "User name not available in body", http.StatusUnauthorized)
	ErrAuthWrongUseridOrPwd      = newNamedError(16008, "Invalid username or password", http.StatusUnauthorized)
	ErrAuthUserIdNotAvailbleInDB = newNamedError(16009, "Username is not configured in db", http.StatusNotFound)
	ErrAuthTokenNotAvailable     = newNamedError(16010, "Token not present in request", http.StatusUnauthorized)
	ErrAuthTokenError            = newNamedError(16011, "Invalid token", http.StatusUnauthorized)
)

//Search
var (
	ErrSearchMultipleEntity = newNamedError(25001, "Multiple entitites not allowed search parameter", http.StatusForbidden)
	ErrSearchError          = newNamedError(25002, "Error in search module", http.StatusNotFound)
)

/***************************************
	Token: 75000
****************************************/
var (
	ErrTokenDataNotFound = newNamedError(75000, "Token data not found", http.StatusNotFound)
	ErrTokenInsert       = newNamedError(75001, "Token Inserting entity to database failed", http.StatusConflict)
	ErrTokenDelete       = newNamedError(75002, "Token delete failed", http.StatusConflict)
	ErrTokenUpdate       = newNamedError(75003, "Token update failed", http.StatusConflict)
	ErrTokenValidation   = newNamedError(75004, "Token data invalid", http.StatusPreconditionFailed)
)
