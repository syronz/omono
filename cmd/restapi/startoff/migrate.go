package startoff

import (
	"omono/domain/base/basmodel"
	"omono/domain/bill/bilmodel"
	"omono/domain/eaccounting/eacmodel"
	"omono/domain/location/locmodel"
	"omono/domain/material/matmodel"
	"omono/domain/notification/notmodel"
	"omono/domain/sync/synmodel"
	"omono/internal/core"
)

// Migrate the database for creating tables
func Migrate(engine *core.Engine) {
	// Sync Domain
	engine.DB.Table(synmodel.CompanyTable).AutoMigrate(&synmodel.Company{})

	// Base Domain
	engine.DB.Table(basmodel.SettingTable).AutoMigrate(&basmodel.Setting{})
	engine.DB.Exec("ALTER TABLE `bas_settings` ADD UNIQUE `idx_bas_settings_companyID_property`(`company_id`, property(50))")

	engine.DB.Table(basmodel.RoleTable).AutoMigrate(&basmodel.Role{})
	engine.DB.Exec("ALTER TABLE `bas_roles` ADD UNIQUE `idx_bas_roles_company_id_name`(`company_id`, name(40))")

	engine.DB.Table(basmodel.AccountTable).AutoMigrate(&basmodel.Account{})
	engine.DB.Exec("ALTER TABLE bas_accounts ADD CONSTRAINT `fk_bas_accounts_self` FOREIGN KEY (parent_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(basmodel.UserTable).AutoMigrate(&basmodel.User{})
	engine.DB.Exec("ALTER TABLE bas_users ADD CONSTRAINT `fk_bas_users_bas_roles` FOREIGN KEY (role_id) REFERENCES bas_roles(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bas_users ADD CONSTRAINT `fk_bas_users_bas_accounts` FOREIGN KEY (id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.ActivityDB.Table(basmodel.ActivityTable).AutoMigrate(&basmodel.Activity{})
	engine.DB.Table(basmodel.PhoneTable).AutoMigrate(&basmodel.Phone{})

	engine.DB.Table(basmodel.AccountPhoneTable).AutoMigrate(&basmodel.AccountPhone{})
	engine.DB.Exec("ALTER TABLE bas_account_phones ADD CONSTRAINT `fk_bas_accounts_phones_bas_accounts` FOREIGN KEY (account_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bas_account_phones ADD CONSTRAINT `fk_bas_accounts_phones_bas_phones` FOREIGN KEY (phone_id) REFERENCES bas_phones(id) ON DELETE CASCADE ON UPDATE CASCADE;")

	engine.DB.Table(basmodel.CityTable).AutoMigrate(&basmodel.City{})
	engine.DB.Exec("ALTER TABLE `bas_cities` ADD UNIQUE `idx_bas_cities_company_id_name`(`company_id`, name(20))")

	// Notification Domain
	engine.DB.Table(notmodel.MessageTable).AutoMigrate(&notmodel.Message{})
	engine.DB.Exec("ALTER TABLE not_messages ADD CONSTRAINT `fk_not_messages_created_by_bas_users` FOREIGN KEY (created_by) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE not_messages ADD CONSTRAINT `fk_not_messages_recepient_id_bas_users` FOREIGN KEY (recepient_id) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	// EAccounting Domain
	engine.DB.Table(eacmodel.CurrencyTable).AutoMigrate(&eacmodel.Currency{})

	engine.DB.Table(eacmodel.RateTable).AutoMigrate(&eacmodel.Rate{})
	engine.DB.Exec("ALTER TABLE eac_rates ADD CONSTRAINT `fk_eac_rates_bas_cities` FOREIGN KEY (city_id) REFERENCES bas_cities(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_rates ADD CONSTRAINT `fk_eac_rates_eac_currencies` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(eacmodel.TransactionTable).AutoMigrate(&eacmodel.Transaction{})
	engine.DB.Exec("ALTER TABLE eac_transactions ADD CONSTRAINT `fk_eac_transactions_eac_currencies` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_transactions ADD CONSTRAINT `fk_eac_transactions_bas_users` FOREIGN KEY (created_by) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE `eac_transactions` ADD UNIQUE `idx_eac_transactions_invoice_company_id_type`(`company_id`, invoice(20), type(20), post_date(4) )")

	engine.DB.Table(eacmodel.SlotTable).AutoMigrate(&eacmodel.Slot{})
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_bas_accounts` FOREIGN KEY (account_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_eac_transactions` FOREIGN KEY (transaction_id) REFERENCES eac_transactions(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_eac_currency` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(eacmodel.TempSlotTable).AutoMigrate(&eacmodel.TempSlot{})
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_bas_accounts` FOREIGN KEY (account_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_eac_transactions` FOREIGN KEY (transaction_id) REFERENCES eac_transactions(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_slots ADD CONSTRAINT `fk_eac_slot_eac_currency` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(eacmodel.BalanceTable).AutoMigrate(&eacmodel.Balance{})
	engine.DB.Exec("ALTER TABLE eac_balances ADD CONSTRAINT `fk_eac_balances_bas_accounts` FOREIGN KEY (account_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE eac_balances ADD CONSTRAINT `fk_eac_balances_eac_currencies` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	// Material Domain
	engine.DB.Table(matmodel.CompanyTable).AutoMigrate(&matmodel.Company{})
	engine.DB.Table(matmodel.ColorTable).AutoMigrate(&matmodel.Color{})

	engine.DB.Table(matmodel.GroupTable).AutoMigrate(&matmodel.Group{})

	engine.DB.Table(matmodel.UnitTable).AutoMigrate(&matmodel.Unit{})

	engine.DB.Table(matmodel.TagTable).AutoMigrate(&matmodel.Tag{})
	engine.DB.Exec("ALTER TABLE mat_tags ADD CONSTRAINT `fk_mat_tags_self` FOREIGN KEY (parent_id) REFERENCES mat_tags(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(matmodel.ProductTable).AutoMigrate(&matmodel.Product{})
	engine.DB.Exec("ALTER TABLE mat_products ADD CONSTRAINT `fk_mat_products_mat_units` FOREIGN KEY (unit_id) REFERENCES mat_units(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(matmodel.ProductTagTable).AutoMigrate(&matmodel.ProductTag{})
	engine.DB.Exec("ALTER TABLE mat_product_tags ADD CONSTRAINT `fk_mat_products_tag_mat_tags` FOREIGN KEY (tag_id) REFERENCES mat_tags(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE mat_product_tags ADD CONSTRAINT `fk_mat_products_tag_mat_products` FOREIGN KEY (product_id) REFERENCES mat_products(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(matmodel.GroupProductTable).AutoMigrate(&matmodel.GroupProduct{})
	engine.DB.Exec("ALTER TABLE mat_group_products ADD CONSTRAINT `fk_mat_group_products_mat_groups` FOREIGN KEY (group_id) REFERENCES mat_groups(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE mat_group_products ADD CONSTRAINT `fk_mat_group_products_mat_products` FOREIGN KEY (product_id) REFERENCES mat_products(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	// Location Domain
	engine.DB.Table(locmodel.StoreTable).AutoMigrate(&locmodel.Store{})
	engine.DB.Exec("ALTER TABLE loc_stores ADD CONSTRAINT `fk_loc_stores_bas_cities` FOREIGN KEY (city_id) REFERENCES bas_cities(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(locmodel.StoreUserTable).AutoMigrate(&locmodel.StoreUser{})
	engine.DB.Exec("ALTER TABLE loc_store_users ADD CONSTRAINT `fk_loc_store_users_loc_stores` FOREIGN KEY (store_id) REFERENCES loc_stores(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE loc_store_users ADD CONSTRAINT `fk_loc_store_users_bas_users` FOREIGN KEY (user_id) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	// Bill Domain
	engine.DB.Table(bilmodel.InvoiceTable).AutoMigrate(&bilmodel.Invoice{})
	engine.DB.Exec("ALTER TABLE bil_invoices ADD CONSTRAINT `fk_bil_invoices_loc_stores` FOREIGN KEY (store_id) REFERENCES loc_stores(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoices ADD CONSTRAINT `fk_bil_invoices_bas_accounts` FOREIGN KEY (account_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoices ADD CONSTRAINT `fk_bil_invoices_eac_currencies` FOREIGN KEY (currency_id) REFERENCES eac_currencies(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoices ADD CONSTRAINT `fk_bil_invoices_bas_users` FOREIGN KEY (created_by) REFERENCES bas_users(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")

	engine.DB.Table(bilmodel.InvoiceProductTable).AutoMigrate(&bilmodel.InvoiceProduct{})
	engine.DB.Exec("ALTER TABLE bil_invoice_products ADD CONSTRAINT `fk_bil_invoice_products_bil_invoices` FOREIGN KEY (invoice_id) REFERENCES bil_invoices(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoice_products ADD CONSTRAINT `fk_bil_invoice_products_bas_accounts1` FOREIGN KEY (source_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoice_products ADD CONSTRAINT `fk_bil_invoice_products_bas_accounts2` FOREIGN KEY (dest_id) REFERENCES bas_accounts(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.Exec("ALTER TABLE bil_invoice_products ADD CONSTRAINT `fk_bil_invoice_products_mat_products` FOREIGN KEY (product_id) REFERENCES mat_products(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
}
