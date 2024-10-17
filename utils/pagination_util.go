package utils

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PaginatedOpts(page, size int64) *options.FindOptions {
	skip := page*size - size

	fOpt := options.FindOptions{Limit: &size, Skip: &skip}

	return &fOpt
}

func CalculatePaginatedSkip(page, size int64) int64 {
	skip := page*size - size
	return skip
}

func CalculatePageCount(recordCount int64, pageSize int64) int64 {
	pageCount := float64(recordCount) / float64(pageSize)
	_pageCount := int64(math.Ceil(pageCount))
	return _pageCount
}

func GetPaginationParams(c *fiber.Ctx) (int64, int64, int64, int64, bson.D) {
	pageStr := c.Query("page", "1")
	sizeStr := c.Query("size", "10")
	sortStr := c.Query("sort", "")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size < 1 {
		size = 10
	}
	firstIdx := (page - 1) * size
	lastIdx := firstIdx + size - 1

	sortParams := bson.D{}
	for _, v := range strings.Split(sortStr, ";") {
		sortElement := strings.Split(v, ",")

		if len(sortElement) >= 2 && sortElement[0] != "" {
			if sortElement[1] == "asc" {
				sortParams = append(sortParams, bson.E{sortElement[0], 1})
			} else if sortElement[1] == "desc" {
				sortParams = append(sortParams, bson.E{sortElement[0], -1})
			}
		}
	}

	return page, size, firstIdx, lastIdx, sortParams
}
