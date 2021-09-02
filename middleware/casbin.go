package middleware

import (
	"net/http"

	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CasbinModel casbin model in string
const CasbinModel string = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && regexMatch(r.obj, p.obj) && r.act == p.act
`

// CasbinEnforcer global casbin instance
var CasbinEnforcer *casbin.Enforcer

// InitCasbin init casbin instance
func InitCasbin() *casbin.Enforcer {
	a := fileadapter.NewAdapter("conf/casbin_policy.csv")
	m := casbin.NewModel(CasbinModel)
	CasbinEnforcer = casbin.NewEnforcer(m, a)
	CasbinEnforcer.LoadPolicy()
	return CasbinEnforcer
}

// CasbinMiddleware Casbin middleware, load authorize policy and enforce authorize check
func CasbinMiddleware() gin.HandlerFunc {

	if CasbinEnforcer == nil {
		CasbinEnforcer = InitCasbin()
	}
	return func(c *gin.Context) {

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		user := GetUserBySession(c)
		var sub string
		if user != nil {
			sub = user.UserID
		} else {
			sub = "anonymous"
		}
		log.Println(sub)
		//判断策略中是否存在
		if CasbinEnforcer.Enforce(sub, obj, act) {
			log.Info("通过权限")
			c.Next()
		} else {
			log.Warn("权限没有通过")
			c.JSON(http.StatusUnauthorized, "权限不足")
			c.Abort()
		}
	}
}
