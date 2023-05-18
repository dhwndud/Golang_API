/*
	[`mappings` 패키지의 코드 리뷰] Markdown.
		Framework: 'gin' 프레임워크 사용
		Middleware: 'Recovery'와 'CORS'를 설정
		+ '/users' 경로에 대한 GET, POST, PUT, DELETE 메서드에 대한 컨트롤러
*/

package mappings

import (
	"chatbot/controllers"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// 'Welcome' message to check if server is running well !
// 서버가 잘 돌아가는 지 확인하는 메세지.
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{"Welcome": "Server is running well."})
}

// 'LimitHandler' blocks too many request at once.
// 'LimitHandler'는 한 번에 너무 많은 요청을 차단한다.
func LimitHandler(lmt *limiter.limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmtGetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// Router ...
// 'Router' 변수는 `gin.Engine` 타입의 라우터를 저장하는 변수
var Router *gin.Engine

/*
'CreateURLMAppings' 함수는 URL 매핑을 설정하는 함수로,
'Router'를 초기화하고 미들웨어와 라우팅을 설정한다.
*/
func CreateURLMappings() {
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)
	Router = gin.New()

	// Create a limiter struct.
	// Allow only 1 request per 1 second
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	Router.Use(LimitHandler(limiter))

	// Apply the middleware to the router (works with groups too)
	Router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, POST",
		RequestHeaders: "*",
		Credentials:    true,
	}))

	v1 := Router.Group("/v1")
	{
		v1.GET("/", Welcome)
		v1.GET("/notices/:num", controllers.GetAllNotices)
		v1.POST("/last", controllers.GetLastNotice)
		v1.POST("/today", controllers.GetTodayNotices)
		v1.POST("/today2", controllers.GetTodayMoreNotices)

		v1.POST("/yesterday", controllers.GetYesterdayNotices)
		v1.POST("/ask", controllers.AskCategory)
		v1.POST("/ask/category", controllers.ShowCategory)
		v1.POST("/schedule", controllers.GetSchedule)
		v1.POST("/search", controllers.SearchKeyword)
		// Infomation
		// v1.POST("/info/weather", controllers.AskWeather)
		v1.POST("/info/weather2", controllers.AskWeatherInCard)
		v1.POST("/info/prof", controllers.SearchProf)
		v1.POST("/info/library", controllers.GetSeatsAvailable)
		v1.POST("/info/meal", controllers.AskMeal)
		v1.POST("/info/job", controllers.AskJob)
	}
}
