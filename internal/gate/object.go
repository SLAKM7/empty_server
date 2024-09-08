package gate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Upload(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	if header == nil {
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	err = os.WriteFile(header.Filename, data, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"code":    "0000",
		"message": "上传成功",
		"data":    header.Filename,
	}
	rb, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	w.Write(rb)
	return
}

func Download(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	n := pathParams["name"]
	fi, err := os.Stat(n)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	data, err := os.ReadFile(n)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fi.Name())
	w.Write(data)
	return

}
