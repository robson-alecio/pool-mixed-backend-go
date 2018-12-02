timestamp=$(date +%s)
echo --$1 down > ./migrations/${timestamp}_$1.down.sql
echo --$1 up > ./migrations/${timestamp}_$1.up.sql