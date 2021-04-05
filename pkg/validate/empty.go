package validate

import (
	"fmt"

	"github.com/spf13/viper"
)

const errEmpty = "%s cannot be empty"

// NotEmpty will validate that the given keys with the viper objects
// are not empty
func NotEmpty(v *viper.Viper, keys []string) error {
	if v == nil {
		return fmt.Errorf("viper cannot be nil")
	}
	var err error
	for _, key := range keys {
		switch t := v.Get(key).(type) {
		case string:
			if t == "" {
				err = constructError(err, fmt.Sprintf(errEmpty, key))
				continue
			}
		case int:
			if t == 0 {
				err = constructError(err, fmt.Sprintf(errEmpty, key))
				continue
			}
		case []string:
			if len(t) == 0 {
				err = constructError(err, fmt.Sprintf(errEmpty, key))
				continue
			}
		case []int:
			if len(t) == 0 {
				err = constructError(err, fmt.Sprintf(errEmpty, key))
				continue
			}
		default:
			err = constructError(err, fmt.Sprintf("unsupported type %T", t))
			continue
		}
	}
	return err
}

func constructError(err error, newErrorMessage string) error {
	if err != nil {
		return fmt.Errorf("%w: "+newErrorMessage, err)
	}
	return fmt.Errorf(newErrorMessage)
}

func in(s string, keys []string) bool {
	for _, key := range keys {
		if s == key {
			return true
		}
	}
	return false
}
