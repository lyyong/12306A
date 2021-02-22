# 用户服务

# API

## 用户接口

### 注册

POST /user/api/v1/register

**请求**

```json
{
    "username": "",
    "password": "",
    "certificateType": 0,
    "name": "",
    "certificateNumber": "",
    "phoneNumber": "",
    "email": "",
    "passengerType": 0
}
```

**响应**

```json
{
    "code":"",
    "msg": "",
    "data":{}
}
```

### 登录

POST /user/api/v1/login

**请求**

```json
{
    "username": "",
    "password": ""
}
```

**响应**

```json
{
    "code":"",
    "msg": "",
    "data":{
        "token":""
    }
}
```

### 添加乘车人

POST /user/api/v1/passenger

**Header**

```json
{
    "token": ""
}
```

**请求**

```json
{
    "data": [
        {
            "name": "",
            "certificateType": 0,
            "certificateNumber": "",
            "passengerType": 0
        },
        {
            "name": "",
            "certificateType": 0,
            "certificateNumber": "",
            "passengerType": 0
        }
    ]
}
```

**响应**

```json
{
    "code":"",
    "msg": "",
    "data":{}
}
```

### 修改乘车人

PUT /user/api/v1/passenger

**Header**

```json
{
    "token": ""
}
```

**请求**

```json
{
    "data": [
        {
            "name": "",
            "certificateType": 0,
            "certificateNumber": "",
            "passengerType": 0
        },
        {
            "name": "",
            "certificateType": 0,
            "certificateNumber": "",
            "passengerType": 0
        }
    ]
}
```

**响应**

```json
{
    "code":"",
    "msg": "",
    "data":{}
}
```

### 查询乘车人

GET /user/api/v1/passenger

**Header**

```json
{
    "token": ""
}
```

**请求**

```json
{
}
```

**响应**

```json
{
    "code":"",
    "msg": "",
    "data":{
        "passenger": [
            {
                "id": 0,
                "name": "",
                "certificateType": 0,
                "certificateNumber": "",
                "passengerType": 0
            },
            {
                "id": 0,
                "name": "",
                "certificateType": 0,
                "certificateNumber": "",
                "passengerType": 0
            }
        ]
    }
}
```

