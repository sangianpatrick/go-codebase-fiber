package status

const (
	// success
	OK      string = "OK"
	CREATED string = "CREATED"

	// client error
	BAD_REQUEST          string = "BAD_REQUEST"
	UNAUTHORIZED         string = "UNAUTHORIZED"
	FORBIDDEN            string = "FORBIDDEN"
	NOT_FOUND            string = "NOT_FOUND"
	UNPROCESSABLE_ENTITY string = "UNPROCESSABLE_ENTITY"

	// server error
	INTERNAL_SERVER_ERROR string = "INTERNAL_SERVER_ERROR"

	// custom
	ALREADY_EXIST string = "ALREADY_EXIST"
)
