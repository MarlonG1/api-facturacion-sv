package constants

import (
	"fmt"
	"github.com/MarlonG1/api-facturacion-sv/config"
	"reflect"
	"strings"
)

// ErrorCodesEn es un mapa que contiene los mensajes de error en inglés
var ErrorCodesEn = map[string]string{
	// Common errors
	"RequiredField":                "The field %s is required",
	"WithoutParams":                "Parameters are expected but none were received",
	"ExceededParameters":           "Error %s The required parameters have been exceeded, check the number of parameters sent, actual: %d",
	"InvalidLength":                "The field %s does not meet the required length, it must be %s characters, received %s",
	"InvalidNumberRange":           "The field %s does not meet the required length, it must be %s digits, received %s",
	"InvalidFormat":                "The field %s does not meet the required format, it must be %s, received %s",
	"InvalidPattern":               "The field %s does not meet the required pattern, it must be %s, received %s",
	"InvalidField":                 "The field %s does not meet the defined business rules",
	"InvalidDocumentNumber":        "Invalid document number, must be: DUI with hyphen, NIT (with or without hyphen), or a document of 3-20 characters. Received: %s",
	"NegativeDiscount":             "The discount %s is negative",
	"ExcessiveDiscount":            "The discount %s is greater than the subtotal %s",
	"NegativeTaxedAmount":          "The taxed sale %s is negative",
	"ExcessiveTaxedAmount":         "The taxed sale %f is greater than the price * quantity %f",
	"InvalidValue":                 "The value %d is not valid, it must be %s for field %s",
	"InvalidEmail":                 "The email %s is not valid. It must be a valid, existing email address and must be at least 3 characters long and no longer than 100",
	"InvalidRetentionCode":         "The retention code '%s' is not valid. It must be one of the allowed retention codes: (22 -> IVA 1 percent, C4 -> IVA 13 percent, C9 -> Other Cases)",
	"InvalidRetentionIVA":          "For item %d, the retention IVA %f does not match the expected value %f. ",
	"InvalidTotalSubjectRetention": "Total subject retention %f does not match the expected value %f",
	"InvalidTotalIVARetention":     "Total IVA retention %f does not match the expected value %f",
	"InconsistentReceiverNRC":      "The receiver NRC on document %s (%s) does not match the NRC on the first document (%s)",
	"InconsistentReceiverNIT":      "The receiver NIT in document %s (%s) does not match the NIT in the first document (%s)",
	"InvalidRetentionReceiver":     "The receiver information for the retention document could not be determined.",
	"DateOutOfAllowedRange":        "The item %d with document number %s has a date that is out of the allowed range. The document must be from the current period or the immediate previous period, and the retention document must be issued within the first 10 business days of the following month",

	// Identification errors
	"InvalidVersion":     "The version %s is not valid, it must be a number between 1 and 3",
	"InvalidAmbientCode": "The ambient code %s is not valid, it must be: 00 -> (Testing) and 01 -> (Production)",
	"InvalidDateTime":    "The date %s is not valid, it must be a date less than or equal to the current date",

	// Temporal (Time) errors
	"InvalidEmissionTime":     "The emission time %s is not valid, it must be a time less than or equal to the current time",
	"InvalidTransmissionType": "The transmission type %s is not valid, it must be: 1 -> (Normal) and 2 -> (Contingency)",

	// Monetary errors
	"InvalidAmount":   "The amount %s is not valid, it must be a number between 0 and 99999999999.99",
	"InvalidQuantity": "The quantity %s is not valid, it must be a number between 1 and 99999999999.99",
	"InvalidDiscount": "The discount %s is not valid, it must be a number between 0 and 100",
	"InvalidTax":      "The tax %s is not valid, it must be a number between 0 and 99999999999.99",
	"InvalidCurrency": "The currency %s is not valid, it must be USD",
	"InvalidItemType": "The item type %s is not valid, it must be: 1 -> (Product), 2 -> (Service), 3 -> (Both) and 4 -> (Tax)",
	"InvalidTaxType":  "The tax type %s is not valid, it must be within the allowed tax catalog",

	// Location errors
	"InvalidMunicipality":       "Invalid municipality code %s for department %s. Municipality code must be a two-digit number that follows the official catalog pattern. For example, valid codes for this department include: %s. Please refer to the official municipality catalog.",
	"InvalidEstablishmentType":  "The establishment type %s is not valid, it must be: 01 -> (Headquarters), 02 -> (Branch), 04 -> (Warehouse), 07 -> (Property or Yard) and 20 -> (Other)",
	"InvalidDocumentNumberItem": "The document number %s is not valid, when document type is '1' (physical), it must be a number between 1 and 20 characters, when document type is '2' (electronic), it must be a valid UUID",

	// Electronic invoice
	"InvalidServiceType":            "The services type %s is not valid, it must be a number between 1 and 6",
	"InvalidAssociatedDocumentCode": "The associated document code %s is not valid, it must be a number between 1 and 4",
	"InvalidPaymentTerms":           "When the operation condition type is '2' (credit), payment terms and period are required",
	"InvalidPaymentTypeOP2":         "When the operation condition type is '2' (credit), the payment type '01' (cash) is not allowed",
	"InvalidPaymentTermsOF":         "When the operation condition type is '1' (cash), payment terms and period are not allowed",
	"InvalidCreditNoteTransaction":  "The total %s is greater than the remaining balance for the document %s. Got total:%f, Remaining balance for this total: %f",

	"UnknownError": "No documentation found for the error, please contact the support team",
	"ServerError":  "An unexpected error has occurred",
}

