package types

import "fmt"

type ActionType string
type EntityType string
type StateTag string

const (
	ActionTypeCreate ActionType = "CREATE"
	ActionTypeDrop   ActionType = "DROP"

	EntityTypeSchema     EntityType = "SCHEMA"
	EntityTypeTable      EntityType = "TABLE"
	EntityTypeColumn     EntityType = "COLUMN"
	EntityTypeConstraint EntityType = "CONST"
	EntityTypeView       EntityType = "VIEW"
	EntityTypeSequence   EntityType = "SEQ"

	StateTagCS StateTag = "CS"
	StateTagPS StateTag = "PS"
)

type MigAction struct {
	ActionType ActionType
	StateTag   StateTag
	Entity     Entity

	/*
		For table, dep id is "schema.table", for col it's "schema.table.colname", for constraint/type it's "schema.name"
	*/
	DepArr []string
}

type DiffKey struct {
	Entity   Entity
	StateTag StateTag
}

type Entity struct {
	Type EntityType
	Id   []string

	// pointer to the real entity
	Ptr interface{}
}

func (e *Entity) GetEntityId() []string {
	switch e.Type {
	case EntityTypeSchema:
		schemaEntity := e.Ptr.(*Schema)
		e.Id = []string{schemaEntity.Name}
	case EntityTypeTable, EntityTypeView:
		tableEntity := e.Ptr.(*TableSchema)
		e.Id = []string{tableEntity.Schema, tableEntity.Name, tableEntity.DefSyntax}
	case EntityTypeColumn:
		colEntity := e.Ptr.(*Column)
		e.Id = []string{
			colEntity.TableSchema,
			colEntity.TableName,
			colEntity.Name,
			colEntity.DataType,
			colEntity.Default,
			fmt.Sprint(colEntity.Nullable),
		}
	case EntityTypeConstraint:
		conEntity := e.Ptr.(*DbConstraint)
		e.Id = []string{
			conEntity.Type,
			conEntity.ConSchema,
			conEntity.ConName,
			conEntity.TableSchema,
			conEntity.TableName,
			conEntity.ColName,
			conEntity.Syntax,
			conEntity.UpdateActionType,
			conEntity.DeleteActionType,
		}
	case EntityTypeSequence:
		seqEntity := e.Ptr.(*DbSequence)
		e.Id = []string{
			seqEntity.Schema,
			seqEntity.Name,
			seqEntity.DataType,
			fmt.Sprint(seqEntity.MinValue),
			fmt.Sprint(seqEntity.MaxValue),
			fmt.Sprint(seqEntity.IncrementBy),
			fmt.Sprint(seqEntity.CacheSize),
			fmt.Sprint(seqEntity.Cycle),
		}
	}

	return e.Id
}

func NewEntity(entityType EntityType, entityPtr interface{}) Entity {
	entity := Entity{
		Type: entityType,
		Ptr:  entityPtr,
	}

	entity.GetEntityId()

	return entity
}
