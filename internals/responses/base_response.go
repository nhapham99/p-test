package responses

//import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Status  int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	Data    any    `json:"result"`
}

type PaginationResponse struct {
	RecordCount int64 `json:"recordCount" bson:"recordCount"`
	PageCount   int64 `json:"pageCount" bson:"pageCount"`
	CurrentPage int64 `json:"currentPage" bson:"currentPage"`
	PageSize    int64 `json:"pageSize" bson:"pageSize"`
	Records     any   `json:"records" bson:"records"`
}
