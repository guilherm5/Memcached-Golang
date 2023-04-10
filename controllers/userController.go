package controllers

import (
	"log"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/guilherm5/memcachedGorm/database"
	"github.com/guilherm5/memcachedGorm/models"
	"gorm.io/gorm"
)

var DB = database.Init()
var MC = memcache.New("127.0.0.1:11211")

func UsersFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Users

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"erro ao decodificar body para realizar insert": err})
			return
		}

		//adicionando usuario no postgresql
		if err := DB.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"erro ao realizar insert": err})
			log.Println("erro ao realizar insert", err)
			return
		}

		//adicionando user no memcached
		item := &memcache.Item{
			Key:        "user_" + strconv.Itoa(int(user.ID)),
			Value:      []byte(user.Name),
			Expiration: 3600,
		}

		if err := MC.Add(item); err != nil {
			c.JSON(500, gin.H{"erro ao adicionar dado no mccached": err})
			log.Println("erro ao adicionar dado no mccached", err)
			return
		}
		c.JSON(200, user)
	}
}

func UsersFuncID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		item, err := MC.Get("user_" + id)
		if err == nil {
			c.JSON(200, gin.H{
				"name": string(item.Value),
			})
			log.Println("VIM DO CACHE")
			return
		}

		//get user postgres
		var user models.Users
		if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"erro": "Usuario nao existe"})
			} else {
				c.JSON(500, gin.H{"erro": err})
			}
			return
		}

		//adicionando item no memcached
		item = &memcache.Item{
			Key:        "user_" + strconv.Itoa(int(user.ID)),
			Value:      []byte(user.Name),
			Expiration: 10,
		}
		if err := MC.Add(item); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, user)
		log.Println("VIM DO BANCO")
	}
}