// DTEErrorCodesEn es un mapa que contiene los mensajes de error DTE en inglés
var DTEErrorCodesEn = map[string]string{
	// Totals errors
	"InvalidTotalAmount":   "The total amounts do not match. Calculated total: %f, declared total: %f",
	"InvalidPaymentTotal":  "The total payments (%f) do not match the operation amount (%f)",
	"InvalidTotalIVA":      "The total IVA %f does not match the sum of IVA taxes %f",
	"InvalidTotalTaxed":    "The total taxed %f does not match the sum of taxed taxes %f",
	"InvalidTaxedAmount":   "The taxed sale %f does not match the taxable base %f",
	"InvalidTotalDiscount": "The total discounts %f do not match the sum of discounts %f",
	"UnsupportedTaxCode":   "Unsupported tax code %s, must be in the allowed tax catalog",

	// Document errors
	"ExceededItemsLimit":    "The number of items (%d) exceeds the allowed limit of 2000",
	"RequiredFieldMissing":  "Required field '%s' is not present for the document type %s",
	"InvalidDocumentState":  "Invalid state '%s' for the document",
	"RequiredSummary":       "Summary field required for the document type %s",
	"InvalidSubTotal":       "The subtotal %f does not match the sum of taxed, exempt, and non-subject totals %f",
	"InvalidItemNumber":     "The item number %d is not valid, it must be a number between 1 and 2000",
	"InvalidTaxRules":       "The tax rules do not match the taxed amount",
	"InvalidTaxRulesCCF":    "The tax rules do not match the taxed amount, only IVA tax is allowed for type 4 items",
	"InvalidTaxCode":        "The tax code %s is not valid, %s",
	"InvalidUnitMeasure":    "The unit of measure %d is not valid, it must be 99 for type 4 items",
	"InvalidNonTaxedAmount": "The Summary total_non_taxed cannot be greater than 0 when items non_taxed is 0",
	"InvalidTotalNonTaxed":  "The the sum of summary non-taxed amount items %f does not match the total_non-taxed sales %f",
	"DocumentInvalid":       "The document %s is already invalid, it cannot be invalidated again",

	// Business errors
	"MissingNRC":                         "NRC required for %s when the document type is %s",
	"InvalidTaxCalculation":              "The tax calculation %s (total_taxed - total_discount * tax) is incorrect. Expected: %f, Actual: %f",
	"MissingPaymentCondition":            "Payment condition required for the document type %s",
	"InvalidTotalToPay":                  "The total to pay (%f) does not match the total operation amount (%f)",
	"InvalidTotalToPayNonTaxedCCF":       "The total to pay (%f) must be equal to the total operation amount (%f) plus the sum of non-taxed amounts (%f) plus iva perception (%f), expected total to pay is %f",
	"InvalidTotalOperation":              "The total operation (%f) does not match the sum of totals (%f)",
	"InvalidDUIFormat":                   "The DUI format %s is not valid, it must be XXXXXXXX-X, must contain a hyphen",
	"InvalidNITFormat":                   "The NIT format %s is not valid, it must be 14 digits for legal entities and 9 digits for natural persons homologation",
	"MixedDocumentTypesNotAllowed":       "Mixed document types are not allowed, a different type than expected was found. Related document types must be the same, do not mix document types",
	"InvalidItemTotal":                   "The item total %f does not match the total sales %f",
	"MissingTaxes":                       "In items exists taxes but in summary dont exists no one tax",
	"MissingTaxesItem":                   "The item %d has taxed sale but no taxes are present, tax (20) IVA must be present for taxed sales",
	"MissingItemUnitPrice":               "The item %d, unit price is required when taxed sale is greater than zero",
	"InvalidRetentionAmount":             "The retention amount %s does not match the taxed total %s",
	"InvalidPerceptionAmount":            "The perception amount %s does not match the taxed total %s",
	"MixedSalesNotAllowed":               "Mixed sales (own and third party) are not allowed",
	"InvalidDTETypeForInvalidation":      "The document type %s is not valid for invalidation, only (01) -> Electronic Invoice, (03) -> Electronic Tax Credit Voucher and (04) -> Electronic Remission Note are allowed",
	"TransmissionFailed":                 "The transmission failed, the document was not received by Hacienda",
	"ExceededRelatedDocsLimit":           "The number of related documents (%d) exceeds the allowed limit of 50",
	"InvalidRelatedDocType":              "The related document type %s is not valid, it must be (04) -> Nota de Crédito, (08) -> Factura de Sujeto Excluido or (09) -> Nota de Débito",
	"InvalidRelatedDocDate":              "The related document date %s is not valid, it must be a date less than or equal to the current date",
	"InvalidRelatedDocNumberContingency": "The related document number %s is not valid, it must be a number between 1 and 36 characters",
	"InvalidRelatedDocNumberNormal":      "The related document number %s is not valid, it must be a number between 1 and 20 characters",
	"InvalidOtherDocsCount":              "The number of associated documents must be between 1 and 10, found %d",
	"InvalidMedicalDocFields":            "For medical documents (code 3) the description and detail fields should not be sent",
	"MutuallyExclusiveFields":            "The fields %s and %s are mutually exclusive, only one of them can be sent",
	"InvalidUnitPrice":                   "The unit price must equal taxed sale when non_taxed > 0. Item: %d, UnitPrice: %f, TaxedSale: %f",
	"InvalidUnitPriceZero":               "The unit price must be greater than zero taxed sale is greater than zero. Item: %d, UnitPrice: %f, TaxedSale: %f",
	"InvalidDecimals":                    "The amount must be multiple of 0.01, got: %s",
	"InvalidTotalTypeCCF":                "CCF cannot include non-subject totals as they do not generate tax credit",
	"InvalidSaleType":                    "CCF cannot include non-subject sales on item, as they do not generate tax credit",
	"ExcessiveItemTotal":                 "The item total %f is greater than the total sales %f",
	"InvalidMultipleSaleCategories":      "The item %d has multiple sale categories, only one",
	"InvalidTotalExempt":                 "The total exempt %f does not match the sum of exempt sales %f",
	"InvalidTotalNonSubject":             "The total non-subject %f does not match the sum of non-subject sales %f",
	"InvalidMixedSalesWithNonTaxed":      "The item %d has mixed sales with non-taxed, taxes field must be empty",
	"InvalidTaxesWithNonTaxed":           "The item %d has taxes with non-taxed, taxes field must be empty",
	"MissingRelatedDocWithNonTaxed":      "The item %d, related document is required for non-taxed items",
	"InvalidUnitPriceWithNonTaxed":       "The item %d, unit price must be 0 because it is a non-taxed item",
	"InvalidMixedSalesWithExempt":        "The item %d has mixed sales with exempt, only one",
	"InvalidTaxesWithExempt":             "The item %d has taxes with exempt, taxes field must be empty",
	"InvalidTaxesWithNonSubject":         "The item %d has taxes with non-subject, taxes field must be empty",
	"ValidationFailed":                   "Validation failed for the document",
	"InvalidSubTotalSales":               "The subtotal sales %f does not match the sum of all sales types %f",
	"InvalidTaxes":                       "Tax amounts present with zero taxed sale",
	"MixedSalesTypesNotAllowed":          "The item %d has mixed sales types, only one type allowed per item",
	"InvalidUnitPriceWithTaxedSale":      "The item %d has unit price zero with non-zero taxed sale",
	"InvalidMonetaryAmount":              "The %s amount %s is invalid",
	"InvalidUnitPriceForNonTaxed":        "Item %d with non-taxed amount must have unit price 0",
	"MissingItemRelatedDoc":              "The item %d, related document is required when Related documents section is present",
	"InvalidItemRelatedDoc":              "The item %d, related document %s is not valid",
	"InvalidSubTotalCalculation":         "The subtotal calculation with discounts (exempt_discount - non_subject_discount - taxed_discount) is incorrect. Expected: %f, Actual: %f",
	"InvalidIVACalculation":              "The IVA calculation considering discounts is incorrect. Expected: %f, Actual: %f",
	"InvalidTotalToPayCalculation":       "The total to pay calculation is incorrect. Expected: %f, Actual: %f",
	"InvalidMixedSalesWithNonSubject":    "The item %d has mixed sales with non-subject, only one",
	"InvalidIVAItemWithoutTaxedSale":     "Item %d has IVA item but no taxed sale",
	"InvalidIVAItemCalculation":          "Item %d has invalid IVA item calculation. Expected: %f, Actual: %f",

	"DiscountExceedsSubtotal": "The %s, with value %f, exceeds the subtotal %f",

	// Contingency errors
	"MissingContingencyType":   "Contingency type required when the operation type is contingency",
	"MissingContingencyReason": "Contingency reason required when the contingency type is 'Other reason'",
	"InvalidContingencyReason": "The contingency reason must be between 5 and 150 characters",

	// Model errors
	"InvalidModelType": "Invalid model type %d for normal transmission",

	// Extension errors
	"RequiredExtension": "The extension is required when the total amount is greater than or equal to $1,095.00, received amount: %f",

	// Operation errors
	"InvalidTransmissionModel": "Invalid billing model for the transmission type",

	// DTE-Receiver relationship errors
	"RequiredReceiver": "The receiver is required for the document type %s",

	"UnknownError": "Unknown error in the electronic tax document",
	"ServerError":  "Error in the electronic tax document",
}

