package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/liyue201/goqr"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

//文件最大长度2M
var maxLenght int64 = 2 << 20
//设置上传目录
var dir = os.TempDir()

func main() {

	app := fiber.New()

	app.Post("/decode", func(c *fiber.Ctx) error {
		//接收文件
		file, err := c.FormFile("img")
		if err !=nil {
			return c.Status(200).JSON(fiber.Map{"code":500, "msg":err, "data": fiber.Map{}})
		}

		//尺寸验证
		if file.Size > maxLenght{
			return c.Status(200).JSON(fiber.Map{
				"code":404,
				"msg":"Maximum length：2M",
				"data": fiber.Map{},
			})
		}

		//类型验证

		//设置文件名
		tmpName := fmt.Sprintf("/%d", time.Now().UnixNano())
		fileName := dir + tmpName + file.Filename

		//删除图片
		defer os.Remove(fileName)

		//保存图片
		if err:= c.SaveFile(file, fileName); err != nil{
			return c.Status(500).JSON(fiber.Map{
				"code":500,
				"msg":err,
				"data": fiber.Map{},
			})
		}

		return decodeFile(fileName, c)



	})

	app.Get("/url_decode", func(c *fiber.Ctx) error {
		http_url := c.Query("http_url")
		if http_url == "" || len(http_url) < 2 {
			return c.JSON(fiber.Map{"code":404, "msg":"Parameter Error", "data": fiber.Map{}})
		}

		isurl, err := regexp.MatchString(`^([hH][tT]{2}[pP]:\/\/|[hH][tT]{2}[pP][sS]:\/\/|www\.)(([A-Za-z0-9-~]+)\.)+([A-Za-z0-9-~\/])+`, http_url)
		if err != nil{
			return c.JSON(fiber.Map{"code":405, "msg": fmt.Sprintf("url error：%s", err), "data": fiber.Map{}})
		}

		if isurl {
			res, err := http.Get(http_url)
			if err !=nil{
				return c.JSON(fiber.Map{"code":406, "msg":fmt.Sprintf("url error：%s", err), "data": fiber.Map{}})
			}
			//请求成功
			if res.Status == "200 OK" {
				//长度验证
				if res.ContentLength > maxLenght {
					return c.Status(200).JSON(fiber.Map{
						"code":404,
						"msg":"Maximum length：2M",
						"data": fiber.Map{},
					})
				}


				//设置文件名
				tmpName := fmt.Sprintf("/%d", time.Now().UnixNano())
				fileName := dir + tmpName


				f, err := os.Create(fileName)
				//关闭文件
				defer f.Close()
				//删除图片
				defer os.Remove(fileName)
				if err != nil {
					return c.Status(500).JSON(fiber.Map{
						"code":500,
						"msg":err,
						"data": fiber.Map{},
					})

				}
				_,e := io.Copy(f, res.Body)
				if e !=nil{
					return c.Status(500).JSON(fiber.Map{
						"code":500,
						"msg":err,
						"data": fiber.Map{},
					})
				}


				//解析文件
				return decodeFile(fileName, c)


			}else{
				return c.Status(200).JSON(fiber.Map{
					"code":404,
					"msg":"error"+res.Status,
					"data": fiber.Map{},
				})
			}


		}else{
			return c.JSON(fiber.Map{"code":405, "msg": fmt.Sprintf("url error：%s", err), "data": fiber.Map{}})
		}

		return nil
	})

	app.Listen(":9900")
}


func decodeFile(path string, c *fiber.Ctx) error {

	imgdata, err := ioutil.ReadFile(path)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"code" : 404,
			"msg" : err.Error(),
			"data" : fiber.Map{},
		})
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"code" : 404,
			"msg" : "image.Decode error:" + err.Error(),
			"data" : fiber.Map{},
		})

	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"code" : 404,
			"msg" : "Recognize failed:" + err.Error(),
			"data" : fiber.Map{},
		})
	}
	qrArr := make([]string,0)
	for _, qrCode := range qrCodes {

		//fmt.Println(qrCode.Payload)
		text := fmt.Sprintf("%s", qrCode.Payload)
		qrArr = append(qrArr, text)
	}

	return c.Status(200).JSON(fiber.Map{
		"code" : 200,
		"msg" : "Decode Success！",
		"data" : fiber.Map{"qrText":qrArr},
	})
}