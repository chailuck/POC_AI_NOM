// internal/models/models.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Individual struct {
	gorm.Model
	ID               string    `json:"id" gorm:"primaryKey"`
	HREF             string    `json:"href,omitempty"`
	Title            string    `json:"title,omitempty"`
	GivenName        string    `json:"givenName"`
	FamilyName       string    `json:"familyName"`
	MaritalStatus    string    `json:"maritalStatus,omitempty"`
	Gender           string    `json:"gender,omitempty"`
	NameType         string    `json:"nameType,omitempty"`
	Nationality      string    `json:"nationality,omitempty"`
	CreationDate     time.Time `json:"creationDate"`
	ModificationDate time.Time `json:"modificationDate"`
	CreatedBy        string    `json:"createdBy"`
	ModifiedBy       string    `json:"modifiedBy"`

	ContactMedium            []ContactMedium            `json:"contactMedium,omitempty" gorm:"foreignKey:IndividualID"`
	ExternalReference        []ExternalReference        `json:"externalReference,omitempty" gorm:"foreignKey:IndividualID"`
	IndividualIdentification []IndividualIdentification `json:"individualIdentification,omitempty" gorm:"foreignKey:IndividualID"`
	PartyCharacteristic      []PartyCharacteristic      `json:"partyCharacteristic,omitempty" gorm:"foreignKey:IndividualID"`
}

type ContactMedium struct {
	gorm.Model
	ID           string `json:"id" gorm:"primaryKey"`
	IndividualID string
	Type         string `json:"@type"`
	MediumType   string `json:"mediumType"`
	Preferred    bool   `json:"preferred"`

	// For PhoneContactMedium
	PhoneNumber string `json:"phoneNumber,omitempty"`

	// For GeographicAddressContactMedium
	Street1         string `json:"street1,omitempty"`
	Street2         string `json:"street2,omitempty"`
	City            string `json:"city,omitempty"`
	StateOrProvince string `json:"stateOrProvince,omitempty"`
	Country         string `json:"country,omitempty"`
	PostCode        string `json:"postCode,omitempty"`
}

type ExternalReference struct {
	gorm.Model
	ID                     string `json:"id" gorm:"primaryKey"`
	IndividualID           string
	Name                   string `json:"name"`
	ExternalIdentifierType string `json:"externalIdentifierType"`
	Type                   string `json:"@type"`
}

type IndividualIdentification struct {
	gorm.Model
	ID                 string `json:"id" gorm:"primaryKey"`
	IndividualID       string
	IdentificationType string    `json:"identificationType"`
	IdentificationId   string    `json:"identificationId"`
	ValidForEnd        time.Time `json:"validFor.endDateTime"`
}

type PartyCharacteristic struct {
	gorm.Model
	ID           string `json:"id" gorm:"primaryKey"`
	IndividualID string
	Name         string `json:"name"`
	Value        string `json:"value"`
	ValueType    string `json:"valueType"`
	Type         string `json:"@type"`
}
