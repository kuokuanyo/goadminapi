// Base is base model structure.
type Base struct {
	TableName string

	Conn db.Connection
	Tx   *sql.Tx
}