# Arkor - Object Storage Service

## Introduction 

Arkor is an object storage service system with high availability and high performance. It is best suited for storing blobs such as container/VM images, software packages, log files, videos, photos and so on. As a distributed system, Arkor can be easily scaled out by simply adding more servers horizontally.

## Features

* Compatible with [Amazon S3](http://aws.amazon.com/s3/)
* Compatible with [Git-LFS (Git Large File Storage)](https://git-lfs.github.com/)
* High Performance
	* Efficient data storage server written in C language
	* Split large file into small pieces and upload/download concurrently
	* Object server upload/download file directly to/from data server
* High Availability&Scalability
	* Stateless Object&Registry server
	* 3 replica data storage 
	* Dynamic scaled out

## Quick Start
	TODO

## Architecture

Arkor include three modules: Object Server Group, Registration Center and Data Groups:

![arch](docs/images/Arch.jpg "architecture")

### Object Server Group

`Object Server Group` is a statelless front-end proxy server group that designed to provide RESTFUL API for user to upload/download files and manage the file/metadata. When a user wants to upload/download a file, it retrives data server info and file storage info(e.g fileID) from  `Registration Center` and execute the upload/download operation directly with related `Data Server`.

### Registration Center

`Registration Center` maintains the server and file informations. It has three main functions: A. recive, store and send the basic information and status of all `Data Servers`. B. recive, store and send the metadata of files storaged in arkor. C. allocate fileID to `Object Server Group` when user wants to upload files.

### Data Groups

`Data Groups` are composed of many data groups. Each group contains multiple (usually three) highly optimized storage engine for performance and space efficiency called `Data Server`. It appends single small file to large file to keep the IO overhead at a low level. When user starts a UPLOAD process, the `Object Server Group` will upload the file to all `Data Servers` in the group. But in the DOWNLOAD process, the `Object Server Group` retrive file from any of it.  

## Developers Guide

See the document of [APIs](docs/specs/API.md), [FileID spec](docs/specs/FileID.md), [Internal APIs](docs/specs/Internal.md)

## License
The code is licensed under the MIT licence






