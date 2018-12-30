#!/usr/bin/env bash

channnel_ids=(UCtEorrVfo4GQsN82HsrnKyA UCu0r_ub9cXm7UaMMK8fSlrg)
addr=$1

for i in ${channnel_ids[@]}; do
	body='{"hub.callback":"'"$addr"'/pubsubhubbub/feed", "hub.topic":"https://www.youtube.com/xml/feeds/videos.xml?channel_id='"$i"'"}'
	url="$addr/pubsubhubbub/subscribe" 
	curl $url -XPOST -d "$body"
done
