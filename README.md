# [![Actions Status](https://github.com/SarthakMakhija/goselect/workflows/GoSelectCI/badge.svg)](https://github.com/SarthakMakhija/goselect/actions) goselect
SQL like 'select' interface for files

# Some examples
- select fName from /home/apps
- select fSize, fName from /home/apps where fSize > 10MB

# Initial thoughts

The idea is to provide support for `selecting` files based on various conditions. At this stage some features that are planned include:
- Support for `where` clause
  - [X] select * from /home/apps where eq(add(2,3), 5)
  - [X] select * from /home/apps where eq(lower(ext), .log)
  - [X] select * from /home/apps where ne(lower(ext), .log)
  - [X] `where` clause should support functions for comparison like `eq`, `le`, `lt`, `ge`, `gt`, `ne`, `contains`, `or`, `and` and `not`
- Support for searching text within a file
- Support for some shortcuts for common files types like images, videos, text:
  - select * from /home/apps where fileType = 'image' and fileSize > 10MB
  - select * from /home/apps where fileType = 'text' and textContains = 'get('
- Support for projections
  - [X] projections with attribute name: `name`, `size`
  - [X] projections with alias in attribute name: `fName` instead of `name`
  - [X] projections with scalar functions: `contains`, `lower`
  - [X] projections with alias in scalar functions: `low` instead of `lower`
  - [ ] projections with aggregate functions: `min`, `max`
  - [X] projections with expression: `1 + 2`
    - supports by giving functions like `add`, `sub`, `mul` and `div`
- Support for `order by` clause
  - [X] order by with positions: `order by 1`
  - [X] order by with descending order: `order by 1 desc`
  - [X] order by with optional ascending order: `order by 1 asc`
- Support for `limit` clause
  - [X] limit clause with a value: `limit 10`
- Support for `aggregation functions`
  - [X] min
  - [X] max
  - [X] avg
  - [X] sum
  - [X] count
  - [ ] median
  - [X] countDistinct
- Support for various `scalar functions`
  - [X] add
  - [X] subtract
  - [X] multiply
  - [X] divide
  - [X] equal (eq)
  - [X] lessThan (lt)
  - [X] greaterThan (gt)
  - [X] lessEqual (le)
  - [X] greaterEqual (ge)
  - [X] notEqual (ne)
  - [X] or
  - [X] and 
  - [X] not 
  - [X] like 
  - [X] lower
  - [X] upper
  - [X] title
  - [X] base64
  - [X] length
  - [X] lTrim
  - [X] rTrim
  - [X] trim
  - [X] now
  - [X] date
  - [X] day
  - [X] month
  - [X] year
  - [X] dayOfWeek
  - [X] working directory (wd)
  - [X] concat
  - [X] concat with separator (concatWs)
  - [X] contains
  - [X] substr
  - [X] replace
  - [X] replaceAll
  - [ ] formatSize
- Support for formatting the results
  - [X] Json formatter
  - [X] Html formatter
  - [X] Table formatter
- Support for exporting the formatted result
  - [X] Console
  - [X] File
- Design consideration for searching files in a directory with a huge number of files
