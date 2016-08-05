# Specification of internal APIs (version v1.0)

## APIs of Registration Center

### List

|Method|Path|Request Source|Entity|Description|
|------|----|--------------|------|-----------|
|POST|`/internal/v1/dataserver`|Administrator|Data Server|Register a new data server or a batch of data servers|
|DELETE|`/internal/v1/:dataserver`|Administrator|Data Server|Delete a data server|
|GET|`/internal/v1/:dataserver`|Administrator|Data Server|Get data server info by a given data server ID|
|PUT|`/internal/v1/dataserver`|Data Groups|Data Server|Recieve status from data server|
|GET|`/internal/v1/groups/:group`|Object Server|Group|Get group info by a given group_id |
|GET|`/internal/v1/groups`|Object Server|Group|Get all available groups for Object Server|
|GET|`/internal/v1/object/id`|Object Server|Object|Returns a FileID allocated by Registration Center|
|PUT|`/internal/v1/object`|Object Server|Object|Save/Update object info|
|GET|`/internal/v1/object/:object`|Object Server|Object|Query object info|


### POST  /internal/v1/dataserver

Register a new data server or a batch of data servers
#### Request 

- **Syntax**
```http
POST /internal/v1/dataserver   HTTP/1.1
Host: <Arkor Host>
Content-Type: application/json
```
```json
[
  {
    "group_id": 1,
    "ip": "10.1.0.6",
    "port": 7654
  },
  {
    "group_id": 2,
    "ip": "10.1.0.8",
    "port": 7664
  }
]
```

- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`group_id`|String|Identifier of the group|
|`ip`|String|IP address of data server|
|`port`|Int|Port of data server listenning on |

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
Content-Type: application/json
```

```json
[
  {
    "data_server_id":"as89ik",
    "group_id": 1,
    "ip": "10.1.0.6",
    "port": 7654
  },
  {
    "data_server_id":"zd32hg",
    "group_id": 2,
    "ip": "10.1.0.8",
    "port": 7664
  }
]
```

- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`data_server_id`|String|Unique ID assigned by registration center|
|`group_id`|String|Same value with Request|
|`ip`|String|Same value with Request|
|`port`|Int|Same value with Request|

#### Response On Failure

* 400 - Bad Request (Invalid Parameters or Incorrect json content)
* 409 - Conflict (Data Server already registered)
* 500 - Internal Server Error

### DELETE  /internal/v1/:dataserver

Delete a data server

#### Request 

- **Syntax**
```http
DELETE  /internal/v1/:dataserver   HTTP/1.1
Host: <Arkor Host>
```

- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`dataserver`|URL|Data server ID to be delete|

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
```

#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Data server is NOT registered)
* 409 - Conflict (Data Server already deleted)
* 500 - Internal Server Error

### GET  /internal/v1/:data_server_id

Get data server info by a given data server ID
#### Request 

- **Syntax**
```http
GET /internal/v1/:dataserver   HTTP/1.1
Host: <Arkor Host>
```

- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`dataserver`|URL|Data server ID to be query|

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
Content-Type: application/json
```
```json
{
  "data_server_id": "as89ik",
  "group_id": 1,
  "ip": "10.1.0.6",
  "port": 7654,
  "status": 2,
  "total_free_space": 4294966273,
  "max_free_space": 2147483593,
  "pend_writes": 123,
  "data_path": "/root/ossdata",
  "reading_count": 123,
  "total_chunks": 2,
  "conn_counts": 0,
  "deleted": 0,
  "create_time": "2016-06-05 19:35:36",
  "update_time": "2016-07-04 16:09:38"
}
```

- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`data_server_id`|String|The same with request|
|`group_id`|String|Identifier of the group|
|`ip`|String|IP address of data server|
|`port`|Int|Port of data server listenning on|
|`status`|Int|Current status of dataserver, must be one of: **INIT_STATUS(0)**, **RW_STATUS(1)**, **RO_STATUS(2)** or **ERR_STATUS(3)**|
|`total_free_space`|Int64|Total free space of the data server|
|`max_free_space`|Int64|Max remain free space of all chunks on the data server|
|`pend_writes`|Int|Pending queue writes count|
|`data_path`|String|File path of chunks to store|
|`reading_count`|Int|Reading connections connectted to the data server|
|`total_chunks`|Int|Chunk number in the data server|
|`conn_counts`|Int|Connections connectted to the data server|
|`deleted`|Int|The flag of whethere the data server is deleted|
|`create_time`|Date|Date of data server regitered|
|`update_time`|Date|Date of data server info updated|


#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Data server is NOT registered)
* 409 - Conflict (Data Server already deleted)
* 500 - Internal Server Error


### PUT /internal/v1/dataserver

Recieve status from data server

#### Request 

- **Syntax**
```http
PUT /internal/v1/dataserver   HTTP/1.1
Host: <Arkor Host>
Date: date
```

```json
{
  "group_id": 1,
  "ip": "10.1.0.6",
  "port": 7654,
  "status": 2,
  "total_free_space": 4294966273,
  "max_free_space": 2147483593,
  "pend_writes": 123,
  "data_path": "/root/ossdata",
  "writing_count": 123,
  "reading_count": 123,
  "total_chunks": 2,
  "conn_counts": 123
}
```
- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`group_id`|String|Identifier of the group|
|`ip`|String|IP address of data server|
|`port`|Int|Port of data server listenning on|
|`status`|Int|Current status of dataserver, including: **INIT_STATUS:0**, **RW_STATUS:1**, **RO_STATUS:2**, **ERR_STATUS:3**|
|`total_free_space`|Int64|Total free space of the data server|
|`max_free_space`|Int64|Max remain free space of all chunks on the data server|
|`pend_writes`|Int|Pending queue writes count|
|`data_path`|String|File path of chunks to store|
|`total_chunks`|Int|Chunk number in the data server|
|`conn_counts`|Int|Connections connectted to the data server|

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
```


#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Data server is NOT registered)
* 409 - Conflict (Data Server already deleted)
* 500 - Internal Server Error



### GET  /internal/v1/groups/:group

Get group info by a given group_id

#### Request 

- **Syntax**
```http
GET /internal/v1/groups/:group  HTTP/1.1
Host: <Arkor Host>
```

- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`group`|URL|Group ID to be query|

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
Content-Type: application/json
```
```json
{
  "id": "1",
  "group_status": "0",
  "servers": [
    {
      "data_server_id": "as89ik",
      "group_id": 1,
      "ip": "10.1.0.6",
      "port": 7654,
      "status": 2,
      "total_free_space": 4294966273,
      "max_free_space": 2147483593,
      "pend_writes": 123,
      "data_path": "/root/ossdata",
      "reading_count": 123,
      "total_chunks": 2,
      "conn_counts": 0,
      "update_time": "2016-07-04 16:09:38"
    },
    {
      "data_server_id": "zd32hg",
      "group_id": 1,
      "ip": "10.1.0.8",
      "port": 7654,
      "status": 2,
      "total_free_space": 4294966273,
      "max_free_space": 2147483593,
      "pend_writes": 123,
      "data_path": "/root/ossdata",
      "reading_count": 123,
      "total_chunks": 2,
      "conn_counts": 0,
      "update_time": "2016-07-04 16:09:38"
    }
  ]
}

