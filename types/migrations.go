package types

type MigrationStatus struct {
	Migrations      []DbMigration
	ActiveMigration *DbMigration
	IsInit          bool
}
