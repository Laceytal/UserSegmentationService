package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"UserSegmentationServise/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ======================
// RegisterRoutes
// ======================

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/segments", CreateSegment(db))
	r.GET("/segments", ListSegments(db))
	r.PUT("/segments/:segment_id", UpdateSegment(db))
	r.DELETE("/segments/:segment_id", DeleteSegment(db))

	r.POST("/users/:user_id/segments/:segment_id", AddUserSegment(db))
	r.DELETE("/users/:user_id/segments/:segment_id", RemoveUserSegment(db))
	r.GET("/users/:user_id/segments", GetUserSegments(db))

	r.POST("/segments/:segment_id/assign_pct", AssignPct(db))
}

// ======================
// Создание сегмента
// ======================

func CreateSegment(db *gorm.DB) gin.HandlerFunc {
	type req struct{ Key, Description string }
	return func(c *gin.Context) {

		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		seg := models.Segment{Key: body.Key, Description: body.Description}
		if err := db.Create(&seg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, seg)
	}
}

// ======================
// Список всех сегментов
// ======================

func ListSegments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var segments []models.Segment
		db.Find(&segments)
		c.JSON(http.StatusOK, segments)
	}
}

// ======================
// Обновление сегмента
// ======================

func UpdateSegment(db *gorm.DB) gin.HandlerFunc {
	type req struct{ Description string }
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("segment_id"))
		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var seg models.Segment
		if err := db.First(&seg, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "segment not found"})
			return
		}
		seg.Description = body.Description
		db.Save(&seg)
		c.JSON(http.StatusOK, seg)
	}
}

// ======================
// Удаление сегмента
// ======================

func DeleteSegment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("segment_id"))
		db.Delete(&models.Segment{}, id)
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

// ======================
// Добавление сегмента пользователю
// ======================

func AddUserSegment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := strconv.Atoi(c.Param("user_id"))
		segID, _ := strconv.Atoi(c.Param("segment_id"))

		us := models.UserSegment{
			UserID:     uint(userID),
			SegmentID:  uint(segID),
			AssignedAt: time.Now(),
		}
		if err := db.Create(&us).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "assigned"})
	}
}

// ======================
// Удаление сегмента у пользователя
// ======================

func RemoveUserSegment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := strconv.Atoi(c.Param("user_id"))
		segID, _ := strconv.Atoi(c.Param("segment_id"))

		db.Delete(&models.UserSegment{}, "user_id = ? AND segment_id = ?", userID, segID)
		c.JSON(http.StatusOK, gin.H{"status": "removed"})
	}
}

// ======================
// Назначение правила распределения X% пользователей
// ======================

func AssignPct(db *gorm.DB) gin.HandlerFunc {
	type req struct{ Pct int }
	return func(c *gin.Context) {
		segID, _ := strconv.Atoi(c.Param("segment_id"))
		var body req
		if err := c.ShouldBindJSON(&body); err != nil || body.Pct < 1 || body.Pct > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pct должен быть 1..100"})
			return
		}

		rule := models.SegmentRule{
			SegmentID: uint(segID),
			Pct:       body.Pct,
			Seed:      "SEG_" + strconv.Itoa(segID), // генерация seed для детерминизма
		}

		// Upsert: INSERT or UPDATE
		db.Clauses(
		// если используешь Postgres ≥13, применяй on conflict do update
		).Create(&rule)

		c.JSON(http.StatusOK, rule)
	}
}

// ======================
// Получение сегментов пользователя (ручных + X%)
// ======================

func GetUserSegments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := strconv.Atoi(c.Param("user_id"))

		// 1. Явные назначения
		var manual []models.Segment
		db.Model(&models.Segment{}).
			Joins("JOIN user_segments us ON us.segment_id = segments.id").
			Where("us.user_id = ?", userID).
			Find(&manual)

		// 2. Автоматические назначения по правилам
		var rules []models.SegmentRule
		db.Find(&rules)

		result := make([]string, 0)
		for _, s := range manual {
			result = append(result, s.Key)
		}

		// Проверка inPct для каждого правила
		for _, r := range rules {
			if inPct(userID, r) {
				var seg models.Segment
				db.First(&seg, r.SegmentID)
				result = append(result, seg.Key)
			}
		}

		c.JSON(http.StatusOK, gin.H{"segments": result})
	}
}

// ======================
// Хелпер: проверка входит ли user_id в X%
// ======================

func inPct(userID int, rule models.SegmentRule) bool {
	// SHA-256(user_id:seed)
	h := sha256.Sum256([]byte(
		strconv.Itoa(userID) + ":" + rule.Seed,
	))
	hexVal := hex.EncodeToString(h[:])
	// первые 8 символов → число
	v, _ := strconv.ParseUint(hexVal[:8], 16, 32)
	return int(v%100) < rule.Pct
}
