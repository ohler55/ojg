// Copyright (c) 2021, Peter Ohler, All rights reserved.

/*

Package asm provides a means of building JSON or simple types using JSON
encoded scripts. The assembly scripts are encapsuled in the Plan type.

An assembly plan is described by a JSON document or a SEN document. The format
is much like LISP but with brackets instead of parenthesis. A plan is
evaluated by evaluating the plan function which is usually an 'asm'
function. The plan operates on a data map which is the root during
evaluation. The source data is in the $.src and the expected assembled output
should be in $.asm.

An example of a plan in SEN format is (the first asm is optional):

  [ asm
    [set $.asm {good: bye}]  // set output to {good: bad}
    [set $.asm.hello world]  // output is now {good: bad, hello: world}
  ]

The functions available are:

          !=: Returns true if any the argument are not equal. An alias is !==.

           *: Returns the product of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised.

           +: Returns the sum of all arguments. All arguments must be numbers
              or strings. If any argument is a string then the result will be
              a string otherwise the result will be a number. If any of the
              arguments are not a number or a string an error is raised.

           -: Returns the difference of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised.

           /: Returns the quotient of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised. If an attempt is made to divide by zero and error will
              be raised.

           <: Returns true if each argument is less than any subsequent
              argument. An alias is lt.

          <=: Returns true if each argument is less than or equal to any
              subsequent argument. An alias is lte.

          ==: Returns true if all the argument are equal. Aliases are eq, ==,
              and equal.

           >: Returns true if each argument is greater than any subsequent
              argument. An alias is gt.

          >=: Returns true if each argument is greater than or equal to any
              subsequent argument. An alias is gte.

         and: Returns true if all argument evaluate to true. Any arguments
              that do not evaluate to a boolean or null (false) raise an error.

      append: Appends the second argument to the first argument which must be
              an array.

      array?: Returns true if the single required argumement is an array
              otherwise false is returned.

         asm: Processes all arguments in order using the return of each as
              input for the next.

          at: Forms a path starting with @. The remaining string arguments are
              joined with a '.' and parsed to form a jp.Expr.

       bool?: Returns true if the single required argumement is a boolean
              otherwise false is returned.

        cond: A conditional construct modeled after the LISP cond. All
              arguments must be array of two elements. The first element must
              evaluate to a boolean and the second can be any value. The value
              of the first true first argument is returned. If none match nil
              is returned.

         del: Deletes the first matching value in either the root ($) or
              local (@) data. Exactly one argument is required and it must be
              a path. The jp.DelOne() function is used to delete the value.
              The local (@) value is returned.

      delall: Deletes the all matching values in either the root ($) or
              local (@) data. Exactly one argument is required and it must be
              a path. The jp.DelOne() function is used to delete the value.
              The local (@) value is returned.

         dif: Returns the difference of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised.

        each: Each .

          eq: Returns true if all the argument are equal. Aliases are eq, ==,
              and equal.

       equal: Returns true if all the argument are equal. Aliases are eq, ==,
              and equal.

       float: Converts a value into a float if possible. I no conversion is
              possible nil is returned.

         get: Gets the first matching value in either the root ($), local (@),
              or if present, the second argument. The required first argument
              must be a path and the option second argument is the
              data to apply the path to. The jp.First() function is used to
              get the results

      getall: Gets all matching values in either the root ($), or local (@),
              or if present, the second argument. The required first argument
              must be a path and the option second argument is the
              data to apply the path to. The jp.Get() function is used to get
              the results

          gt: Returns true if each argument is greater than any subsequent
              argument. An alias is >.

         gte: Returns true if each argument is greater than or equal to any
              subsequent argument. An alias is >=.

     include: Returns true if a list first argument includes the second
              argument. It will also return true if the first argument is a
              string and the second string argument is included in the first.

     inspect: Print the arguments as JSON unless the argument is an integer.
              Integers are assumed to be the indentation for the arguments
              that follow.

         int: Converts a value into a integer if possible. I no conversion is
              possible nil is returned.

        join: Join an array of strings with the provided separator. If a
              separator is not provided as the second argument then an empty
              string is used.

        list: Creates a list from all the argument and return that list.

          lt: Returns true if each argument is less than any subsequent
              argument. An alias is <.

         lte: Returns true if each argument is less than or equal to any
              subsequent argument. An alias is <=.

        map?: Returns true if the single required argumement is a map
              otherwise false is returned.

         mod: Returns the remainer of a modulo operation on the first two
              argument. Both arguments must be integers and are both required.
              An error is raised if the wrong argument types are given.

         neq: Returns true if any the argument are not equal. An alias is !==.

        nil?: Returns true if the single required argumement is null (JSON)
              or nil (golang) otherwise false is returned.

         not: Returns the boolean NOT of the argument. Exactly one argument
              is expected and it must be a boolean.

         nth: Returns a nth element of an array. The second argument must be
              an integer that indicates the element of the array to return.
              If the index is less than 0 then the index is from the end of
              the array.

       null?: Returns true if the single required argumement is null (JSON)
              or nil (golang) otherwise false is returned.

        num?: Returns true if the single required argumement is number
              otherwise false is returned.

          or: Returns true if any of the argument evaluate to true. Any
              arguments that do not evaluate to a boolean or null (false)
              raise an error.

     product: Returns the product of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised.

       quote: Does not evaluate arguments. One argument is expected. Null is
              returned if no arguments are given while any arguments other
              than the first are ignored. An example for use would be to
              treats "@.x" as a string instead of as a path.

    quotient: Returns the quotient of all arguments. All arguments must be
              numbers. If any of the arguments are not a number an error is
              raised. If an attempt is made to divide by zero and error will
              be raised.

     replace: Replace an occurrences the second argument with the third
              argument. All three arguments must be strings.

     reverse: Reverse the items in an array and return a copy of it.

        root: Forms a path starting with @. The remaining string arguments are
              joined with a '.' and parsed to form a jp.Expr.

         set: Sets a single value in either the root ($) or local (@) data. Two
              arguments are required, the first must be a path and the second
              argument is evaluate to a value and inserted using the
              jp.SetOne() function.

      setall: Sets multiple values in either the root ($) or local (@) data.
              Two arguments are required, the first must be a path and the
              second argument is evaluate to a value and inserted using the
              jp.Set() function.

        size: Returns the size or length of a string, array, or object (map).
              For all other types zero is returned

        sort: Sort the items in an array and return a copy of the array. Valid
              types for comparison are strings, numbers, and times. Any other
              type returned or a type mismatch will raise an error.

       split: Split a string on using a specified separator.

      string: Converts a value into a string.

     string?: Returns true if the single required argumement is a string
              otherwise false is returned.

      substr: Returns a substring of the input string. The second argument
              must be an integer that marks the start of the substring. The
              third integer argument indicates the length of the substring
              if provided. If the length argument is not provided the end of
              the substring is the end of the input string.

         sum: Returns the sum of all arguments. All arguments must be numbers
              or strings. If any argument is a string then the result will be
              a string otherwise the result will be a number. If any of the
              arguments are not a number or a string an error is raised.

        time: Converts the first argument to a time if possible otherwise
              an error is raised. The first argument can be a integer, float,
              or string and are converted as follows:
                integer < 10^10:  time in seconds since 1970-01-01 UTC
                integer >= 10^10: time in nanoseconds 1970-01-01 UTC
                decimal (float):  time in seconds 1970-01-01 UTC
                string:           assumed to be formated as RFC3339 unless a
                                  format argument is provided

       time?: Returns true if the single required argumement is a time
              otherwise false is returned.

       title: Convert a string to capitalized string. There must be exactly
              one string argument.

     tolower: Convert a string to lowercase. There must be exactly one
              string argument.

     toupper: Convert a string to uppercase. There must be exactly one
              string argument.

        trim: Trim white space from both ends of a string unless a second
              argument provides an alternative cut set.

        zone: Changes the timezone on a time to the location specified in the
              second argument. Raises an error if the first argument does not
              evaluate to a time or the location can not be determined.
              Location can be either a string or the number of minutes offset
              from UTC.

*/
package asm
