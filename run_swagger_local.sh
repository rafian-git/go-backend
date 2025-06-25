#!/bin/bash

echo "running swagger"
#
#docker stop swagger-local
#
#docker rm swagger-local
#
#docker run -d --name swagger-local -p 8080:8080 -e API_HOST="localhost:8080" -e URLS="[
#              { url: \"./api/auth.swagger.json\", name: \"Auth\" },
#              { url: \"./api/me.swagger.json\", name: \"Me\" },
#              { url: \"./api/sms-pusher.swagger.json\", name: \"SmsPusher\" },
#              { url: \"./api/admin-portal.swagger.json\", name: \"AdminPortal\" },
#                      ]" -v $(pwd)/integration/swagger/:/usr/share/nginx/html/api/ swaggerapi/swagger-ui



docker stop swagger-local; docker rm swagger-local;
docker run -d --name swagger-local -p 8280:8080 -e API_HOST="localhost:8080" -e URLS="[
              { url: \"./api/auth.swagger.json\", name: \"Auth\" },
              { url: \"./api/oms-auth.swagger.json\", name: \"OmsAuth\" },
              { url: \"./api/oms-user-management.swagger.json\", name: \"OmsUserManagement\" },
              { url: \"./api/oms-admin-portal.swagger.json\", name: \"OMSAdminPortal\" },
              { url: \"./api/admin-portal.swagger.json\", name: \"AdminPortal\" },
              { url: \"./api/order-manager.swagger.json\", name: \"OMSOrderManager\" },
              { url: \"./api/trade-capture.swagger.json\", name: \"OMSTradeCapture\" },
              { url: \"./api/order-executor.swagger.json\", name: \"OrderExecutor\" },
              { url: \"./api/market_data_feed.swagger.json\", name: \"MarketDataFeed\" },
              { url: \"./api/mercurius.swagger.json\", name: \"Mercurius\" },
                      ]" gcr.io/stock-x-342909/swagger:84ed03a