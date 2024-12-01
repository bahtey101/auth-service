package dbconverter

import "fmt"

func ConvertToString(dest *string, src any) error {
	if dest == nil {
		return fmt.Errorf("error nil pointer")
	}

	switch s := src.(type) {
	case string:
		*dest = s
		return nil
	case []byte:
		*dest = string(s)
		return nil
	default:
		return fmt.Errorf("unknown type for string")
	}
}
