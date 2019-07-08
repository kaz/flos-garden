resp=$(
	curl "http://localhost:9000/api/database/query" -s \
	-X POST \
	-u k4mu1:k4mu1 \
	-H "Content-Type: application/json" \
	-d "$(printf '"%s"' "$(cat PZ.sql | tr "\n" " "  | sed -E "s/ +/ /g")")"
)

for row in $(echo $resp | jq -c ".rows[]")
do
	row_raw=$(echo $row | jq -r ".[]")
	host=$(echo $row_raw | awk '{ print $1 }')
	id=$(echo $row_raw | awk '{ print $2 }')
	path=$(echo $row_raw | awk '{ print $3 }')

	mkdir -p $(dirname ".$path")
	echo Downloading: "$path"
	curl "http://localhost:9000/api/database/blob/$host/$id" -s -u k4mu1:k4mu1 -o ".$path"
done
