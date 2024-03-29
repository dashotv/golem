package tasks

import (
	"bytes"
	"os"
	"regexp"

	"github.com/pkg/errors"

	"github.com/dashotv/golem/templates"
)

func File(template, output string, data interface{}) error {
	buf := bytes.NewBufferString("")
	err := templates.New(template).Execute(buf, data)
	if err != nil {
		return errors.Wrap(err, "execute template")
	}

	err = os.WriteFile(output, buf.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "writing template output")
	}

	return nil
}

func FileDoesntExist(template, output string, data interface{}) error {
	if PathExists(output) {
		return nil
	}
	return File(template, output, data)
}

func RawFile(output string, data string) error {
	err := os.WriteFile(output, []byte(data), 0644)
	if err != nil {
		return errors.Wrap(err, "writing raw output")
	}

	return nil
}

func AppendFile(template, output string, data interface{}) error {
	buf, err := Buffer(template, data)
	if err != nil {
		return errors.Wrap(err, "execute template")
	}

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}
	defer f.Close()

	if _, err := f.Write([]byte(buf)); err != nil {
		return errors.Wrap(err, "writing to file")
	}

	return nil
}

func Buffer(template string, data interface{}) (string, error) {
	buf := bytes.NewBufferString("")
	err := templates.New(template).Execute(buf, data)
	if err != nil {
		return "", errors.Wrap(err, "execute template")
	}

	return buf.String(), nil
}

func Modify(output string, data interface{}) error {
	rx, err := regexp.Compile(`//golem:template:(.*)`)
	if err != nil {
		return errors.Wrap(err, "compiling regex")
	}

	file, err := os.ReadFile(output)
	if err != nil {
		return errors.Wrap(err, "reading file")
	}

	if matches := rx.FindAllStringSubmatch(string(file), -1); matches != nil {
		found := map[string]bool{}
		for _, match := range matches {
			if found[match[0]] {
				continue
			}

			buf, err := Buffer(match[1], data)
			if err != nil {
				return errors.Wrap(err, "execute template")
			}

			re, err := regexp.Compile(`(?ms)` + match[0] + `.*` + match[0])
			if err != nil {
				return errors.Wrap(err, "compiling replace regex")
			}

			file = re.ReplaceAll(file, []byte(buf))

			found[match[0]] = true
		}
	}

	if err := os.WriteFile(output, file, 0644); err != nil {
		return errors.Wrap(err, "writing contents")
	}

	return nil
}
