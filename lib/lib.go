package lib

import (
	"context"
	loggerInterfaces "github.com/fajarardiyanto/flt-go-logger/interfaces"
	log "github.com/fajarardiyanto/flt-go-logger/lib"
	"github.com/fajarardiyanto/flt-go-router/interfaces"
	"github.com/fajarardiyanto/flt-go-router/lib/tree"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	Logger       loggerInterfaces.Logger
	Prefix       string
	Middleware   []interfaces.MiddlewareFunc
	Trees        map[string]*tree.Tree
	NotFound     interfaces.Handler
	PanicHandler func(w http.ResponseWriter, r *http.Request, err interface{})
}

func New() interfaces.Routers {
	logger := log.NewLib()
	logger.Init("HTTP Router")

	return &Router{
		Logger: logger,
		Trees:  make(map[string]*tree.Tree),
	}
}

func (r *Router) GET(path string, handle interfaces.Handler) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *Router) POST(path string, handle interfaces.Handler) {
	r.Handle(http.MethodPost, path, handle)
}

func (r *Router) DELETE(path string, handle interfaces.Handler) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *Router) PUT(path string, handle interfaces.Handler) {
	r.Handle(http.MethodPut, path, handle)
}

func (r *Router) PATCH(path string, handle interfaces.Handler) {
	r.Handle(http.MethodPatch, path, handle)
}

func (r *Router) Group(prefix string) interfaces.Routers {
	return &Router{
		Prefix:     prefix,
		Trees:      r.Trees,
		Middleware: r.Middleware,
	}
}

func (r *Router) Run(addr ...string) {
	address := interfaces.ResolveAddress(addr, r.Logger)
	if address == ":" {
		r.Logger.Error("port can't be empty").Quit()
	}

	r.Logger.Info("Listening and serving HTTP on %s", address)

	if err := http.ListenAndServe(address, r); err != nil {
		r.Logger.Error(err).Quit()
	}
}

func (r *Router) RunTLS(addr, certFile, keyFile string) (err error) {
	err = http.ListenAndServeTLS(addr, certFile, keyFile, r)
	return
}

func (r *Router) NotFoundFunc(handler interfaces.Handler) {
	r.NotFound = handler
}

func (r *Router) Handle(method string, path string, handle interfaces.Handler) {
	if _, ok := interfaces.METHOD[method]; !ok {
		r.Logger.Error("invalid method").Quit()
	}

	treeMethod, ok := r.Trees[method]
	if !ok {
		treeMethod = tree.NewTree()
		r.Trees[method] = treeMethod
	}
	if r.Prefix != "" {
		path = r.Prefix + "/" + path
	}

	treeMethod.Add(path, handle, r.Middleware...)
}

func (r *Router) Use(middleware ...interfaces.MiddlewareFunc) {
	if len(middleware) > 0 {
		r.Middleware = append(r.Middleware, middleware...)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestUrl := req.URL.Path

	if r.PanicHandler != nil {
		defer func() {
			if err := recover(); err != nil {
				r.PanicHandler(w, req, err)
			}
		}()
	}

	if _, ok := r.Trees[req.Method]; !ok {
		r.HandleNotFound(w, req, r.Middleware)
		return
	}

	nodes := r.Trees[req.Method].Find(requestUrl, false)
	if len(nodes) > 0 {
		node := nodes[0]
		if node.Handle != nil {
			if node.Path == requestUrl {
				r.Logger.Trace("[%s][%s]", req.Method, node.Path)
				handle(w, req, node.Handle, node.Middleware)
				return
			}
			if node.Path == requestUrl[1:] {
				r.Logger.Trace("[%s][%s]", req.Method, node.Path)
				handle(w, req, node.Handle, node.Middleware)
				return
			}
		}
	}

	if len(nodes) == 0 {
		res := strings.Split(requestUrl, "/")
		prefix := res[1]
		node := r.Trees[req.Method].Find(prefix, true)
		for _, n := range node {
			if handler := n.Handle; handler != nil && n.Path != requestUrl {
				if matchParamsMap, ok := r.MatchAndParse(requestUrl, n.Path); ok {
					r.Logger.Trace("[%s][%s]", req.Method, res[1])
					ctx := context.WithValue(req.Context(), interfaces.ContextKey, matchParamsMap)
					req = req.WithContext(ctx)
					handle(w, req, handler, n.Middleware)
					return
				}
			}
		}
	}

	r.Logger.Trace("[%s][%s]", req.Method, requestUrl)
	r.HandleNotFound(w, req, r.Middleware)
}

func (r *Router) HandleNotFound(w http.ResponseWriter, req *http.Request, middleware []interfaces.MiddlewareFunc) {
	if r.NotFound != nil {
		handle(w, req, r.NotFound, middleware)
		return
	}
	http.NotFound(w, req)
}

func handle(w http.ResponseWriter, req *http.Request, handler interfaces.Handler, middleware []interfaces.MiddlewareFunc) {
	var baseHandler = handler
	for _, m := range middleware {
		baseHandler = m(baseHandler)
	}
	err := baseHandler(w, req)
	if err != nil {
		return
	}
}

func (r *Router) Match(requestUrl string, path string) bool {
	_, ok := r.MatchAndParse(requestUrl, path)
	return ok
}

func (r *Router) MatchAndParse(requestUrl string, path string) (matchParams interfaces.ParamsMapType, b bool) {
	var (
		matchName []string
		pattern   string
	)

	b = true
	matchParams = make(interfaces.ParamsMapType)
	res := strings.Split(path, "/")
	for _, str := range res {
		if str == "" {
			continue
		}

		strLen := len(str)
		firstChar := str[0]
		lastChar := str[strLen-1]
		if string(firstChar) == "{" && string(lastChar) == "}" {
			matchStr := str[1 : strLen-1]
			res := strings.Split(matchStr, ":")
			matchName = append(matchName, res[0])
			pattern = pattern + "/" + "(" + res[1] + ")"
		} else if string(firstChar) == ":" {
			matchStr := str
			res := strings.Split(matchStr, ":")
			matchName = append(matchName, res[1])
			if res[1] == interfaces.IDKey {
				pattern = pattern + "/" + "(" + interfaces.IDPattern + ")"
			} else {
				pattern = pattern + "/" + "(" + interfaces.DefaultPattern + ")"
			}
		} else {
			pattern = pattern + "/" + str
		}
	}
	if strings.HasSuffix(requestUrl, "/") {
		pattern = pattern + "/"
	}
	re := regexp.MustCompile(pattern)
	if subMatch := re.FindSubmatch([]byte(requestUrl)); subMatch != nil {
		if string(subMatch[0]) == requestUrl {
			subMatch = subMatch[1:]
			for k, v := range subMatch {
				matchParams[matchName[k]] = string(v)
			}
			return
		}
	}
	return nil, false
}
