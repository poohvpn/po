package po

import "github.com/hashicorp/go-multierror"

func Errors(i ...interface{}) error {
	errs := new(multierror.Error)
	for _, v := range i {
		if e, ok := v.(error); ok && e != nil {
			errs = multierror.Append(errs, e)
		}
	}
	return errs.ErrorOrNil()
}
