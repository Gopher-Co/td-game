package replay

import (
	"bufio"
	"os"
)

func Save(filename string, w *Watcher) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o666)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	buf := bufio.NewWriter(f)
	if err := w.Write(buf); err != nil {
		return err
	}

	_ = buf.Flush()

	return nil
}
