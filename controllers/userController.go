package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/guilherm5/memcachedGorm/database"
	"github.com/guilherm5/memcachedGorm/models"
)

var DB = database.Init()
var MC = memcache.New("127.0.0.1:11211")
var user models.Users

func UsersFunc() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Erro ao decodificar body para realizar insert": err})
			log.Println("Erro ao decodificar body para realizar insert", err)
			return
		}

		//adicionando usuario no postgresql
		if err := DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro ao realizar insert": err})
			log.Println("Erro ao realizar insert", err)
			return
		}

		//adicionando user ao memcached
		item := &memcache.Item{
			Key:        "user_" + strconv.Itoa(int(user.ID)),
			Value:      []byte(user.Name),
			Expiration: 3600,
		}

		if err := MC.Add(item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Erro ao adicionar dado no mccached": err})
			log.Println("Erro ao adicionar dado no mccached", err)
			return
		}

		cookie := &http.Cookie{
			Name:  "user_id",
			Value: strconv.Itoa(user.ID),
		}
		http.SetCookie(c.Writer, cookie)
		log.Println("Armazenamos cookie ID:", cookie)
	}
}

func UsersFuncID() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("user_id")
		if err != nil {
			// cookie não encontrado
			c.JSON(http.StatusBadRequest, gin.H{
				"Cookie não encontrado": err.Error(),
			})
			log.Println("Cookie não encontrado", err)
			return
		}

		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			// valor do cookie não é um número
			c.JSON(http.StatusBadRequest, gin.H{
				"Valor do cookie não é um inteiro": err.Error(),
			})
			log.Println("Valor do cookie não é um inteiro", err)
			return
		}

		// verificando se o usuario esta no memcached
		item, err := MC.Get("user_" + strconv.Itoa(int(userID)))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Usuario nao encontrado no memcached": err,
			})
			log.Println("Usuario nao encontrado no memcached", err)
		}

		valueItemMemcached := string(item.Value[:])
		c.JSON(http.StatusOK, gin.H{
			"name": string(item.Value),
		})
		log.Println(valueItemMemcached, "Vim do Memcached")

		//get user postgres
		var userPostgres []models.Users

		if err := DB.Find(&userPostgres).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Erro ao realizar select": err.Error(),
			})
			log.Println("Erro ao realizar select", err)
			return
		}

		userJSON, err := json.MarshalIndent(userPostgres, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Erro ao deserializar dados": err.Error(),
			})
			log.Println("Erro ao deserializar dados", err)
			return
		}

		c.Data(http.StatusOK, "application/json", userJSON)
		log.Println("Vim do Banco de Dados", string(userJSON))

	}
}