```

- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`id`|string|ID of the group|
|`global_status`|Int|Global status of whole group, must be one of: **GROUP_STATUS_NORMAL(0)**, or **GROUP_STATUS_UNNORMAL(1)**|
|`data_server_id`|String|The same with request|
|`group_id`|String|Identifier of the group|
|`ip`|String|IP address of data server|
|`port`|Int|Port of data server listenning on|
|`status`|Int|Current status of dataserver, must be one of: **INIT_STATUS(0)**, **RW_STATUS(1)**, **RO_STATUS(2)** or **ERR_STATUS(3)**|
|`total_free_space`|Int64|Total free space of the data server|
|`max_free_space`|Int64|Max remain free space of all chunks on the data server|
|`pend_writes`|Int|Pending queue writes count|
|`data_path`|String|File path of chunks to store|
|`reading_count`|Int|Reading connections connectted to the data server|
|`total_chunks`|Int|Chunk number in the data server|
|`conn_counts`|Int|Connections connectted to the data server|
|`update_time`|Date|Date of data server info updated|

#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Group not found)
* 500 - Internal Server Error

### GET  /internal/v1/groups

Get all available groups for Object Server

#### Request 

- **Syntax**
```http
GET /internal/v1/groups   HTTP/1.1
Host: <Arkor Host>
```

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
Content-Type: application/json
```
```json
[
  {
    "id": "1",
    "group_status": "0",
    "servers": [
      {
        "data_server_id": "as89ik",
        "group_id": 1,
        "ip": "10.1.0.6",
        "port": 7654,
        "status": 2,
        "group_status": 0,
        "total_free_space": 4294966273,
        "max_free_space": 2147483593,
        "pend_writes": 123,
        "data_path": "/root/ossdata",
        "reading_count": 123,
        "total_chunks": 2,
        "conn_counts": 0,
        "update_time": "2016-07-04 16:09:38"
      },
      {
        "data_server_id": "zd32hg",
        "group_id": 1,
        "ip": "10.1.0.8",
        "port": 7654,
        "status": 2,
        "global_status": 0,
        "total_free_space": 4294966273,
        "max_free_space": 2147483593,
        "pend_writes": 123,
        "data_path": "/root/ossdata",
        "reading_count": 123,
        "total_chunks": 2,
        "conn_counts": 0,
        "update_time": "2016-07-04 16:09:38"
      }
    ]
  },
  {
    "id": "2",
    "group_status": "0",
    "servers": [
      {
        "data_server_id": "asadsk",
        "group_id": 2,
        "ip": "10.1.0.6",
        "port": 7654,
        "status": 2,
        "global_status": 0,
        "total_free_space": 4294966273,
        "max_free_space": 2147483593,
        "pend_writes": 123,
        "data_path": "/root/ossdata",
        "reading_count": 123,
        "total_chunks": 2,
        "conn_counts": 0,
        "update_time": "2016-07-04 16:09:38"
      },
      {
        "data_server_id": "12e2hg",
        "group_id": 2,
        "ip": "10.1.0.8",
        "port": 7654,
        "status": 2,
        "global_status": 0,
        "total_free_space": 4294966273,
        "max_free_space": 2147483593,
        "pend_writes": 123,
        "data_path": "/root/ossdata",
        "reading_count": 123,
        "total_chunks": 2,
        "conn_counts": 0,
        "update_time": "2016-07-04 16:09:38"
      }
    ]
  }
]

```

- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`id`|string|ID of the group|
|`global_status`|Int|Global status of whole group, must be one of: **GROUP_STATUS_NORMAL(0)**, or **GROUP_STATUS_UNNORMAL(1)**|
|`data_server_id`|String|The same with request|
|`group_id`|String|Identifier of the group|
|`ip`|String|IP address of data server|
|`port`|Int|Port of data server listenning on|
|`status`|Int|Current status of dataserver, must be one of: **INIT_STATUS(0)**, **RW_STATUS(1)**, **RO_STATUS(2)** or **ERR_STATUS(3)**|
|`total_free_space`|Int64|Total free space of the data server|
|`max_free_space`|Int64|Max remain free space of all chunks on the data server|
|`pend_writes`|Int|Pending queue writes count|
|`data_path`|String|File path of chunks to store|
|`reading_count`|Int|Reading connections connectted to the data server|
|`total_chunks`|Int|Chunk number in the data server|
|`conn_counts`|Int|Connections connectted to the data server|
|`update_time`|Date|Date of data server info updated|

