# bidl
`bild downloader = bidl`

A fast, multi-threaded downloader for images. Provides an interface for creating handlers for websites.

## Adding new handlers 

Routing of site to image urls is done by downloaders.GetImUrls, which refers to rules inited in InitRules.
Main route is based on host, then rest of the url is matched with a regex with named groups, which then are zipped to a map, passed to the provided handler.
