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
