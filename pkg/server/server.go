package server

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/user/b3i/pkg/device"
)

//go:embed all:dist
var staticAssets embed.FS

type Server struct {
	Port       int
	DeviceIP   string
	Password   string
	Insecure   bool
	Client     *device.Client
	AdminToken string
}

func NewServer(port int, deviceIP string, password string, insecure bool) *Server {
	// Generate an ephemeral admin token for the session
	token := make([]byte, 16)
	rand.Read(token)
	adminToken := hex.EncodeToString(token)

	return &Server{
		Port:       port,
		DeviceIP:   deviceIP,
		Password:   password,
		Insecure:   insecure,
		AdminToken: adminToken,
	}
}

func (s *Server) Start() error {
	// Initialize device client
	s.Client = device.NewClient(s.DeviceIP, s.Password, s.Insecure)
	if s.Password != "" {
		if err := s.Client.Login(); err != nil {
			fmt.Printf("Warning: device login failed: %v\n", err)
		}
	}

	fmt.Printf("\n--- B3i Web UI Session Token: %s ---\n", s.AdminToken)
	fmt.Printf("Access the UI at: http://localhost:%d/ui/?token=%s\n\n", s.Port, s.AdminToken)

	r := gin.Default()

	// Simple Token-based Auth Middleware
	authMiddleware := func(c *gin.Context) {
		token := c.GetHeader("X-B3i-Token")
		if token == "" {
			token = c.Query("token")
		}
		if token != s.AdminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}

	// Serve static files from embedded FS
	sub, _ := fs.Sub(staticAssets, "dist")
	r.StaticFS("/ui", http.FS(sub))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})

	// API Routes
	api := r.Group("/api")
	api.Use(authMiddleware) // Protect all API endpoints
	{
		api.GET("/device", s.handleGetDeviceInfo)
		api.GET("/apps", s.handleListApps)
		api.POST("/install", s.handleInstallApp)
		api.DELETE("/apps/:id", s.handleUninstallApp)
		api.POST("/apps/:id/:action", s.handleManageApp)
	}

	return r.Run(fmt.Sprintf(":%d", s.Port))
}

func (s *Server) handleGetDeviceInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ip":     s.DeviceIP,
		"status": "connected",
	})
}

func (s *Server) handleListApps(c *gin.Context) {
	apps, err := s.Client.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func (s *Server) handleInstallApp(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	tempPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temp file"})
		return
	}
	defer os.Remove(tempPath)

	if err := s.Client.InstallApp(tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) handleUninstallApp(c *gin.Context) {
	id := c.Param("id")
	if err := s.Client.UninstallApp(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) handleManageApp(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")
	if action != "launch" && action != "terminate" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	if err := s.Client.ManageApp(id, action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
