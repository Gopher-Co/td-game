package updater

// Updater is a struct that represents an updater.
type Updater struct {
	funcs []func()
}

// Append appends a function to the updater.
func (u *Updater) Append(f func()) {
	u.funcs = append(u.funcs, f)
}

// Update updates the updater.
func (u *Updater) Update() {
	for _, f := range u.funcs {
		f()
	}
}
