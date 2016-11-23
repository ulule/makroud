package sqlxx

// Preloader is a custom preloader.
type Preloader func(d Driver) (Driver, error)

// Preload preloads related fields.
func Preload(driver Driver, out Model, related ...string) error {
	return nil
}

// PreloadFuncs preloads with the given preloader functions.
func PreloadFuncs(driver Driver, out Model, preloaders ...Preloader) error {
	return nil
}
