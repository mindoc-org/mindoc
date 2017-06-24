package migrate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/mindoc/models"
)

type MigrationVersion04 struct {
	isValid bool
	tables  []string
}

func NewMigrationVersion04() *MigrationVersion04 {
	return &MigrationVersion04{isValid: false, tables: make([]string, 0)}
}

func (m *MigrationVersion04) Version() int64 {
	return 201706210611
}

func (m *MigrationVersion04) ValidUpdate(version int64) error {
	if m.Version() > version {
		m.isValid = true
		return nil
	}
	m.isValid = false
	return errors.New("The target version is higher than the current version.")
}

func (m *MigrationVersion04) ValidForBackupTableSchema() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	var err error
	m.tables, err = ExportDatabaseTable()

	return err
}

func (m *MigrationVersion04) ValidForUpdateTableSchema() error {
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

func (m *MigrationVersion04) MigrationOldTableData() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	return nil
}

func (m *MigrationVersion04) MigrationNewTableData() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	o := orm.NewOrm()

	_, err := o.Raw("UPDATE md_members SET nickname = account WHERE nickname = '' ").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("UPDATE md_books SET link_id = 0 WHERE link_id IS NULL ").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("UPDATE md_attachment SET description = file_name WHERE description IS NULL ").Exec()
	if err != nil {
		return err
	}
	return nil
}

func (m *MigrationVersion04) AddMigrationRecord(version int64) error {
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

func (m *MigrationVersion04) MigrationCleanup() error {

	return nil
}

func (m *MigrationVersion04) RollbackMigration() error {
	if !m.isValid {
		return errors.New("The current version failed to verify.")
	}
	o := orm.NewOrm()
	_, err := o.Raw("ALTER TABLE md_members DROP COLUMN nickname").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("ALTER TABLE md_books DROP COLUMN link_id").Exec()
	if err != nil {
		return err
	}
	_, err = o.Raw("ALTER TABLE md_attachment DROP COLUMN description").Exec()
	if err != nil {
		return err
	}

	return nil
}
