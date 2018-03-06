for GOOS in darwin linux openbsd windows; do
    EXT=""
    for GOARCH in 386 amd64; do
        echo "Compiling $GOOS-$GOARCH..."
        if [ "$GOOS" == "windows" ]; then
            EXT=".exe"
        fi
        GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/pgpi-$GOOS-$GOARCH$EXT
    done
done