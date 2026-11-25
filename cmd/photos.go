package cmd

import (
 "crypto/rand"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "log"
 "net/http"
 "os"
 "path/filepath"
 "strconv"
 "strings"
 "inventory/api/providers/database_provider"
 "inventory/pkg/models"
)

func photoDir() string {
 if config != nil && config.PhotosDir != "" { return config.PhotosDir }
 return "photos"
}
const maxPhotoBytes int64 = 10 << 20

func photoItemID(r *http.Request) (int,error) {
 id,err:=strconv.Atoi(r.PathValue("id"))
 if err!=nil || id<=0 {return 0,errors.New("invalid item ID")}
 return id,nil
}

func listPhotosHandler(provider database_provider.DatabaseProviderInterface) http.HandlerFunc {
 return func(w http.ResponseWriter,r *http.Request) {
  id,err:=photoItemID(r);if err!=nil {http.Error(w,err.Error(),400);return}
  photos,err:=provider.ListPhotos(id);if err!=nil {http.Error(w,"unable to list photos",500);return}
  if photos==nil {photos=[]models.ItemPhoto{}}
  w.Header().Set("Content-Type","application/json; charset=utf-8")
  json.NewEncoder(w).Encode(photos)
 }
}

func uploadPhotoHandler(provider database_provider.DatabaseProviderInterface) http.HandlerFunc {
 return func(w http.ResponseWriter,r *http.Request) {
  id,err:=photoItemID(r);if err!=nil {http.Error(w,err.Error(),400);return}
  // Check the item exists before accepting its photo.
  item,err:=provider.Get(id);if err!=nil || item.ID==0 {http.Error(w,"item not found",404);return}
  r.Body=http.MaxBytesReader(w,r.Body,maxPhotoBytes+(1<<20))
  if err=r.ParseMultipartForm(maxPhotoBytes);err!=nil {http.Error(w,"invalid or oversized upload (10 MB max)",400);return}
  file,_,err:=r.FormFile("photo");if err!=nil {http.Error(w,"photo is required",400);return};defer file.Close()
  header:=make([]byte,512);n,err:=io.ReadFull(file,header);if err!=nil && err!=io.ErrUnexpectedEOF && err!=io.EOF {http.Error(w,"unable to read photo",400);return}
  kind:=http.DetectContentType(header[:n])
  extensions:=map[string]string{"image/jpeg":".jpg","image/png":".png","image/gif":".gif","image/webp":".webp"}
  ext,ok:=extensions[kind];if !ok {http.Error(w,"JPEG, PNG, GIF or WebP required",400);return}
  if _,err=file.Seek(0,io.SeekStart);err!=nil {http.Error(w,"unable to read photo",500);return}
  if err=os.MkdirAll(photoDir(),0700);err!=nil {http.Error(w,"unable to create photo directory",500);return}
  b:=make([]byte,16);if _,err=rand.Read(b);err!=nil {http.Error(w,"unable to create filename",500);return}
  filename:=hex.EncodeToString(b)+ext
  path:=filepath.Join(photoDir(),filename)
  out,err:=os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600);if err!=nil {http.Error(w,"unable to save photo",500);return}
  size,copyErr:=io.Copy(out,io.LimitReader(file,maxPhotoBytes+1));closeErr:=out.Close()
  if copyErr!=nil || closeErr!=nil || size>maxPhotoBytes || size==0 {
   os.Remove(path);http.Error(w,"unable to save photo (10 MB max)",400);return
  }
  photo,err:=provider.SavePhoto(&models.ItemPhoto{ItemID:id,Filename:filename})
  if err!=nil {os.Remove(path);log.Printf("save photo record: %v",err);http.Error(w,"unable to save photo record",500);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8");w.WriteHeader(http.StatusCreated);json.NewEncoder(w).Encode(photo)
 }
}

func servePhotoHandler(provider database_provider.DatabaseProviderInterface) http.HandlerFunc {
 return func(w http.ResponseWriter,r *http.Request) {
  id,err:=photoItemID(r);if err!=nil {http.Error(w,err.Error(),400);return}
  name:=r.PathValue("filename")
  if name=="" || name!=filepath.Base(name) || strings.Contains(name,"..") {http.NotFound(w,r);return}
  photos,err:=provider.ListPhotos(id);if err!=nil {http.Error(w,"unable to list photos",500);return}
  found:=false;for _,p:=range photos {if p.Filename==name {found=true;break}}
  if !found {http.NotFound(w,r);return}
  w.Header().Set("X-Content-Type-Options","nosniff")
  http.ServeFile(w,r,filepath.Join(photoDir(),name))
 }
}

func deletePhotoHandler(provider database_provider.DatabaseProviderInterface) http.HandlerFunc {
 return func(w http.ResponseWriter,r *http.Request) {
  id,err:=photoItemID(r);if err!=nil {http.Error(w,err.Error(),400);return}
  photoID,err:=strconv.Atoi(r.PathValue("photoID"));if err!=nil || photoID<=0 {http.Error(w,"invalid photo ID",400);return}
  photos,err:=provider.ListPhotos(id);if err!=nil {http.Error(w,"unable to list photos",500);return}
  var filename string;for _,p:=range photos {if p.ID==photoID {filename=p.Filename;break}}
  if filename=="" {http.NotFound(w,r);return}
  if err=provider.DeletePhoto(photoID);err!=nil {http.Error(w,"unable to delete photo",500);return}
  if err=os.Remove(filepath.Join(photoDir(),filepath.Base(filename)));err!=nil && !os.IsNotExist(err) {log.Printf("remove photo file %s: %v",filename,err)}
  w.WriteHeader(http.StatusNoContent)
 }
}
