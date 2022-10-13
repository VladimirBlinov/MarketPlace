package apiserver

//type server struct {
//	logger  *logrus.Logger
//	handler handler.Handler
//}
//
//func newServer(handler handler.Handler) *server {
//	s := &server{
//		logger:  logrus.New(),
//		handler: handler,
//	}
//
//	s.handler.InitHandler()
//
//	return s
//}
//
//func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	s.handler.Router.ServeHTTP(w, r)
//}
