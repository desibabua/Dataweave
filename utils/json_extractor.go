package utils

import (
	"database/sql"
	"fmt"
	"strconv"
)

func MakeStructJSON(rows sql.Rows) (map[string][]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	masterData := make(map[string][]interface{})

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for i, v := range values {
			x := v.([]byte)
			if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
				masterData[columns[i]] = append(masterData[columns[i]], nx)
			} else if b, ok := strconv.ParseBool(string(x)); ok == nil {
				masterData[columns[i]] = append(masterData[columns[i]], b)
			} else if "string" == fmt.Sprintf("%T", string(x)) {
				masterData[columns[i]] = append(masterData[columns[i]], string(x))
			} else {
				return nil, fmt.Errorf("Failed on if for type %T of %v\n", x, x)
			}
		}
	}

	return masterData, err
}
