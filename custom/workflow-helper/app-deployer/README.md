## About Git-App-Deployer
Git-App-Deployer has been used for the installation of applications.

At first, the user is asked to give the namespace, typeName, operation, app, scope, and timeout.

<table>
  <tr>
    <th> Environment </th>
    <th> Description </th>
    <th> Example </th>
  </tr>
  <tr>
    <td> -namespace=<"app-namespace"> </td>
    <td> Application namespace has to pass </td>
    <td> -namespace="sock-shop" </td>
  </tr>
  <tr>
    <td> -typeName=<"scenario"> </td>
    <td> It can be weak/resilient in apllication </br>
     Note: </br>
     1. In a weak scenario, it will create a single replica and Deployments for all. </br> 
     2. In a resilient scenario, it will create two replicas of pods with Statefulsets for databases and Deployments for others. </td>
     <td> -typeName="resilent"</td>
  </tr>
  <tr>
    <td> -timeout=<"timeout"> </td>
    <td> Timeout is used for the termination of application </td>
    <td> -timeout="400" </td>
  </tr>
  <tr>
    <td> -operation=<"operation"> </td>
    <td> Operation will be Kubernetes CRUD operations 
         on these resources.
    </td>
    <td> -operation="apply" </td>
  </tr>
  </tr>
    <td> -app=<"app-name"> </td>
    <td> App is the name of application, which will be used as labels </td>
    <td> -app="sock-shop" </td>
  </tr>
    <td> -scope=<"scope"> </td>
    <td> Scope for an application </td>
    <td> -scope="cluster" </td>
  </tr>
</table>
</br>

It creates a namespace and then installs the required application based on the given -namespace and -typeName.
If namespace already exists, then it shows log and starts installing the application.

`[Status]: Namespace already exists!`

## Load-Test:
The load test packages a test script in a container for Locust that simulates user traffic to the application. Please run it against the front-end service. The address and port of the frontend will be different and depend on which platform you've deployed to. See the notes for each deployment.
It has been used parallelly with a chaos engine, which loads against the catalog front-end service.
In the manifest, it is written as:
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=loadtest"] 
```

## Now let’s see how the app-deployer works in the workflow:
At first, the installation of App-Deployer(application installation) is performed.
```
- name: install-application
      container:
        image: litmuschaos/litmus-app-deployer:latest
        args: ["-namespace=sock-shop","-typeName=resilient","-operation=apply","-timeout=400", "-app=sock-shop","-scope=cluster"] 
```
 ##### Note: 
  - by default it will be resilient for weak provide type flagName as resilient(-typeName=weak)
  - To delete application or loadtest correspondence operation need to pass ```-operation=delete```

## Application details
As of now Sock-shop and Potato-Head application have been added in a pre-defined workflow.
Where weak and resilient cases are present.

### In a weak scenario, 
All the services of the application are deployments, in both sock-shop and podtato-head. 
Single replica for every service. 

#### For Databases:
MongoDB is run as a single replica deployment.
MySQL DB is hosted on ephemeral storage with a single replica.

Note: Only one replica of the pod is present. After chaos injection, it will be down, and therefore accessibility will not be there, and eventually, it will fail due to lack of resources availability.

### In a resilient scenario, 
All the databases(MongoDB and MySQL) services of the application are statefulsets and other services are deployments in sock-shop, whereas deployments for podtato-head. 
For each service two replicas are present. 

#### For Databases:
MongoDB multi-replica statefulset with persistent volumes. 
MySQL DB is Persistent, which can be dynamically extensible.

Note: Two replicas of pods are present. After chaos injection, one will be down. Therefore one pod is still up and accessibility. 