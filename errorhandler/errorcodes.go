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

//authentication
var (
	ErrAuthGenerateToken         = newNamedError(1000, "Generate token failed", http.StatusUnauthorized)
	ErrUnAuthorizedUserForAPI    = newNamedError(1001, "User is not authorized to use current API", http.StatusUnauthorized)
	ErrAuthUserNameNotAvailable  = newNamedError(1002, "User name not available in body", http.StatusUnauthorized)
	ErrAuthWrongUseridOrPwd      = newNamedError(1003, "Invalid username or password", http.StatusUnauthorized)
	ErrAuthUserIdNotAvailbleInDB = newNamedError(1004, "Username is not configured in db", http.StatusNotFound)
	ErrAuthTokenNotAvailable     = newNamedError(1005, "Token not present in request", http.StatusUnauthorized)
	ErrAuthTokenError            = newNamedError(1006, "Invalid token", http.StatusUnauthorized)
)

var (
	ErrCompanyInsert       = newNamedError(1101, "Company Inserting entity to database failed", http.StatusConflict)
	ErrCompanyDataNotFound = newNamedError(1102, "Company Data not found", http.StatusNotFound)
	ErrCompanyUpdate       = newNamedError(1103, "Update Company failed", http.StatusConflict)
	ErrCompanyDelete       = newNamedError(1104, "Company delete failed", http.StatusConflict)
)

var (
	ErrAuthUserInsert       = newNamedError(1201, "Auth user  Inserting entity to database failed", http.StatusConflict)
	ErrAuthUserDataNotFound = newNamedError(1202, "User Data not found", http.StatusNotFound)
	ErrAuthUserUpdate       = newNamedError(1203, "Update authuser failed", http.StatusConflict)
	ErrAuthUserDelete       = newNamedError(1204, "AuthUser  delete failed", http.StatusConflict)
)
var (
	ErrUserInsert       = newNamedError(1301, "User Inserting entity to database failed", http.StatusConflict)
	ErrUserDataNotFound = newNamedError(1302, "User Data not found", http.StatusNotFound)
	ErrUserUpdate       = newNamedError(1303, "Update User failed", http.StatusConflict)
	ErrUserDelete       = newNamedError(1304, "User delete failed", http.StatusConflict)
)

var (
	ErrEmployeeInsert       = newNamedError(1401, "Employee Inserting entity to database failed", http.StatusConflict)
	ErrEmployeeDataNotFound = newNamedError(1402, "Employee Data not found", http.StatusNotFound)
	ErrEmployeeUpdate       = newNamedError(1403, "Update Employee failed", http.StatusConflict)
	ErrEmployeeDelete       = newNamedError(1404, "Employee delete failed", http.StatusConflict)
	ErrEmployeeUpload       = newNamedError(1405, "Employee upload  failed", http.StatusConflict)
	ErrEmployeeSearch       = newNamedError(1406, "Employee search  failed", http.StatusConflict)
)
