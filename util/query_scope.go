package util

import (
	"strings"
	"github.com/jinzhu/gorm"
	"fmt"
)

type Parameter struct {
	Raw 	string
	Key 	string
	Operate Operator
	Function OperateFunction
	Value 	string
}

func (param Parameter) IsOperator() bool {
	return param.Operate != "" && param.Function == ""
}

func (param Parameter) IsFunction() bool {
	return param.Function != "" && param.Operate == ""
}

func (param Parameter) IsQueryOperatable() bool {
	return param.Operate != "" || param.Function != ""
}

type ParameterMapper map[string][]*Parameter


type QueryOperatable interface {
	Must(string) bool
}

type Operator string

func (operator Operator) Must(paramString string) bool {
	if len(strings.Split(paramString, string(operator))) == 2 {
		return true
	}
	return false
}

func (operator Operator) Is(otherOp Operator) bool {
	return string(operator) == string(otherOp)
}

type OperateFunction string

func (operatorFunc OperateFunction) Must(paramString string) bool {
	if OPERATOR_EQUAL.Must(paramString) {
		if strings.Split(paramString, "=")[0] == string(operatorFunc) {
			return true
		}
	}
	return false
}

const (
	OPERATOR_EQUAL         Operator = "="
	OPERATOR_NOT_EQUAL              = "_ne_"
	OPERATOR_GREATER_THAN           = "_gt_"
	OPERATOR_GREATER_EQUAL          = "_ge_"
	OPERATOR_LESS_THAN              = "_lt_"
	OPERATOR_LESS_EQUAL             = "_le_"
	OPERATOR_NOT_IN                 = "_ni_"		   //example: type<ni>dressing_room,zone
	OPERATOR_IN                     = "_in_"		   //example: status<in>created,published
	OPERATOR_LIKE 					= "_lk_" 		   //example: name<lk>%线下店%

	OPERATOR_FUNCTION_ORDER OperateFunction = "order"  //example: order=+price,-publish_date
	OPERATOR_FUNCTION_FIELD_SCANNER = "fields"		   //example: fields=type,status,name,title
	OPERATOR_FUNCTION_PAGINATION = "flag"			   //example: flag=1_100_18
															// meaning: version_(flag_format), flag_format example: offset_limit
)

var (
	currentSupportedOperators []Operator = []Operator {
		OPERATOR_EQUAL,
		OPERATOR_NOT_EQUAL,
		OPERATOR_GREATER_THAN,
		OPERATOR_GREATER_EQUAL,
		OPERATOR_LESS_THAN,
		OPERATOR_LESS_EQUAL,
		OPERATOR_NOT_IN,
		OPERATOR_IN,
		OPERATOR_LIKE,
	}
	currentSupportedFunctions []OperateFunction = []OperateFunction {
		OPERATOR_FUNCTION_ORDER,
		OPERATOR_FUNCTION_FIELD_SCANNER,
		OPERATOR_FUNCTION_PAGINATION,
	}
)

func determineOperatable(paramString string) QueryOperatable {
	for _, operator := range currentSupportedOperators {
		if operator.Must(paramString) {
			if operator.Is(OPERATOR_EQUAL) {
				for _, function := range currentSupportedFunctions {
					if function.Must(paramString) {
						return function
					}
				}
			}
			return operator
		}
	}
	return nil
}

func parseParam(paramString string) *Parameter {
	operator := determineOperatable(paramString)
	switch operator.(type) {
	case Operator:
		values := strings.Split(paramString, string(operator.(Operator)))
		return &Parameter{
			Raw: paramString,
			Operate: operator.(Operator),
			Key: values[0],
			Value: values[1],
		}
	case OperateFunction:
		values := strings.Split(paramString, string(OPERATOR_EQUAL))
		return &Parameter{
			Raw: paramString,
			Function: operator.(OperateFunction),
			Key: values[0],
			Value: values[1],
		}
	}
	return nil
}

func ParseQuery(db *gorm.DB, query string) *gorm.DB {
	mapper := parseQuery(query)
	for _, params := range mapper {
		for _, param := range params {
			db = bindParameter(db, param)
		}
	}
	return db
}

func parseQuery(query string) ParameterMapper {
	mapper := ParameterMapper{}
	paramStringList := strings.Split(query, "&")
	for _, paramString := range paramStringList {
		param := parseParam(paramString)
		if param == nil {
			continue
		}
		key := param.Key
		mapper[key] = append(mapper[key], param)
	}
	return mapper
}

func bindParameter(db *gorm.DB, param *Parameter) *gorm.DB {
	if param.IsQueryOperatable() {
		if param.IsFunction() {
			switch param.Function {
			case OPERATOR_FUNCTION_ORDER:
				orderConditions := strings.Split(param.Value, ",")
				for _, condition := range orderConditions {
					if strings.HasPrefix(condition, "-") {
						db = db.Order(fmt.Sprintf("%s %s", strings.TrimPrefix(condition, "-"), "desc"))
					} else if strings.HasPrefix(condition, "+") {
						db = db.Order(fmt.Sprintf("%s %s", strings.TrimPrefix(condition, "+"), "asc"))
					}
				}
			case OPERATOR_FUNCTION_PAGINATION:
				paginationVersion := strings.Split(param.Value, "_")[0]
				if paginationVersion == "1" {
					parts := strings.Split(param.Value, "_")
					db = db.Offset(parts[1]).Limit(parts[2])
				}
			case OPERATOR_FUNCTION_FIELD_SCANNER:
				db = db.Select(strings.Split(param.Value, ","))
			}
		} else if param.IsOperator() {
			switch param.Operate {
			case OPERATOR_EQUAL:
				db = db.Where(fmt.Sprintf("%s = ?", param.Key), param.Value)
			case OPERATOR_GREATER_THAN:
				db = db.Where(fmt.Sprintf("%s > ?", param.Key), param.Value)
			case OPERATOR_GREATER_EQUAL:
				db = db.Where(fmt.Sprintf("%s >= ?", param.Key), param.Value)
			case OPERATOR_LESS_THAN:
				db = db.Where(fmt.Sprintf("%s < ?", param.Key), param.Value)
			case OPERATOR_LESS_EQUAL:
				db = db.Where(fmt.Sprintf("%s <= ?", param.Key), param.Value)
			case OPERATOR_LIKE:
				db = db.Where(fmt.Sprintf("%s like ?", param.Key), fmt.Sprint("%", param.Value, "%"))
			case OPERATOR_IN:
				db = db.Where(fmt.Sprintf("%s in (?)", param.Key),  strings.Split(param.Value, ","))
			case OPERATOR_NOT_EQUAL:
				fallthrough
			case OPERATOR_NOT_IN:
				db = db.Not(param.Key, strings.Split(param.Value, ","))
			}
		}
	}
	return db
}
