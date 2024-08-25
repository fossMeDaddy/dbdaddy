package types

type ActionType string
type EntityType string
type StateTag string

const (
	MigActionTypeCreate ActionType = "CREATE"
	MigActionTypeDrop   ActionType = "DROP"

	EntityTypeSchema EntityType = "SCHEMA"
	EntityTypeType   EntityType = "TYPE"
	EntityTypeTable  EntityType = "TABLE"
	EntityTypeColumn EntityType = "COLUMN"
	EntityTypeView   EntityType = "VIEW"

	StateTagCS StateTag = "CS"
	StateTagPS StateTag = "PS"
)

type MigAction struct {
	ActionType       ActionType
	ActionEntityType EntityType
	StateTag         StateTag
	EntityId         []string
}

type DiffKey struct {
	EntityType EntityType
	StateTag   StateTag
	EntityId   []string
}
