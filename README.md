# httputil

helpers to create http rest services

## features

- bind json body
- bind form body
- bind xml body
- bind query values
- bind `url.Values`
- custom **error response writer**
- write json response
- write xml response
- write text response
- write blob response
- pprof handler and middleware

## examples

- handler with json body

```go
var ErrBadRequest = httputil.NewHTTPError(http.StatusBadRequest)

func CreateTaskHandler(svc TaskService) http.HandlerFunc {
 return httputil.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
  var req CreateTaskRequest
  if err := httputil.BindJSON(r, &req); err != nil {
   return ErrBadRequest.WithCause(err)
  }

  res, err := svc.Create(r.Context(), req)
  if err != nil {
   return err
  }

  return httputil.WriteJSON(w, http.StatusCreated, res)
 })
}
```

- handler with query params

```go
func SearchTaskHandler(svc TaskService) http.HandlerFunc {
 return httputil.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
  var req SearchTaskRequest
  if err := httputil.BindQuery(r, &req, httputil.WithTagName("json")); err != nil {
   return ErrBadRequest.WithCause(err)
  }

  res, err := svc.Search(r.Context(), req)
  if err != nil {
   return err
  }

  return httputil.WriteJSON(w, http.StatusOK, res)
 })
}
```

- custom error writer

```go

var errorWriter = httputil.ErrorWriterFunc(func(w http.ResponseWriter, err error) {
 _ = httputil.WriteString(w, http.StatusInternalServerError, err.Error())
})

func CreateTaskHandler(svc TaskService) http.HandlerFunc {
 return httputil.NewHandler(func(w http.ResponseWriter, r *http.Request) error {
  // ...
 }, httputil.WithErrorWriter(errorWriter))
}

```

- pprof handler

```go
r := chi.NewRouter()
r.Handle("/debug/*", pprof.Handler())
```

- pprof middleware

```go
r := chi.NewRouter()
r.Use(pprof.Middleware)
```