#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 500 - Internal Server Error

### GET /internal/v1/object/id

Returns a FileID allocated by Registration Center

#### Request 

- **Syntax**
```http
GET /internal/v1/object/id   HTTP/1.1
Host: <Arkor Host>
Date: date
```
#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
```
```json
{
  "file_id": "a50ec4ff5f0ea",
}
```
- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`fileID`|String|FileID allocated by Registration Center, the allocation rule see [FileID.md](FileID.md)|



#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Object NOT found)
* 409 - Conflict (Object already deleted)
* 500 - Internal Server Error


### PUT /internal/v1/object/info

Save/Update object info

#### Request 

- **Syntax**
```http
PUT /internal/v1/dataserver   HTTP/1.1
Host: <Arkor Host>
Date: date
```

```json
{
  "object_id": "b06bd",
  "object_key": "test.tar",
  "md5_key": "f1ddfa66735ab83451b8ec3bf9c0fe46",
  "fragments": [
    {
      "index": 0,
      "start": 0,
      "end": 4194304,
      "group_id": 1,
      "fileId": 240001,
      "isLast": false,
      "modTime": "0001-01-01T00:00:00Z"
    },
    {
      "index": 1,
      "start": 4194304,
      "end": 8388608,
      "group_id": 1,
      "fileId": 240002,
      "isLast": false,
      "modTime": "0001-01-01T00:00:00Z"
    }
  ]
}
```
- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`object_id`|String|ID of the object|
|`object_key`|String|Key(name) of the object|
|`md5_key`|String|MD5 value of the object key|
|`fragments`|container|Fragements info of the object|
|`index`|Int|The sequence number of the fragment|
|`start`|Int64|Where this fragement range starts|
|`end`|Int64|Where this fragement range ends|
|`group_id`|String|Group ID of which the fragment is located|
|`isLast`|Boolean|If this fragement is the last fragement|
|`modTime`|Date|Date of modification|


#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
```

#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Object NOT found)
* 409 - Conflict (Object already deleted)
* 500 - Internal Server Error

### GET /internal/v1/object/:object

Query object info

#### Request 

- **Syntax**
```http
GET /internal/v1/object/:object   HTTP/1.1
Host: <Arkor Host>
Date: date
```
- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`object`|URL|Object ID to be query|

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
Date: date
```
```json
{
  "object_id": "b06bd",
  "object_key": "test.tar",
  "md5_key": "f1ddfa66735ab83451b8ec3bf9c0fe46",
  "fragments": [
    {
      "index": 0,
      "start": 0,
      "end": 4194304,
      "group_id": 1,
      "fileId": 240001,
      "isLast": false,
      "modTime": "0001-01-01T00:00:00Z"
    },
    {
      "index": 1,
      "start": 4194304,
      "end": 8388608,
      "group_id": 1,
      "fileId": 240002,
      "isLast": false,
      "modTime": "0001-01-01T00:00:00Z"
    }
  ]
}
```
- **Elements:**

|Name|Type|Description|
|----|----|-----------|
|`object_id`|String|ID of the object|
|`object_key`|String|Key(name) of the object|
|`md5_key`|String|MD5 value of the object key|
|`fragments`|container|Fragements info of the object|
|`index`|Int|The sequence number of the fragment|
|`start`|Int64|Where this fragement range starts|
|`end`|Int64|Where this fragement range ends|
|`group_id`|String|Group ID of which the fragment is located|
|`isLast`|Boolean|If this fragement is the last fragement|
|`modTime`|Date|Date of modification|


#### Response On Failure

* 400 - Bad Request (Invalid Parameters)
* 404 - Not Found (Object NOT found)
* 409 - Conflict (Object already deleted)
* 500 - Internal Server Error


