# StockApp

Barebones app and k8s deployment.

Run `kubectl create secret generic apikey --from-literal=apikey=$APIKEY` to create the apikey and then `make deploy`

To build `make build`, to build the docker container `make image`.



Improvements I would make for productionalising app
- Use http library with more features: logging, tracing etc
- Secret management for apikey, add a serviceaccount to allow app to auth and get creds
- Improve logging, add metrics endpoint
- Dashboards, Runbooks and Alerting- start with RED or Golden Signals and iterate as we go
- Make HA (topologyspread, more replicas, PDB)
- Template k8s manifests to allow to deploy into different envs with different paramaters (kustomize, helm, etc)
- Publish docs and schema output- switch to openapi spec to generate handlers and make this easier to manage/support
- TLS! Everything is unencrypted- generate certs for e2e communications (or use a service mesh)
- Use a better json parsing library- unmarshalling json into maps means i had to double handle the array to get the last nDays- though bonus points for the challenge of using only go stdlib? :D
- Caching the data: either store that data in an external cache or a DB so that we dont have to query the service directly- and if the container goes up or down it can quickly retrieve the cache rather than store the data in-memory.
- Separation of Concerns: To the above point also separate the data retrieval and serving into different services- that way the data retrieval is abstracted by the way we serve the data and we can scale those independently. If we decide to use a different data source in the future (for example because of licensing cost) we can jus swap out the retrieval service.
- I realise my implementation means that container needs to be restarted everyday to get the last nDays of prices- if we didnt split the data retrieval and service components, I would change this to use a channel or mutex system where the app starts and the retrieval and service pieces operate independently of each other on different Goroutines, the data retrieval can loop every hour or so and see if the data has been updated upstream, RW lock a mutex, update the data and then the server side can have an RO lock and serve
