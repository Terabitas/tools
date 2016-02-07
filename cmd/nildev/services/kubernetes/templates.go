package kubernetes

var (
	SecreteTemaplate = `{
  "kind": "Secret",
  "apiVersion": "v1",
  "metadata": {
    "name": "{{.Name}}"
  },
  "data": {
  	{{range $j, $e := .Data}}
	"{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Data | not}},{{end}}
	{{end}}
  }
}`

	ReplicationControllerTemplate = `{
   "kind":"ReplicationController",
   "apiVersion":"v1",
   "metadata":{
      "name":"{{.Name}}",
      "labels":{
      	 {{range $j, $e := .Labels}}
		 "{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Labels | not}},{{end}}
	 	 {{end}}
      }
   },
   "spec":{
      "replicas":2,
      "selector":{
         {{range $j, $e := .Labels}}
		 "{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Labels | not}},{{end}}
	 	 {{end}}
      },
      "template":{
         "metadata":{
            "labels":{
               {{range $j, $e := .Labels}}
		       "{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Labels | not}},{{end}}
	 	       {{end}}
            }
         },
         "spec":{
            "volumes":[
               {{range $j, $e := .SecretVolumes}}
		       {
				 "name":"{{$e.Name}}",
				 "secret":{
				 	"secretName": "{{$e.Name}}"
				 }
			   },
	 	       {{end}}
               {
				 "name":"{{.EmptyDBVolume.Name}}",
				 "emptyDir":{}
			   }
            ],
            "containers":[
               {
                  "name":"{{.Name}}",
                  "image":"{{.Image}}",
                  "env":[
                  	 {{range $j, $e := .Env}}
					 {
					   "name":"{{$e.Key}}",
					   "value":"{{$e.Value}}"
					 }{{if last $j $.Env | not}},{{end}}
					 {{end}}
                  ],
                  "volumeMounts":[
                     {{range $j, $e := .SecretVolumes}}
					 {
						 "name":"{{$e.Name}}",
						 "mountPath":"{{$e.MountPath}}"
					 },
					 {{end}}
                     {
						 "name":"{{.EmptyDBVolume.Name}}",
						 "mountPath":"{{.EmptyDBVolume.MountPath}}"
					 }
                  ],
                  "ports":[
                     {
                        "name":"{{.Name}}",
                        "containerPort":8080
                     }
                  ]
               }
            ]
         }
      }
   }
}
`

	ServiceTemplate = `{
  "kind":"Service",
  "apiVersion":"v1",
  "metadata":{
    "name":"{{.Name}}",
    "labels":{
      {{range $j, $e := .Labels}}
	  "{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Labels | not}},{{end}}
	  {{end}}
    }
  },
  "spec":{
    "ports": [
      {
        "port":80,
        "targetPort":"{{.Name}}"
      }
    ],
    "selector":{
      {{range $j, $e := .Labels}}
	  "{{$e.Key}}":"{{$e.Value}}"{{if last $j $.Labels | not}},{{end}}
	  {{end}}
    },
    "externalIPs" : [
      {{range $j, $e := .ExternalIPs}}
	  "{{$e}}"{{if last $j $.ExternalIPs | not}},{{end}}
	  {{end}}
    ]
  }
}
`

	IngressTemplate = `apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: nildev
spec:
  rules:
  - host: empty.nildev.nil.services
    http:
      paths:
      - backend:
          serviceName: empty-nildev
          servicePort: 80
  - host: account.nildev.nil.services
    http:
      paths:
      - backend:
          serviceName: account-nildev
          servicePort: 80`
)
