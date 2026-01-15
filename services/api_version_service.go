package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	APIVersionV1   = "v1"
	APIVersionV2   = "v2"
	CurrentVersion = APIVersionV1
)

type APIVersion struct {
	Version    string `json:"version"`
	Deprecated bool   `json:"deprecated"`
	SunsetDate string `json:"sunset_date,omitempty"`
	Migration  string `json:"migration,omitempty"`
}

type VersionHandler interface {
	GetVersionInfo(version string) *APIVersion
	IsDeprecated(version string) bool
	GetMigrationPath(version, from, to string) string
}

type DefaultVersionHandler struct{}

func (h *DefaultVersionHandler) GetVersionInfo(version string) *APIVersion {
	info := &APIVersion{
		Version:    version,
		Deprecated: false,
	}

	switch version {
	case APIVersionV1:
		info.Deprecated = false
	case APIVersionV2:
		info.Deprecated = false
	}

	return info
}

func (h *DefaultVersionHandler) IsDeprecated(version string) bool {
	return version == APIVersionV1 // Example: v1 is deprecated
}

func (h *DefaultVersionHandler) GetMigrationPath(version, from, to string) string {
	if from == "" || to == "" {
		return ""
	}

	if strings.HasPrefix(version, "v") {
		return fmt.Sprintf("/migrate/%s-to-%s", from, to)
	}

	return fmt.Sprintf("/migrate/%s", from)
}

func APIVersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := extractVersionFromPath(c.Request.URL.Path)
		c.Set("api_version", version)
		c.Set("version_info", GetVersionHandler().GetVersionInfo(version))
		c.Next()
	}
}

func extractVersionFromPath(path string) string {
	// Extract version from path like /v1/products or /v2/orders
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && strings.HasPrefix(parts[0], "v") {
		return parts[0]
	}
	return CurrentVersion
}

func SetupVersionedRoutes(r *gin.Engine) {
	// V1 Routes (current)
	v1 := r.Group("/v1")
	{
		setupV1Routes(v1)
	}

	// V2 Routes (future)
	v2 := r.Group("/v2")
	{
		setupV2Routes(v2)
	}

	// Legacy routes (deprecated)
	legacy := r.Group("/api")
	{
		setupLegacyRoutes(legacy)
	}
}

func setupV1Routes(rg *gin.RouterGroup) {
	rg.GET("/products", func(c *gin.Context) {
		version := c.GetString("api_version")
		c.JSON(http.StatusOK, gin.H{
			"message": "V1 Products API",
			"version": version,
		})
	})
}

func setupV2Routes(rg *gin.RouterGroup) {
	rg.GET("/products", func(c *gin.Context) {
		version := c.GetString("api_version")
		c.JSON(http.StatusOK, gin.H{
			"message": "V2 Products API - Enhanced",
			"version": version,
		})
	})
}

func setupLegacyRoutes(rg *gin.RouterGroup) {
	rg.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Legacy API - Consider migrating to v1",
			"deprecated": true,
		})
	})
}

func GetDeprecationHeaders(version string) map[string]string {
	headers := make(map[string]string)

	if GetVersionHandler().IsDeprecated(version) {
		headers["Deprecation"] = "true"
		headers["Sunset"] = "2024-12-31" // Example sunset date
		headers["Link"] = fmt.Sprintf("</%s>; rel=\"successor-version\"",
			strings.Replace(version, "v", "v", 1))
	}

	return headers
}
