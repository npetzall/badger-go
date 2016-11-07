## Badger-go

Making badges in go based of https://github.com/odino/nodejs-badges  
and the blog-post http://odino.org/generating-badges-slash-shields-with-nodejs/  

Shield/badge template is from https://github.com/odino/nodejs-badges    
which is inspired from shields.io's shield.  

### Latest version in artifactory

`/artifactory/latestVersion?g=[groupId]&a=[artifactId]&repos=[repoId]`

g = groupId, a= artifactId, repos= repository

### Latest version in nexus


### Sonarqube overall coverage (is going to be version specific)
`/sonarqube/overall-coverage?id=[componentId/projectId]`

id = projectId/componentId not really sure

## TODO

### Bamboo buildStatus
`/bamboo/buildstatus?planKey=[planKey]`

planKey =  [the plankey in bamboo]

### Jenkins buildStatus
`/jenkins/buildstatus?`
