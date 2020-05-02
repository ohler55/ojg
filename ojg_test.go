// Copyright (c) 2020, Peter Ohler, All rights reserved.

package ojg_test

import (
	"testing"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/tt"
)

type data struct {
	src string
	// Empty means no error expected while non empty should be compared
	// err.Error().
	expect    string
	value     interface{}
	onlyOne   bool
	noComment bool
}

func TestValidateString(t *testing.T) {
	for _, d := range []data{
		{src: "null"},
		{src: "true"},
		{src: "false"},
		{src: "[]"},
		{src: "[true]"},
		{src: "[true,false]"},
		{src: "[[],[true],false]"},
		{src: "[[],[true]false]", expect: "expected a comma or close, not 'f' at 1:11"},

		{src: "123"},
		{src: "-1.23"},
		{src: "[1,2]"},

		{src: "[]"},
		{src: "null {}"},
		{src: "null {}", expect: "extra characters after close, '{' at 1:6", onlyOne: true},

		{src: "-1.23"},
		{src: "+1.23", expect: "unexpected character '+' at 1:1"},
		{src: "1.23e+3"},
		{src: "1.23e-3"},
		{src: "1.23e3"},
		{src: "1.2e3e3", expect: "invalid number '1.2e3e' at 1:6"},
		{src: "0.3"},
		{src: "03", expect: "invalid number '03' at 1:2"},

		{src: `""`},
		{src: `"abc"`},
		{src: `"a\tb\nc\b\"\\d\f\r"`},
		{src: "\"bass \U0001D122\""},
		{src: `"a \u2669"`},
		{src: `"bad \uabcz"`, expect: "invalid JSON unicode character 'z' at 1:11"},

		{src: "[\n  // a comment\n]"},
		{src: "[\n  // a comment\n]", expect: "comments not allowed at 2:3", noComment: true},
		{src: "[\n  / a comment\n]", expect: "unexpected character ' ' at 2:4"},

		{src: "{}"},
		{src: `{"a":3}`},
		{src: `{"a": 3, "b": true}`},
		{src: `{"a":{"b":{"c":true}}}`},
		{src: `{x}`, expect: "expected a string start or object close, not 'x' at 1:2"},
	} {
		var err error
		if d.onlyOne || d.noComment {
			p := ojg.Validator{OnlyOne: d.onlyOne, NoComment: d.noComment}
			err = p.Validate([]byte(d.src))
		} else {
			err = ojg.Validate([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}

func TestDev(t *testing.T) {
	for _, d := range []data{
		{src: "0"},
		{src: "-0"},
		{src: "0.3"},
		{src: "03", expect: "invalid number at 1:2"},
		{src: "123"},
		{src: "-123"},
		{src: "-", expect: "incomplete JSON at 1:1"},
		{src: "12.3"},
		{src: "12.3456"},
		{src: "+1.23", expect: "unexpected character '+' at 1:1"},
		{src: "1.23e+3"},
		{src: "1.23e-3"},
		{src: "1.23e3"},
		{src: "1.2e3e3", expect: "invalid number at 1:6"},
	} {
		var err error
		if d.onlyOne || d.noComment {
			p := ojg.Validator{OnlyOne: d.onlyOne, NoComment: d.noComment}
			err = p.Validate([]byte(d.src))
		} else {
			err = ojg.Validate([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
		}
	}
}

func TestParseString(t *testing.T) {
	for _, d := range []data{

		{src: "null", value: nil},
		{src: "true", value: true},
		{src: "false", value: false},
		{src: "123", value: 123},
		{src: "-321", value: -321},
		{src: "12.3", value: 12.3},
		{src: `"xyz"`, value: "xyz"},

		{src: "[]", value: []interface{}{}},
		{src: "[true]", value: []interface{}{true}},
		{src: "[true,false]", value: []interface{}{true, false}},
		{src: "[[]]", value: []interface{}{[]interface{}{}}},
		{src: "[[true]]", value: []interface{}{[]interface{}{true}}},

		{src: "{}", value: map[string]interface{}{}},
		{src: `{"abc":true}`, value: map[string]interface{}{"abc": true}},
		{src: `{"abc":{"def":3}}`, value: map[string]interface{}{"abc": map[string]interface{}{"def": 3}}},
	} {
		var err error
		var v interface{}
		if d.onlyOne || d.noComment {
			p := ojg.Parser{OnlyOne: d.onlyOne, NoComment: d.noComment}
			v, err = p.Parse([]byte(d.src))
		} else {
			v, err = ojg.Parse([]byte(d.src))
		}
		if 0 < len(d.expect) {
			tt.NotNil(t, err, d.src)
			tt.Equal(t, d.expect, err.Error(), d.src)
		} else {
			tt.Nil(t, err, d.src)
			tt.Equal(t, d.value, v, d.src)
		}
	}
}

const patientJSON = `{
  "resourceType": "Patient",
  "id": "example",
  "text": {
    "status": "generated",
    "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">\n\t\t\t<table>\n\t\t\t\t<tbody>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Name</td>\n\t\t\t\t\t\t<td>Peter James \n              <b>Chalmers</b> (&quot;Jim&quot;)\n            </td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Address</td>\n\t\t\t\t\t\t<td>534 Erewhon, Pleasantville, Vic, 3999</td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Contacts</td>\n\t\t\t\t\t\t<td>Home: unknown. Work: (03) 5555 6473</td>\n\t\t\t\t\t</tr>\n\t\t\t\t\t<tr>\n\t\t\t\t\t\t<td>Id</td>\n\t\t\t\t\t\t<td>MRN: 12345 (Acme Healthcare)</td>\n\t\t\t\t\t</tr>\n\t\t\t\t</tbody>\n\t\t\t</table>\n\t\t</div>"
  },
  "identifier": [
    {
      "use": "usual",
      "system": "urn:oid:1.2.36.146.595.217.0.3",
      "value": "54321"
    },
    {
      "use": "official",
      "type": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v2-0203",
            "code": "MR"
          }
        ]
      },
      "system": "urn:oid:1.2.36.146.595.217.0.1",
      "value": "12345",
      "period": {
        "start": "2001-05-06"
      },
      "assigner": {
        "display": "Acme Healthcare"
      }
    }
  ],
  "active": true,
  "name": [
    {
      "use": "official",
      "family": "Chalmers",
      "given": [
        "Peter",
        "James"
      ]
    },
    {
      "use": "usual",
      "given": [
        "Jim"
      ]
    },
    {
      "use": "maiden",
      "family": "Windsor",
      "given": [
        "Peter",
        "James"
      ],
      "period": {
        "end": "2002"
      }
    }
  ],
  "telecom": [
    {
      "use": "home"
    },
    {
      "system": "phone",
      "value": "(03) 5555 6473",
      "use": "work",
      "rank": 1
    },
    {
      "system": "phone",
      "value": "(03) 3410 5613",
      "use": "mobile",
      "rank": 2
    },
    {
      "system": "phone",
      "value": "(03) 5555 8834",
      "use": "old",
      "period": {
        "end": "2014"
      }
    }
  ],
  "gender": "male",
  "birthDate": "1974-12-25",
  "_birthDate": {
    "extension": [
      {
        "url": "http://hl7.org/fhir/StructureDefinition/patient-birthTime",
        "valueDateTime": "1974-12-25T14:35:45-05:00"
      }
    ]
  },
  "deceasedBoolean": false,
  "address": [
    {
      "use": "home",
      "type": "both",
      "text": "534 Erewhon St PeasantVille, Rainbow, Vic  3999",
      "line": [
        "534 Erewhon St"
      ],
      "city": "PleasantVille",
      "district": "Rainbow",
      "state": "Vic",
      "postalCode": "3999",
      "period": {
        "start": "1974-12-25"
      }
    }
  ],
  "contact": [
    {
      "relationship": [
        {
          "coding": [
            {
              "system": "http://terminology.hl7.org/CodeSystem/v2-0131",
              "code": "N"
            }
          ]
        }
      ],
      "name": {
        "family": "du Marché",
        "_family": {
          "extension": [
            {
              "url": "http://hl7.org/fhir/StructureDefinition/humanname-own-prefix",
              "valueString": "VV"
            }
          ]
        },
        "given": [
          "Bénédicte"
        ]
      },
      "telecom": [
        {
          "system": "phone",
          "value": "+33 (237) 998327"
        }
      ],
      "address": {
        "use": "home",
        "type": "both",
        "line": [
          "534 Erewhon St"
        ],
        "city": "PleasantVille",
        "district": "Rainbow",
        "state": "Vic",
        "postalCode": "3999",
        "period": {
          "start": "1974-12-25"
        }
      },
      "gender": "female",
      "period": {
        "start": "2012"
      }
    }
  ],
  "managingOrganization": {
    "reference": "Organization/1"
  }
}
`
