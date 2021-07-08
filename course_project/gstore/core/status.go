package main

const (
	StatusOK = 100

	StatusQueryErr         = 201
	StatusJsonDecodeErr    = 202
	StatusParamsErr        = 203
	StatusRequestDecodeErr = 204
	StatusReadBodyErr      = 205

	// code for query3
	StatusNoCycle = 301

	// gstore api error
	StatusGStoreClientErr = 401
	StatusGStoreReadBodyErr = 402
	StatusGStoreJsonDecodeErr = 403
)

var statusText = map[int]string{
	StatusOK:               "OK",
	StatusQueryErr:         "Query Error",
	StatusJsonDecodeErr:    "Json Decode Error",
	StatusParamsErr:        "Request Params Error",
	StatusRequestDecodeErr: "Request Decode Error",
	StatusReadBodyErr:      "Read Request Error",

	StatusNoCycle: "No Cycle",

	StatusGStoreClientErr:     "connect gstore err",
	StatusGStoreReadBodyErr:   "read gstore data err",
	StatusGStoreJsonDecodeErr: "decode gstore data err",
}

func StatusText(code int) string {
	return statusText[code]
}
