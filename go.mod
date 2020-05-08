module github.com/Yesterday17/pug-backend

go 1.14

// TODO: Remove local path when PUG is stable

replace github.com/Yesterday17/pug => ../pug

replace github.com/Yesterday17/pug/api => ../pug/api

require (
	github.com/Yesterday17/pug v0.0.0
	github.com/Yesterday17/pug/api v1.0.6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.9
	github.com/satori/go.uuid v1.2.0
)
