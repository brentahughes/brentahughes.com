APP_NAME=brentahughes

SITE_ROOT=""
RUN_FLAG="-d --restart=always"
if [ "$1" == "debug" ]; then
    RUN_FLAG="--rm"
    SITE_ROOT="-v ${PWD}/public:/www"
fi

echo "Building $APP_NAME image"
docker build -t $APP_NAME .

echo "Removing $APP_NAME container if it exists"
docker rm -f $APP_NAME

echo "Running $APP_NAME container"
docker run $RUN_FLAG --name $APP_NAME \
    -p 8800:80 \
    $SITE_ROOT \
    $APP_NAME