// GetErrorMessage obtiene el mensaje de error según el idioma configurado
func GetErrorMessage(errorCode string, params ...interface{}) string {
	var message string
	var template string

	// Por si se genera un panic en la obtención del mensaje se manda el codigo error "ServerError"
	defer func() {
		if r := recover(); r != nil {
			message = fmt.Sprintf("%s", ErrorCodesEn["ServerError"])
		}
	}()

	template = getTemplate(errorCode)
	template = validateTemplate(template, errorCode, params...)

	if template == ErrorCodesEn["ServerError"] {
		return fmt.Sprintf("[%s] %s", "ServerError", template)
	}

	if template == ErrorCodesEn["WithoutParams"] {
		return fmt.Sprintf("[%s] %s", "WithoutParams", template)
	}

	message = fmt.Sprintf(template, params...)

	// Si el modo debug está activado, se muestra el código de error
	if config.Server.Debug {
		return fmt.Sprintf("[%s] %s", errorCode, message)
	}

	return message
}

// getTemplate Obtiene la plantilla de error según el código de error
func getTemplate(errorCode string) string {
	var template string

	template = ErrorCodesEn[errorCode]
	if template == "" {
		template = DTEErrorCodesEn[errorCode]
	}

	return template
}

// validateTemplate Obtiene la plantilla de error validada con los parámetros enviados
func validateTemplate(template string, errorCode string, params ...interface{}) string {
	switch {
	case template == "":
		return ErrorCodesEn["UnknownError"]
	case len(params) == 0 && strings.Contains(template, "%"):
		return ErrorCodesEn["WithoutParams"]
	case strings.Count(template, "%") != len(params):
		return fmt.Sprintf(ErrorCodesEn["ExceededParameters"],
			errorCode,
			len(params))
	case !validateParams(template, params...):
		return ErrorCodesEn["ServerError"]
	}

	return template
}

