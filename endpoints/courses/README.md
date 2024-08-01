# Endpoints for courses

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik"}%
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik

## Upload a large file
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/chunk/upload \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "chunkNumber=0" \
  -F "totalChunks=1" \
  -F "file=@/home/agust/Videos/test_edit.mp4"
```

## Response for Upload large file
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/1cd35233-c3c5-4d6f-b155-1c1a74a1a133/test.mkv"}%
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"}%
```

## Create course
```bash
# this endpoint need to be authenticated as admin
curl -X POST http://localhost:8081/courses/create \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "title=Data Structures" \
  -F "description=Description for Data Structures" \
  -F "author=agustfricke" \
  -F "duration=4 hours, 20 minutes" \
  -F "is_active=true" \
  -F "thumbnail=@/home/agust/Pictures/test.png" \
  -F "preview_tmp=/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"
```

## Get courses
```bash
curl -X GET "http://localhost:8081/courses/?cursor=0&q="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" | jq
```

## Get admin courses
```bash
# this endpoint need to be authenticated as admin
curl -X GET "http://localhost:8081/courses/admin?cursor=0&q=&a="  \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" | jq
```

## Add course to user
```bash
# this endpoint need to be authenticated as admin
curl -X PUT "http://localhost:8081/courses/add/user" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ2ODI3MzAsImlhdCI6MTcyMjA5MDczMCwibmJmIjoxNzIyMDkwNzMwLCJzdWIiOjF9.FPOoBntSbQNs8klEuNOYzGD-XRB07buGMACGofcK7mY" \
  -d '{"user_id": 1, "course_id": 1}'
```

## Update course
First you need to upload the video to http://localhost:8081/courses/chunk/upload and get 
the response that looks like this
```json
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/cc557fd0-239b-48f1-8f0c-a6d297bfca9c/test.mp4"}%
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/762716dc-f053-4150-86c0-e6678c0e263a/test_edit.mp4"}%
```

```bash
# this endpoint need to be authenticated as admin
curl -X PUT http://localhost:8081/courses/update \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjUxMjQ0OTYsImlhdCI6MTcyMjUzMjQ5NiwibmJmIjoxNzIyNTMyNDk2LCJzdWIiOjF9.ENH-zsDg-s1Z4aKOMP6tnV7Wg91-qaRJHlXvKhc_Uik" \
  -H "Content-Type: multipart/form-data" \
  -F "id=1" \
  -F "sort_order=4" \
  -F "title=Data Structures updated" \
  -F "description=Description for Data Structures updated" \
  -F "author=agustfricke updated" \
  -F "duration=6 hours, 09 minutes" \
  -F "is_active=false" \
  -F "thumbnail=@/home/agust/Pictures/edit.png" \
  -F "previewTmpDir=/home/agust/work/iron-stream/backend/web/uploads/tmp/762716dc-f053-4150-86c0-e6678c0e263a/test_edit.mp4"
```


