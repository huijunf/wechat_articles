package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"your-app-name/models"
)

func Index(c *gin.Context) {
	accounts, err := models.GetAllPublicAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func FetchArticles(c *gin.Context) {
	bizID := c.Param("bizID")
	account, err := models.GetPublicAccountByBizID(bizID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Public Account not found"})
		return
	}

	// Get the HTML of the account's page
	doc, err := goquery.NewDocument(fmt.Sprintf("https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=%s==#wechat_redirect", bizID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Get the article list items and extract their information
	articleList := doc.Find("div#js_articlelist div.rich_media_area_primary")
	articleList.Each(func(i int, article *goquery.Selection) {
		link := article.Find("h4 a")
		url, _ := link.Attr("href")
		title := link.Text()
		createTimeStr := article.Find("div.rich_media_meta_text").Eq(1).Text()
		createTime, _ := time.Parse("2006-01-02", createTimeStr)
		id := xid.New().String()

		// Check if the article already exists in the database
		exists, err := models.ArticleExists(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		if exists {
			return
		}

		// Download the article's HTML
		articleDoc, err := goquery.NewDocument(url)
		if err != nil {
			return
		}

		// Save the article to the database
		articleHTML, _ := articleDoc.Find("div.rich_media_content").Html()
		err = models.CreateArticle(models.Article{
			ID:         id,
			URL:        url,
			AccountID:  account.ID,
			CreateTime: createTime,
			HTML:       articleHTML,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
