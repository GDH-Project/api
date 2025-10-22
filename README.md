# GDH API 서버 입니다.

## 환경변수 
```dotenv
AUTH_GRPC_SERVER="서버주소:포트"

DB_URL="데이터베이스주소"

# proxy 제한 위해 필요
HOST_URL="배포주소"

# cors 설정을 위한 호스트 주소 다중 호스트의 경우 , 로 구분한다.
CORS_HOST_LIST="http://localhost:3000,http://localhost:5173"
```