// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/articles": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "list all articles",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page num",
                        "name": "p",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ArticleList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article": {
            "post": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create article",
                "parameters": [
                    {
                        "description": "title and content",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateArticle"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ArticleNoReply"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article/{article-id}": {
            "get": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "show article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ArticleForSingle"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "delete article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "article id",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article/{article-id}/reply": {
            "post": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create rely",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "comment",
                        "name": "article",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateReply"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.ReplyNoNested"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article/{article-id}/reply/{reply-id}": {
            "delete": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "delete reply",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "reply's id",
                        "name": "reply-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "reply id",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article/{article-id}/reply/{reply-id}/nested-reply": {
            "post": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create nested rely",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "reply's id",
                        "name": "reply-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "comment",
                        "name": "reply",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateReply"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.NestedReply"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/article/{article-id}/reply/{reply-id}/nested-reply/{nested-reply-id}": {
            "delete": {
                "security": [
                    {
                        "JWTToken": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "delete nested reply",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "article's id",
                        "name": "article-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "reply's id",
                        "name": "reply-id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "nested reply's id",
                        "name": "nested-reply-id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "nested reply id",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/list": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "list articles where article id is in ids",
                "parameters": [
                    {
                        "enum": [
                            "\"1",
                            "2",
                            "3\""
                        ],
                        "type": "string",
                        "description": "article id list separated by comma",
                        "name": "ids",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.ArticleNoReply"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        },
        "/v1/articles/search": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "list articles search by q",
                "parameters": [
                    {
                        "type": "string",
                        "description": "some word in article content",
                        "name": "q",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.ArticleList"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.HttpError": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string",
                    "example": "Some error comment"
                }
            }
        },
        "domain.Article": {
            "type": "object",
            "properties": {
                "claps": {
                    "type": "integer",
                    "example": 123
                },
                "content": {
                    "type": "string",
                    "example": "example content..."
                },
                "created_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "replies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Reply"
                    }
                },
                "title": {
                    "type": "string",
                    "example": "example title"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "domain.ArticleForSingle": {
            "type": "object",
            "properties": {
                "article": {
                    "$ref": "#/definitions/domain.Article"
                }
            }
        },
        "domain.ArticleList": {
            "type": "object",
            "properties": {
                "article_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ArticleNoReply"
                    }
                }
            }
        },
        "domain.ArticleNoReply": {
            "type": "object",
            "properties": {
                "claps": {
                    "type": "integer",
                    "example": 123
                },
                "content": {
                    "type": "string",
                    "example": "example content..."
                },
                "created_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "title": {
                    "type": "string",
                    "example": "example title"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "domain.CreateArticle": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string",
                    "example": "example content..."
                },
                "title": {
                    "type": "string",
                    "example": "example title"
                }
            }
        },
        "domain.CreateReply": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string",
                    "example": "example comment..."
                }
            }
        },
        "domain.NestedReply": {
            "type": "object",
            "properties": {
                "claps": {
                    "type": "integer",
                    "example": 2
                },
                "comment": {
                    "type": "string",
                    "example": "example nested comment..."
                },
                "created_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "reply_id": {
                    "type": "integer",
                    "example": 1
                },
                "updated_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "domain.Reply": {
            "type": "object",
            "properties": {
                "article_id": {
                    "type": "integer",
                    "example": 1
                },
                "claps": {
                    "type": "integer",
                    "example": 5
                },
                "comment": {
                    "type": "string",
                    "example": "example comment..."
                },
                "created_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "nested_replies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.NestedReply"
                    }
                },
                "updated_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "test"
                }
            }
        },
        "domain.ReplyNoNested": {
            "type": "object",
            "properties": {
                "article_id": {
                    "type": "integer",
                    "example": 1
                },
                "claps": {
                    "type": "integer",
                    "example": 5
                },
                "comment": {
                    "type": "string",
                    "example": "example comment..."
                },
                "created_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "updated_at": {
                    "type": "string",
                    "example": "2021-01-15T09:44:35.151+09:00"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                },
                "user_name": {
                    "type": "string",
                    "example": "test"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWTToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Medium Rare Article Service",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
