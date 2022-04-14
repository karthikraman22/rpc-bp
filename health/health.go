package api

import "github.com/gin-gonic/gin"

type Probe struct {
	live  *gin.HandlerFunc
	ready *gin.HandlerFunc
}

var p Probe

func RegLivenessProbe(h gin.HandlerFunc) {
	p.live = &h
}

func RegReadinessProbe(h gin.HandlerFunc) {
	p.ready = &h
}

func LivenessProbe() gin.HandlerFunc {
	if p.live == nil {
		return func(c *gin.Context) {
			c.String(200, "ok")
		}
	}
	return *p.live
}

func ReadinessProbe() gin.HandlerFunc {
	if p.ready == nil {
		return func(c *gin.Context) {
			c.String(200, "ok")
		}
	}
	return *p.ready
}
