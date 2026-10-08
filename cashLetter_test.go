package imagecashletter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCashLetterPanics(t *testing.T) {
	var cl *CashLetter

	require.Nil(t, cl.GetBundles())
	require.Nil(t, cl.GetRoutingNumberSummary())
	require.Nil(t, cl.GetCreditItems())
}

// TestCashLetterNoBundle validates no Bundle when CashLetterHeader.RecordTypeIndicator = "N"
func TestCashLetterNoBundle(t *testing.T) {
	// Create CheckDetail
	cd := mockCheckDetail()
	cd.AddCheckDetailAddendumA(mockCheckDetailAddendumA())
	cd.AddCheckDetailAddendumB(mockCheckDetailAddendumB())
	cd.AddCheckDetailAddendumC(mockCheckDetailAddendumC())
	cd.AddImageViewDetail(mockImageViewDetail())
	cd.AddImageViewData(mockImageViewData())
	cd.AddImageViewAnalysis(mockImageViewAnalysis())
	bundle := NewBundle(mockBundleHeader())
	bundle.AddCheckDetail(cd)

	// Create CashLetter
	cl := NewCashLetter(mockCashLetterHeader())
	cl.GetHeader().RecordTypeIndicator = "N"
	cl.AddBundle(bundle)
	err := cl.Create()
	var e *CashLetterError
	require.ErrorAs(t, err, &e)
	require.Equal(t, "RecordTypeIndicator", e.FieldName)
}

// TestCashLetterNoRoutingNumberSummary validates no Bundle when CashLetterHeader.CollectionTypeIndicator is not
// 00, 01, 02
func TestCashLetterRoutingNumberSummary(t *testing.T) {
	// Create CheckDetail
	cd := mockCheckDetail()
	cd.AddCheckDetailAddendumA(mockCheckDetailAddendumA())
	cd.AddCheckDetailAddendumB(mockCheckDetailAddendumB())
	cd.AddCheckDetailAddendumC(mockCheckDetailAddendumC())
	cd.AddImageViewDetail(mockImageViewDetail())
	cd.AddImageViewData(mockImageViewData())
	cd.AddImageViewAnalysis(mockImageViewAnalysis())
	bundle := NewBundle(mockBundleHeader())
	bundle.AddCheckDetail(cd)

	// Create CashLetter
	cl := NewCashLetter(mockCashLetterHeader())
	cl.GetHeader().CollectionTypeIndicator = "03"
	cl.AddBundle(bundle)
	rns := mockRoutingNumberSummary()
	cl.AddRoutingNumberSummary(rns)
	err := cl.Create()
	var e *CashLetterError
	require.ErrorAs(t, err, &e)
	require.Equal(t, "CollectionTypeIndicator", e.FieldName)
}

