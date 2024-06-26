package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/agungdhewe/dwpath"
	"github.com/agungdhewe/dwtpl"
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/go-chi/chi/v5"
)

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	ws := GetWebservice()
	setupAllowedAsset(ws)

	ctx := r.Context()
	device := ctx.Value(appsmodel.DeviceKeyName).(appsmodel.Device)

	pathparam := chi.URLParam(r, "*")
	filename := filepath.Base(pathparam)
	extension := filepath.Ext(filename)

	// cek apakah asset boleh diakses
	_, allowed := ws.AllowedAsset[extension]
	if !allowed {
		fmt.Fprintf(w, "Akses ke asset %s tidak diperbolehkan", filename)
		return
	}

	var exist bool
	var path string

	// apakah ada asset yang dibuat untuk mobile/tablet ?
	if device.Type == dwtpl.DeviceMobile || device.Type == dwtpl.DeviceTablet {
		dir := filepath.Dir(pathparam)
		filename = fmt.Sprintf("%s~%s", device.Type, filename)
		path = filepath.Join(ws.RootDir, ws.Configuration.Template.Dir, dir, filename)
		exist, _, _ = dwpath.IsFileExists(path)
		if exist {
			pathparam = filepath.Join(dir, filename)
		}
	}

	// cek apakah asset ada
	path = filepath.Join(ws.RootDir, ws.Configuration.Template.Dir, pathparam)
	exist, _, _ = dwpath.IsFileExists(path)
	if !exist {
		fmt.Fprintf(w, "Asset %s tidak ditemukan", pathparam)
		return
	}

	// muat asset
	// filesize := fileinfo.Size()
	// contenttype := (*fct)[0]
	filedatasource, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(w, "Error memuat asset %s", pathparam)
		return
	}
	defer filedatasource.Close()

	http.ServeFile(w, r, path)

}
