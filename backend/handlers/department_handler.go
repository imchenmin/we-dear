package handlers

import (
	"net/http"
	"time"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

// storage 初始化函数, 后续要改成单例模式
func initStorage() *storage.DepartmentStorage {
	return storage.GetDepartmentStorage()
}

func GetAllDepartments(c *gin.Context) {
	departmentStorage := initStorage()
	departments, err := departmentStorage.GetAllDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, departments)
}

func GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")
	departmentStorage := initStorage()
	department, err := departmentStorage.GetDepartmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, department)
}

func CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置创建时间等基础字段
	now := time.Now()
	department.BaseModel = models.BaseModel{
		ID:        utils.GenerateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	departmentStorage := initStorage()
	if err := departmentStorage.CreateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, department)
}

func UpdateDepartment(c *gin.Context) {
	// id := c.Param("id")
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initStorage().UpdateDepartment(&department); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, department)
}

func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	if err := initStorage().DeleteDepartment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