func TestCashLetter_customSequenceNumber(t *testing.T) {
	// Create a forward check bundle
	checkBundleHeader := mockBundleHeader()
	checkBundleHeader.SetBundleSequenceNumber(564)
	checkBundle := NewBundle(checkBundleHeader)
	cd := mockCheckDetail()
	cd.SetEceInstitutionItemSequenceNumber(283)
	cd.AddendumCount = 4
	firstAddendumA := mockCheckDetailAddendumA()
	firstAddendumA.RecordNumber = 1
	cd.AddCheckDetailAddendumA(firstAddendumA)
	secondAddendumA := mockCheckDetailAddendumA()
	secondAddendumA.RecordNumber = 2
	cd.AddCheckDetailAddendumA(secondAddendumA)
	firstAddendumC := mockCheckDetailAddendumC()
	firstAddendumC.RecordNumber = 1
	firstAddendumC.EndorsingBankItemSequenceNumber = "7"
	cd.AddCheckDetailAddendumC(firstAddendumC)
	secondAddendumC := mockCheckDetailAddendumC()
	secondAddendumC.RecordNumber = 2
	secondAddendumC.EndorsingBankItemSequenceNumber = ""
	cd.AddCheckDetailAddendumC(secondAddendumC)
	checkBundle.AddCheckDetail(cd)

	// Create a return bundle
	returnBundleHeader := mockBundleHeader()
	returnBundleHeader.BundleSequenceNumber = "" // test auto-increment behavior
	returnBundle := NewBundle(returnBundleHeader)
	rd := mockReturnDetail()
	rd.SetEceInstitutionItemSequenceNumber(4923)
	rd.AddendumCount = 2
	rd.AddReturnDetailAddendumA(mockReturnDetailAddendumA())
	rd.AddReturnDetailAddendumD(mockReturnDetailAddendumD())
	returnBundle.AddReturnDetail(rd)

	clh := mockCashLetterHeader()
	cl := NewCashLetter(clh)
	cl.AddBundle(checkBundle)
	cl.AddBundle(returnBundle)
	require.NoError(t, cl.Create())

	require.Len(t, cl.Bundles, 2)
	require.Equal(t, "0564", cl.Bundles[0].BundleHeader.BundleSequenceNumber)
	require.Equal(t, "0565", cl.Bundles[1].BundleHeader.BundleSequenceNumber)

	require.Len(t, cl.Bundles[0].Checks, 1)
	wantCheckSeq := "000000000000283"
	require.Equal(t, wantCheckSeq, cl.Bundles[0].Checks[0].EceInstitutionItemSequenceNumber)

	// CheckDetailAddendumA
	require.Len(t, cl.Bundles[0].Checks[0].CheckDetailAddendumA, 2)
	require.Equal(t, 1, cl.Bundles[0].Checks[0].CheckDetailAddendumA[0].RecordNumber)
	require.Equal(t, 2, cl.Bundles[0].Checks[0].CheckDetailAddendumA[1].RecordNumber)
	require.Equal(t, "1              ", cl.Bundles[0].Checks[0].CheckDetailAddendumA[0].BOFDItemSequenceNumber, "should not have overwritten custom sequence number")

	// CheckDetailAddendumC
	require.Len(t, cl.Bundles[0].Checks[0].CheckDetailAddendumC, 2)
	require.Equal(t, 1, cl.Bundles[0].Checks[0].CheckDetailAddendumC[0].RecordNumber)
	require.Equal(t, 2, cl.Bundles[0].Checks[0].CheckDetailAddendumC[1].RecordNumber)
	require.Equal(t, "7", cl.Bundles[0].Checks[0].CheckDetailAddendumC[0].EndorsingBankItemSequenceNumber, "should not have overwritten custom sequence number")
	require.Equal(t, wantCheckSeq, cl.Bundles[0].Checks[0].CheckDetailAddendumC[1].EndorsingBankItemSequenceNumber, "should have populated empty sequence number")

	require.Len(t, cl.Bundles[1].Returns, 1)
	wantReturnSeq := "000000000004923"
	require.Equal(t, wantReturnSeq, cl.Bundles[1].Returns[0].EceInstitutionItemSequenceNumber)
	require.Equal(t, "1              ", cl.Bundles[1].Returns[0].ReturnDetailAddendumA[0].BOFDItemSequenceNumber, "should not have overwritten custom sequence number")
	require.Equal(t, wantReturnSeq, cl.Bundles[1].Returns[0].ReturnDetailAddendumD[0].EndorsingBankItemSequenceNumber)
}

// TestCashLetterValidate_NilHeader verifies basic structural check for CashLetterHeader
// is enforced even when SkipAll is set (to avoid later nil derefs in build/processing).
func TestCashLetterValidate_NilHeader(t *testing.T) {
	cl := CashLetter{} // no header set
	err := cl.Validate()
	require.Error(t, err)
	require.Equal(t, "nil CashLetterHeader", err.Error())
}

