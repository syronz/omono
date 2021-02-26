package transactiontype

import "omono/internal/types"

// Transactions type
const (
	Manual         types.Enum = "manual"
	JournalEntry   types.Enum = "journal-entry"
	JournalUpdate  types.Enum = "journal-update"
	VoucherApprove types.Enum = "voucher-Approve"
	JournalVoucher types.Enum = "journal-voucher"
	VoucherUpdate  types.Enum = "voucher-update"
	OpeningEntry   types.Enum = "opening-entry"
	OpeningVocuher types.Enum = "opening-voucher"
	PaymentEntry   types.Enum = "payment-entry"
	PaymentVoucher types.Enum = "payment-voucher"
	ReceiptEntry   types.Enum = "receipt-entry"
	ReceiptVoucher types.Enum = "receipt-voucher"
)

var List = []types.Enum{
	Manual,
	JournalEntry,
	JournalUpdate,
	JournalVoucher,
	OpeningEntry,
	OpeningVocuher,
	PaymentEntry,
	PaymentVoucher,
	ReceiptEntry,
	ReceiptVoucher,
	VoucherApprove,
	VoucherUpdate,
}

// Join make a string for showing in the api
func Join() string {
	return types.JoinEnum(List)
}

//SimilarType is used for fetching the similar type for transaction
//currentlu used in invoice generation. journal entry -> returns journal voucher, etc...
func SimilarType(tType types.Enum) (similarType types.Enum) {
	switch tType {
	case JournalEntry:
		return JournalVoucher
	case JournalVoucher:
		return JournalEntry
	case ReceiptEntry:
		return ReceiptVoucher
	case ReceiptVoucher:
		return JournalEntry
	case PaymentEntry:
		return PaymentVoucher
	case Manual:
		return Manual
	default:
		return ""
	}

}
