package lib

import (
	"html/template"
	"io"
)

// CreateReport creates html report from HashList and
// saves html file to io.Writer w.
func CreateReport(w io.Writer, pathList PathList, result HashList) error {
	report := template.New("report")
	b, err := Asset("lib/templates/index.html")
	if err != nil {
		return err
	}

	report, err = report.Funcs(
		template.FuncMap{
			"passFail": passFail,
			"rowAttr":  rowAttr,
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			}}).Parse(string(b))
	if err != nil {
		return err
	}

	dir, file := result.GetDirectoryInfo()
	type DirectoryInfo struct {
		Directories int
		Files       int
	}
	values := struct {
		PathList      PathList
		Result        HashList
		DirectoryInfo DirectoryInfo
	}{
		PathList: pathList,
		Result:   result,
		DirectoryInfo: DirectoryInfo{
			Directories: dir,
			Files:       file,
		},
	}
	err = report.Execute(w, values)
	if err != nil {
		return err
	}
	return nil
}

func passFail(result bool) string {
	var message string
	if result {
		message = "Pass"
	} else {
		message = "Fail"
	}
	return message
}

func rowAttr(result bool) template.HTMLAttr {
	var attr template.HTMLAttr
	if result {
		attr = template.HTMLAttr("success")
	} else {
		attr = template.HTMLAttr("danger")
	}
	return attr
}
