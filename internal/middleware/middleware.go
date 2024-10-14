package middleware

import (
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strings"
	"time"
)

func RequestLogger(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Infof("| %3d | %13v | %15s | %-7s %s",
			status,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func GzipDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(400, gin.H{"error": "Bad Request"})
				return
			}
			defer gz.Close()
			c.Request.Body = io.NopCloser(gz)
		}
		c.Next()
	}
}

func GzipCompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			// Добавьте логирование
			fmt.Println("Gzip compression is enabled for this response")

			c.Writer.Header().Set("Content-Encoding", "gzip")
			c.Writer.Header().Set("Vary", "Accept-Encoding")

			gz := gzip.NewWriter(c.Writer)
			c.Writer = &gzipResponseWriter{Writer: gz, ResponseWriter: c.Writer}

			c.Next()

			gz.Close()
		} else {
			c.Next()
		}
	}
}
