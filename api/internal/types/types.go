package types

type ContextKey string
type TxKey struct{}

const (
	ContextID = ContextKey("id")
)
