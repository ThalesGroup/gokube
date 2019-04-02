# GoKube Release Notes

## Version 1.5.0 - 02/08/2019
* Bump to helm-spray v3.2.0 and stern v1.10.0 [#6](https://github.com/gemalto/gokube/pull/6)

## Version 1.4.0 - 11/21/2018
* Bump to monocular v1.2.0 (that support http proxy)
* Bump to docker v18.06
* Bump to k8s v1.10.9
* Add docker image 'mongodb' in cache for monocular
* Add support of any-proxy (for transparent proxy)
* Add support of --cache-images
* Patch kubernetes-dashboard service to expose it on nodePort: 30000

## Version 1.3.0 - 10/15/2018
* Add option for using a specific version for Minikube
* Add option for using a specific version for Helm
* Add option for using a specific fork for Minikube

## Version 1.2.0 - 10/09/2018
* Updated Minikube to 0.30.0
* Updated Kubectl to 1.12.0
* Remove directories: '.kube' and '.docker' before upgrading gokube
