package elasticsearch

type BulkAction string

const (
	BulkActionIndex  BulkAction = "index"
	BulkActionCreate BulkAction = "create"
	BulkActionDelete BulkAction = "delete"
	BulkActionUpdate BulkAction = "update"
)

type BulkOperation struct {
	Action BulkAction
	Index  string
	ID     string
	Doc    map[string]interface{}
}

type BulkRequestBody struct {
	Operations []BulkOperation
}
