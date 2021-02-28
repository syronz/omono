package startoff

import (
	"omono/domain/base/basmodel"
	"omono/domain/notification/notmodel"
	"omono/domain/subscriber/submodel"
	"omono/internal/core"
)

// Migrate the database for creating tables
func Migrate(engine *core.Engine) {

	// Base Domain
	engine.DB.Table(basmodel.SettingTable).AutoMigrate(&basmodel.Setting{})

	engine.DB.Table(basmodel.RoleTable).AutoMigrate(&basmodel.Role{})

	engine.DB.Table(basmodel.UserTable).AutoMigrate(&basmodel.User{})
	engine.DB.Exec("ALTER TABLE bas_users ADD CONSTRAINT `fk_bas_users_bas_roles` FOREIGN KEY (role_id) REFERENCES bas_roles(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.ActivityDB.Table(basmodel.ActivityTable).AutoMigrate(&basmodel.Activity{})

	engine.DB.Table(basmodel.CityTable).AutoMigrate(&basmodel.City{})

	// Subscriber Domain
	engine.DB.Table(submodel.AccountTable).AutoMigrate(&submodel.Account{})
	engine.DB.Exec("ALTER TABLE sub_accounts ADD CONSTRAINT `fk_sub_accounts_self` FOREIGN KEY (parent_id) REFERENCES sub_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(submodel.PhoneTable).AutoMigrate(&submodel.Phone{})

	engine.DB.Table(submodel.AccountPhoneTable).AutoMigrate(&submodel.AccountPhone{})
	engine.DB.Exec("ALTER TABLE sub_account_phones ADD CONSTRAINT `fk_sub_accounts_phones_sub_accounts` FOREIGN KEY (account_id) REFERENCES sub_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE sub_account_phones ADD CONSTRAINT `fk_sub_accounts_phones_sub_phones` FOREIGN KEY (phone_id) REFERENCES sub_phones(id) ON DELETE CASCADE ON UPDATE CASCADE;")

	// Notification Domain
	engine.DB.Table(notmodel.MessageTable).AutoMigrate(&notmodel.Message{})
	engine.DB.Exec("ALTER TABLE not_messages ADD CONSTRAINT `fk_not_messages_created_by_bas_users` FOREIGN KEY (created_by) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE not_messages ADD CONSTRAINT `fk_not_messages_recipient_id_bas_users` FOREIGN KEY (recipient_id) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

}
