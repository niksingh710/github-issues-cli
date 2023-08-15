// This package handles everything related to the editor.
package editor

import (
	"os"
	"os/exec"
)

const defaultEditor = "vim"

// This gives the editor name if GIT_EDITOR variable is setted in user's env
func getEditor() string {
	editor := os.Getenv("GIT_EDITOR")
	if editor == "" {
		editor = defaultEditor
	}
	return editor
}

func openFileWithEditor(filename string) error {
	editor := getEditor()
	bin, err := exec.LookPath(editor)
	if err != nil {
		return err
	}
	cmd := exec.Command(bin, filename)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func Edit(content string) ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "github-issue-edit-")
	if err != nil {
		return []byte{}, err
	}
	filename := file.Name()
	defer os.Remove(filename)
	file.WriteString(content)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}
	if err = openFileWithEditor(filename); err != nil {
		return []byte{}, err
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
