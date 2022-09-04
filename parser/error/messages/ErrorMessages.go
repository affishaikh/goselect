package messages

const (
	ErrorMessageEmptyQuery                                = "expected query to be non-empty"
	ErrorMessageNonSelectQuery                            = "expected a select query statement"
	ErrorMessageLimitValue                                = "expected a limit value"
	ErrorMessageLimitValueInt                             = "expected limit to be a positive integer"
	ErrorMessageLimitValueIntWithExistingError            = "expected limit to be a positive integer, %v"
	ErrorMessageMissingBy                                 = "expected 'by' after order"
	ErrorMessageMissingCommaOrderBy                       = "expected a comma after 'order by' in attribute positions"
	ErrorMessageMissingOrderByAttributes                  = "expected an attribute position after 'order by'. attribute positions start with 1"
	ErrorMessageNonZeroPositivePositions                  = "expected non-zero & positive 'order by' positions"
	ErrorMessageNonZeroPositivePositionsWithExistingError = "expected non-zero & positive 'order by' positions, %v"
	ErrorMessageOrderByPositionOutOfRange                 = "expected 'order by' position to be between %v and %v, both inclusive"
	ErrorMessageMissingSource                             = "expected a source path after 'from`"
	ErrorMessageInaccessibleSource                        = "expected directory path %v to exist. please check the path, also ensure that it is accessible"
	ErrorMessageMissingCommaProjection                    = "expected a comma in the projection list. please check the spellings, supported attributes and supported functions as well"
	ErrorMessageOpeningParenthesesProjection              = "expected an opening parentheses in the projection list"
	ErrorMessageInvalidProjection                         = "invalid projection list, please check the opening and closing parentheses for all the functions"
	ErrorMessageExpectedExpressionInProjection            = "expected atleast one expression in the projection list"
	ErrorMessageExpectedExpressionInWhere                 = "expected one expression in the where clause, or remove 'where' keyword"
	ErrorMessageExpectedSingleExpressionInWhere           = "expected only a single expression in the where clause"
	ErrorMessageInvalidWhere                              = "invalid where clause, please check opening and closing parentheses for all the functions"
	ErrorMessageInvalidWhereFunctionUsed                  = "invalid where clause, 'where' supports only functions. please check all the function supported in 'where' clause"
	ErrorMessageAggregateFunctionInsideWhere              = "invalid where clause, aggregate functions are not supported in the where clause"
	ErrorMessageMissingParameterInScalarFunctions         = "expected %v parameters in the function %v but did not receive all"
	ErrorMessageIncorrectValueType                        = "expected a %v value type but was not"
	ErrorMessageIncorrectEndIndexInSubstring              = "expected the end index to be greater than the from index in substr"
	ErrorMessageIllegalFromToIndexInSubstring             = "expected the from and to index to be positive integers"
	ErrorMessageExpectedNumericArgument                   = "expected numeric type argument value"
	ErrorMessageExpectedNonZeroInDivide                   = "expected a non zero denominator in divide operation"
	ErrorMessageFunctionNamePrefixWithExistingError       = "[Function %v], %s"
)