.
├── previews
│   ├── 1
│   │   ├── master-1080p0.ts
│   │   ├── master-1080p10.ts
│   │   ├── master-1080p11.ts
│   │   ├── master-1080p12.ts
│   │   ├── master-1080p13.ts
│   │   ├── master-1080p14.ts
│   │   ├── master-1080p15.ts
│   │   ├── master-1080p16.ts
│   │   ├── master-1080p17.ts
│   │   ├── master-1080p18.ts
│   │   ├── master-1080p19.ts
│   │   ├── master-1080p1.ts
│   │   ├── master-1080p20.ts
│   │   ├── master-1080p21.ts
│   │   ├── master-1080p22.ts
│   │   ├── master-1080p23.ts
│   │   ├── master-1080p24.ts
│   │   ├── master-1080p2.ts
│   │   ├── master-1080p3.ts
│   │   ├── master-1080p4.ts
│   │   ├── master-1080p5.ts
│   │   ├── master-1080p6.ts
│   │   ├── master-1080p7.ts
│   │   ├── master-1080p8.ts
│   │   ├── master-1080p9.ts
│   │   ├── master-1080p.m3u8
│   │   ├── master-360p0.ts
│   │   ├── master-360p10.ts
│   │   ├── master-360p11.ts
│   │   ├── master-360p12.ts
│   │   ├── master-360p13.ts
│   │   ├── master-360p14.ts
│   │   ├── master-360p15.ts
│   │   ├── master-360p16.ts
│   │   ├── master-360p17.ts
│   │   ├── master-360p18.ts
│   │   ├── master-360p19.ts
│   │   ├── master-360p1.ts
│   │   ├── master-360p20.ts
│   │   ├── master-360p21.ts
│   │   ├── master-360p22.ts
│   │   ├── master-360p23.ts
│   │   ├── master-360p24.ts
│   │   ├── master-360p2.ts
│   │   ├── master-360p3.ts
│   │   ├── master-360p4.ts
│   │   ├── master-360p5.ts
│   │   ├── master-360p6.ts
│   │   ├── master-360p7.ts
│   │   ├── master-360p8.ts
│   │   ├── master-360p9.ts
│   │   ├── master-360p.m3u8
│   │   ├── master-480p0.ts
│   │   ├── master-480p10.ts
│   │   ├── master-480p11.ts
│   │   ├── master-480p12.ts
│   │   ├── master-480p13.ts
│   │   ├── master-480p14.ts
│   │   ├── master-480p15.ts
│   │   ├── master-480p16.ts
│   │   ├── master-480p17.ts
│   │   ├── master-480p18.ts
│   │   ├── master-480p19.ts
│   │   ├── master-480p1.ts
│   │   ├── master-480p20.ts
│   │   ├── master-480p21.ts
│   │   ├── master-480p22.ts
│   │   ├── master-480p23.ts
│   │   ├── master-480p24.ts
│   │   ├── master-480p2.ts
│   │   ├── master-480p3.ts
│   │   ├── master-480p4.ts
│   │   ├── master-480p5.ts
│   │   ├── master-480p6.ts
│   │   ├── master-480p7.ts
│   │   ├── master-480p8.ts
│   │   ├── master-480p9.ts
│   │   ├── master-480p.m3u8
│   │   ├── master-720p0.ts
│   │   ├── master-720p10.ts
│   │   ├── master-720p11.ts
│   │   ├── master-720p12.ts
│   │   ├── master-720p13.ts
│   │   ├── master-720p14.ts
│   │   ├── master-720p15.ts
│   │   ├── master-720p16.ts
│   │   ├── master-720p17.ts
│   │   ├── master-720p18.ts
│   │   ├── master-720p19.ts
│   │   ├── master-720p1.ts
│   │   ├── master-720p20.ts
│   │   ├── master-720p21.ts
│   │   ├── master-720p22.ts
│   │   ├── master-720p23.ts
│   │   ├── master-720p24.ts
│   │   ├── master-720p2.ts
│   │   ├── master-720p3.ts
│   │   ├── master-720p4.ts
│   │   ├── master-720p5.ts
│   │   ├── master-720p6.ts
│   │   ├── master-720p7.ts
│   │   ├── master-720p8.ts
│   │   ├── master-720p9.ts
│   │   ├── master-720p.m3u8
│   │   └── master.m3u8
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a
│       ├── master-1080p0.ts
│       ├── master-1080p10.ts
│       ├── master-1080p11.ts
│       ├── master-1080p12.ts
│       ├── master-1080p13.ts
│       ├── master-1080p14.ts
│       ├── master-1080p15.ts
│       ├── master-1080p16.ts
│       ├── master-1080p17.ts
│       ├── master-1080p18.ts
│       ├── master-1080p19.ts
│       ├── master-1080p1.ts
│       ├── master-1080p20.ts
│       ├── master-1080p21.ts
│       ├── master-1080p22.ts
│       ├── master-1080p23.ts
│       ├── master-1080p24.ts
│       ├── master-1080p25.ts
│       ├── master-1080p26.ts
│       ├── master-1080p27.ts
│       ├── master-1080p28.ts
│       ├── master-1080p29.ts
│       ├── master-1080p2.ts
│       ├── master-1080p3.ts
│       ├── master-1080p4.ts
│       ├── master-1080p5.ts
│       ├── master-1080p6.ts
│       ├── master-1080p7.ts
│       ├── master-1080p8.ts
│       ├── master-1080p9.ts
│       ├── master-1080p.m3u8
│       ├── master-360p0.ts
│       ├── master-360p10.ts
│       ├── master-360p11.ts
│       ├── master-360p12.ts
│       ├── master-360p13.ts
│       ├── master-360p14.ts
│       ├── master-360p15.ts
│       ├── master-360p16.ts
│       ├── master-360p17.ts
│       ├── master-360p18.ts
│       ├── master-360p19.ts
│       ├── master-360p1.ts
│       ├── master-360p20.ts
│       ├── master-360p21.ts
│       ├── master-360p22.ts
│       ├── master-360p23.ts
│       ├── master-360p24.ts
│       ├── master-360p25.ts
│       ├── master-360p26.ts
│       ├── master-360p27.ts
│       ├── master-360p28.ts
│       ├── master-360p29.ts
│       ├── master-360p2.ts
│       ├── master-360p3.ts
│       ├── master-360p4.ts
│       ├── master-360p5.ts
│       ├── master-360p6.ts
│       ├── master-360p7.ts
│       ├── master-360p8.ts
│       ├── master-360p9.ts
│       ├── master-360p.m3u8
│       ├── master-480p0.ts
│       ├── master-480p10.ts
│       ├── master-480p11.ts
│       ├── master-480p12.ts
│       ├── master-480p13.ts
│       ├── master-480p14.ts
│       ├── master-480p15.ts
│       ├── master-480p16.ts
│       ├── master-480p17.ts
│       ├── master-480p18.ts
│       ├── master-480p19.ts
│       ├── master-480p1.ts
│       ├── master-480p20.ts
│       ├── master-480p21.ts
│       ├── master-480p22.ts
│       ├── master-480p23.ts
│       ├── master-480p24.ts
│       ├── master-480p25.ts
│       ├── master-480p26.ts
│       ├── master-480p27.ts
│       ├── master-480p28.ts
│       ├── master-480p29.ts
│       ├── master-480p2.ts
│       ├── master-480p3.ts
│       ├── master-480p4.ts
│       ├── master-480p5.ts
│       ├── master-480p6.ts
│       ├── master-480p7.ts
│       ├── master-480p8.ts
│       ├── master-480p9.ts
│       ├── master-480p.m3u8
│       ├── master-720p0.ts
│       ├── master-720p10.ts
│       ├── master-720p11.ts
│       ├── master-720p12.ts
│       ├── master-720p13.ts
│       ├── master-720p14.ts
│       ├── master-720p15.ts
│       ├── master-720p16.ts
│       ├── master-720p17.ts
│       ├── master-720p18.ts
│       ├── master-720p19.ts
│       ├── master-720p1.ts
│       ├── master-720p20.ts
│       ├── master-720p21.ts
│       ├── master-720p22.ts
│       ├── master-720p23.ts
│       ├── master-720p24.ts
│       ├── master-720p25.ts
│       ├── master-720p26.ts
│       ├── master-720p27.ts
│       ├── master-720p28.ts
│       ├── master-720p29.ts
│       ├── master-720p2.ts
│       ├── master-720p3.ts
│       ├── master-720p4.ts
│       ├── master-720p5.ts
│       ├── master-720p6.ts
│       ├── master-720p7.ts
│       ├── master-720p8.ts
│       ├── master-720p9.ts
│       ├── master-720p.m3u8
│       └── master.m3u8
├── thumbnails
│   ├── 25590623-ff44-4338-8cba-5ce4aa024d78.png
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a.png
└── tmp
    ├── 762716dc-f053-4150-86c0-e6678c0e263a
    │   └── test_edit.mp4
    └── cc557fd0-239b-48f1-8f0c-a6d297bfca9c
        └── test.mp4

