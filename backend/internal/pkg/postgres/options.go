package postgres

type Options struct {
	driverName string
}

var defaultOptions = Options{
	driverName: "pgx",
}

type OptionFunc func(*Options)
