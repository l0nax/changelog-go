package errors

/// -=-=-=-=-=-=-=-=-=] EntryTypeExistsError [-=-=-=-=-=-=-=-=-=

// EntryTypeExistsError occurs when the Application tries to add/ register a new
// Entry Type but it already Exists a equal Entry Type
type EntryTypeExistsError struct { message string }

func NewEntryTypeExistsError() *EntryTypeExistsError {
	return &EntryTypeExistsError{
		message: "Trying to register a new Entry Type, but it already exists. (No more information's :( )",
	}
}

func NewEntryTypeExistsErrorMsg(message string) *EntryTypeExistsError {
	return &EntryTypeExistsError{
		message: message,
	}
}

func (e *EntryTypeExistsError) Error() string {
	return e.message
}

