package models

type InputDataCommon struct {
	Identification *Identification
	Issuer         *Issuer
	Receiver       *Receiver
	Extension      *Extension        `json:"extension,omitempty"`      // opcional
	RelatedDocs    []RelatedDocument `json:"relatedDocs,omitempty"`    // opcional
	OtherDocs      []OtherDocument   `json:"otherDocs,omitempty"`      // opcional
	ThirdPartySale *ThirdPartySale   `json:"thirdPartySale,omitempty"` // opcional
	Appendixes     []Appendix        `json:"appendixes,omitempty"`     // opcional
}