7 directories, 234 files

.
├── previews
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a
│       ├── master-1080p0.ts
│       ├── master-1080p10.ts
│       ├── master-1080p11.ts
│       ├── master-1080p12.ts
│       ├── master-1080p13.ts
│       ├── master-1080p14.ts
│       ├── master-1080p15.ts
│       ├── master-1080p16.ts
│       ├── master-1080p17.ts
│       ├── master-1080p18.ts
│       ├── master-1080p19.ts
│       ├── master-1080p1.ts
│       ├── master-1080p20.ts
│       ├── master-1080p21.ts
│       ├── master-1080p22.ts
│       ├── master-1080p23.ts
│       ├── master-1080p24.ts
│       ├── master-1080p25.ts
│       ├── master-1080p26.ts
│       ├── master-1080p27.ts
│       ├── master-1080p28.ts
│       ├── master-1080p29.ts
│       ├── master-1080p2.ts
│       ├── master-1080p3.ts
│       ├── master-1080p4.ts
│       ├── master-1080p5.ts
│       ├── master-1080p6.ts
│       ├── master-1080p7.ts
│       ├── master-1080p8.ts
│       ├── master-1080p9.ts
│       ├── master-1080p.m3u8
│       ├── master-360p0.ts
│       ├── master-360p10.ts
│       ├── master-360p11.ts
│       ├── master-360p12.ts
│       ├── master-360p13.ts
│       ├── master-360p14.ts
│       ├── master-360p15.ts
│       ├── master-360p16.ts
│       ├── master-360p17.ts
│       ├── master-360p18.ts
│       ├── master-360p19.ts
│       ├── master-360p1.ts
│       ├── master-360p20.ts
│       ├── master-360p21.ts
│       ├── master-360p22.ts
│       ├── master-360p23.ts
│       ├── master-360p24.ts
│       ├── master-360p25.ts
│       ├── master-360p26.ts
│       ├── master-360p27.ts
│       ├── master-360p28.ts
│       ├── master-360p29.ts
│       ├── master-360p2.ts
│       ├── master-360p3.ts
│       ├── master-360p4.ts
│       ├── master-360p5.ts
│       ├── master-360p6.ts
│       ├── master-360p7.ts
│       ├── master-360p8.ts
│       ├── master-360p9.ts
│       ├── master-360p.m3u8
│       ├── master-480p0.ts
│       ├── master-480p10.ts
│       ├── master-480p11.ts
│       ├── master-480p12.ts
│       ├── master-480p13.ts
│       ├── master-480p14.ts
│       ├── master-480p15.ts
│       ├── master-480p16.ts
│       ├── master-480p17.ts
│       ├── master-480p18.ts
│       ├── master-480p19.ts
│       ├── master-480p1.ts
│       ├── master-480p20.ts
│       ├── master-480p21.ts
│       ├── master-480p22.ts
│       ├── master-480p23.ts
│       ├── master-480p24.ts
│       ├── master-480p25.ts
│       ├── master-480p26.ts
│       ├── master-480p27.ts
│       ├── master-480p28.ts
│       ├── master-480p29.ts
│       ├── master-480p2.ts
│       ├── master-480p3.ts
│       ├── master-480p4.ts
│       ├── master-480p5.ts
│       ├── master-480p6.ts
│       ├── master-480p7.ts
│       ├── master-480p8.ts
│       ├── master-480p9.ts
│       ├── master-480p.m3u8
│       ├── master-720p0.ts
│       ├── master-720p10.ts
│       ├── master-720p11.ts
│       ├── master-720p12.ts
│       ├── master-720p13.ts
│       ├── master-720p14.ts
│       ├── master-720p15.ts
│       ├── master-720p16.ts
│       ├── master-720p17.ts
│       ├── master-720p18.ts
│       ├── master-720p19.ts
│       ├── master-720p1.ts
│       ├── master-720p20.ts
│       ├── master-720p21.ts
│       ├── master-720p22.ts
│       ├── master-720p23.ts
│       ├── master-720p24.ts
│       ├── master-720p25.ts
│       ├── master-720p26.ts
│       ├── master-720p27.ts
│       ├── master-720p28.ts
│       ├── master-720p29.ts
│       ├── master-720p2.ts
│       ├── master-720p3.ts
│       ├── master-720p4.ts
│       ├── master-720p5.ts
│       ├── master-720p6.ts
│       ├── master-720p7.ts
│       ├── master-720p8.ts
│       ├── master-720p9.ts
│       ├── master-720p.m3u8
│       └── master.m3u8
├── thumbnails
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a.png
└── tmp
    └── cc557fd0-239b-48f1-8f0c-a6d297bfca9c
        └── test.mp4

