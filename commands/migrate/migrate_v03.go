package migrate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/mindoc-org/mindoc/models"
)

type MigrationVersion03 struct {
	isValid bool
	tables  []string
}

func NewMigrationVersion03() *MigrationVersion03 {
	return &MigrationVersion03{isValid: false, tables: make([]string, 0)}
}

func (m *MigrationVersion03) Version() int64 {
	return 201705271114
}

func (m *MigrationVersion03) ValidUpdate(version int64) error {
	if m.Version() > version {
		m.isValid = true
		return nil
	}
	m.isValid = false
	return errors.New("The target version is higher than the current version.")
}

func (m *MigrationVersion03) ValidForBackupTableSchema() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	var err error
	m.tables, err = ExportDatabaseTable()

	return err
}

func (m *MigrationVersion03) ValidForUpdateTableSchema() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}

	err := orm.RunSyncdb("default", false, true)

	if err != nil {
		return err
	}

	//_,err = o.Raw("ALTER TABLE md_members ADD auth_method VARCHAR(50) DEFAULT 'local' NULL").Exec()

	return err
}

func (m *MigrationVersion03) MigrationOldTableData() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	return nil
}

func (m *MigrationVersion03) MigrationNewTableData() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	o := orm.NewOrm()

	_, err := o.Raw("UPDATE md_members SET auth_method = 'local'").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用文档历史','ENABLE_DOCUMENT_HISTORY','true' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLE_DOCUMENT_HISTORY');").Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *MigrationVersion03) AddMigrationRecord(version int64) error {
	o := orm.NewOrm()
	tables, err := ExportDatabaseTable()

	if err != nil {
		return err
	}
	migration := models.NewMigration()
	migration.Version = version
	migration.Status = "update"
	migration.CreateTime = time.Now()
	migration.Name = fmt.Sprintf("update_%d", version)
	migration.Statements = strings.Join(tables, "\r\n")

	_, err = o.Insert(migration)

	return err
}

func (m *MigrationVersion03) MigrationCleanup() error {

	return nil
}

func (m *MigrationVersion03) RollbackMigration() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	o := orm.NewOrm()
	_, err := o.Raw("ALTER TABLE md_members DROP COLUMN auth_method").Exec()
	if err != nil {
		return err
	}

	_, err = o.Raw("DROP TABLE md_document_history").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("DELETE md_options WHERE option_name = 'ENABLE_DOCUMENT_HISTORY'").Exec()

	if err != nil {
		return err
	}

	return nil
}
