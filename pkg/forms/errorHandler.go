package forms

type errors map[string][]string

func (err errors) Add(formField, errDesc string) {
	err[formField] = append(err[formField], errDesc)
}

func (err errors) Get(formField string) string {
	if len(err[formField]) == 0 {
		return ""
	}

	return err[formField][0]
}
