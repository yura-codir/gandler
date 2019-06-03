package handler

type Storage interface {
  Get(key string) interface{}
  Set(key string, value interface{})
}

type Handler interface {
  Run(handlers ...func(h Storage) error) Handler
  OnFail(handlers ...func(h Storage, err error)) Handler
  OnSuccess(handlers ...func(h Storage)) Handler
}

type MapStorage struct {
  storage map[string]interface{}
}

func (m MapStorage) Get(key string) interface{} {
  return m.storage[key]
}
func (m MapStorage) Set(key string, value interface{}) {
  m.storage[key] = value
}

type SimpleHandler struct {
  Storage Storage
  Error   error
}

func (h *SimpleHandler) Run(handlers ...func(h Storage) error) Handler {
  if h.Error != nil {
    return h
  }
  for _, fun := range handlers {
    if err := fun(h.Storage); err != nil {
      h.Error = err
      break
    }
  }
  return h
}

func (h *SimpleHandler) OnFail(handlers ...func(h Storage, err error)) Handler {
  if h.Error == nil {
    return h
  }
  for _, fun := range handlers {
    fun(h.Storage, h.Error)
  }
  return h
}

func (h *SimpleHandler) OnSuccess(handlers ...func(h Storage)) Handler {
  if h.Error != nil {
    return h
  }
  for _, fun := range handlers {
    fun(h.Storage)
  }
  return h
}
