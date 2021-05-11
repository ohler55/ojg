// Copyright (c) 2021, Peter Ohler, All rights reserved.

package main

// Patient is a struct used for Marshal and Unmarshal benchmarks.
type Patient struct {
	ResourceType         string
	ID                   string
	Text                 Text
	Identifier           []*Identifier
	Active               bool
	Name                 []*Name
	Telecom              []*Telecom
	Gender               string
	BirthDate            string
	XBirthDate           X `json:"_birthDate"`
	DeceasedBoolean      bool
	Address              []*Address
	Contact              []*Contact
	Communication        []*Communication
	ManagingOrganization Ref
	Meta                 Meta
}

// Text is a struct used for Marshal and Unmarshal benchmarks.
type Text struct {
	Status string
	Div    string
}

// Name is a struct used for Marshal and Unmarshal benchmarks.
type Name struct {
	Given   []string
	Family  string
	XFamily X `json:"_family"`
	Use     string
	Period  Period
}

// Ref is a struct used for Marshal and Unmarshal benchmarks.
type Ref struct {
	Reference string
	Display   string
}

// Identifier is a struct used for Marshal and Unmarshal benchmarks.
type Identifier struct {
	Use      string
	Type     CC
	System   string
	Value    string
	Period   Period
	Assigner Ref
}

// CC is a struct used for Marshal and Unmarshal benchmarks.
type CC struct {
	Coding []*Tag
	Text   string
}

// Period is a struct used for Marshal and Unmarshal benchmarks.
type Period struct {
	Start string
	End   string
}

// Meta is a struct used for Marshal and Unmarshal benchmarks.
type Meta struct {
	Tag []*Tag
}

// Tag is a struct used for Marshal and Unmarshal benchmarks.
type Tag struct {
	System string
	Code   string
}

// X is a struct used for Marshal and Unmarshal benchmarks.
type X struct {
	Extension []Extension
}

// Extension is a struct used for Marshal and Unmarshal benchmarks.
type Extension struct {
	URL           string
	ValueDateTime string
}

// Address is a struct used for Marshal and Unmarshal benchmarks.
type Address struct {
	Use        string
	Type       string
	Text       string
	Line       []string
	City       string
	District   string
	State      string
	PostalCode string
	Country    string
	Period     Period
}

// Telecom is a struct used for Marshal and Unmarshal benchmarks.
type Telecom struct {
	Use    string
	System string
	Value  string
	Rank   int
	Period Period
}

// Contact is a struct used for Marshal and Unmarshal benchmarks.
type Contact struct {
	Relationship []*CC
	Name         Name
	Telecom      []*Telecom
	Address      Address
	Gender       string
	Period       Period
}

// Communication is a struct used for Marshal and Unmarshal benchmarks.
type Communication struct {
	Language  CC
	Preferred bool
}

// Types for the citm_catalog.json file.

// Catalog is a struct used for Marshal and Unmarshal benchmarks.
type Catalog struct {
	AreaNames                map[string]string
	AudienceSubCategoryNames map[string]string
	BlockNames               map[string]string
	Events                   map[string]*Event
	Performances             []*Performance
	SeatCategoryNames        map[string]string
	SubTopicNames            map[string]string
	SubjectNames             map[string]string
	TopicNames               map[string]string
	TopicSubTopics           map[string][]int
	VenueNames               map[string]string
}

// Event is a struct used for Marshal and Unmarshal benchmarks.
type Event struct {
	Description *string
	ID          int `json:"id"`
	Logo        *string
	Name        string
	SubTopicIds []int
	SubjectCode *string
	Subtitle    *string
	TopicIds    []int
}

// Performance is a struct used for Marshal and Unmarshal benchmarks.
type Performance struct {
	EventID        int `json:"eventId"`
	ID             int `json:"id"`
	Logo           *string
	Name           *string
	Prices         []*Price
	SeatCategories []*SeatCategory
	SeatMapImage   *string
	Start          int64
	VenueCode      string
}

// Price is a struct used for Marshal and Unmarshal benchmarks.
type Price struct {
	Amount                int
	AudienceSubCategoryID int `json:"audienceSubCategoryId"`
	SeatCategoryID        int
}

// SeatCategory is a struct used for Marshal and Unmarshal benchmarks.
type SeatCategory struct {
	SeatCategoryID int `json:"seatCategoryId"`
	Areas          []*Area
}

// Area is a struct used for Marshal and Unmarshal benchmarks.
type Area struct {
	AreaID   int
	BlockIDs []int
}
