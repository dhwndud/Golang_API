package controllers

import (
	"chatbot/models"
	"fmt"
	"strings"
	"unicode/utf8"

	k "github.com/Alfex4936/kakao"

	"github.com/gin-gonic/gin"
)

// AskCategory :POST /ask, MUST: "cate": 카테고리 이름
// 메시지 종류: SimpleText
func AskCategory(c *gin.Context) {
	categories := []string{"공지사항", "회사공유일정", "주간메뉴표",
		"건의함", "이달의 우수사원", "업무공유", "경조사 알림", "기타", "행사",
	}

	var replies k.Kakao
	for _, cate := range categories {
		replies.Add(k.QuickReply{}.New(cate, cate))
	}

	c.PureJSON(200, k.SimpleText{}.Build("무슨 공지를 보고 싶으신가요?", replies))
}

// ShowCategory :POST /ask/category, MUST: "cate": 카테고리 이름
// 메시지 종류: SimpleText | ListCard
func ShowCategory(c *gin.Context) {
	// JSON request parse
	var kjson k.Request
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("[ERROR] 오류가 발생했습니다.\n:( [Try Again] 다시 시도해 주세요!", nil)) // http.StatusBadRequest
		return
	}

	categories := map[string]int{
		"공지사항":     1,
		"회사공유일정":   2,
		"주간메뉴표":    3,
		"건의함":      4,
		"이달의 우수사원": 5,
		"업무공유":     6,
		"경조사 알림":   7,
		"기타":       8,
		"행사":       166,
	}

	// Cast to string as cate parameter is an interface
	userCategory := strings.Replace(kjson.Action.Params["cate"].(string), " ", "", 1)
	url := fmt.Sprintf("%v?mode=list&srCategoryId=%v&srSearchKey=&srSearchVal=&articleLimit=5&article.offset=0", models.AjouLink, categories[userCategory])

	var notices []models.Notice = models.Parse(url, 5)
	if len(notices) == 0 {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(" Somma.Inc 홈페이지 서버 반응이 늦고 있네요. 잠시 후 다시 시도해보세요. ", nil))
		return
	}

	listCard := k.ListCard{}.New(false)
	listCard.Title = fmt.Sprintf("%v 공지", userCategory)

	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))
	listCard.Buttons.Add(k.LinkButton{}.New(userCategory, fmt.Sprintf("%v?mode=list&srCategoryId=%v", models.AjouLink, categories[userCategory])))

	for _, notice := range notices {
		if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
			notice.Title = string([]rune(notice.Title)[0:32]) + "..."
		}
		description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])

		listCard.Items.Add(k.ListItemLink{}.New(notice.Title, description, "", notice.Link))
	}

	c.JSON(200, listCard.Build())
}
