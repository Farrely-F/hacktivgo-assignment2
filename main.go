package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type Item struct {
    ItemID      uint   `gorm:"primary_key"`
    ItemCode    string `json:"itemCode"`
    Description string `json:"description"`
    Quantity    int    `json:"quantity"`
    OrderID     int    `json:"orderId"`
}

type Order struct {
    OrderID      uint      `gorm:"primary_key"`
    CustomerName string    `json:"customerName"`
    OrderedAt    time.Time `json:"orderedAt"`
    Items        []Item    `gorm:"foreignkey:OrderID"`
}

var (
    db *gorm.DB
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    dbPassword := os.Getenv("DB_PASSWORD")

    connectionString := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"

    db, err = gorm.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&Item{}, &Order{})
}

func main() {
    defer db.Close()

    router := gin.Default()

    router.POST("/orders", createOrder)
    router.GET("/orders", getOrders)
    router.PUT("/order/:orderId", updateOrder)
    router.DELETE("/order/:orderId", deleteOrder)

    router.Run(":8080")
}

func createOrder(c *gin.Context) {
    var req struct {
        CustomerName string `json:"customerName"`
        Items        []struct {
            ItemCode    string `json:"itemCode"`
            Description string `json:"description"`
            Quantity    int    `json:"quantity"`
        } `json:"items"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order := Order{
        CustomerName: req.CustomerName,
        OrderedAt:    time.Now(), 
    }

    for _, item := range req.Items {
        order.Items = append(order.Items, Item{
            ItemCode:    item.ItemCode,
            Description: item.Description,
            Quantity:    item.Quantity,
        })
    }

    db.Create(&order)
    c.JSON(http.StatusCreated, order)
}

func getOrders(c *gin.Context) {
    var orders []Order
    db.Preload("Items").Find(&orders)
    c.JSON(http.StatusOK, orders)
}

func updateOrder(c *gin.Context) {
    orderIdStr := c.Param("orderId")
    orderId, err := strconv.Atoi(orderIdStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
        return
    }

    var req struct {
        CustomerName string `json:"customerName"`
        Items        []struct {
            LineItemID  int    `json:"lineItemId"`
            ItemCode    string `json:"itemCode"`
            Description string `json:"description"`
            Quantity    int    `json:"quantity"`
        } `json:"items"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var order Order
    if err := db.Preload("Items").First(&order, orderId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }

    order.CustomerName = req.CustomerName

    for _, item := range req.Items {
        for i, existingItem := range order.Items {
            if existingItem.ItemID == uint(item.LineItemID) {
                order.Items[i] = Item{
                    ItemID:      existingItem.ItemID,
                    ItemCode:    item.ItemCode,
                    Description: item.Description,
                    Quantity:    item.Quantity,
                }
                break
            }
        }
    }

    db.Save(&order)
    c.JSON(http.StatusOK, order)
}

func deleteOrder(c *gin.Context) {
    orderIdStr := c.Param("orderId")
    orderId, err := strconv.Atoi(orderIdStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
        return
    }

    var order Order
    if err := db.First(&order, orderId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }

    db.Delete(&order)
    c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