5 directories, 127 files
{"message":"Archivo cargado con éxito","path":"/home/agust/work/iron-stream/backend/web/uploads/tmp/762716dc-f053-4150-86c0-e6678c0e263a/test_edit.mp4"}%
.
├── previews
│   ├── 1
│   │   ├── master-1080p0.ts
│   │   ├── master-1080p1.ts
│   │   ├── master-1080p.m3u8
│   │   ├── master-360p0.ts
│   │   ├── master-360p1.ts
│   │   ├── master-360p.m3u8
│   │   ├── master-480p0.ts
│   │   ├── master-480p1.ts
│   │   ├── master-480p.m3u8
│   │   ├── master-720p0.ts
│   │   ├── master-720p1.ts
│   │   ├── master-720p.m3u8
│   │   └── master.m3u8
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a
│       ├── master-1080p0.ts
│       ├── master-1080p10.ts
│       ├── master-1080p11.ts
│       ├── master-1080p12.ts
│       ├── master-1080p13.ts
│       ├── master-1080p14.ts
│       ├── master-1080p15.ts
│       ├── master-1080p16.ts
│       ├── master-1080p17.ts
│       ├── master-1080p18.ts
│       ├── master-1080p19.ts
│       ├── master-1080p1.ts
│       ├── master-1080p20.ts
│       ├── master-1080p21.ts
│       ├── master-1080p22.ts
│       ├── master-1080p23.ts
│       ├── master-1080p24.ts
│       ├── master-1080p25.ts
│       ├── master-1080p26.ts
│       ├── master-1080p27.ts
│       ├── master-1080p28.ts
│       ├── master-1080p29.ts
│       ├── master-1080p2.ts
│       ├── master-1080p3.ts
│       ├── master-1080p4.ts
│       ├── master-1080p5.ts
│       ├── master-1080p6.ts
│       ├── master-1080p7.ts
│       ├── master-1080p8.ts
│       ├── master-1080p9.ts
│       ├── master-1080p.m3u8
│       ├── master-360p0.ts
│       ├── master-360p10.ts
│       ├── master-360p11.ts
│       ├── master-360p12.ts
│       ├── master-360p13.ts
│       ├── master-360p14.ts
│       ├── master-360p15.ts
│       ├── master-360p16.ts
│       ├── master-360p17.ts
│       ├── master-360p18.ts
│       ├── master-360p19.ts
│       ├── master-360p1.ts
│       ├── master-360p20.ts
│       ├── master-360p21.ts
│       ├── master-360p22.ts
│       ├── master-360p23.ts
│       ├── master-360p24.ts
│       ├── master-360p25.ts
│       ├── master-360p26.ts
│       ├── master-360p27.ts
│       ├── master-360p28.ts
│       ├── master-360p29.ts
│       ├── master-360p2.ts
│       ├── master-360p3.ts
│       ├── master-360p4.ts
│       ├── master-360p5.ts
│       ├── master-360p6.ts
│       ├── master-360p7.ts
│       ├── master-360p8.ts
│       ├── master-360p9.ts
│       ├── master-360p.m3u8
│       ├── master-480p0.ts
│       ├── master-480p10.ts
│       ├── master-480p11.ts
│       ├── master-480p12.ts
│       ├── master-480p13.ts
│       ├── master-480p14.ts
│       ├── master-480p15.ts
│       ├── master-480p16.ts
│       ├── master-480p17.ts
│       ├── master-480p18.ts
│       ├── master-480p19.ts
│       ├── master-480p1.ts
│       ├── master-480p20.ts
│       ├── master-480p21.ts
│       ├── master-480p22.ts
│       ├── master-480p23.ts
│       ├── master-480p24.ts
│       ├── master-480p25.ts
│       ├── master-480p26.ts
│       ├── master-480p27.ts
│       ├── master-480p28.ts
│       ├── master-480p29.ts
│       ├── master-480p2.ts
│       ├── master-480p3.ts
│       ├── master-480p4.ts
│       ├── master-480p5.ts
│       ├── master-480p6.ts
│       ├── master-480p7.ts
│       ├── master-480p8.ts
│       ├── master-480p9.ts
│       ├── master-480p.m3u8
│       ├── master-720p0.ts
│       ├── master-720p10.ts
│       ├── master-720p11.ts
│       ├── master-720p12.ts
│       ├── master-720p13.ts
│       ├── master-720p14.ts
│       ├── master-720p15.ts
│       ├── master-720p16.ts
│       ├── master-720p17.ts
│       ├── master-720p18.ts
│       ├── master-720p19.ts
│       ├── master-720p1.ts
│       ├── master-720p20.ts
│       ├── master-720p21.ts
│       ├── master-720p22.ts
│       ├── master-720p23.ts
│       ├── master-720p24.ts
│       ├── master-720p25.ts
│       ├── master-720p26.ts
│       ├── master-720p27.ts
│       ├── master-720p28.ts
│       ├── master-720p29.ts
│       ├── master-720p2.ts
│       ├── master-720p3.ts
│       ├── master-720p4.ts
│       ├── master-720p5.ts
│       ├── master-720p6.ts
│       ├── master-720p7.ts
│       ├── master-720p8.ts
│       ├── master-720p9.ts
│       ├── master-720p.m3u8
│       └── master.m3u8
├── thumbnails
│   ├── 25590623-ff44-4338-8cba-5ce4aa024d78.png
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a.png
└── tmp
    ├── 762716dc-f053-4150-86c0-e6678c0e263a
    │   └── test_edit.mp4
    └── cc557fd0-239b-48f1-8f0c-a6d297bfca9c
        └── test.mp4