// validateParams Valida que los parámetros enviados coincidan con los especificadores de formato de la plantilla
func validateParams(template string, params ...interface{}) bool {
	specifiers := extractFormatSpecifiers(template)
	if len(specifiers) != len(params) {
		return false
	}

	for i, spec := range specifiers {
		expectedType := mapSpecifierToType(spec)
		if expectedType == nil || params[i] == nil {
			return false
		}

		// Revisa si el tipo del parámetro coincide con el especificador de formato
		paramKind := reflect.ValueOf(params[i]).Kind()
		if (spec == "s" && paramKind != reflect.String) ||
			((spec == "d" || spec == "i" || spec == "f") && !isNumericKind(paramKind)) {
			return false
		}
	}
	return true
}

// mapSpecifierToType Mapea los especificadores de formato a tipos de datos de reflect para su validación
func mapSpecifierToType(spec string) reflect.Type {
	switch spec {
	case "s":
		return reflect.TypeOf("")
	case "d", "i":
		return reflect.TypeOf(0)
	case "f":
		return reflect.TypeOf(0.0)
	default:
		return nil
	}
}

// extractFormatSpecifiers Extrae los especificadores de formato de la plantilla de error (%s, %d, %f, etc.)
func extractFormatSpecifiers(template string) []string {
	var specifiers []string
	for _, part := range strings.Split(template, "%")[1:] {
		if len(part) > 0 {
			specifiers = append(specifiers, string(part[0]))
		}
	}
	return specifiers
}

// isNumericKind Valida si el tipo de dato es numérico
func isNumericKind(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Float64
}
