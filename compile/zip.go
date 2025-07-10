package compile

import (
	"archive/zip"
	"encoding/json"
	"os"
)

func (self *Artifact) WriteZip(path string) error {
	out_fd, err := os.OpenFile(path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out_fd.Close()

	w := zip.NewWriter(out_fd)
	defer w.Close()

	artifact, err := self.Yaml()
	if err != nil {
		return err
	}

	f, err := w.Create("SQLiteHunter.yaml")
	_, err = f.Write([]byte(artifact))
	if err != nil {
		return err
	}

	serialized, err := json.MarshalIndent(self.Spec, " ", " ")
	if err != nil {
		return err
	}

	f, err = w.Create("spec.json")
	_, err = f.Write(serialized)
	if err != nil {
		return err
	}

	f, err = w.Create("definitions.json")
	_, err = f.Write(self.BuildIndex())
	if err != nil {
		return err
	}

	return err
}
