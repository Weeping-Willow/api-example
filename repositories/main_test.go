package repositories

import (
	"os"
	"testing"

	"github.com/Weeping-Willow/api-example/testt"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	testt.LoadConfig(m, "../.env")
	testt.GetConnection = GetConnection
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
