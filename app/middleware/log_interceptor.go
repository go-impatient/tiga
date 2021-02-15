package middleware

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"moocss.com/tiga/internal/service/dto"
)

type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.Buffer.Write(data)
	return g.out.Write(data)
}

func LogInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		w := bufio.NewWriter(ctx.Writer)
		buff := bytes.Buffer{}
		newWriter := &bufferedWriter{ctx.Writer, w, buff}
		ctx.Writer = newWriter
		defer func() {
			logrus.Infof("response status : %d; body : %s", ctx.Writer.Status(), newWriter.Buffer.Bytes())
			w.Flush()
		}()
		body, err := ctx.GetRawData()
		if err != nil {
			dto.HandleErrorF(ctx, http.StatusOK, err)
			return
		}
		logrus.Infof("request [%s%s] body : %v", ctx.Request.Host, ctx.Request.URL, string(body))
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		ctx.Next()
	}
}
