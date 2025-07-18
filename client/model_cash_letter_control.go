/*
 * ImageCashLetter API
 *
 * Moov Image Cash Letter (ICL) implements an HTTP API for creating, parsing, and validating ImageCashLetter files.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

// CashLetterControl struct for CashLetterControl
type CashLetterControl struct {
	// CashLetterControl ID
	Id string `json:"id,omitempty"`
	// CashLetterBundleCount identifies the total number of bundles within the cash letter.
	CashLetterBundleCount int32 `json:"cashLetterBundleCount,omitempty"`
	// CashLetterItemsCount identifies the total number of items within the cash letter.
	CashLetterItemsCount int32 `json:"cashLetterItemsCount"`
	// CashLetterTotalAmount identifies the total dollar value of all item amounts within the cash letter.
	CashLetterTotalAmount int32 `json:"cashLetterTotalAmount"`
	// CashLetterImagesCount identifies the total number of ImageViewDetail(s) within the CashLetter.
	CashLetterImagesCount int32 `json:"cashLetterImagesCount,omitempty"`
	// ECEInstitutionName identifies the short name of the institution that creates the CashLetterControl.
	EceInstitutionName string `json:"eceInstitutionName,omitempty"`
	// SettlementDate identifies the date that the institution that creates the cash letter expects settlement.
	SettlementDate time.Time `json:"settlementDate"`
	// CreditTotalIndicator is a code that indicates whether Credit Items are included in this record’s totals. If so, they will be included in TotalItemCount and FileTotalAmount. TotalRecordCount includes all records of all types regardless of the value of this field. * ` ` - No Credit Items * `0` - Credit Items are not included in totals * `1` - Credit Items are included in totals
	CreditTotalIndicator int32 `json:"creditTotalIndicator,omitempty"`
}
