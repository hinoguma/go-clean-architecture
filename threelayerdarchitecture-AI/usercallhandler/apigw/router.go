package main

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"app/threelayerdarchitecture/setup"
	"github.com/aws/aws-lambda-go/events"
)

type handlerFunc func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type route struct {
	method   string
	segments []routeSegment
	factory  handlerFactory
	once     sync.Once
	handler  handlerFunc
}

type routeSegment struct {
	rawValue   string
	matchValue string
	isParam    bool
}

type Router struct {
	routes   []route
	fallback handlerFunc
}

type handlerFactory func() handlerFunc

func NewRouter(container *setup.Container) *Router {
	r := &Router{
		fallback: func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error":"not_found"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		},
	}

	r.addRoute("POST", "/transfer", func() handlerFunc {
		handler := NewTransferHandler(container)
		return handler.Handle
	})

	return r
}

func (r *Router) Route(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := strings.ToUpper(event.RequestContext.HTTP.Method)
	path := event.RawPath
	if path == "" && event.RequestContext.HTTP.Path != "" {
		path = event.RequestContext.HTTP.Path
	}

	for _, rt := range r.routes {
		if rt.method != method {
			continue
		}

		if rt.matches(path) {
			return rt.resolveHandler()(ctx, event)
		}
	}

	return r.fallback(ctx, event)
}

func (r *Router) addRoute(method, pattern string, factory handlerFactory) {
	r.routes = append(r.routes, route{
		method:   strings.ToUpper(method),
		segments: parsePattern(pattern),
		factory:  factory,
	})
}

func (rt *route) resolveHandler() handlerFunc {
	rt.once.Do(func() {
		if rt.factory != nil {
			rt.handler = rt.factory()
		}
	})
	return rt.handler
}
func (rt route) matches(path string) bool {
	pathSegments := splitPath(path)
	if len(pathSegments) != len(rt.segments) {
		return false
	}

	for idx, segment := range rt.segments {
		if segment.isParam {
			continue
		}

		incoming := pathSegments[idx]
		if decoded, err := url.PathUnescape(incoming); err == nil && decoded != "" {
			incoming = decoded
		}

		if segment.matchValue != incoming {
			return false
		}
	}

	return true
}

func parsePattern(pattern string) []routeSegment {
	var segments []routeSegment
	for _, rawPart := range splitPath(pattern) {
		part := strings.TrimSpace(rawPart)
		if part == "" {
			continue
		}
		isParamCandidate := strings.HasPrefix(part, ":") || (strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}"))
		if isParamCandidate {
			paramName := part
			if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
				paramName = strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
			} else {
				paramName = strings.TrimPrefix(part, ":")
			}

			if paramName != "" {
				segments = append(segments, routeSegment{
					rawValue:   paramName,
					matchValue: "",
					isParam:    true,
				})
				continue
			}
		}

		matchValue := part
		if decoded, err := url.PathUnescape(part); err == nil && decoded != "" {
			matchValue = decoded
		}

		segments = append(segments, routeSegment{
			rawValue:   part,
			matchValue: matchValue,
			isParam:    false,
		})
	}
	return segments
}

func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}
