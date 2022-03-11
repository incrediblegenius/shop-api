package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-web/global"
	"shop-api/user-web/models"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.Id)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token is not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//zap.S().Debugf(token.SigningString())
	return token.SignedString(j.SigningKey)

}

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}

		}
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() // 默认token有效期为1个小时
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}

//// 生成令牌
//func GenerateToken(c *gin.Context, user *models.SysUser) (token string, msg string, ok bool) {
//	var grade_list []string
//	var class_id_list []int
//	if user.RoleId == 1 {
//		dao.DB.Model(&models.SysGrade{}).Pluck("name", &grade_list)
//		dao.DB.Model(&models.SysClass{}).Pluck("id", &class_id_list)
//	} else {
//		grade_list = strings.Split(user.Grades, ",")
//		fmt.Println(grade_list)
//		dao.DB.Model(&models.SysClass{}).Joins("left join sys_grades on sys_grades.id = sys_classes.grade_id").Where("major_id = ? and sys_grades.name in (?)", user.MajorID, grade_list).Pluck("sys_classes.id", &class_id_list)
//	}
//	fmt.Println(grade_list)
//	j := &JWT{[]byte(SignKey)}
//	claims := request.CustomClaims{
//		user.ID,
//		user.Username,
//		user.AvatarUrl,
//		user.RoleId,
//		user.MajorID,
//		grade_list,
//		class_id_list,
//		jwt.StandardClaims{
//			NotBefore: int64(time.Now().Unix() - 1000),          // 签名生效时间
//			ExpiresAt: int64(time.Now().Unix() + TokenExpireAt), // 过期时间 一小时
//			Issuer:    Issuer,                                   //签名的发行者
//		},
//	}
//	token, err := j.CreateToken(claims)
//	if err != nil {
//		log.Println(err)
//		return token, "创建token失败", false
//	} else {
//		return token, "登录成功！", true
//	}
//}