7 directories, 142 files

.
├── previews
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a
│       ├── master-1080p0.ts
│       ├── master-1080p10.ts
│       ├── master-1080p11.ts
│       ├── master-1080p12.ts
│       ├── master-1080p13.ts
│       ├── master-1080p14.ts
│       ├── master-1080p15.ts
│       ├── master-1080p16.ts
│       ├── master-1080p17.ts
│       ├── master-1080p18.ts
│       ├── master-1080p19.ts
│       ├── master-1080p1.ts
│       ├── master-1080p20.ts
│       ├── master-1080p21.ts
│       ├── master-1080p22.ts
│       ├── master-1080p23.ts
│       ├── master-1080p24.ts
│       ├── master-1080p25.ts
│       ├── master-1080p26.ts
│       ├── master-1080p27.ts
│       ├── master-1080p28.ts
│       ├── master-1080p29.ts
│       ├── master-1080p2.ts
│       ├── master-1080p3.ts
│       ├── master-1080p4.ts
│       ├── master-1080p5.ts
│       ├── master-1080p6.ts
│       ├── master-1080p7.ts
│       ├── master-1080p8.ts
│       ├── master-1080p9.ts
│       ├── master-1080p.m3u8
│       ├── master-360p0.ts
│       ├── master-360p10.ts
│       ├── master-360p11.ts
│       ├── master-360p12.ts
│       ├── master-360p13.ts
│       ├── master-360p14.ts
│       ├── master-360p15.ts
│       ├── master-360p16.ts
│       ├── master-360p17.ts
│       ├── master-360p18.ts
│       ├── master-360p19.ts
│       ├── master-360p1.ts
│       ├── master-360p20.ts
│       ├── master-360p21.ts
│       ├── master-360p22.ts
│       ├── master-360p23.ts
│       ├── master-360p24.ts
│       ├── master-360p25.ts
│       ├── master-360p26.ts
│       ├── master-360p27.ts
│       ├── master-360p28.ts
│       ├── master-360p29.ts
│       ├── master-360p2.ts
│       ├── master-360p3.ts
│       ├── master-360p4.ts
│       ├── master-360p5.ts
│       ├── master-360p6.ts
│       ├── master-360p7.ts
│       ├── master-360p8.ts
│       ├── master-360p9.ts
│       ├── master-360p.m3u8
│       ├── master-480p0.ts
│       ├── master-480p10.ts
│       ├── master-480p11.ts
│       ├── master-480p12.ts
│       ├── master-480p13.ts
│       ├── master-480p14.ts
│       ├── master-480p15.ts
│       ├── master-480p16.ts
│       ├── master-480p17.ts
│       ├── master-480p18.ts
│       ├── master-480p19.ts
│       ├── master-480p1.ts
│       ├── master-480p20.ts
│       ├── master-480p21.ts
│       ├── master-480p22.ts
│       ├── master-480p23.ts
│       ├── master-480p24.ts
│       ├── master-480p25.ts
│       ├── master-480p26.ts
│       ├── master-480p27.ts
│       ├── master-480p28.ts
│       ├── master-480p29.ts
│       ├── master-480p2.ts
│       ├── master-480p3.ts
│       ├── master-480p4.ts
│       ├── master-480p5.ts
│       ├── master-480p6.ts
│       ├── master-480p7.ts
│       ├── master-480p8.ts
│       ├── master-480p9.ts
│       ├── master-480p.m3u8
│       ├── master-720p0.ts
│       ├── master-720p10.ts
│       ├── master-720p11.ts
│       ├── master-720p12.ts
│       ├── master-720p13.ts
│       ├── master-720p14.ts
│       ├── master-720p15.ts
│       ├── master-720p16.ts
│       ├── master-720p17.ts
│       ├── master-720p18.ts
│       ├── master-720p19.ts
│       ├── master-720p1.ts
│       ├── master-720p20.ts
│       ├── master-720p21.ts
│       ├── master-720p22.ts
│       ├── master-720p23.ts
│       ├── master-720p24.ts
│       ├── master-720p25.ts
│       ├── master-720p26.ts
│       ├── master-720p27.ts
│       ├── master-720p28.ts
│       ├── master-720p29.ts
│       ├── master-720p2.ts
│       ├── master-720p3.ts
│       ├── master-720p4.ts
│       ├── master-720p5.ts
│       ├── master-720p6.ts
│       ├── master-720p7.ts
│       ├── master-720p8.ts
│       ├── master-720p9.ts
│       ├── master-720p.m3u8
│       └── master.m3u8
├── thumbnails
│   └── 3867a629-fab2-4d6a-8d54-89cbe86a1c1a.png
└── tmp
    └── cc557fd0-239b-48f1-8f0c-a6d297bfca9c
        └── test.mp4

5 directories, 127 files
