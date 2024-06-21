package updater

// Updater is a struct that represents an updater.
type Updater struct {
	funcs []func()
}

func (u *Updater) Append(f func()) {
	u.funcs = append(u.funcs, f)
}

func (u *Updater) Update() {
	for _, f := range u.funcs {
		f()
	}
}
