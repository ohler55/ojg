// Copyright (c) 2021, Peter Ohler, All rights reserved.

package main

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

type Text struct {
	Status string
	Div    string
}

type Name struct {
	Given   []string
	Family  string
	XFamily X `json:"_family"`
	Use     string
	Period  Period
}

type Ref struct {
	Reference string
	Display   string
}

type Identifier struct {
	Use      string
	Type     CC
	System   string
	Value    string
	Period   Period
	Assigner Ref
}

type CC struct {
	Coding []*Tag
	Text   string
}

type Period struct {
	Start string
	End   string
}

type Meta struct {
	Tag []*Tag
}

type Tag struct {
	System string
	Code   string
}

type X struct {
	Extension []Extension
}

type Extension struct {
	URL           string
	ValueDateTime string
}

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

type Telecom struct {
	Use    string
	System string
	Value  string
	Rank   int
	Period Period
}

type Contact struct {
	Relationship []*CC
	Name         Name
	Telecom      []*Telecom
	Address      Address
	Gender       string
	Period       Period
}

type Communication struct {
	Language  CC
	Preferred bool
}

// Types for the citm_catalog.json file.
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

type Price struct {
	Amount                int
	AudienceSubCategoryId int
	SeatCategoryID        int
}

type SeatCategory struct {
	SeatCategoryId int
	Areas          []*Area
}

type Area struct {
	AreaID   int
	BlockIDs []int
}
