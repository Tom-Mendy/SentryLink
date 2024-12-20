// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "schemes": {{ marshal .Schemes }},






















































    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/register": {
            "post": {
                "description": "Register a new user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "parameters": [
                    {
                        "description": "string",
                        "in": "formData",
                        "name": "username",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "description": "string",
                        "in": "formData",
                        "name": "email",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "description": "string",
                        "in": "formData",
                        "name": "password",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered successfully",
                        "schema": {
                            "$ref": "#/definitions/schemas.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/schemas.Response"
                        }
                    },
                    "409": {
                        "description": "User already exists",
                        "schema": {
                            "$ref": "#/definitions/schemas.Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Authenticates a user and provides a JWT to Authorize API calls",
                "produces": ["application/json"],
                "tags": ["auth"],
                "parameters": [
                    {
                        "description": "string",
                        "in": "formData",
                        "name": "username",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "description": "string",
                        "in": "formData",
                        "name": "password",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT",
                        "schema": {
                        "$ref": "#/definitions/schemas.JWT"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                        "$ref": "#/definitions/schemas.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.LinkApi": {
            "type": "object"
        },
        "schemas.JWT": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "schemas.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "bearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:		"",
	Host:			"",
	BasePath:		"",
	Schemes:		[]string{},
	Title:			"",
	Description:		"",
	InfoInstanceName:	"swagger",
	SwaggerTemplate:	docTemplate,
	LeftDelim:		"{{",
	RightDelim:		"}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
