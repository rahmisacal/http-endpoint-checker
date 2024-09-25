##HttpResponseCheck

        Network problems experienced in applications running on cloud native infrastructures can be challenging compared to other problems. 

	I decided to use Golang for this small project because I knew it wouldn’t be a massive amount of code, and because it is the preferred language of our Infra people. The main constraint I had was to be limited only to Go’s standard library: no external dependencies.

        Here, a service that will run on kubernetes and openshift infrastructure, send a get request to the given endpoint url address every 15 seconds and print the http status code and response time will make the job easier.

        with this app, coredns, container network interface, ingress, egress, service and route objects can be tested and their performance can be observed from time to time with response time. the related infrastructure can be examined by assigning it to worker nodes suspected of having problems. it can be deployed to two different projects or clusters and the response time between the two can be compared.

        you can easily implement and test.