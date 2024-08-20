CORE:

backend -- camps at gateway, communicates (auth -> key change) with store
store -- encrypted at rest access req keys, within internal network

desktop + mobile app frontend:
keys stored locally, used for encrypt + decrypt of blobs to/from backend
