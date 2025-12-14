package models

type ExportRequest struct {
	Data   interface{} `json:"data"`
	Title  string      `json:"title"`
	Format string      `json:"format"`
}
