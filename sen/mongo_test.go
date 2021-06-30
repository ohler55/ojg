// Copyright (c) 2021, Peter Ohler, All rights reserved.

package sen_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/ojg/tt"
)

func TestParserMongo(t *testing.T) {
	p := sen.Parser{}
	p.AddMongoFuncs()

	src := `{
  _id: ObjectId("60c02af61528f028e174d95c"),
  date: ISODate("2021-06-29T02:03:04.005Z"),
  date2: ISODate(1624932184005),
  long: NumberLong("1234567890")
  decimal: NumberDecimal("12345.67890")
  badInt: NumberInt("1234zz")
  badDecimal: NumberDecimal("1e2e3")
}`
	v := p.MustParse([]byte(src)).(map[string]interface{})

	tt.Equal(t, "60c02af61528f028e174d95c", v["_id"])
	tm := time.Unix(0, 1624932184005000000).UTC()
	tt.Equal(t, tm, v["date"])
	tt.Equal(t, tm, v["date2"])
	tt.Equal(t, int64(1234567890), v["long"])
	tt.Equal(t, 12345.67890, v["decimal"])
}