func TestCashLetterValidate_NilHeaderSkipAll(t *testing.T) {
	cl := CashLetter{}
	cl.SetValidation(&ValidateOpts{SkipAll: true})
	// Nil header check happens before SkipAll for structural safety
	err := cl.Validate()
	require.Error(t, err)
	require.Equal(t, "nil CashLetterHeader", err.Error())
}

func TestCashLetterBuild_NilHeader(t *testing.T) {
	cl := NewCashLetter(nil)
	err := cl.build()
	require.Error(t, err)
	require.Equal(t, "nil CashLetterHeader", err.Error())
}

// TestCashLetter_preservesBOFDItemSequenceNumber verifies that a caller-supplied
// BOFDItemSequenceNumber survives Create().
//
// The BOFD item sequence number identifies the item at the bank of first deposit. On a
// return received from the Fed it carries the sequence number of the *original forward*
// item, which is what lets the depositing bank tie the return back to what it sent. It
// is therefore not interchangeable with the detail record's own ECE institution item
// sequence number, and overwriting it with that value destroys the linkage.
//
// Matches the behaviour already implemented for EndorsingBankItemSequenceNumber on
// CheckDetailAddendumC: a value set by the caller is left alone, an empty one is
// populated from the detail's sequence number.
func TestCashLetter_preservesBOFDItemSequenceNumber(t *testing.T) {
	const originalForwardISN = "000078366005826"

	// Return whose addendum A carries the forward item's sequence number.
	returnBundle := NewBundle(mockBundleHeader())
	rd := mockReturnDetail()
	rd.SetEceInstitutionItemSequenceNumber(4923)
	rd.AddendumCount = 1
	rdAddendumA := mockReturnDetailAddendumA()
	rdAddendumA.BOFDItemSequenceNumber = originalForwardISN
	rd.AddReturnDetailAddendumA(rdAddendumA)
	returnBundle.AddReturnDetail(rd)

	// Return whose addendum A leaves it unset, so it should be populated.
	rdEmpty := mockReturnDetail()
	rdEmpty.SetEceInstitutionItemSequenceNumber(4924)
	rdEmpty.AddendumCount = 1
	rdEmpty.AddReturnDetailAddendumA(mockReturnDetailAddendumAWithoutBOFDItemSequenceNumber())
	returnBundle.AddReturnDetail(rdEmpty)

	// Forward check with the same contract.
	checkBundle := NewBundle(mockBundleHeader())
	cd := mockCheckDetail()
	cd.SetEceInstitutionItemSequenceNumber(283)
	cd.AddendumCount = 1
	cdAddendumA := mockCheckDetailAddendumA()
	cdAddendumA.BOFDItemSequenceNumber = originalForwardISN
	cd.AddCheckDetailAddendumA(cdAddendumA)
	checkBundle.AddCheckDetail(cd)

	cl := NewCashLetter(mockCashLetterHeader())
	cl.AddBundle(returnBundle)
	cl.AddBundle(checkBundle)
	require.NoError(t, cl.Create())

	require.Equal(t, originalForwardISN, cl.Bundles[0].Returns[0].ReturnDetailAddendumA[0].BOFDItemSequenceNumber,
		"return: should not have overwritten the caller's BOFD item sequence number")
	require.Equal(t, "000000000004924", cl.Bundles[0].Returns[1].ReturnDetailAddendumA[0].BOFDItemSequenceNumber,
		"return: should have populated the empty BOFD item sequence number")
	require.Equal(t, originalForwardISN, cl.Bundles[1].Checks[0].CheckDetailAddendumA[0].BOFDItemSequenceNumber,
		"forward: should not have overwritten the caller's BOFD item sequence number")

	// The detail records' own sequence numbers are unaffected.
	require.Equal(t, "000000000004923", cl.Bundles[0].Returns[0].EceInstitutionItemSequenceNumber)
	require.Equal(t, "000000000000283", cl.Bundles[1].Checks[0].EceInstitutionItemSequenceNumber)
}
