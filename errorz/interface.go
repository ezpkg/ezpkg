package errorz

type zUnwrapI interface{ Unwrap() error }
type zUnwraps interface{ Unwrap() []error }
type zErrorsI interface{ Errors() []error }
type zWrErrorsI interface{ WrappedErrors() []error }
